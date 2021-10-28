package fin

import (
	"Backend-Review/constants"
	"Backend-Review/model"
	"Backend-Review/service"
	"Backend-Review/storage"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type FinService struct{}

func NewFinService() service.FinServiceInterface {
	return &FinService{}
}

func (s *FinService) AddFund(userID string, amount int64) (int, interface{}, []string) {
	if len(userID) == 0 || amount == 0 {
		return 0, "", []string{}
	}
	bodyReq := model.FinOrderRequest{
		UserID: userID,
		Amount: amount,
	}

	bodyJSON, err := json.Marshal(bodyReq)
	if err != nil {
		return 0, "", []string{"unable marshal json"}
	}

	req := strings.NewReader(string(bodyJSON))

	rawRes, err := http.Post(constants.FinURL+"/add-fund", "application/json", req)
	if err != nil {

		storage.Redis.LPush(constants.AddFundJobKey, string(bodyJSON))
		return 0, "", []string{err.Error()}
	}

	defer rawRes.Body.Close()
	resData, err := ioutil.ReadAll(rawRes.Body)
	if err != nil {
		log.Println(err)
		return 0, "", []string{err.Error()}
	}
	var response model.StandardResponse
	err = json.Unmarshal(resData, &response)
	if err != nil {
		log.Println(err)
		return 0, "", []string{err.Error()}
	}

	if len(response.Error) != 0 || rawRes.StatusCode != http.StatusOK {
		storage.Redis.LPush(constants.AddFundJobKey, string(bodyJSON))
	}

	return rawRes.StatusCode, response.Response, response.Error
}

func (s *FinService) Withdraw(userID string, amount int64) (int, interface{}, []string) {

	if len(userID) == 0 || amount == 0 {
		return 0, "", []string{}
	}
	bodyReq := model.FinOrderRequest{
		UserID: userID,
		Amount: amount,
	}

	bodyJSON, err := json.Marshal(bodyReq)
	if err != nil {
		return 0, "", []string{"unable marshal json"}
	}

	req := strings.NewReader(string(bodyJSON))

	rawRes, err := http.Post(constants.FinURL+"/withdraw", "application/json", req)
	if err != nil {
		storage.Redis.LPush(constants.WithdrawJobKey, string(bodyJSON))
		return 0, "", []string{err.Error()}
	}

	defer rawRes.Body.Close()
	resData, err := ioutil.ReadAll(rawRes.Body)
	if err != nil {
		log.Println(err)
		return 0, "", []string{err.Error()}
	}
	var response model.StandardResponse
	err = json.Unmarshal(resData, &response)
	if err != nil {
		log.Println(err)
		return 0, "", []string{err.Error()}
	}

	if len(response.Error) != 0 || rawRes.StatusCode != http.StatusOK {
		storage.Redis.LPush(constants.WithdrawJobKey, string(bodyJSON))
	}

	return rawRes.StatusCode, response.Response, response.Error
}
