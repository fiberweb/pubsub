# pubsub
Google Cloud's [PubSub](https://cloud.google.com/pubsub) request middleware for [Fiber](https://github.com/gofiber/fiber). This middleware handle the validation of PubSub request and its payload decoding.

## Install

```
go get -u github.com/fiberweb/pubsub
```

## Usage

```
package main

import (
  "encoding/json"
  "fmt"
  
  "github.com/gofiber/fiber"
  "github.com/fiberweb/pubsub"
)

// our PubSub data structure
type User struct {
  Name  string `json:"name"`
  Age   int    `json:"age"`
}

func main() {
  app := fiber.New()
  
  // use the middleware
  app.Use(pubsub.New())
  app.Post("/", func(c *fiber.Ctx) {
    msg := c.Locals(pubsub.LocalsKey).(*pubsub.Message)
    
    var user User
    if err := json.Unmarshal(msg.Message.Data, &user); err != nil {
      c.SendStatus(400)
      return
    }
    
    fmt.Println(user.Name, user.Age)
    c.Send("Ok")
  })
  app.Listen("8080")
}
```

When the middleware successfully decode the message, PubSub data will be available for the next handlers inside the Fiber context `Locals` called `PubSubMessage`.

## Configuration
You could also initialize the middleware with a config:
```
app.Use(pubsub.New(pubsub.Config{
  Debug: false,
}))
```

This middleware has only two configuration options:

```
type Config struct {
	// Skip this middleware
	Skip func(*fiber.Ctx) bool
  
	// Debug true will log any errors during validation 
        // and marshalling and log the unmarshalled payload, default: true
	Debug bool
}
```

#### Example of using Skip
```
app.Use(pubsub.New(pubsub.Config{
    Skip: func(c *fiber.Ctx) bool {
      // add your logic here
      return true // returning true will skip this middleware.
    }
}))
```
