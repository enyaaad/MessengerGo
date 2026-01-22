package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const maxBodyBytes = 1 << 20

func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			return writeError(w, http.StatusBadRequest, "invalid json")
		case errors.Is(err, io.ErrUnexpectedEOF):
			return writeError(w, http.StatusBadRequest, "invalid json")
		case errors.As(err, &unmarshalTypeErr):
			return writeError(w, http.StatusBadRequest, "invalid json field type")
		case errors.Is(err, io.EOF):
			return writeError(w, http.StatusBadRequest, "empty body")
		default:
			return writeError(w, http.StatusBadRequest, "invalid json")
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return writeError(w, http.StatusBadRequest, "invalid json")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) error {
	writeJSON(w, status, map[string]string{"error": message})
	return errors.New(message)
}
