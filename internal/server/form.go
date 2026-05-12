package server

import (
	"errors"
	"net/http"
)

// maxFormBodyBytes is a hard limit for URL-encoded form bodies to prevent
// memory exhaustion when parsing request data.
//
// 1 MiB is plenty for these pages (a list of selected IDs + a few fields) while
// still preventing accidental/hostile oversized submissions.
const maxFormBodyBytes int64 = 1 << 20 // 1 MiB

// parseFormWithLimit applies a maximum request body size before parsing form data.
//
// This is required for gosec (and good practice) because (*http.Request).ParseForm
// reads the whole body into memory.
func parseFormWithLimit(w http.ResponseWriter, r *http.Request) error {
	if r.Body != nil {
		r.Body = http.MaxBytesReader(w, r.Body, maxFormBodyBytes)
	}

	if err := r.ParseForm(); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return StatusError(http.StatusRequestEntityTooLarge)
		}
		return StatusError(http.StatusBadRequest)
	}

	return nil
}

