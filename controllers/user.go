package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/helpers"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		fmt.Println("Failed to hash password")
		return "", errors.New("failed to hash password")
	}
	return string(byte), nil
}

// var modelsuser *models.User

func Verifypassword(dbpassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password)); err != nil {
		log.Println("failed to compare password ", err)
		return false
	}
	return true
}

func SignUp(c *gin.Context) {
	c.HTML(200, "SignUp.html", nil)
}

func PostSignUp(c *gin.Context) {

	name := c.Request.FormValue("name")
	phone := c.Request.FormValue("number")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	confpassword := c.Request.FormValue("confpassword")

	err := namevalidator(name)
	if err != nil {
		c.HTML(400, "SignUp.html", gin.H{
			"error": err,
		})
		return
	}

	err = numbervalidator(phone)
	if err != nil {
		c.HTML(400, "SignUp.html", gin.H{
			"error": err,
		})
		return
	}

	err = emailvalidator(email)

	if err != nil {
		c.HTML(400, "SignUp.html", gin.H{
			"error": err,
		})
		return
	}

	err = passwordvalidator(password)
	if err != nil {
		c.HTML(400, "SignUp.html", gin.H{
			"error": err,
		})
		return
	}
	if confpassword != password {
		c.HTML(400, "SignUp.html", gin.H{
			"error": "Incorrect confirm password",
		})
		return
	}

	var dtuser models.User

	database.DB.Where("email=?", email).First(&dtuser)
	if email == dtuser.Email {
		c.HTML(http.StatusBadRequest, "SignUp.html", gin.H{
			"error": email + " has already been used",
		})
		return
	}

	otp := helpers.VerifyOtp(c, email)

	hashPassword, err := hashPassword(password)

	if err != nil {
		log.Println(err)
		return
	}

	err = database.DB.Create(&models.User{
		Name:       name,
		Password:   hashPassword,
		Email:      email,
		Phone:      phone,
		Otp:        otp,
		User_type:  "user",
		Created_at: time.Now(),
	}).Error
	if err != nil {
		log.Println(err)
		return
	}
	c.Redirect(303, "/user/varifyotp")
}

func VarifyOtp(c *gin.Context) {
	c.HTML(200, "otp.html", nil)
}
func PostVarifyOtp(c *gin.Context) {
	eMail := c.Request.FormValue("email")
	otp := c.Request.FormValue("otp")
	var user models.User
	err := database.DB.First(&user, "email =?", eMail).Error

	if user.Otp == otp && err == nil {

		err = database.DB.Model(&models.User{}).Where("email =?", eMail).Update("validate", true).Error

		if err != nil {
			log.Println(err)
			return
		}

		c.Redirect(303, "/user/login")
	} else {
		err = database.DB.Where("validate=?", false).Delete(&models.User{}).Error

		if err != nil {
			log.Println(err)
			return
		}

		c.HTML(200, "otp.html", gin.H{
			"error": "Invalid Email or otp",
		})
	}
}

func Login(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func Postlogin(c *gin.Context) {

	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	//finding with username in database

	var user models.User
	get := database.DB.Where("email=?", email).First(&user)

	if get.Error != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "This user is not signed up",
		})
		return
	}
	checkpassword := Verifypassword(user.Password, password)
	if !checkpassword {
		c.HTML(401, "login.html", gin.H{
			"error": "Incorrect password",
		})
		return
	}

	if user.IsBlocked {
		c.HTML(401, "login.html", gin.H{
			"error": "Unautharized access user is blocked",
		})
		return
	}

	tokenstring, err := helpers.GenerateJWTToken(user.Email, user.User_type, int(user.User_id))
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(tokenstring)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_user", tokenstring, 3600*24*30, "", "", false, false)

	c.Redirect(303, "/")
}
func Logout(c *gin.Context) {
	c.SetCookie("jwt_user", "", -1, "", "", false, false)
	c.Redirect(303, "/user/login")
}
func Home(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
func About(c *gin.Context) {
	c.HTML(200, "about.html", nil)
}
func Gallery(c *gin.Context) {
	c.HTML(200, "gallery.html", nil)
}
func Testimonial(c *gin.Context) {
	c.HTML(200, "testimonial.html", nil)
}
func Contact(c *gin.Context) {
	c.HTML(200, "contact.html", nil)
}
func News(c *gin.Context) {
	c.HTML(200, "news.html", nil)
}
