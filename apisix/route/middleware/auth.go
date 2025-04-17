package middleware

import (
	"apisix-api/domain/auth/client"
	auth_gen "apisix-api/proto/gen"

	"apisix-api/util"
	"context"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
)

func CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		if ctx == nil {
			ctx = context.Background()
		}
		var req auth_gen.CheckTokenRequest
		headerToken := c.Request().Header.Get(echo.HeaderAuthorization)
		req.Token = strings.Replace(headerToken, "Bearer ", "", -1)
		log.Println(req.Token)
		res, err := client.CheckToken(ctx, &req)
		log.Println(res)
		if err != nil {
			st, _ := status.FromError(err)
			resp := &util.Response{
				Code:    st.Code(),
				Message: st.Message(),
				Errors:  []string{st.Message()},
			}
			return resp.JSON(c)
		}
		payload := auth_gen.CheckTokenResponse{
			Email:  res.Email,
			Status: res.Status,
		}
		// TODO wrapping context
		c.Set(util.ContextTokenValueKey, payload)
		return next(c)
	}
}
