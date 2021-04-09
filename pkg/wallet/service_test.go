package wallet

import (
	"testing"
	"reflect"
	"github.com/darkside1809/wallet/pkg/types"
)

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}

	account, _ := svc.RegisterAccount("+992000000001")

	acc, err := svc.FindAccountByID(account.ID)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(account, acc) {
		t.Error("Accounts not found")
	}
}
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

func Test_FindAccountByID_accountNotFound(t *testing.T) {
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
				Balance: 6_000,
			},
			{
				ID:      3,
				Phone:   "992918632026",
				Balance: 10_000,
			},
		},
	}

	var expected *types.Account = nil
	account, err := svc.FindAccountByID(5)

	if !reflect.DeepEqual(expected, account) {
		t.Errorf("invalid result, expected: %v, actual %v", expected, account)
	}

	if ErrAccountNotFound != err {
		t.Errorf("invalid result, expected: %v, actual %v", ErrAccountNotFound, err)
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

	if err != ErrPaymentNotFound {
		t.Error("payment already exists")
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