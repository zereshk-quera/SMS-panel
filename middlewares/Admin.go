package middlewares

import (
	database "SMS-panel/database"
	"SMS-panel/models"
	"os"

	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		tokenString := req.Header.Get("Authorization")

		//Account Doesn't have Token
		if tokenString == "" {
			return echo.ErrConflict
		}

		//Parse Token
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			//Wrong Algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected Sigining Method %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil

		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			//Check Expiration Time
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				return echo.ErrUnauthorized
			}

			//Connect To Database
			db, err := database.GetConnection()
			if err != nil {
				return echo.ErrInternalServerError
			}
			//account isn't an admin
			if claims["admin"].(bool) == false {
				return echo.ErrUnauthorized
			}

			//Find Account
			var account models.Account
			db.First(&account, claims["id"])

			//Token And Id are not For Same Accounts
			if account.ID == 0 {
				return echo.ErrUnauthorized
			}

			//account isn't active
			if account.Token == "" && account.IsActive == false {
				return echo.ErrUnauthorized
			}

			return next(c)

		} else {
			return echo.ErrUnauthorized
		}
	}
}
