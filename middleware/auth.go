package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func UserAuthentication(c *gin.Context) {

	tokenString, err := c.Cookie("jwt_user")
	if err != nil {
		log.Println("error in jwt_user : ", err)
		c.Redirect(303, "/user/login")
		c.AbortWithStatus(400)
		return
	}
	//Decode / validate it
	// Parse takes the token string and a function for looking up the key. The latter is especially

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.Redirect(303, "/user/login")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		c.Redirect(303, "/user/login")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("exp")
			c.Redirect(303, "/user/login")
			return
		}
		//Find the user with token sub
		var user models.User
		database.DB.First(&user, claims["sub"])
		if user.User_id == 0 {
			c.Redirect(303, "/user/login")
			return
		}
		//Attach to req
		c.Set("user", user)
		//Continue
		c.Next()
	} else {
		fmt.Println("Failed \n @62")
		c.Redirect(303, "/user/login")
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func AdminAuthentication(c *gin.Context) {
	tokenString, err := c.Cookie("jwt_admin")
	if err != nil {
		log.Println("unatherized access")
		c.Redirect(303, "/admin-login")
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		log.Println("Failed to admin generate token  ", err)
		c.HTML(404, "adminLogin.html", gin.H{
			"error": "error occurse while token generation",
		})
		c.AbortWithStatus(404)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	} else {
		c.AbortWithStatus(401)
	}
}
