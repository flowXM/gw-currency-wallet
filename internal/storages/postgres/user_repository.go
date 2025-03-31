package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/pkg/client/postgresql"
	"gw-currency-wallet/pkg/utils"
)

type userRepository struct{}

func NewUserRepository() storages.UserRepository {
	return &userRepository{}
}

func (u userRepository) Register(username, password, email string) error {
	db, err := postgresql.NewClient()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	salt := utils.GenerateRandomSalt()
	hash := utils.HashPassword(password, salt)
	id := uuid.New()

	_, err = tx.Exec("INSERT INTO users VALUES ($1, $2, $3, $4, $5);", id, username, email, hash, salt)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO wallets VALUES ($1, $2, 100, 100, 100)", uuid.New(), id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (u userRepository) Login(username, password string) (token string, err error) {
	db, err := postgresql.NewClient()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var salt []byte
	var password_hash string
	var user_id string

	res := db.QueryRow("SELECT salt, password_hash, user_id FROM users WHERE username = $1", username)
	err = res.Scan(&salt, &password_hash, &user_id)
	if err != nil {
		return "", err
	}

	if utils.ConfirmPassword(password, salt, password_hash) {
		token, err := utils.GenerateToken(user_id)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	return "", fmt.Errorf("invalid username or password")
}
