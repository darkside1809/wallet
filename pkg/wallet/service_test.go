package wallet

import (
	"fmt"
	"testing"
	"reflect"
	"github.com/darkside1809/wallet/pkg/types"
)


var defaultTestAccount = testAccount{
	phone:   "+992938151007",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{{
		amount:   1000_00,
		category: "auto",
	}},
}
type testAccount struct {
	phone 	types.Phone
	balance 	types.Money
	payments []struct {
		amount	types.Money
		category types.PaymentCategory
	}
}
type testService struct {
	*Service
}
func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("cant register account %v = ", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("cant deposit account %v = ", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("cant make payment %v = ", err)
		}
	}

	return account, payments, nil
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}
func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := &Service{}

	_, err := svc.FindAccountByID(123)

	if err != ErrAccountNotFound {
		t.Error(err)
	}
}
func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "992918632026",
				Balance: 1_000,
			},
			{
				ID:      2,
				Phone:   "992918632026",
				Balance: 2_000,
			},
			{
				ID:      3,
				Phone:   "992918632026",
				Balance: 3_000,
			},
		},
	}

	expected := &types.Account{
		ID:      2,
		Phone:   "992918632026",
		Balance: 2_000,
	}
	account, err := svc.FindAccountByID(2)
	if !reflect.DeepEqual(expected, account) {
		t.Errorf("invalid result, expected: %v, actual %v", expected, account)
	}

	if exErr != err {
		t.Errorf("invalid result, expected: %v, actual %v", exErr, err)
	}
}
func TestService_FindPaymentByID_success(t *testing.T) {
	svc := &Service{}

	account, errorReg := svc.RegisterAccount("+992000000001")

	if errorReg != nil {
		t.Error("error on register account")
	}
	_, err := svc.Pay(account.ID, 1000, "auto")
	if err == ErrAmountMustBePositive {
		t.Error(ErrAmountMustBePositive)
	}

	if err == nil {
		t.Error("error on pay")
	}
}

func TestService_FindPaymentByID_notFound(t *testing.T) {
	svc := &Service{}

	_, err := svc.FindPaymentByID("aaa")

	if err == ErrPaymentNotFound {
		t.Error(ErrAccountNotFound)
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)

	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)

	if err != nil {
		t.Errorf("Reject(): error = %v", err)
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): cant't find payment by id, error = %v", err)
	}

	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", err)
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)

	if err != nil {
		t.Errorf("Reject(): cant't find account by id, error = %v", err)
	}

	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't changed")
	}


}
func TestService_Reject_paymentNotFound(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAccount("992918632026")

	if err == ErrPhoneRegistered {
		t.Error(ErrPhoneRegistered)
	}

	_, errPay := svc.Pay(account.ID, 2000, "auto")

	if errPay != ErrNotEnoughBalance {
		t.Error(ErrNotEnoughBalance)
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	newPayment, nil := s.Repeat(payment.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if payment.ID == newPayment.ID {
		t.Error("repeated payment id not different")
		return
	}

	if payment.AccountID != newPayment.AccountID ||
		payment.Status != newPayment.Status ||
		payment.Category != newPayment.Category ||
		payment.Amount != newPayment.Amount {
		t.Error("some field is not equal the original")
	}
}