package sri

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"io"
	"strings"
)

type Integrity struct {
	// Algorithm is the algorithm used to generate the hash.
	Algorithm
	// Hash is the hash of the payload.
	Hash []byte
}

func FromString(s string) (Integrity, error) {
	var algorithm Algorithm
	switch {
	case strings.HasPrefix(s, "sha256-"):
		algorithm = SHA256
	case strings.HasPrefix(s, "sha384-"):
		algorithm = SHA384
	case strings.HasPrefix(s, "sha512-"):
		algorithm = SHA512
	default:
		return Integrity{}, errors.New("invalid algorithm")
	}
	hash, err := base64.StdEncoding.DecodeString(s[len(algorithm)+1:])
	if err != nil {
		return Integrity{}, fmt.Errorf("decoding hash: %w", err)
	}
	if len(hash) != algorithm.ByteLen() {
		return Integrity{}, fmt.Errorf("invalid hash length: %d", len(hash))
	}
	return Integrity{Algorithm: algorithm, Hash: hash}, nil
}

func FromReader(algorithm Algorithm, payload io.Reader) (Integrity, error) {
	algorithm, err := AlgorithmFromString(string(algorithm))
	if err != nil {
		return Integrity{}, err
	}
	hash, err := algorithm.Hash(payload)
	if err != nil {
		return Integrity{}, err
	}
	return Integrity{Algorithm: algorithm, Hash: hash}, nil
}

func (i Integrity) Validate(payload io.Reader) error {
	payloadHash, err := i.Algorithm.Hash(payload)
	if err != nil {
		return err
	}
	if !bytes.Equal(i.Hash, payloadHash) {
		return errors.New("hash mismatch")
	}
	return nil
}

func (i Integrity) String() string {
	return string(i.Algorithm) + "-" + base64.StdEncoding.EncodeToString(i.Hash)
}

type Algorithm string

func AlgorithmFromString(s string) (Algorithm, error) {
	switch s {
	case "sha256":
		return SHA256, nil
	case "sha384":
		return SHA384, nil
	case "sha512":
		return SHA512, nil
	default:
		return "", errors.New("invalid algorithm")
	}
}

const (
	SHA256 Algorithm = "sha256"
	SHA384 Algorithm = "sha384"
	SHA512 Algorithm = "sha512"
)

func (a Algorithm) ByteLen() int {
	switch a {
	case SHA256:
		return 32
	case SHA384:
		return 48
	case SHA512:
		return 64
	}
	return 0 // unreachable
}

func (a Algorithm) Hash(in io.Reader) ([]byte, error) {
	var hasher hash.Hash
	switch a {
	case SHA256:
		hasher = sha256.New()
	case SHA384:
		hasher = sha512.New384()
	case SHA512:
		hasher = sha512.New()
	default:
		return nil, errors.New("hashing: invalid algorithm")
	}
	if _, err := io.Copy(hasher, in); err != nil {
		return nil, fmt.Errorf("hashing: %w", err)
	}
	return hasher.Sum(nil), nil
}
