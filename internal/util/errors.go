package util

import (
	"errors"
	"net/http"
)

// IsBodyTooLargeError true jika error berasal dari http.MaxBytesReader
// (request body melebihi batas ukuran yang dikonfigurasi).
func IsBodyTooLargeError(err error) bool {
	if err == nil {
		return false
	}
	var mbe *http.MaxBytesError
	return errors.As(err, &mbe)
}