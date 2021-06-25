package controller

import (
	"net/http"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/pooladkhay/tinier-shortener-service/domain"
	"github.com/pooladkhay/tinier-shortener-service/service"
)

type Url interface {
	Shorten(c *fiber.Ctx) error
	GetByHash(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type url struct {
	service service.Url
}

func NewUrl(s service.Url) Url {
	return &url{service: s}
}

func (u *url) Shorten(c *fiber.Ctx) error {
	reqBody := new(domain.ShortenRequest)

	if err := c.BodyParser(reqBody); err != nil {
		return err
	}
	if err := reqBody.Validate(); err != nil {
		return c.Status(err.Status).JSON(err)
	}

	reqBody.User = ""
	if c.Locals("user") != nil {
		if usr := c.Locals("user").(string); len(usr) > 0 {
			reqBody.User = usr
		}
	}

	resp, err := u.service.Shorten(reqBody)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (u *url) GetByHash(c *fiber.Ctx) error {
	if c.Params("hash") == "" {
		c.Status(http.StatusBadRequest).JSON("invalid hash param")
	}

	resp, err := u.service.GetByHash(c.Params("hash"), c.Query("secret"))
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func (u *url) GetAll(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	usr := claims["user"].(string)

	resp, err := u.service.GetAll(usr)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(resp)
}

func (u *url) Delete(c *fiber.Ctx) error {
	reqBody := new(domain.DeleteHashRequest)

	if err := c.BodyParser(reqBody); err != nil {
		return err
	}
	if err := reqBody.Validate(); err != nil {
		return c.Status(err.Status).JSON(err)
	}

	reqBody.User = ""
	if c.Locals("user") != nil {
		if usr := c.Locals("user").(string); len(usr) > 0 {
			reqBody.User = usr
		}
	}

	err := u.service.Delete(reqBody)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusOK)
}
