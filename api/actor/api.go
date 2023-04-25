package actor

import (
	chttp "challenge/http"
	"challenge/model"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
)

type AddActorRequest struct {
	FullName string `json:"full_name"`
}

type UpdateActorRequest struct {
	FullName string `json:"full_name"`
}

type GetActorsResponse struct {
	*model.ActorCollection
}

type GetActorResponse struct {
	Actor model.Actor `json:"actor"`
}

func AddActor(m *model.ApplicationModels) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request AddActorRequest
		decodeJsonErr := chttp.DecodeJSONPayload(w, r, &request)
		if decodeJsonErr != nil {
			rp, _ /* << criminal */ := json.Marshal(decodeJsonErr.Msg)
			w.WriteHeader(decodeJsonErr.Status)
			_, _ /* << criminal */ = w.Write(rp)

			return
		}

		err := m.Actors.Add(&model.Actor{
			FullName: request.FullName,
		})

		if err != nil {
			_, _ = chttp.CreateResponseInternalServerError(w, err)
			return
		}

		_, _ = chttp.CreateResponseCreated(w)
	}
}

func UpdateActor(m *model.ApplicationModels) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UpdateActorRequest
		decodeJsonErr := chttp.DecodeJSONPayload(w, r, &request)
		if decodeJsonErr != nil {
			rp, _ /* << criminal */ := json.Marshal(decodeJsonErr.Msg)
			w.WriteHeader(decodeJsonErr.Status)
			_, _ /* << criminal */ = w.Write(rp)

			return
		}

		code := mux.Vars(r)["uuid"]
		if !govalidator.IsUUIDv4(code) {
			_, _ = chttp.CreateResponseBadRequestUuid(w)

			return
		}

		actor, err := m.Actors.Find(code)
		if err != nil {
			_, _ = chttp.CreateResponseInternalServerError(w, err)

			return
		}

		if actor == nil {
			_, _ = chttp.CreateResponseNotFound(w)

			return
		}

		actor.FullName = request.FullName
		err = m.Actors.Update(actor)

		if err != nil {
			_, _ = chttp.CreateResponseInternalServerError(w, err)
			return
		}

		_, _ = chttp.CreateResponseNoContent(w)
	}
}

func GetActors(m *model.ApplicationModels) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actors, err := m.Actors.FindAll()

		if err != nil {
			_, _ = chttp.CreateResponseInternalServerError(w, err)
			return
		}

		_, _ = chttp.CreateResponseOk(w, GetActorsResponse{actors})
	}
}

func GetActor(m *model.ApplicationModels) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := mux.Vars(r)["uuid"]
		if !govalidator.IsUUIDv4(code) {
			_, _ = chttp.CreateResponseBadRequestUuid(w)

			return
		}

		actor, err := m.Actors.Find(code)
		if err != nil {
			_, _ = chttp.CreateResponseInternalServerError(w, err)

			return
		}

		if actor == nil {
			_, _ = chttp.CreateResponseNotFound(w)

			return
		}
		_, _ = chttp.CreateResponseOk(w, GetActorResponse{Actor: *actor})
	}
}
