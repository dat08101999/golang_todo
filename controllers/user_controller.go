package controllers

import (
	"api_fiber/dbquery"
	"api_fiber/middlewares"
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
	keyTime := time.Now().Add(time.Minute * 30).Unix()
	token, errCreate := Create(user.UserName, string(keyTime))
	refreshToken, _ := refreshtTokenCreate(string(keyTime), user.UserName)
	userModel := dbquery.GetUser(user.UserName)
	if userModel == nil {
		return res.Response(c, models.ErrorResponse{}, errors.New("user does not exitst "), "")
	}
	userModel.Refresh_Token = refreshToken
	go func() {
		errorUpdate := dbquery.UpdateUser(*userModel)
		if errorUpdate != nil {
			print(errorUpdate.Error())
			// return res.Response(c, models.ErrorResponse{}, errorUpdate, "")
		}
	}()
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
	token, err := middlewares.ParseToken(tokenReq.Refresh_Token)
	if err != nil {
		return res.Response(c, models.ErrorResponse{
			Message: "error",
		}, err, "")

	}
	///
	var headertemp string
	var userHeader string
	///
	tempHeader, err := middlewares.ParseToken(tokenHeader)
	if claims, ok := tempHeader.Claims.(jwt.MapClaims); ok && token.Valid {
		headertemp = fmt.Sprintf("%v", claims["key"])
		userHeader = fmt.Sprintf("%v", claims["username"])
	}
	var tokenRe string
	var userName string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenRe = fmt.Sprintf("%v", claims["key"])
		userName = fmt.Sprintf("%v", claims["user"])
		if tokenRe != headertemp || userName != userHeader {
			return res.Response(c, models.ErrorResponse{}, errors.New("invalid token "), "")
		}
	}
	userModel := dbquery.GetUser(userName)
	if userModel == nil {
		return res.Response(c, models.ErrorResponse{}, errors.New("invalid refresh token "), "")
	}
	if userModel.Refresh_Token != tokenReq.Refresh_Token {
		return res.Response(c, models.ErrorResponse{}, errors.New("invalid refresh token token "), "")
	}
	keyTime := time.Now().Add(time.Minute * 30).Unix()
	newToken, err := Create(userName, string(keyTime))
	refreshToken, _ := refreshtTokenCreate(string(keyTime), userName)
	userModel.Refresh_Token = refreshToken
	errorUpdate := dbquery.UpdateUser(*userModel)

	if errorUpdate != nil {
		return res.Response(c, models.ErrorResponse{}, errors.New("invalid refresh token "), "")
	}
	return res.Response(c, fiber.Map{
		"Token":         newToken,
		"Refresh_token": refreshToken,
	}, err, "Success")
}

func Create(username string, key string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["sub"] = 1
	claims["key"] = key
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() //Token hết hạn sau 12 giờ
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

func refreshtTokenCreate(key string, userName string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rtClaims["key"] = key
	rtClaims["user"] = userName
	/// error

	return refreshToken.SignedString([]byte(os.Getenv("SECRET_JWT_REFRESH")))

}
