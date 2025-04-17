package handler

import (
	"apisix-api/domain/auth"
	"apisix-api/domain/auth/client"
	auth_gen "apisix-api/proto/gen"
	"apisix-api/util"
	"context"
	"encoding/json"
	"log"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Login struct{}

func (h *Login) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	r := new(auth.LoginRequest)
	err := h.validate(r, c)
	if err != nil {
		log.Println("validate error : ", err.Error())
		return err
	}
	var req auth_gen.LoginRequest
	bytes, _ := json.Marshal(&r)
	_ = json.Unmarshal(bytes, &req)
	log.Println("username ", req.GetUsername())
	log.Println("email ", req.GetEmail())
	log.Println("password ", req.GetPassword())
	grpcResp, err := client.Login(ctx, &req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Println("response", err.Error())
		resp, err := h.buildErrorResponse(ctx, grpcResp, c, st.Code(), st.Message())
		if err != nil {
			return err
		}
		return resp.JSON(c)
	}
	resp, err := h.buildResponse(ctx, grpcResp, c)
	if err != nil {
		return err
	}
	return resp.JSON(c)
}
func (h *Login) buildResponse(ctx context.Context, response *auth_gen.LoginResponse, c echo.Context) (*util.Response, error) {
	resp := &util.Response{
		Code:    util.Success,
		Message: util.StatusMessage[util.Success],
		Data: map[string]interface{}{
			"token":             response.Token,
			"expired_timestamp": response.ExpiredTimestamp,
		},
	}
	return resp, nil
}
func (h *Login) buildErrorResponse(ctx context.Context, response *auth_gen.LoginResponse, c echo.Context, errorCode codes.Code, message string) (*util.Response, error) {
	resp := &util.Response{
		Code:    errorCode,
		Message: util.StatusMessage[errorCode],
		Data: map[string]interface{}{
			"message": message,
		},
	}
	return resp, nil
}
func (h *Login) validate(r *auth.LoginRequest, c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	return c.Validate(r)
}
func NewLogin() *Login {
	return &Login{}
}
