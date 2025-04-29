package toolshed

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid"
)

func CreateULID(prefix string, timestamp time.Time) (string, error) {
	var newUlid ulid.ULID

	entropy := rand.Reader
	ms := ulid.Timestamp(timestamp)
	newUlid, err := ulid.New(ms, entropy)
	if err != nil {
		return "", err
	}

	ulidString := newUlid.String()
	if prefix != "" {
		ulidString = prefix + "_" + ulidString
	}

	return ulidString, nil
}

func RandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = letterBytes[int(b[i])%len(letterBytes)]
	}
	return string(b)
}
