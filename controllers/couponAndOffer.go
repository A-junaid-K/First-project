package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func AddCoupon(c *gin.Context) {
	couponCode, _ := strconv.Atoi(c.Param("coupon_code"))
	var dtcoupon models.Coupon
	database.DB.Where("coupon_code=?", couponCode).First(&dtcoupon)
	if couponCode == dtcoupon.CouponId {
		c.HTML(400, "payment", gin.H{
			"error": "This coupon code already exist",
		})
		return
	}

	if len(dtcoupon.Coupon_Code) < 5 || len(dtcoupon.Coupon_Code) > 10 {
		c.HTML(400, "payment.html", gin.H{
			"error": "Coupon code must be lenght between 5 to 10",
		})
		return
	}

	if dtcoupon.Type == "fixed" || dtcoupon.Type == "percentage" {
		database.DB.Create(&models.Coupon{
			Coupon_Code:   dtcoupon.Coupon_Code,
			Starting_Time: time.Now(),
			// Ending_Time:   time.Now().Add(time.Hour * 24 * time.Duration(dtcoupon.Days)),
			Value:        dtcoupon.Value,
			Type:         dtcoupon.Type,
			Max_Discount: dtcoupon.Max_Discount,
			Min_Discount: dtcoupon.Min_Discount,
		})
		c.HTML(200, "payment.html", gin.H{
			"success": "successfully created coupon",
		})

	} else {
		c.HTML(400, "payment.html", gin.H{
			"error": "This type not applicable",
		})
		return
	}
}
func Offer(c *gin.Context) {
	dtcategory := DtTables()
	c.HTML(200, "offer.html", dtcategory)
}
func PostAddOffer(c *gin.Context) {
	layout := "2006-02-01"
	offer_name := c.Request.FormValue("offer_name")

	raw_starting_date := c.PostForm("startingDate")
	starting_date, err := time.Parse(layout, raw_starting_date)
	if err != nil {
		log.Println("faild to get starting date : ", err)
	}
	raw_expiry_date := c.PostForm("expiryDate")
	expiry_date, _ := time.Parse(layout, raw_expiry_date)

	percentage, err := strconv.Atoi(c.Request.FormValue("percentage"))

	log.Println("starting time : ", starting_date, "\n ending time : ", expiry_date)

	if err != nil {
		log.Println("failed form value : ", err)
	}

	// recieving the chosen category
	category := c.PostForm("chosencategory")

	var dtoffer models.Category_Offer
	err = database.DB.Where("offer_name=?", offer_name).First(&dtoffer).Error

	if err != nil {
		log.Println("database err  : ", err)
	}

	log.Println("dtoffer : ", dtoffer)
	if offer_name == dtoffer.Offer_Name {
		log.Println("This offer already exist : ")
		c.HTML(http.StatusBadRequest, "offer.html", gin.H{
			"error": "This offer already exist",
		})
		return
	}

	// Adding offer
	result := database.DB.Create(&models.Category_Offer{
		Offer_Name:    offer_name,
		Category_Name: category,
		Starting_Time: starting_date,
		Expiry_date:   expiry_date,
		Percentage:    uint(percentage),
		Offer:         true,
	})
	database.DB.Table("products").Where("category_name=?", category).Updates(map[string]interface{}{
		"offer_name": offer_name,
		"percentage": percentage,
	})
	err = database.DB.Table("categories").Where("name=?", category).Updates(map[string]interface{}{
		"offer_name": offer_name,
		"percentage": percentage,
	}).Error
	if err != nil {
		log.Println("errrrrrr : ", err)
	}

	if result.Error != nil {
		log.Println("Failed to add offer : ", result.Error)
		c.HTML(http.StatusBadRequest, "offer.html", gin.H{
			"error": "Failed to add offer",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin-offer")
}
func RemoveOffer(c *gin.Context) {
	db := database.DB
	delete_offer := c.Query("offer_name")
	// Find and update the rows where 'offer_name' is same
	result := db.Model(&models.Category{}).Where("offer_name = ?", delete_offer).Updates(map[string]interface{}{
		"offer_name": nil,
		"percentage": 0,
	})
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete data"})
		return
	}

	// Find and update the rows where 'offer_name' is same
	result2 := db.Model(&models.Product{}).Where("offer_name = ?", delete_offer).Updates(map[string]interface{}{
		"offer_name": nil,
		"percentage": 0,
	})
	if result2.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete data"})
		return
	}

	// Find and update the rows where 'offer_name' is "Casual sale"
	// result3 := db.Model(&models.Category_Offer{}).Where("offer_name = ?", delete_offer).Updates(map[string]interface{}{
	// 	"offer_name": nil,
	// 	"percentage": 0,
	// })
	db.Where("offer_name=?", delete_offer).Delete(&models.Category_Offer{})

	c.Redirect(303, "/admin-category")
}
