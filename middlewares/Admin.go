package middlewares

import (
	"net/http"
	"os"

	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
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

			//account isn't an admin
			if claims["admin"].(bool) == false {
				return echo.ErrUnauthorized
			}

			return next(c)

		} else {
			return echo.ErrUnauthorized
		}
	}
}
