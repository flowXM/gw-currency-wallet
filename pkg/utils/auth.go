package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		jt, err := ValidateToken(token[7:])
		if err != nil || !jt.Valid {
			jsonBytes, _ := json.Marshal(
				struct {
					Error string `json:"error"`
				}{Error: "Incorrect token"})

			http.Error(w, string(jsonBytes), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GenerateRandomSalt() []byte {
	var salt = make([]byte, saltSize)
	rand.Read(salt[:])
	return salt
}

func HashPassword(password string, salt []byte) string {
	passwordWithSalt := append([]byte(password), salt...)
	var hash = sha512.Sum512(passwordWithSalt)
	return hex.EncodeToString(hash[:])
}

func ConfirmPassword(password string, salt []byte, expectedHash string) bool {
	return HashPassword(password, salt) == expectedHash
}

func GenerateToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}

func GetUserId(r *http.Request) string {
	token := r.Header.Get("Authorization")
	jt, _ := jwt.Parse(token[7:], func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	claims := jt.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(string)
	return userId
}
