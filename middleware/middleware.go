package middleware

import (
	"Backend-Review/model"
	"Backend-Review/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func BasicAuthen(s *service.Service) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			token := c.Request().Header.Get("Authorization")
			if len(token) == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized, model.StandardResponse{
					Success:  false,
					Response: "Must have access token",
					Error:    []string{"Unauthorized"},
				})
			}

			exist := s.AuthenService.CheckIfTokenExist(token)
			if !exist {
				return echo.NewHTTPError(http.StatusUnauthorized, model.StandardResponse{
					Success:  false,
					Response: "Invalid access token",
					Error:    []string{"Invalid access token"},
				})
			}
			c.Set("token", token)
			return hf(c)
		}
	}
}
