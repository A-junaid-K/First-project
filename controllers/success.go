package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context) {

	id := c.Query("id")

	pid, err := strconv.Atoi(id)

	if err != nil {
		log.Println("err in success : ", err)
	}

	// database.DB.Model(&models.Payment{}).Select("payment_id").Where("user_id=?", userid).Scan(&pid)

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})

	fmt.Println("success payment id  : ", pid)

}
