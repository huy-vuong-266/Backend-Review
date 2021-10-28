package user

import (
	"Backend-Review/constants"
	"Backend-Review/model"
	"Backend-Review/service"
	"time"

	"Backend-Review/util"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

func GetUserInfoHandler(s *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {

		token := c.Get("token")
		tok, ok := token.(string)
		if !ok {

			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant parse token"},
			})
		}

		userID, err := s.AuthenService.GetUserIDByAccesstoken(tok)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant get user id"},
			})
		}

		user, err := s.UserService.GetUserByID(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant get user info"},
			})
		}

		return c.JSON(http.StatusOK, model.StandardResponse{
			Success:  true,
			Response: user,
			Error:    []string{},
		})
	}
}

func LoginHandler(s *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type loginRequest struct {
			Phone    string `json:"phone" validate:"required,min=10,max=12"`
			Password string `json:"password"  validate:"required,min=8,max=50"`
		}
		req := new(loginRequest)
		c.Bind(req)
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, model.StandardResponse{
				Success:  false,
				Response: "",
				Error:    []string{err.Error()},
			})
		}

		isValidPhone := util.ValidatePhone(req.Phone)
		if !isValidPhone {
			return echo.NewHTTPError(http.StatusBadRequest, model.StandardResponse{
				Success:  false,
				Response: "",
				Error:    []string{constants.ErrorInvalidPhoneNum},
			})
		}

		userRes, err := s.UserService.GetUserByPhone(req.Phone)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: "Get User by phone failed",
				Error:    []string{err.Error()},
			})
		}

		user, ok := userRes.(*model.User)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: "Parse user failed",
				Error:    []string{err.Error()},
			})
		}

		if encryptPw := util.Hashing(req.Password, user.Salt); encryptPw != user.Encrypted_PW {
			return echo.NewHTTPError(http.StatusForbidden, model.StandardResponse{
				Success:  false,
				Response: "Wrong phone/password",
				Error:    []string{"Wrong phone/password"},
			})
		}
		token := ""
		for token == "" {
			token = util.GenerateToken(user.UserID.String(), time.Now().Unix())
			exist := s.AuthenService.CheckIfTokenExist(token)
			if exist {
				token = ""
			}
		}

		t := &model.Token{
			UserID:    user.UserID,
			Token:     token,
			CreatedAt: time.Now().Unix(),
		}

		err = s.AuthenService.CreateToken(t)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: "Create token failed",
				Error:    []string{err.Error()},
			})
		}

		return c.JSON(http.StatusOK, model.StandardResponse{
			Success:  true,
			Response: token,
			Error:    []string{},
		})
	}
}

func RegisterHandler(s *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type registerRequest struct {
			Phone    string `json:"phone" validate:"required,min=10,max=12"`
			Name     string `json:"name"  validate:"required,max=100"`
			Email    string `json:"email" validate:"omitempty,max=255,email"`
			Password string `json:"password"  validate:"required,min=8,max=50"`
		}
		req := new(registerRequest)
		c.Bind(req)
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, model.StandardResponse{
				Success:  false,
				Response: "",
				Error:    []string{err.Error()},
			})
		}

		isValidPhone := util.ValidatePhone(req.Phone)
		if !isValidPhone {
			return echo.NewHTTPError(http.StatusBadRequest, model.StandardResponse{
				Success:  false,
				Response: "",
				Error:    []string{constants.ErrorInvalidPhoneNum},
			})
		}

		_, err := s.UserService.GetUserByPhone(req.Phone)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: "Get User by phone failed",
				Error:    []string{err.Error()},
			})
		}
		if gorm.IsRecordNotFoundError(err) {

			salt := util.GenerateSalt()

			encryptPw := util.Hashing(req.Password, salt)

			newUser := &model.User{
				UserID:       uuid.NewV4(),
				CreatedAt:    time.Now().Unix(),
				UpdatedAt:    time.Now().Unix(),
				Fullname:     req.Name,
				Phone:        req.Phone,
				Email:        req.Email,
				Encrypted_PW: encryptPw,
				Budget:       0,
				Salt:         salt,
				Status:       constants.UserStatus.StatusEnable,
			}
			err := s.UserService.CreateUser(newUser)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
					Success:  false,
					Response: "Create user failed",
					Error:    []string{err.Error()},
				})
			}

			return c.JSON(http.StatusOK, model.StandardResponse{
				Success:  true,
				Response: "register success",
				Error:    []string{},
			})
		}

		return echo.NewHTTPError(http.StatusUnprocessableEntity, model.StandardResponse{
			Success:  false,
			Response: "Phone number has been used",
			Error:    []string{"Phone number has been used"},
		})
	}
}

