package routes

import (
	"github.com/first_project/controllers"
	"github.com/first_project/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	r.GET("/admin-login", controllers.Adminlogin)
	r.POST("/admin-login", controllers.PostAdminlogin)

	// Dashboard
	r.GET("/admin-dashboard", middleware.AdminAuthentication, controllers.AdminDashboard)

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
	r.GET("/admin-brands", middleware.AdminAuthentication, controllers.DisplayBrands)
	r.POST("/admin-addbrand", middleware.AdminAuthentication, controllers.AddBrand)
	r.GET("/admin/remove-brands/:brand_id", middleware.AdminAuthentication, controllers.RemoveBrand)

	// Add Category
	r.GET("/admin-category", middleware.AdminAuthentication, controllers.DisplayCategory)
	r.POST("/admin-addcategory", middleware.AdminAuthentication, controllers.AddCategory)
	r.GET("/admin-unlistcategory/:category_id", middleware.AdminAuthentication, controllers.UnlistCategory)
	r.GET("/admin-list-category/:category_id", middleware.AdminAuthentication, controllers.ListCategory)

	// Coupoun
	r.GET("/admin-coupon", middleware.AdminAuthentication, controllers.Coupon)
	r.POST("/admin-addcoupon", middleware.AdminAuthentication, controllers.PostAddCoupon)
	r.GET("/admin/cancel-coupon", middleware.AdminAuthentication, controllers.CancelCoupon)
	r.GET("/admin/approve-coupon", middleware.AdminAuthentication, controllers.ApproveCoupon)
	r.GET("/admin/remove-coupon", middleware.AdminAuthentication, controllers.RemoveCoupon)
	r.GET("/admin/apply-coupon", middleware.AdminAuthentication, controllers.ApplyCoupon)

	//offer
	r.GET("/admin-offer", middleware.AdminAuthentication, controllers.Offer)
	r.POST("/admin/add-offer", middleware.AdminAuthentication, controllers.PostAddOffer)
	r.GET("/admin/remove-offer", middleware.AdminAuthentication, controllers.RemoveOffer)

	//Order
	r.GET("/admin-order", middleware.AdminAuthentication, controllers.Order)
	r.POST("/admin/order-status/:order_id", middleware.AdminAuthentication, controllers.PostOrder)

	//Sales
	r.GET("/admin-sales", middleware.AdminAuthentication, controllers.Sales)
	r.POST("/admin-sales-post", middleware.AdminAuthentication, controllers.Salesreport)
	r.GET("/admin/salesreport/xlsx", middleware.AdminAuthentication, controllers.DownloadExel)
	r.GET("/admin/salesreport/pdf", middleware.AdminAuthentication, controllers.Downloadpdf)

}
