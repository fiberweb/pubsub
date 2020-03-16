package pubsub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber"
)

// pubsubMessage is a format of the pubsub payload
type pubsubMessage struct {
	Message struct {
		Data []byte `json:"data"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

// Config is this middleware configuration
type Config struct {
	// Skip this middleware
	Skip func(*fiber.Ctx) bool
	// Debug true will log the unmarshalled payload
	Debug bool
}

// New returns the middleware
func New(config ...Config) func(*fiber.Ctx) {
	var cfg Config
	if len(config) == 0 {
		cfg = Config{
			Debug: true,
		}
	} else {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) {
		if cfg.Skip != nil && cfg.Skip(c) {
			c.Next()
			return
		}
		// validates request method
		if http.MethodPost != c.Method() {
			println(cfg, "PubSub middleware error: request method != POST")
			c.SendStatus(http.StatusMethodNotAllowed)
			return
		}
		// unmarshal PubSub message
		var msg *pubsubMessage
		err := json.Unmarshal([]byte(c.Body()), &msg)
		if err != nil {
			println(cfg, fmt.Sprintf("PubSub middleware error: %s", err))
			c.SendStatus(http.StatusBadRequest)
			return
		}
		println(cfg, fmt.Sprintf("PubSub data: %s, msgId: %s, subId: %s", string(msg.Message.Data), msg.Message.ID, msg.Subscription))
		c.Locals("PubSubData", msg.Message.Data)
		c.Next()
	}
}

func println(cfg Config, msg string) {
	if !cfg.Debug {
		return
	}
	log.Println(msg)
}
