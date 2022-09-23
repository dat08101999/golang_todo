package middlewares

import (
	"api_fiber/models"
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func CheckToken(p func(c *fiber.Ctx, userName string) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bearer := string(c.Request().Header.Peek("authorization"))
		var username string
		fmt.Println(bearer)
		if bearer == "" {
			err := c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Message: "missing authorization",
			})
			return err
		}
		tokenString := strings.Split(bearer, " ")[1]
		/// parse token to username
		token, errToken := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_JWT")), nil
		})
		/// error
		if errToken != nil {
			err := c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Message: errToken.Error(),
			})
			return err
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			username = fmt.Sprintf("%v", claims["username"])
			if claims["username"] == nil {
				fmt.Println("c")
				return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
					Message: "Token is not valid",
				})
			}
		}
		fmt.Println(username)

		return p(c, username)
	}
}
