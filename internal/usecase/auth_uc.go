package usecase

import (
	"context"

	"github.com/ffajarpratama/boiler-api/internal/http/request"
	"github.com/ffajarpratama/boiler-api/internal/model"
	"github.com/ffajarpratama/boiler-api/pkg/hash"
	"github.com/ffajarpratama/boiler-api/pkg/jwt"
	"github.com/google/uuid"
)

// TODO: do properly
// RegisterUser implements IFaceUsecase.
func (u *Usecase) RegisterUser(ctx context.Context, req *request.ReqRegisterUser) error {
	tx := u.DB.Begin()
	defer tx.Rollback()

	pwd, err := hash.HashAndSalt([]byte(req.Password))
	if err != nil {
		return err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: pwd,
	}

	err = u.Repo.CreateUser(ctx, user, tx)
	if err != nil {
		return err
	}

	return tx.Commit().Error
}

// Login implements IFaceUsecase.
func (u *Usecase) Login(ctx context.Context, req *request.ReqLoginUser) (*model.User, error) {
	res, err := u.Repo.FindOneUser(ctx, "email = ?", req.Email)
	if err != nil {
		return nil, err
	}

	err = hash.Compare(res.Password, []byte(req.Password))
	if err != nil {
		return nil, err
	}

	claims := &jwt.CustomUserClaims{
		ID:   res.UserID.String(),
		Role: "",
	}

	res.Token, err = jwt.GenerateToken(claims, u.Cnf.JWTConfig.User)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserProfile implements IFaceUsecase.
func (u *Usecase) GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	res, err := u.Repo.FindOneUser(ctx, "user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
