package usecase

import (
	"context"
	"net/http"

	"github.com/ffajarpratama/boiler-api/internal/http/request"
	"github.com/ffajarpratama/boiler-api/internal/model"
	"github.com/ffajarpratama/boiler-api/internal/repository"
	"github.com/ffajarpratama/boiler-api/lib/custom_error"
	"github.com/ffajarpratama/boiler-api/lib/hash"
	"github.com/ffajarpratama/boiler-api/lib/jwt"
	"github.com/google/uuid"
)

// Register implements IFaceUsecase.
func (u *Usecase) Register(ctx context.Context, req *request.Register) (*model.User, error) {
	pwd, err := hash.HashAndSalt([]byte(req.Password))
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: pwd,
	}

	err = u.repo.CreateUser(ctx, user, u.db)
	if err != nil {
		if repository.IsDuplicateErr(err) {
			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusConflict,
				Message:  "email sudah digunakan",
			})

			return nil, err
		}

		return nil, err
	}

	claims := &jwt.CustomClaims{
		UserID: user.UserID.String(),
	}

	user.AccessToken, err = jwt.GenerateToken(claims, u.cnf.JWT.Secret)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login implements IFaceUsecase.
func (u *Usecase) Login(ctx context.Context, req *request.Login) (*model.User, error) {
	user, err := u.repo.FindOneUser(ctx, "email = ?", req.Email)
	if err != nil {
		return nil, err
	}

	claims := &jwt.CustomClaims{
		UserID: user.UserID.String(),
	}

	user.AccessToken, err = jwt.GenerateToken(claims, u.cnf.JWT.Secret)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetProfile implements IFaceUsecase.
func (u *Usecase) GetProfile(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return u.repo.FindOneUser(ctx, "user_id = ?", userID)
}
