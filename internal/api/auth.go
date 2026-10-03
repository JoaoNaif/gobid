package api

import (
	"net/http"

	"github.com/JoaoNaif/Go-bid/internal/utils"
	"github.com/gorilla/csrf"
)

func (api *Api) HandleGetCSRFtoken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	utils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"csrf_token": token,
	})
}

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticateUserId") {
			utils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
				"message": "must be logged in",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
