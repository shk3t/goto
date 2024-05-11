package handler

import (
	"context"
	"goto/src/config"
	q "goto/src/database/query"
	m "goto/src/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getJwtToken(user *m.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.Id,
		"login": user.Login,
		"exp":   time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString([]byte(config.SecretKey))
	return encodedToken, err
}

func Register(fctx *fiber.Ctx) error {
	ctx := context.Background()
	body := m.User{}
	if err := fctx.BodyParser(&body); err != nil {
		return err
	}

	if body.Login == "" || body.Password == "" {
		return fctx.Status(fiber.StatusBadRequest).
			SendString("Login and password must be provided")
	}
	if len(body.Password) < 8 {
		return fctx.Status(fiber.StatusBadRequest).SendString("Password is too short")
	}
	if q.IsLoginInUse(ctx, body.Login) {
		return fctx.Status(fiber.StatusBadRequest).SendString("Login is already in use")
	}

	passwordHash, err := hashPassword(body.Password)
	if err != nil {
		return fctx.SendStatus(fiber.StatusInternalServerError)
	}

	user := &m.User{Login: body.Login, Password: passwordHash}
	user, err = q.CreateUser(ctx, user)
	if err != nil {
		return fctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	encodedToken, err := getJwtToken(user)
	if err != nil {
		return fctx.SendStatus(fiber.StatusInternalServerError)
	}
	return fctx.JSON(fiber.Map{"user": user, "token": encodedToken})
}

func Login(fctx *fiber.Ctx) error {
	ctx := context.Background()
	body := m.User{}
	if err := fctx.BodyParser(&body); err != nil {
		return err
	}

	if body.Login == "" || body.Password == "" {
		return fctx.Status(fiber.StatusBadRequest).
			SendString("Login and password must be provided")
	}

	user, err := q.GetUserByLogin(ctx, body.Login)
	if err != nil {
		return fctx.Status(fiber.StatusUnauthorized).SendString("Login or password is not valid")
	}
	if valid := checkPasswordHash(body.Password, user.Password); !valid {
		return fctx.Status(fiber.StatusUnauthorized).SendString("Login or password is not valid")
	}

	encodedToken, err := getJwtToken(user)
	if err != nil {
		return fctx.SendStatus(fiber.StatusInternalServerError)
	}
    return fctx.JSON(fiber.Map{"user": user, "token": encodedToken})
}