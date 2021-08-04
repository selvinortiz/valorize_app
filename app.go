package main

import (
	"net/http"
	"valorize-app/config"
	"valorize-app/handlers"
	appmiddleware "valorize-app/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stripe/stripe-go/v72"
)

func main() {

	cfg := config.NewConfig()
	s := handlers.NewServer(cfg)
	stripe.Key = "sk_test_51JGBbjBhSkl0qU1AdCzBjVv6N0Z2xyYqHTfYPOECkuFdl4lA9fyLIz6lHrKP702wlybuwcfh1rB7ljG8zUzzta7k00ytyRYt2d"
	e := *s.Echo
	payment := handlers.NewPaymentHandler(s)
	auth := handlers.NewAuthHandler(s)
	eth := handlers.NewEthHandler(s)
	user := handlers.NewUserHandler(s)
	wallet := handlers.NewWalletHandler(s)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://valorize.local:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	e.Static("/*", "app/dist")
	e.Static("/static/images", "dist/images")

	e.GET("/success", func(c echo.Context) error {
		return c.String(http.StatusOK, "Success")
	})
	e.GET("/cancel", func(c echo.Context) error {
		return c.String(http.StatusOK, "Payment error")
	})
	e.POST("/login", auth.Login)
	e.GET("/logout", auth.Logout)
	e.POST("/register", auth.Register)
	e.GET("/create-checkout-session", payment.CreateCheckoutSession)
	e.GET("/eth", eth.Ping)
	e.POST("/payments/successhook", payment.OnPaymentAccepted)

	api := e.Group("/api/v0")

	me := api.Group("/me", appmiddleware.AuthMiddleware)
	me.GET("", auth.ShowUser)
	me.PUT("/picture", auth.UpdatePicture)
	me.PUT("/profile", auth.UpdateProfile)

	userGroup := api.Group("/users")
	userGroup.GET("/:username", user.Show)
	userGroup.GET("/:username/wallets", wallet.Index)

	r := api.Group("/admin", appmiddleware.AuthMiddleware)
		r.POST("/wallet", eth.CreateWalletFromRequest)
		r.POST("/deploy", eth.DeployCreatorToken)
	e.Logger.Fatal(e.Start(":1323"))
}
