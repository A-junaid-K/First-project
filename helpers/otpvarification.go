package helpers

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateOtp() string {
	seed := time.Now().UnixNano()
	randomGenerator := rand.New(rand.NewSource(seed))

	return strconv.Itoa(randomGenerator.Intn(8999) + 1000)
}

func VerifyOtp(c *gin.Context, email string) string {
	otp := GenerateOtp()
	email1 := []string{email}

	SendOtp(otp, email1)
	// if err != nil {
	// 	log.Panic(err)
	// }
	return otp
}

func SendOtp(otp string, to []string) {
	auth := smtp.PlainAuth(
		"",
		"junaidkaidakath@gmail.com",
		"rkjsjsaqbxwuitvl",
		"smtp.gmail.com",
	)

	msg := "Subject:" + " otp verification" + "\n" + otp

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		os.Getenv("USERNAME"),
		to,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)

	}

}

//-----------------------------------------------//

// var eMail string
// var otp string

// func sendMail(c *gin.Context, email string) {
// 	eMail = email
// 	user := os.Getenv("USERNAME")
// 	//	password := os.Getenv("PASSWORD")
// 	host := os.Getenv("HOST")
// 	to := email
// 	sub := os.Getenv("eSub")

// 	// auth := smtp.PlainAuth("", user, password, host)
// 	auth := smtp.PlainAuth("", "junaidkaidakath@gmail.com", "fsdfnzhhvwwqfytf", "smtp.gmail.com")

// 	body := fmt.Sprintf("OTP for varification is: %s", otp)
// 	fmt.Println("OTP:", otp)

// 	msg := []byte("To: " + to + "\r\n" +
// 		"Subject: " + sub + "\r\n" +
// 		"\r\n" + body + "\r\n")

// 	err := smtp.SendMail(host+":587", auth, user, []string{to}, msg)
// 	if err != nil {
// 		// c.JSON(http.StatusInternalServerError, gin.H{
// 		// 	"err": err,
// 		// })
// 		c.HTML(400, "signup.html", gin.H{
// 			"message": "invalid email or password",
// 		})
// 		c.Abort()
// 		return
// 	}
// }

//----------------------------------------------------//

// func SendOtp(otp, email string) {
// 	auth := smtp.PlainAuth("", os.Getenv("USERNAME"), os.Getenv("PASSWORD"), "smtp.gmail.com")
// 	to := []string{email}
// 	message := "Subject: Otp verification\nyour verification otp is " + otp
// 	err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("USERNAME"), to, []byte(message))
// 	if err != nil {
// 		fmt.Println("failed to send otp")
// 		return
// 	}
// }
