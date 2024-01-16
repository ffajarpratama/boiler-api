package usecase

import (
	"context"
)

// Ping implements IFaceUsecase.
func (u *Usecase) Ping(ctx context.Context) (string, error) {
	res, err := u.Repo.Ping(ctx)
	if err != nil {
		return "", err
	}

	return res, nil
}
