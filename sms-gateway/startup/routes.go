package startup

import (
	"sms-gateway/handler/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func newHTTPApplication(h handlers) *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{
			validate: validator.New(),
		},
	})
	initRoutes(app, h)

	return app
}

func initRoutes(app *fiber.App, h handlers) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{fiber.HeaderContentType, fiber.HeaderContentLength},
		AllowMethods: []string{fiber.MethodGet, fiber.MethodPost},
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

	// message
	message := v1.Group("/message")
	message.Post("/send", h.message.SendMessage)

	// report
	v1.Get("/report", h.message.GetUserReport)
}
