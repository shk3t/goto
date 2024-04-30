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
		"login": user.Login,
		"exp":   time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString([]byte(config.SecretKey))
	return encodedToken, err
}

func Register(c *fiber.Ctx) error {
	ctx := context.Background()
	user := model.User{}
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	if user.Login == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).
			SendString("Login and password must be provided")
	}
	if len(user.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).SendString("Password is too short")
	}
	if query.IsLoginInUse(ctx, user.Login) {
		return c.Status(fiber.StatusBadRequest).SendString("Login is already in use")
	}

	passwordHash, err := hashPassword(user.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user = model.User{Login: user.Login, Password: passwordHash}
	if err := query.CreateUser(ctx, &user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	encodedToken, err := getJwtToken(&user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": encodedToken})
}

func Login(c *fiber.Ctx) error {
	ctx := context.Background()
	user := model.User{}
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	if user.Login == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).
			SendString("Login and password must be provided")
	}

	targetUser, err := query.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Login or password is not valid")
	}
	if valid := checkPasswordHash(user.Password, targetUser.Password); !valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Login or password is not valid")
	}

	encodedToken, err := getJwtToken(&user)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": encodedToken})
}

// func restricted(c fiber.Ctx) error {
// 	user := c.Locals("user").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	name := claims["name"].(string)
// 	return c.SendString("Welcome " + name)
// }