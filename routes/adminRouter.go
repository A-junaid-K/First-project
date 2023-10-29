package routes

import (
	"github.com/first_project/controllers"
	"github.com/first_project/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	// router := r.Group("/admin")
	// {
	r.GET("/admin-login", controllers.Adminlogin)
	r.POST("/admin-login", controllers.PostAdminlogin)

	// user
	r.GET("/users-list", middleware.AdminAuthentication, controllers.Listusers)
	r.GET("/block/:user_id", middleware.AdminAuthentication, controllers.Blockuser)
	r.GET("/unblock/:user_id", middleware.AdminAuthentication, controllers.Unblockuser)

	// products
	r.GET("/add-product", middleware.AdminAuthentication, controllers.Addproducts)
	r.POST("/add-product", middleware.AdminAuthentication, controllers.PostAddproducts)
	r.GET("/admin-products-list", middleware.AdminAuthentication, controllers.AdminListproducts)
	r.GET("/edit-product/:id", middleware.AdminAuthentication, controllers.Editproduct)
	r.POST("/edit-product/:id", middleware.AdminAuthentication, controllers.PostEditproduct)
	r.GET("/delete-product/:prdctid", middleware.AdminAuthentication, controllers.Deleteproduct)

	// Brand
	r.POST("/admin-products-list-addbrand", middleware.AdminAuthentication, controllers.AddBrand)
	// r.GET("/products-list", middleware.AdminAuthentication, controllers.ListBrand)

	// Add Category
	r.POST("/admin-products-list-addcategory", middleware.AdminAuthentication, controllers.AddCategory)

	// Category Offers

	// Coupen
	//add coupoun
	r.POST("/addcoupon/:coupon_code", middleware.AdminAuthentication, controllers.AddCoupon)
	// r.GET("/listcoupons", middleware.AdminAuthentication, controllers.ListCoupons)
	// r.PUT("/cancelcoupon/:coupon_id", middleware.AdminAuthentication, controllers.CancelCoupon)

	r.GET("/generate-sales-report", middleware.AdminAuthentication, controllers.Getsalesreport)
	r.POST("/generate-sales-reports", middleware.AdminAuthentication, controllers.SalesReport)

}
