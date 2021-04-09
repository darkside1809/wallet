package wallet

import (
	"testing"
	"reflect"
	"github.com/darkside1809/wallet/pkg/types"
)

func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := &Service{}

	_, err := svc.FindAccountByID(123)

	if err != ErrAccountNotFound {
		t.Error(err)
	}
}
func Test_FindAccountByID_success(t *testing.T) {
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
		t.Error(ErrPhoneRegistered)
	}

	_, err := svc.Pay(account.ID, 1000, "auto")

	if err == ErrAmountMustBePositive {
		t.Error(ErrAmountMustBePositive)
	}

	if err == nil {
		t.Error(ErrPaymentNotFound)
	}
}

func TestService_FindPaymentByID_notFound(t *testing.T) {
	svc := &Service{}

	_, err := svc.FindPaymentByID("aaa")

	if err == ErrPaymentNotFound {
		t.Error(ErrAccountNotFound)
	}
}

func Test_Reject_success(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAccount("992918632026")

	if err == ErrPhoneRegistered {
		t.Error(ErrPhoneRegistered)
	}

	_, errPay := svc.Pay(account.ID, 1000, "auto")

	if errPay != ErrNotEnoughBalance {
		t.Error(ErrNotEnoughBalance)
	}

}
func Test_Reject_paymentNotFound(t *testing.T) {
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