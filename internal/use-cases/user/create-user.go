package user

import (
	"context"

	"github.com/JoaoNaif/Go-bid/internal/validator"
)

type CreateUserReq struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "This field cannot be empty.")

	eval.CheckField(validator.NotBlank(req.Email), "email", "This field cannot be empty.")
	eval.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "must be a valid email")

	eval.CheckField(validator.NotBlank(req.Bio), "bio", "This field cannot be empty.")
	eval.CheckField(
		validator.MinChars(req.Bio, 10) &&
			validator.MaxChars(req.Bio, 255),
		"bio", "This field must have a length between 10 and 255",
	)

	eval.CheckField(validator.MinChars(req.Password, 8), "password", "Password must be bigger than 8 chars")

	return eval
}
