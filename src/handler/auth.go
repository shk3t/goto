package handler

import (
	"context"
	"goto/src/config"
	"goto/src/database/query"
	"goto/src/model"
	"time"

	"github.com/gofiber/fiber/v3"
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

// func main() {
//     password := "secret"
//     hash, _ := HashPassword(password) // ignore error for the sake of simplicity
//
//     fmt.Println("Password:", password)
//     fmt.Println("Hash:    ", hash)
//
//     match := CheckPasswordHash(password, hash)
//     fmt.Println("Match:   ", match)
// }

func Register(c fiber.Ctx) error {
	ctx := context.Background()
	user := model.User{}
	c.Bind().Body(&user)

	if user.Login == "" || user.Password == "" {
		return c.Status(400).SendString("`login` or `password` are not specified")
	}
	if len(user.Password) < 8 {
		return c.Status(400).SendString("Password is too short")
	}
	if query.IsLoginInUse(ctx, user.Login) {
		return c.Status(400).SendString("Login is already in use")
	}

	passwordHash, err := hashPassword(user.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user = model.User{Login: user.Login, Password: passwordHash}
	if err := query.CreateUser(ctx, &user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	claims := jwt.MapClaims{
		"login": user.Login,
		"exp":   time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encodedToken, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": encodedToken})
}

// func Login(c fiber.Ctx) error {
// 	user := model.User{}
// 	c.Bind().Body(&user)
//
// 	// Throws Unauthorized error
// 	if user.Login != "john" || user.Password != "doe" {
// 		return c.SendStatus(fiber.StatusUnauthorized)
// 	}
//
// 	// Create the Claims
// 	claims := jwt.MapClaims{
// 		"name":  "John Doe",
// 		"admin": true,
// 		"exp":   time.Now().Add(time.Hour * 72).Unix(),
// 	}
//
// 	// Create token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
// 	// Generate encoded token and send it as response.
// 	t, err := token.SignedString([]byte(config.SecretKey))
// 	if err != nil {
// 		return c.SendStatus(fiber.StatusInternalServerError)
// 	}
//
// 	return c.JSON(fiber.Map{"token": t})
// }

// func restricted(c fiber.Ctx) error {
// 	user := c.Locals("user").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	name := claims["name"].(string)
// 	return c.SendString("Welcome " + name)
// }