package middlewares

import (
	"net/http"
	"os"

	database "SMS-panel/database"
	"SMS-panel/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func IsLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//Get Cookie
		cookies := c.Cookies()
		tokenString := ""
		for _, cookie := range cookies {
			if cookie.Name == "account_token" {
				tokenString = cookie.Value
				break
			}
		}

		//Account Doesn't have Token
		if tokenString == "" {
			return echo.ErrUnauthorized
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
				cookie := &http.Cookie{
					Name:   "account_token",
					Value:  "",
					Path:   "/",
					MaxAge: -1,
				}
				c.SetCookie(cookie)
				c.SetCookie(&http.Cookie{Name: "account_token", MaxAge: -1})
				return echo.ErrUnauthorized
			}

			//Connect To Database
			db, err := database.GetConnection()
			if err != nil {
				return echo.ErrInternalServerError
			}

			//Find Account
			var account models.Account
			db.First(&account, claims["id"])

			//Token And Id are not For Same Accounts
			if account.ID == 0 {
				return echo.ErrUnauthorized
			}
			//account is deactive
			if account.Token == "" && account.IsActive == false {
				cookie := &http.Cookie{
					Name:   "account_token",
					Value:  "",
					Path:   "/",
					MaxAge: -1,
				}
				c.SetCookie(cookie)
				c.SetCookie(&http.Cookie{Name: "account_token", MaxAge: -1})
				return echo.ErrUnauthorized
			}

			//Add Account Object To Context
			c.Set("account", account)
			return next(c)

		} else {
			return echo.ErrUnauthorized
		}
	}
}
