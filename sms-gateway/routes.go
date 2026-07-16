package main

import (
	"sms-gateway/handler"
	"sms-gateway/handler/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type handlers struct {
	wallet *handler.WalletHandler
}

func initRoutes(app *fiber.App, h handlers) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Content-Type, Content-Length"},
		AllowMethods: []string{"POST, GET"},
	}))

	app.Use(logger.New())

	// we use this middleware so it enforce
	// Idempotency header for unsafe requests
	app.Use(middleware.RequireIdempotencyKey)

	// it store data in memory.
	//  we should use redis as storage in production
	app.Use(idempotency.New())

	v1 := app.Group("/api/v1/sms-gateway")

	//wallet
	wallet := v1.Group("/wallet")
	wallet.Post("/add", h.wallet.TopUpWallet)
	wallet.Get("/:user_id", h.wallet.GetWalletBalance)
}
