package handler

import (
	"context"
	"goto/src/config"
	"goto/src/database/query"
	"goto/src/model"
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

func getJwtToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.Id,
		"login": user.Login,
		"exp":   time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString([]byte(config.SecretKey))
	return encodedToken, err
}

func GetCurrentUser(c *fiber.Ctx) *model.User {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return &model.User{
		Id:    int(claims["id"].(float64)),
		Login: claims["login"].(string),
	}
}

func Register(c *fiber.Ctx) error {
	ctx := context.Background()
	body := model.User{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if body.Login == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).
			SendString("Login and password must be provided")
	}
	if len(body.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).SendString("Password is too short")
	}
	if query.IsLoginInUse(ctx, body.Login) {
		return c.Status(fiber.StatusBadRequest).SendString("Login is already in use")
	}

	passwordHash, err := hashPassword(body.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user := &model.User{Login: body.Login, Password: passwordHash}
	if err := query.CreateUser(ctx, user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	user, _ = query.GetUserByLogin(ctx, body.Login)

	encodedToken, err := getJwtToken(user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": encodedToken})
}

func Login(c *fiber.Ctx) error {
	ctx := context.Background()
	body := model.User{}
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if body.Login == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).
			SendString("Login and password must be provided")
	}

	user, err := query.GetUserByLogin(ctx, body.Login)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Login or password is not valid")
	}
	if valid := checkPasswordHash(body.Password, user.Password); !valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Login or password is not valid")
	}

	encodedToken, err := getJwtToken(user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": encodedToken})
}