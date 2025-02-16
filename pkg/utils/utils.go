package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var fixedGUID = "123e4567-e89b-12d3-a456-426614174000" // Substitua por seu GUID fixo

func HashPassword(password string) (string, error) {
	passwordWithGUID := password + fixedGUID
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordWithGUID), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	passwordWithGUID := password + fixedGUID
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordWithGUID))
	return err == nil
}

func IsValidCPF(cpf string) bool {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")
	if len(cpf) != 11 {
		return false
	}
	var sum int
	var secondDigit, firstDigit int
	for i := 0; i < 9; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * (10 - i)
	}
	firstDigit = sum % 11
	if firstDigit < 2 {
		firstDigit = 0
	} else {
		firstDigit = 11 - firstDigit
	}
	sum = 0
	for i := 0; i < 10; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * (11 - i)
	}
	secondDigit = sum % 11
	if secondDigit < 2 {
		secondDigit = 0
	} else {
		secondDigit = 11 - secondDigit
	}
	return string(cpf[9]) == strconv.Itoa(firstDigit) && string(cpf[10]) == strconv.Itoa(secondDigit)
}

var jwtKey = []byte("169f1ae6-4b24-4750-a3c4-1222cc69eee2-69efb153-527c-40a4-8029-12c04c6585ae")

func GenerateJWTWithExpiration(userID interface{}, expirationTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expirationTime).Unix(),
	})
	return token.SignedString(jwtKey)
}

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expira em 24 horas
	})
	return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (string, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	userID := (*claims)["user_id"].(string)
	return userID, nil
}

func HideEmailPart(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	local := parts[0]
	domain := parts[1]
	if len(local) > 2 {
		local = local[:2] + strings.Repeat("*", len(local)-2)
	}
	return local + "@" + domain
}
