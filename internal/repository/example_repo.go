package repository

import "context"

// Ping implements IFaceRepository.
func (r *Repository) Ping(ctx context.Context) (string, error) {
	return "PONG", nil
}
