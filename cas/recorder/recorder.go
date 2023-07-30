package recorder

import (
	"io"

	"github.com/malt3/abstractfs-core/api"
)

// Recorder is a recorder for a CAS.
// It uses the recorder protocol to read cas contents from a stream
// and write them to the CAS.
type Recorder struct {
	cas    api.CASWriter
	reader io.Reader
}

// New creates a new recorder.
func New(cas api.CASWriter, r io.Reader) *Recorder {
	return &Recorder{
		cas:    cas,
		reader: r,
	}
}

// Consume reads from the reader and writes the contents to the CAS.
func (r *Recorder) Consume() error {
	var err error
	for err == nil {
		err = r.consumeOne()
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func (r *Recorder) consumeOne() error {
	sri, body, err := Decode(r.reader)
	if err != nil {
		return err
	}

	defer body.Close()
	if err := r.cas.Write(sri.String(), body); err != nil {
		return err
	}
	return nil
}
