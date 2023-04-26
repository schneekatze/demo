package api

import (
	cjwt "challenge/service/jwt"
	"encoding/json"
	"fmt"
	"net/http"
)

type SingInRequest struct {
	Name     string `json:"name" valid:"minstringlength(2),required"`
	Password string `json:"password" valid:"required"`
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var request SingInRequest
	decodeJsonErr := DecodeJSONPayload(w, r, &request)
	if decodeJsonErr != nil {
		rp, _ /* << criminal */ := json.Marshal(decodeJsonErr.Msg)
		w.WriteHeader(decodeJsonErr.Status)
		_, _ /* << criminal */ = w.Write(rp)

		return
	}

	if request.Name == "foo" && request.Password == "bar" {
		token, err := cjwt.GetToken(request.Name)
		if err != nil {
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)
		_, _ = CreateResponseOk(w, map[string]string{
			"token": token,
		})
	} else {
		_, _ = CreateResponseUnauthorized(w, "Name and password do not match")

		return
	}
}

func ProtectedResource(w http.ResponseWriter, r *http.Request) {
	_, _ = CreateResponse(w, fmt.Sprintf("Hallo, %s", r.Header.Get("X-Name")), http.StatusOK)
}
