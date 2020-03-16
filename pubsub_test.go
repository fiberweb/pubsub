package pubsub

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber"
)

func Test_MethodNotAllowed(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	app := fiber.New()
	app.Use(New(Config{Debug: false}))
	resp, _ := app.Test(req)

	if http.StatusMethodNotAllowed != resp.StatusCode {
		t.Error("should return 405 if method is not POST")
	}
}

func Test_UnmarshalMessageError(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader("invalid body"))
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	app := fiber.New()
	app.Use(New(Config{Debug: false}))
	resp, _ := app.Test(req)

	if http.StatusBadRequest != resp.StatusCode {
		t.Error("should return 400 if method is not POST")
	}
}

func Test_LocalsData(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"message": {"data": "aGVsbG8="}}`)) // aGVsbG8= = base64 of `hello`
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	app := fiber.New()
	app.Use(New(Config{Debug: false}))
	app.Post("/", func(c *fiber.Ctx) {
		data := c.Locals("PubSubData").([]byte)
		if "hello" != string(data) {
			t.Errorf("PubSub data should be `hello` not `%s`", string(data))
		}
		c.SendStatus(http.StatusOK)
	})
	app.Test(req)
}

func Test_Skip(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"message": {"data": "aGVsbG8="}}`)) // aGVsbG8= = base64 of `hello`
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	app := fiber.New()
	app.Use(New(Config{
		Debug: false,
		Skip: func(c *fiber.Ctx) bool {
			return true
		},
	}))
	app.Post("/", func(c *fiber.Ctx) {
		data := c.Locals("PubSubData")
		if data != nil {
			t.Error("skip should not set Locals")
		}
		c.SendStatus(http.StatusOK)
	})
	app.Test(req)
}