func AddFundHandler(s *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type addFundRequest struct {
			Amount int64 `json:"amount"`
		}

		req := new(addFundRequest)
		c.Bind(req)
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, model.StandardResponse{
				Success:  false,
				Response: "",
				Error:    []string{err.Error()},
			})
		}

		token := c.Get("token")
		tok, ok := token.(string)
		if !ok {

			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant parse token"},
			})
		}

		userID, err := s.AuthenService.GetUserIDByAccesstoken(tok)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant get user id"},
			})
		}

		userRes, err := s.UserService.GetUserByID(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant get user info"},
			})
		}
		user, ok := userRes.(*model.User)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant parse user info"},
			})
		}

		newOrder := &model.Order{
			OrderID:       uuid.NewV4(),
			CreatedAt:     time.Now().Unix(),
			UpdatedAt:     time.Now().Unix(),
			UserID:        user.UserID,
			Amount:        int64(req.Amount),
			Type:          constants.TransactionType.AddFund,
			TransactionID: uuid.Nil,
			Status:        constants.OrderStatus.Pending,
		}

		err = s.OrderService.CreateOrder(newOrder)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant create order"},
			})
		}

		code, status, errList := s.FinService.AddFund(userID, int64(req.Amount))
		if code != http.StatusOK || len(errList) != 0 {
			return c.JSON(http.StatusOK, model.StandardResponse{
				Success:  false,
				Response: status,
				Error:    errList,
			})
		}

		return c.JSON(http.StatusOK, model.StandardResponse{
			Success:  true,
			Response: "Add fund success",
			Error:    []string{},
		})
	}
}

func WithdrawHandler(s *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		type addFundRequest struct {
			Amount int64 `json:"amount"`
		}

		req := new(addFundRequest)
		c.Bind(req)
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, model.StandardResponse{
				Success:  false,
				Response: "",
				Error:    []string{err.Error()},
			})
		}

		token := c.Get("token")
		tok, ok := token.(string)
		if !ok {

			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant parse token"},
			})
		}

		userID, err := s.AuthenService.GetUserIDByAccesstoken(tok)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant get user id"},
			})
		}

		userRes, err := s.UserService.GetUserByID(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant get user info"},
			})
		}
		user, ok := userRes.(*model.User)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant parse user info"},
			})
		}

		newOrder := &model.Order{
			OrderID:       uuid.NewV4(),
			CreatedAt:     time.Now().Unix(),
			UpdatedAt:     time.Now().Unix(),
			UserID:        user.UserID,
			Amount:        int64(req.Amount),
			Type:          constants.TransactionType.Withdraw,
			TransactionID: uuid.Nil,
			Status:        constants.OrderStatus.Pending,
		}

		err = s.OrderService.CreateOrder(newOrder)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, model.StandardResponse{
				Success:  false,
				Response: constants.ErrorInternalServerError,
				Error:    []string{"cant create order"},
			})
		}

		code, status, errList := s.FinService.Withdraw(userID, int64(req.Amount))
		if code != http.StatusOK || len(errList) != 0 {
			return c.JSON(http.StatusOK, model.StandardResponse{
				Success:  false,
				Response: status,
				Error:    errList,
			})
		}

		return c.JSON(http.StatusOK, model.StandardResponse{
			Success:  true,
			Response: "Withdraw success",
			Error:    []string{},
		})
	}
}
