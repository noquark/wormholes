package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	Time       = 1
	Memory     = 64 * 1024
	Threads    = 4
	KeyLength  = 32
	SaltLength = 32
)

var ErrInvalidPassword = errors.New("invalid password")

type hashed struct {
	salt, hash              []byte
	memory, time, keyLength uint32
	threads                 uint8
}

// Auth helpers/parsers

func ParseAuth(auth string) (email, password string, ok bool) {
	const prefix = "Basic "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}

	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}

	cs := string(c)
	s := strings.IndexByte(cs, ':')

	if s < 0 {
		return
	}

	return cs[:s], cs[s+1:], true
}

func newFromHash(hashedSec string) (*hashed, error) {
	parts := strings.Split(hashedSec, "$")

	h := &hashed{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &h.memory, &h.time, &h.threads); err != nil {
		return nil, fmt.Errorf("failed to build hash : %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt : %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, fmt.Errorf("failed to decode hash : %w", err)
	}

	h.salt = salt
	h.hash = hash
	h.keyLength = uint32(len(h.hash))

	return h, nil
}

func GenerateFromPassword(password []byte) (string, error) {
	salt := make([]byte, SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt : %w", err)
	}

	key := argon2.IDKey(password, salt, Time, Memory, Threads, KeyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, Memory, Time, Threads, b64Salt, b64Key), nil
}

func CompareHashAndPassword(hashedPassword string, password []byte) error {
	h, err := newFromHash(hashedPassword)
	if err != nil {
		return err
	}

	compHash := argon2.IDKey(password, h.salt, h.time, h.memory, h.threads, h.keyLength)
	if subtle.ConstantTimeCompare(h.hash, compHash) == 1 {
		return nil
	}

	return ErrInvalidPassword
}
