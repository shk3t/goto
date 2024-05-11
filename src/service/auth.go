package service

import (
	m "goto/src/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetCurrentUser(fctx *fiber.Ctx) *m.User {
	token := fctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return &m.User{
		Id:    int(claims["id"].(float64)),
		Login: claims["login"].(string),
	}
}