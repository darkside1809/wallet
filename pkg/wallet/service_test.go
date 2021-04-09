package wallet

import (
	"testing"
	"reflect"
	"github.com/darkside1809/wallet/pkg/types"
)

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
	_, payments, err := s.addAcount(defaultTestAccount)

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

	_, payments, err := s.addAcount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	newPayment, err := s.Repeat(payment.ID)

	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

	if newPayment.AccountID != payment.AccountID {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

	if newPayment.Amount != payment.Amount {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

	if newPayment.Category != payment.Category {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

	if newPayment.Status != payment.Status {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

}