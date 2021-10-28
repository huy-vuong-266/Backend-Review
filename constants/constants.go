package constants

type UserStatusStruct struct {
	StatusEnable    int8
	StatusDisable   int8
	StatusWhitelist int8
}

var UserStatus UserStatusStruct = UserStatusStruct{
	StatusEnable:    1,
	StatusDisable:   2,
	StatusWhitelist: 9,
}

const FinURL string = "http://127.0.0.1:8011"

type TransactionTypeStruct struct {
	AddFund  int8
	Withdraw int8
}

var TransactionType TransactionTypeStruct = TransactionTypeStruct{
	AddFund:  1,
	Withdraw: 2,
}

type OrderStatusStruct struct {
	Success int8
	Fail    int8
	Pending int8
}

var OrderStatus OrderStatusStruct = OrderStatusStruct{
	Success: 1,
	Fail:    2,
	Pending: 9,
}

const AddFundJobKey string = "job_add_fund"
const WithdrawJobKey string = "job_withdraw"
