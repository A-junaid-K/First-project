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
	layout := "02/01/2006"
	offer_name := c.Request.FormValue("offer_name")
	starting_date, _ := time.Parse(layout, c.PostForm("startingDate"))
	expiry_date, _ := time.Parse(layout, c.PostForm("expiryDate"))
	percentage, err := strconv.Atoi(c.Request.FormValue("percentage"))

	log.Println("offer name : ", offer_name)

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

	if result.Error != nil {
		log.Println("Failed to add offer : ", result.Error)
		c.HTML(http.StatusBadRequest, "offer.html", gin.H{
			"error": "Failed to add offer",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin-offer")
}
