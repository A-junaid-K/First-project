package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {
	c.HTML(200, "address.html", nil)
}

func PostAddAddress(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_id

	building_name := c.Request.FormValue("buildingname")
	city := c.Request.FormValue("city")
	state := c.Request.FormValue("state")
	landmark := c.Request.FormValue("landmark")
	zip := c.Request.FormValue("zip")

	if len(zip) != 6 {
		c.HTML(400, "address.html", gin.H{
			"error": "Zip code must be 6 digits",
		})
		return
	}

	//Checking, Did he have an Address
	var adr models.Address
	database.DB.Where("user_id=?", userId).Find(&adr)

	if adr == (models.Address{}) {
		err := database.DB.Create(&models.Address{
			Building_Name: building_name,
			City:          city,
			State:         state,
			Landmark:      landmark,
			Zip_code:      zip,
			User_ID:       userId,
			Primary:       true,
		}).Error
		if err != nil {
			c.HTML(400, "address.html", gin.H{"error": "Failed to first add address"})
			return
		}
		log.Println("Primary Address has added")
	} else {
		err := database.DB.Create(&models.Address{
			Building_Name: building_name,
			City:          city,
			State:         state,
			Landmark:      landmark,
			Zip_code:      zip,
			User_ID:       userId,
		}).Error
		if err != nil {
			c.HTML(400, "address.html", gin.H{
				"error": "Failed to add address",
			})
			return
		}
	}

	c.Redirect(303, "/user/user-details")
}

func EditAddress(c *gin.Context) {
	adrid, _ := strconv.Atoi(c.Param("adrid"))

	var address models.Address
	database.DB.Where("address_id=?", adrid).First(&address)

	c.HTML(200, "editaddress.html", address)
}

func PostEditAddress(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_id
	adrid, err := strconv.Atoi(c.Param("adrid"))

	if err != nil {
		fmt.Println("Failed to get address param")
		c.HTML(400, "editaddress.html", gin.H{
			"error": "Failed to get address param",
		})
		return
	}

	building_name := c.Request.FormValue("buildingname")
	city := c.Request.FormValue("city")
	state := c.Request.FormValue("state")
	landmark := c.Request.FormValue("landmark")
	zip := c.Request.FormValue("zip")

	err = database.DB.Table("addresses").Where("address_id=? AND user_id=?", adrid, userId).Updates(map[string]interface{}{
		"building_name": building_name,
		"city":          city,
		"state":         state,
		"landmark":      landmark,
		"zip_code":      zip,
	}).Error

	if err != nil {
		fmt.Println("Failed to fetch address", err)
		c.HTML(400, "editaddress.html", gin.H{
			"error": "Failed to fetch address",
		})
		return
	}

	if len(zip) != 6 {
		fmt.Println("Zip code must be 6 digits")
		c.HTML(400, "editaddress.html", gin.H{
			"error": "Zip code must be 6 digits",
		})
		return
	}

	err = database.DB.Model(&models.Address{}).Where("address_id=? AND user_id=?", adrid, userId).Updates(map[string]interface{}{
		"building_name": building_name,
		"city":          city,
		"state":         state,
		"zip_code":      zip,
		"landmark":      landmark,
	}).Error

	if err != nil {
		fmt.Println("Failed to edit address")
		c.HTML(400, "editaddress.html", gin.H{
			"error": "Failed to edit address",
		})
		return
	}

	c.Redirect(303, "/user/user-details")
}

func RemoveAddress(c *gin.Context) {
	adrid, _ := strconv.Atoi(c.Param("adrid"))
	database.DB.Where("address_id=?", adrid).Delete(&models.Address{})
	log.Println("adderress Successfully Removed")
	c.Redirect(303, "/user/user-details")
}

func PrimaryAddress(c *gin.Context) {
	adrid, _ := strconv.Atoi(c.Param("adrid"))
	err := database.DB.Model(&models.Address{}).Where("address_id=?", adrid).Update("primary", true).Error
	if err != nil {
		fmt.Println("Failed to make primary address")
		c.HTML(400, "editaddress.html", gin.H{"error": "Failed to make primary address"})
		return
	}
	log.Println("adderress Successfully Updated as primary")
	c.Redirect(303, "/user/user-details")
}

//----------------------------------------user profile-------------------------------------//

func Editprofile(c *gin.Context) {

	user, _ := c.Get("user")
	userId := user.(models.User).User_id

	var userr models.User
	database.DB.Where("user_id=?", userId).First(&userr)

	c.HTML(200, "editprofile.html", userr)
}

func PostEditprofile(c *gin.Context) {

	user, _ := c.Get("user")
	userId := user.(models.User).User_id

	name := c.Request.FormValue("name")
	phone := c.Request.FormValue("number")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	confpassword := c.Request.FormValue("confpassword")

	err := namevalidator(name)
	if err != nil {
		c.HTML(400, "editprofile.html", gin.H{
			"error": err,
		})
		return
	}

	err = numbervalidator(phone)
	if err != nil {
		c.HTML(400, "editprofile.html", gin.H{
			"error": err,
		})
		return
	}

	err = emailvalidator(email)

	if err != nil {
		c.HTML(400, "editprofile.html", gin.H{
			"error": err,
		})
		return
	}

	err = passwordvalidator(password)
	if err != nil {
		c.HTML(400, "editprofile.html", gin.H{
			"error": err,
		})
		return
	}
	if confpassword != password {
		c.HTML(400, "editprofile.html", gin.H{
			"error": "Incorrect confirm password",
		})
		return
	}

	// var dtuser models.User

	// database.DB.Where("email=?", email).First(&dtuser)
	// if email == dtuser.Email {
	// 	c.HTML(http.StatusBadRequest, "SignUp.html", gin.H{
	// 		"error": email + " has already been used",
	// 	})
	// 	return
	// }

	err = database.DB.Table("users").Where("user_id=?", userId).Updates(map[string]interface{}{
		"name":     name,
		"email":    email,
		"phone":    phone,
		"password": password,
	}).Error

	if err != nil {
		fmt.Println("Failed to edit profile")
		c.HTML(400, "editprofile.html", gin.H{
			"error": "Failed to edit profile",
		})
		return
	}

	c.Redirect(303, "/user/user-details")

}

func ListUserDetails(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_id

	var addresses []models.Address
	var users []models.User

	database.DB.Where("user_id=?", userId).Find(&addresses)
	database.DB.Where("user_id=?", userId).Find(&users)
	// database.DB.Find(&addresses)
	// database.DB.Find(&users)
	userdetails := struct {
		Addresses []models.Address
		Users     []models.User
	}{
		Addresses: addresses,
		Users:     users,
	}
	c.HTML(200, "user-details.html", userdetails)
}
