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
	payload := `{
		"message": {
			"attributes": {
			  "attr": "attr"
			},
			"data": "aGVsbG8=",
			"messageId": "1059130155449808",
			"message_id": "1059130155449808",
			"publishTime": "2020-03-22T01:42:29.391Z",
			"publish_time": "2020-03-22T01:42:29.391Z"
		  },
		  "subscription": "projects/test-project/subscriptions/test-sub"
	}`
	req, _ := http.NewRequest("POST", "/", strings.NewReader(payload)) // aGVsbG8= = base64 of `hello`
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))

	app := fiber.New()
	app.Use(New(Config{Debug: false}))
	app.Post("/", func(c *fiber.Ctx) {
		msg := c.Locals(LocalsKey).(*Message)
		if "hello" != string(msg.Message.Data) {
			t.Errorf("PubSub data should be `hello` not `%s`", string(msg.Message.Data))
		}
		if "1059130155449808" != msg.Message.ID {
			t.Error("Wrong message ID received")
		}
		if "2020-03-22T01:42:29.391Z" != msg.Message.PublishTime {
			t.Error("Wrong publish time received")
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
		data := c.Locals(LocalsKey)
		if data != nil {
			t.Error("skip should not set Locals")
		}
		c.SendStatus(http.StatusOK)
	})
	app.Test(req)
}
