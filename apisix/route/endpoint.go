package route

import (
	auth "apisix-api/domain/auth/handler"

	"github.com/labstack/echo/v4"
)

// Handler endpoint to use it later
type Handler interface {
	Handle(c echo.Context) (err error)
}

var endpoint = map[string]Handler{
	//auth
	"login": auth.NewLogin(),
}
