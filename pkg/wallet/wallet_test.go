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
	account, err := svc.FindAccountById(2)

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
	account, err := svc.FindAccountById(5)

	if !reflect.DeepEqual(expected, account) {
		t.Errorf("invalid result, expected: %v, actual %v", expected, account)
	}

	if ErrAccountNotFound != err {
		t.Errorf("invalid result, expected: %v, actual %v", ErrAccountNotFound, err)
	}
}