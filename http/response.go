package http

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CreateResponse(w http.ResponseWriter, payload interface{}, httpStatus int) (int, error) {
	w.WriteHeader(httpStatus)

	if payload == nil {
		return 0, nil
	}

	rp, _ := json.Marshal(payload)

	return w.Write(rp)
}

func CreateResponseCreated(w http.ResponseWriter) (int, error) {
	return CreateResponse(w, nil, http.StatusCreated)
}

func CreateResponseNoContent(w http.ResponseWriter) (int, error) {
	return CreateResponse(w, nil, http.StatusNoContent)
}

func CreateResponseBadRequest(w http.ResponseWriter, payload interface{}) (int, error) {
	return CreateResponse(w, payload, http.StatusBadRequest)
}

func CreateResponseBadRequestUuid(w http.ResponseWriter) (int, error) {
	return CreateResponseBadRequest(
		w,
		map[string]string{
			"error": "Invalid code. Must be a valid uuidv4.",
		},
	)
}

func CreateResponseOk(w http.ResponseWriter, payload interface{}) (int, error) {
	return CreateResponse(w, payload, http.StatusOK)
}

func CreateResponseNotFound(w http.ResponseWriter) (int, error) {
	return CreateResponse(w, map[string]string{
		"error": "Not Found",
	}, http.StatusNotFound)
}

func CreateResponseInternalServerError(w http.ResponseWriter, err error) (int, error) {
	log.Errorf("http: internal server error: %s", err.Error())
	log.Debug(err)

	return CreateResponse(
		w,
		map[string]string{
			"error": "Internal Server Error",
		},
		http.StatusInternalServerError,
	)
}
