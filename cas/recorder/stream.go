package recorder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/malt3/abstractfs-core/sri"
)

// Encode encodes content from a cas emitter.
// The contents are encoded in the following order:
// sri, payload
// The sri is encoded as follows:
// - 1 byte: type of record (0x01)
// - 8 byte: length of record
// - length bytes: sri
// The payload is encoded as follows:
// - 1 byte: type of record (0x02)
// - 8 byte: length of record
// - length bytes: payload
func Encode(w io.Writer, sri sri.Integrity, size int64, payload io.Reader) error {
	if err := encodeSRI(w, sri); err != nil {
		return err
	}
	if err := encodePayload(w, size, payload); err != nil {
		return err
	}
	return nil
}

// Decode decodes content for a cas recorder.
// The contents are decoded in the following order:
// sri, payload
// The components are expected to be encoded as described in Encode.
// The caller is free to skip the payload if it is not needed (i.e when the sri was recorded previously).
// The caller must close the returned body before calling Decode again.
func Decode(r io.Reader) (sri sri.Integrity, body io.ReadCloser, err error) {
	sri, err = decodeSRI(r)
	if err != nil {
		return sri, nil, err
	}
	body, err = decodePayload(r)
	if err != nil {
		return sri, nil, err
	}
	return sri, body, nil
}

func encodeSRI(w io.Writer, sri sri.Integrity) error {
	rawSRI := []byte(sri.String())
	return encodeTLV(w, typeSRI, int64(len(rawSRI)), bytes.NewReader(rawSRI))
}

func encodePayload(w io.Writer, size int64, payload io.Reader) error {
	return encodeTLV(w, typePayload, size, payload)
}

func encodeTLV(w io.Writer, t byte, l int64, v io.Reader) error {
	if err := binary.Write(w, binary.BigEndian, t); err != nil {
		return fmt.Errorf("encoding type: %w", err)
	}
	if err := binary.Write(w, binary.BigEndian, l); err != nil {
		return fmt.Errorf("encoding length: %w", err)
	}
	if _, err := io.CopyN(w, v, l); err != nil {
		return fmt.Errorf("encoding value: %w", err)
	}
	return nil
}

func decodeSRI(r io.Reader) (sri.Integrity, error) {
	lr, err := decodeTLV(r, typeSRI)
	if err != nil {
		return sri.Integrity{}, err
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, lr); err != nil {
		return sri.Integrity{}, fmt.Errorf("reading sri: %w", err)
	}
	return sri.FromString(buf.String())
}

func decodePayload(r io.Reader) (io.ReadCloser, error) {
	lr, err := decodeTLV(r, typePayload)
	if err != nil {
		return nil, err
	}
	return &limitReadCloser{r: lr}, nil
}

func decodeTLV(r io.Reader, expectType byte) (*io.LimitedReader, error) {
	t, l, err := decodeTL(r)
	if err != nil {
		return nil, err
	}
	if t != expectType {
		return nil, fmt.Errorf("expected type %d, got %d", expectType, t)
	}
	return &io.LimitedReader{R: r, N: l}, nil
}

func decodeTL(r io.Reader) (t byte, l int64, err error) {
	if err := binary.Read(r, binary.BigEndian, &t); err != nil {
		if err == io.EOF {
			return 0, 0, io.EOF
		}
		return 0, 0, fmt.Errorf("decoding type: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &l); err != nil {
		return 0, 0, fmt.Errorf("decoding length: %w", err)
	}
	return t, l, nil
}

// limitReadCloser is a io.ReadCloser that limits the number of bytes that can be read.
// On close, the remaining bytes are read and discarded. The underlying reader is left open.
type limitReadCloser struct {
	r *io.LimitedReader
}

func (l *limitReadCloser) Read(p []byte) (int, error) {
	return l.r.Read(p)
}

func (l *limitReadCloser) Close() error {
	var buf [4096]byte
	for l.r.N > 0 {
		_, err := l.r.Read(buf[:])
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	// typeSRI is the type of the sri record.
	typeSRI = 0x01
	// typePayload is the type of the payload record.
	typePayload = 0x02
)
