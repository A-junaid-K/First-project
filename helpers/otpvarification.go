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

