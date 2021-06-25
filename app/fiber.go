package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/pooladkhay/tinier-shortener-service/client"
	"github.com/pooladkhay/tinier-shortener-service/controller"
	"github.com/pooladkhay/tinier-shortener-service/helper/errs"
	"github.com/pooladkhay/tinier-shortener-service/repository"
	"github.com/pooladkhay/tinier-shortener-service/service"
)

func StartFiber() {
	cassandraSession := client.CassandraSession()
	redis := client.NewRedis()

	urlRepo := repository.NewUrl(cassandraSession)
	cacheRepo := repository.NewCache(redis)

	urlService := service.NewUrl(urlRepo, cacheRepo)
	urlController := controller.NewUrl(urlService)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(cors.New())
	// app.Use(logger.New())
	app.Use(recover.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")
	short := v1.Group("/short")

	//  private routes
	short.Get("/all", jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("USER_JWT_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, _ error) error {
			return c.Status(401).JSON(errs.NewUnauthorizedError("unauthorized"))
		}}), urlController.GetAll)

	// public route
	short.Get("/:hash", urlController.GetByHash)

	// semi-private routes
	// these routes work with access token and secret
	short.Post("/", jwtValidarionExtraction, urlController.Shorten)
	short.Delete("/", jwtValidarionExtraction, urlController.Delete)

	app.All("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNotFound)
	})

	p := os.Getenv("PORT")
	fmt.Printf("Auth api is listening on %s...\r\n", p)
	err := app.Listen(fmt.Sprintf(":%s", p))
	if err != nil {
		log.Fatalf("Err listening on %s: %s\r\n", p, err)
	}
}

func jwtValidarionExtraction(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	l := len("Bearer")
	if len(auth) > l+1 && strings.EqualFold(auth[:l], "Bearer") {
		tokenString := auth[l+1:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("USER_JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("user", claims["user"])
		} else {
			fmt.Println(err)
		}
	}
	return c.Next()
}
