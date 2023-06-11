package middlewares

import (
	database "SMS-panel/database"
	"SMS-panel/models"

	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var SECRET = "s89ut8cn4u3bghyn75gy38ghm9g3mgc85g9m" ///should be in env file !!!!!!!!!!

func IsLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//get cookie
		tokenString, err := c.Cookie("username")
		if err != nil {
			return echo.ErrUnauthorized
		}

		//parse token
		token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {

			//wrong algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected Sigining Method %v", token.Header["alg"])
			}
			return []byte(SECRET), nil

		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			//check expiration date
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return echo.ErrUnauthorized
			}

			db, _ := database.GetConnection()

			//find account
			var account models.Account
			accountID := claims["id"].(int)
			db.First(&account, accountID)
			if account.ID == 0 {
				return echo.ErrUnauthorized
			}
			c.Set("account", account)
			return next(c)

		} else {
			return echo.ErrUnauthorized
		}
	}
}
