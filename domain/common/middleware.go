package common

import (
	"os"

	jwtware "github.com/gofiber/jwt/v3"
)

var JwtMiddleware = jwtware.New(jwtware.Config{
	SigningKey: []byte(os.Getenv("SECRET")),
})
