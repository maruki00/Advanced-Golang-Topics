package route

import (
	"apisix-api/util"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Route for mapping from json file
type Route struct {
	Path       string   `json:"path"`
	Method     string   `json:"method"`
	Module     string   `json:"module"`
	Tag        string   `json:"tag"`
	Endpoint   string   `json:"endpoint_filter"`
	Middleware []string `json:"middleware"`
}
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Init gateway router
func Init() *echo.Echo {
	routes := loadRoutes("./route/gate/")
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	// Set Bundle MiddleWare
	e.Use(middleware.RequestID())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowHeaders, echo.HeaderContentDisposition, "X-Request-Id", "device-id", "X-Summary", "X-Account-Number", "X-Business-Name", "client-secret", "X-CSRF-Token"},
		ExposeHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowHeaders, echo.HeaderContentDisposition, "X-Request-Id", "device-id", "X-Summary", "X-Account-Number", "X-Business-Name", "client-secret", "X-CSRF-Token"},
		AllowMethods:  []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	for _, route := range routes {
		e.Add(route.Method, route.Path, endpoint[route.Endpoint].Handle, chainMiddleware(route)...)
	}
	return e
}

func loadRoutes(filePath string) []Route {
	var routes []Route
	files, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatalf("Failed to load file: %v", err)
	}
	for _, file := range files {
		byteFile, err := os.ReadFile(filePath + "/" + file.Name())
		if err != nil {
			log.Fatalf("Failed to load file: %v", err)
		}
		var tmp []Route
		if err := util.Json.Unmarshal(byteFile, &tmp); err != nil {
			log.Fatalf("Failed to marshal file: %v", err)
		}
		routes = append(routes, tmp...)
	}
	return routes
}
func chainMiddleware(route Route) []echo.MiddlewareFunc {
	var mwHandlers []echo.MiddlewareFunc
	// init mw for router ,attach router properties
	// mwHandlers = append(mwHandlers, customMw.CheckToken(util.ContextRouterKey, route.Tag))
	// for _, v := range route.Middleware {
	// 	mwHandlers = append(mwHandlers, middlewareHandler[v])
	// }
	return mwHandlers
}
