package api

import (
	"net/http"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	_, _ = CreateResponse(w, "alive", http.StatusOK)
}
