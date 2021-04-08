package wallet

import (
	"github.com/darkside1809/wallet/pkg/types"
	"errors"
)
var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than 0")
var ErrAccountNotFound = errors.New("account not found")
var exErr error = nil

type Service struct {
	nextAccountId 	int64
	accounts 		[]*types.Account
	payments 		[]*types.Payment
}

func (s *Service)FindAccountById(account int64) (*types.Account, error){
	var accID *types.Account

	for _, acc := range s.accounts{

		if acc.ID == account{
			accID = acc 
			break
		}
	}

	if  accID == nil {
		return nil, ErrAccountNotFound
	}
	
	return accID, nil
}