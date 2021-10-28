package internal

import (
	userapi "Backend-Review/internal/user"
	middleware "Backend-Review/middleware"
	service "Backend-Review/service"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, s *service.Service) http.Handler {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, true)
	})

	e.POST("/login", userapi.LoginHandler(s))
	e.POST("/register", userapi.RegisterHandler(s))

	u := e.Group("/user", middleware.BasicAuthen(s))
	{
		u.GET("/info", userapi.GetUserInfoHandler(s))
		u.POST("/add-fund", userapi.AddFundHandler(s))
		u.POST("/withdraw", userapi.WithdrawHandler(s))
	}

	return e.Server.Handler
}
