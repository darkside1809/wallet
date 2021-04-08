package wallet

import (
	"testing"
	"reflect"
	"github.com/darkside1809/wallet/pkg/types"
)

func Test_FindAccountByID_success(t *testing.T) {
	svc := &Service{
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "992918632026",
				Balance: 1000,
			},
			{
				ID:      2,
				Phone:   "992918632026",
				Balance: 2000,
			},
			{
				ID:      3,
				Phone:   "992918632026",
				Balance: 3000,
			},
		},
	}

	expected := &types.Account{
		ID:      2,
		Phone:   "992918632026",
		Balance: 2000,
	}
	account, err := svc.FindAccountById(2)

	if !reflect.DeepEqual(expected, account) {
		t.Errorf("invalid result, expected: %v, actual %v", expected, account)
	}

	if expectedErr != err {
		t.Errorf("invalid result, expected: %v, actual %v", expectedErr, err)
	}
}

func Test_FindAccountByID_accountNotFound(t *testing.T) {
	svc := &Service{
		accounts: []*types.Account{
			{
				ID:      1,
				Phone:   "992918632026",
				Balance: 1000,
			},
			{
				ID:      2,
				Phone:   "992918632026",
				Balance: 2000,
			},
			{
				ID:      3,
				Phone:   "992918632026",
				Balance: 3000,
			},
		},
	}

	var expected *types.Account = nil
	account, err := svc.FindAccountById(5)

	if !reflect.DeepEqual(expected, account) {
		t.Errorf("invalid result, expected: %v, actual %v", expected, account)
	}

	if ErrAccountNotFound != err {
		t.Errorf("invalid result, expected: %v, actual %v", ErrAccountNotFound, err)
	}
}