package repository

import (
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
	supa "github.com/supabase-community/supabase-go"
)

type AuthRepository struct {
	client *supa.Client
}

func NewAuthRepository(client *supa.Client) *AuthRepository {
	return &AuthRepository{client: client}
}

func (r *AuthRepository) Register(email, password string) (*types.SignupResponse, error) {
	session, err := r.client.Auth.Signup(types.SignupRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *AuthRepository) DeleteUser(userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return r.client.Auth.AdminDeleteUser(types.AdminDeleteUserRequest{UserID: uid})
}

func (r *AuthRepository) Login(email, password string) (*types.TokenResponse, error) {
	session, err := r.client.Auth.SignInWithEmailPassword(email, password)
	if err != nil {
		return nil, err
	}
	return session, nil
}
