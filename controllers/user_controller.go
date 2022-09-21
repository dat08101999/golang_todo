package controllers

import (
	"api_fiber/dbquery"
	"api_fiber/models"
	"api_fiber/res"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{})
	}
	if user.UserName == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{})
	}

	hash, _ := models.Hash(user.Password)
	user.Password = hash
	result := dbquery.InsertUser(user)
	if result == nil {
		return c.JSON(models.SuccessResponseModel{
			Code:    200,
			Message: "Register success",
			Data:    map[string]interface{}{},
		})
	}
	return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
		Message: result.Error(),
	})
}

func LoginUser(c *fiber.Ctx) error {
	var user models.User
	errUser := c.BodyParser(&user)
	if errUser != nil {
		return res.Response(c, user, errors.New("Missing element field"), "")
	}
	err := dbquery.FindUser(user)
	if err != nil {
		return res.Response(c, user, err, "")
	}
	token, errCreate := Create(user.UserName)
	refreshToken, _ := refreshtTokenCreate(token, user.UserName)

	return res.Response(c, map[string]interface{}{
		"Token":         token,
		"Refresh_token": refreshToken,
	}, errCreate, "Success")
}

func RefreshToken(c *fiber.Ctx) error {
	type bodyParams struct {
		Refresh_Token string `bson:"refresh_token"`
	}

	var tokenHeader string
	var tokenReq bodyParams
	errorBody := c.BodyParser(&tokenReq)
	bearer := string(c.Request().Header.Peek("authorization"))
	fmt.Println(tokenReq)
	if len(strings.Split(bearer, " ")) > 1 {
		tokenHeader = strings.Split(bearer, " ")[1]
	} else {
		return res.Response(c, models.ErrorResponse{}, errors.New("invalid token "), "")
	}

	///
	if errorBody != nil {
		return res.Response(c, models.ErrorResponse{
			Message: "error",
		}, errorBody, "")
	}
	///
	token, err := jwt.Parse(tokenReq.Refresh_Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_JWT_REFRESH")), nil
	})
	if err != nil {
		return res.Response(c, models.ErrorResponse{
			Message: "error",
		}, err, "")
	}

	var tokenRe string
	var userName string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenRe = fmt.Sprintf("%v", claims["token"])
		userName = fmt.Sprintf("%v", claims["user"])
		if tokenRe != tokenHeader {
			return res.Response(c, models.ErrorResponse{}, errors.New("invalid token "), "")
		}
	}
	newToken, err := Create(userName)
	refreshToken, _ := refreshtTokenCreate(newToken, userName)
	return res.Response(c, fiber.Map{
		"Token":         newToken,
		"Refresh_token": refreshToken,
	}, err, "Success")
}

func Create(username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["sub"] = 1
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix() //Token hết hạn sau 12 giờ
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

func refreshtTokenCreate(token string, userName string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rtClaims["token"] = token
	rtClaims["user"] = userName
	/// error

	return refreshToken.SignedString([]byte(os.Getenv("SECRET_JWT_REFRESH")))

}
