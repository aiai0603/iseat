package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

func JsonResponse(i interface{}, w http.ResponseWriter, sc int) {
	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(sc)
	enc := json.NewEncoder(w)
	err := enc.Encode(i)
	// Problems encoding
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func Respond(i interface{}, err error, w http.ResponseWriter) {
	RespondWithCode(i, err, w, http.StatusOK)
}

func RespondWithCode(i interface{}, err error, w http.ResponseWriter, code int) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	JsonResponse(i, w, code)
}

func Str2Uint(s string) (uint64, error) {
	u, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		// TODO: Fixme
		return 0, err
	}

	return u, nil
}
