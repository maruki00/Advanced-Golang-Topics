package util

import (
	"fmt"

	"github.com/labstack/echo/v4"
	grpcCode "google.golang.org/grpc/codes"
)

// Response struct
type Response struct {
	Code       grpcCode.Code          `json:"code"`
	Message    string                 `json:"message,omitempty"`
	Data       interface{}            `json:"data,omitempty"`
	Pagination *Pg                    `json:"pagination,omitempty"`
	Errors     []string               `json:"errors,omitempty"`
	Header     map[string]interface{} `json:"-"`
}

// Pg struct
type Pg struct {
	CurrentPage int32  `json:"current_page"`
	PageSize    int32  `json:"page_size"`
	TotalPage   int32  `json:"total_page"`
	TotalResult int32  `json:"total_result"`
	Next        string `json:"next,omitempty"`
	Prev        string `json:"prev,omitempty"`
}

// WithPagination set response with pagination
func (r *Response) WithPagination(c echo.Context, pagination Pg) *Response {
	r.Pagination = &pagination
	page := r.Pagination.CurrentPage
	u := c.Request().URL
	if r.Pagination.NextPage() {
		q := u.Query()
		q.Set("page", fmt.Sprintf("%d", page+1))
		u.RawQuery = q.Encode()
		r.Pagination.Next = u.String()
	}
	if page > 1 {
		q := u.Query()
		q.Set("page", fmt.Sprintf("%d", page-1))
		u.RawQuery = q.Encode()
		r.Pagination.Prev = u.String()
	}
	return r
}

// JSON render response as JSON
func (r *Response) JSON(c echo.Context) error {
	for k, v := range r.Header {
		c.Response().Header().Set(k, fmt.Sprintf("%v,%v", c.Response().Header().Get(k), v))
	}
	return c.JSON(HTTPStatusFromCode(r.Code), r)
}

// NextPage determine is available next page
func (p *Pg) NextPage() bool {
	return (p.CurrentPage * p.PageSize) < p.TotalResult
}
func BindValidate(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	return c.Validate(i)
}
