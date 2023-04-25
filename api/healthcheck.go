package api

import (
	chttp "challenge/http"
	"net/http"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	_, _ = chttp.CreateResponse(w, "alive", http.StatusOK)
}
