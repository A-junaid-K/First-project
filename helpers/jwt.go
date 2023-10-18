// package helpers

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// var Jwtkey = []byte(os.Getenv("SECRET_KEY"))

// // type claims struct {
// // 	Email string
// // 	jwt.StandardClaims
// // }

// func Generatejwt(email, usertype string, id int) (string, error) {
// 	expireAt := time.Now().Add(time.Hour * 24 * 30).Unix()
// 	// jwtclaims := &claims{
// 	// 	Email: email,

// 	// 	StandardClaims: jwt.StandardClaims{
// 	// 		ExpiresAt: expireAt.Unix(),
// 	// 	},
// 	// }

// 	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtclaims)

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
// 		"exp":      expireAt,
// 		"sub":      id,
// 		"email":    email,
// 		"usertype": usertype,
// 	})

// 	//---------------------
// 	tokenString, err := token.SignedString(Jwtkey)
// 	if err != nil {
// 		fmt.Println("error @ 40 ", err)
// 		return "", err
// 	}

// 	return tokenString, nil

// 	// refreshToken := jwt.New(jwt.SigningMethodHS256)
// 	// rtclaims := refreshToken.Claims.(jwt.MapClaims)
// 	// rtclaims["email"] = email
// 	// rtclaims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
// 	// rt, err := refreshToken.SignedString(Jwtkey)

// 	// if err != nil {
// 	// 	fmt.Println("error @ 44 ", err)
// 	// 	return nil, err
// 	// }
// 	// return map[string]string{
// 	// 	"access_token":  tokenString,
// 	// 	"refresh_token": rt,
// 	// }, nil

// }

// // var P string

// // func ValidateToken(singnedToken string) (err error) {
// // 	token, err := jwt.ParseWithClaims(
// // 		singnedToken,
// // 		&claims{},
// // 		func(token *jwt.Token) (interface{}, error) {
// // 			return []byte(Jwtkey), nil
// // 		},
// // 	)
// // 	if err != nil {
// // 		return fmt.Errorf("failed to validate token: %v", err)
// // 	}

// // 	claims, ok := token.Claims.(*claims)
// // 	P = claims.Email
// // 	if !ok {
// // 		err = errors.New("couldn't parse claims")
// // 		return
// // 	}
// // 	if claims.ExpiresAt < time.Now().Local().Unix() {
// // 		err = errors.New("token expired")
// // 		return
// // 	}
// // 	return
// // }

//---------------------code of manaf-----------------------

package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(email string, userType string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"sub":      id,
		"email":    email,
		"usertype": userType,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
