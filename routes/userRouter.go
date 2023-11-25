package routes

import (
	"os"

	"github.com/first_project/controllers"
	"github.com/first_project/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")
	r.GET("/", controllers.Home)
	// Use cookie-based sessions
	store := cookie.NewStore([]byte(os.Getenv("SECRET_KEY")))
	r.Use(sessions.Sessions("mysession", store))

	router := r.Group("/user")
	{
		//   User
		router.GET("/signup", controllers.SignUp)
		router.POST("/signup", controllers.PostSignUp)

		router.GET("/varifyotp", controllers.VarifyOtp)
		router.POST("/varifyotp", controllers.PostVarifyOtp)

		router.GET("/login", controllers.Login)
		router.POST("/login", controllers.Postlogin)

		router.GET("/login/forgot-password", controllers.ForgotPassword)
		router.POST("/login/forgot-password", controllers.PostForgotPassword)

		router.GET("/log-out", middleware.UserAuthentication, controllers.Logout)

		// User details
		router.GET("/user-details", middleware.UserAuthentication, controllers.ListUserDetails)

		// Address
		router.GET("/add-address", middleware.UserAuthentication, controllers.AddAddress)
		router.POST("/add-address", middleware.UserAuthentication, controllers.PostAddAddress)
		router.GET("/edit-address/:adrid", middleware.UserAuthentication, controllers.EditAddress)
		router.POST("/edit-address/:adrid", middleware.UserAuthentication, controllers.PostEditAddress)
		router.GET("/remove-address/:adrid", middleware.UserAuthentication, controllers.RemoveAddress)
		router.GET("/primary-address/:adrid", middleware.UserAuthentication, controllers.PrimaryAddress)

		//edit profile
		router.GET("/edit-profile", middleware.UserAuthentication, controllers.Editprofile)
		router.POST("/edit-profile", middleware.UserAuthentication, controllers.PostEditprofile)

		//Product
		router.GET("/products-list", controllers.Listproducts)
		router.GET("/product-details/:id", controllers.ProductDetails)

		//Filter
		router.GET("/products-list/category", controllers.FilterCategory)
		router.GET("/products-list/brand", controllers.FilterBrand)

		//Wishlist
		router.GET("/wishlist/:id", middleware.UserAuthentication, controllers.AddToWishlist)
		router.GET("/wishlist", middleware.UserAuthentication, controllers.ListWishlist)
		router.GET("/remove-from-wishlist/:productid", middleware.UserAuthentication, controllers.RemoveFromWishlist)

		// Order
		router.GET("/orders", middleware.UserAuthentication, controllers.Userorder)
		router.GET("/cancel-order/:orderitem_id", middleware.UserAuthentication, controllers.CancelOrder)
		router.GET("/return-order-form/:orderitem_id", middleware.UserAuthentication, controllers.Reason)
		router.POST("/return-order/:orderitem_id", middleware.UserAuthentication, controllers.ReturnOrder)

		//Cart
		router.POST("/cart/:id", middleware.UserAuthentication, controllers.AddtoCart)
		router.GET("/cart", middleware.UserAuthentication, controllers.ListCart)
		router.GET("/remove-from-cart/:productid", middleware.UserAuthentication, controllers.RemoveFromCart)

		// Search
		router.POST("/search",middleware.UserAuthentication,controllers.Search)

		// ------------Payment-----------//

		// Checkout
		router.GET("/checkout", middleware.UserAuthentication, controllers.Checkout)
		router.POST("/checkout", middleware.UserAuthentication, controllers.ApplyCoupon, controllers.Wallet, controllers.PostCheckout)

		//Instand Purchase
		router.GET("/buy-now/:product_id", middleware.UserAuthentication, controllers.BuyNow)
		router.POST("/single-checkout", middleware.UserAuthentication, controllers.ApplyCoupon, controllers.Wallet, controllers.PostBuyCheckout)

		router.GET("/payment-single-cod", middleware.UserAuthentication, controllers.GetSingleCod)
		router.GET("/payment-single-cod-success", middleware.UserAuthentication, controllers.SingleCod)

		router.GET("/payment-single-razorpay", middleware.UserAuthentication, controllers.SingleRazorpay)
		router.GET("/payment-single-razorpay-success", middleware.UserAuthentication, controllers.SingleRazorpaySuccess)

		router.GET("/payment-single-wallet", middleware.UserAuthentication, controllers.SinglePaywithWallet)
		router.GET("/payment-single-wallet-success", middleware.UserAuthentication, controllers.SingleWalletSuccess)

		// COD
		router.GET("/payment-cod", middleware.UserAuthentication, controllers.GetCod)
		router.GET("/payment-cod-success", middleware.UserAuthentication, controllers.Cod)

		// Razor Pay
		router.GET("/payment-razorpay", middleware.UserAuthentication, controllers.RazorPay)
		router.GET("/payment-razorpay-success", middleware.UserAuthentication, controllers.RazorpaySuccess)

		// Wallet
		router.GET("/payment-wallet", middleware.UserAuthentication, controllers.PaywithWallet)
		router.GET("/payment-wallet-success", middleware.UserAuthentication, controllers.WalletSuccess)

		// Success
		router.GET("/payment-success", middleware.UserAuthentication, controllers.Success)
	}

}
