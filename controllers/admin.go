package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/first_project/database"
	"github.com/first_project/helpers"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Adminlogin(c *gin.Context) {
	c.HTML(http.StatusOK, "adminLogin.html", nil)
}
func PostAdminlogin(c *gin.Context) {
	type adminDetails struct {
		Email    string
		Password string
	}
	var admin adminDetails
	admin.Email = c.Request.FormValue("email")
	admin.Password = c.Request.FormValue("password")

	Email := os.Getenv("ADMIN_NAME")
	Password := os.Getenv("ADMIN_PASSWORD")
	// ---------cheking email & password----------
	if admin.Email != Email || admin.Password != Password {
		//log.Println("Unauthorized access invalid username or password")
		c.HTML(http.StatusUnauthorized, "adminLogin.html", gin.H{
			"error": "Unauthorized access invalid username or password",
		})
		return
	}
	//-----------generating token------------
	tokenstring, err := helpers.GenerateJWTToken(Email, "admin", 0)
	if err != nil {
		log.Println("Failed to generate jwt", err)
		return
	}

	//----------set token into browser-------
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_admin", tokenstring, 3600*24*30, "", "", true, false)
	c.Redirect(303, "/users-list")
}
func Listusers(c *gin.Context) {
	type user struct {
		User_id   uint
		Name      string
		Email     string
		Phone     string
		IsBlocked bool
	}
	var users []user
	result := database.DB.Table("users").Select("user_id,name,email,phone,is_blocked").Scan(&users)
	if result.Error != nil {
		log.Println(result.Error)
	}
	c.HTML(http.StatusOK, "users-list.html", users)
}
func Blockuser(c *gin.Context) {
	//get the id
	id := c.Param("user_id")
	initid, _ := strconv.Atoi(id)
	var user models.User
	result := database.DB.First(&user, initid)
	if result.Error != nil {
		log.Println(result.Error)
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": result.Error,
		})
		return
	}
	if user.IsBlocked {
		log.Println("This user is already blocked")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": "This user is already blocked",
		})
		return
	}
	//blocking the user
	result = database.DB.Model(&models.User{}).Where("user_id=?", initid).Update("is_blocked", true)
	if result.Error != nil {
		log.Println("Failed to blocking user")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": "Failed to blocking user",
		})
		return
	}
	// c.HTML(http.StatusOK, "users-list.html", gin.H{
	// 	"message": "Successfully blocked" + user.Name,
	// })
	c.Redirect(http.StatusSeeOther, "/users-list")
}
func Unblockuser(c *gin.Context) {

	//get the id
	id := c.Param("user_id")
	initid, _ := strconv.Atoi(id)
	var user models.User
	result := database.DB.First(&user, initid)
	if result.Error != nil {
		log.Println(result.Error)
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": result.Error,
		})
		return
	}
	//checking user blocked or not
	if !user.IsBlocked {
		log.Println("This user is already unblocked")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": "This user is already unblocked",
		})
		return
	}
	//Unblocking the user
	result = database.DB.Model(&models.User{}).Where("user_id=?", initid).Update("is_blocked", false)
	if result.Error != nil {
		log.Println("Failed to unblock user")
		c.HTML(http.StatusBadRequest, "users-list.html", gin.H{
			"error": "Failed to unblock user",
		})
		return
	}
	// c.HTML(http.StatusOK, "users-list.html", gin.H{
	// 	"message": "Successfully unblocked" + user.Name,
	// })
	log.Println("Successfully unblocked " + user.Name)
	c.Redirect(http.StatusSeeOther, "/users-list")
}
