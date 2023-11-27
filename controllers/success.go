package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context) {
	// user, _ := c.Get("user")
	// userid := user.(models.User).User_id

	id := c.Query("id")

	log.Println("id :", id)

	pid, err := strconv.Atoi(id)

	if err != nil {
		log.Println("err in suuccess : ", err)
	}

	// database.DB.Model(&models.Payment{}).Select("payment_id").Where("user_id=?", userid).Scan(&pid)

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})

	fmt.Println("success payment id  : ", pid)

}
