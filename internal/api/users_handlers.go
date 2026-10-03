package api

import (
	"errors"
	"net/http"

	"github.com/JoaoNaif/Go-bid/internal/services"
	"github.com/JoaoNaif/Go-bid/internal/use-cases/user"
	"github.com/JoaoNaif/Go-bid/internal/utils"
)

func (api *Api) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[user.CreateUserReq](r)

	if err != nil {
		_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(),
		data.UserName,
		data.Email,
		data.Password,
		data.Bio,
	)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateEmailOrUsername) {
			_ = utils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "invalid email or username, already exists",
			})
			return
		}

		_ = utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	_ = utils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := utils.DecodeValidJson[user.LoginUserReq](r)

	if err != nil {
		utils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "invalid email or password",
			})
			return
		}

		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticateUserId", id)

	utils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "logged in, successfully",
	})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		utils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), "AuthenticateUserId")

	utils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "logged out, successfully",
	})
}
