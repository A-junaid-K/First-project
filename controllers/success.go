package controllers

import (
	"fmt"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	var pid int

	database.DB.Model(&models.Payment{}).Select("payment_id").Where("user_id=?", userid).Scan(&pid)

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})

	fmt.Println("success payment id  : ", pid)

}
