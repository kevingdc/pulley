package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func ToRequest(c *fiber.Ctx) *http.Request {
	request := &http.Request{}
	fasthttpadaptor.ConvertRequest(c.Context(), request, false)
	return request
}
