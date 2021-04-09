package wallet

import (
	"github.com/darkside1809/wallet/pkg/types"
	"github.com/google/uuid"
	"errors"
	"fmt"
)
var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than 0")
var ErrNotEnoughBalance = errors.New("not enought balance in account")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var exErr error = nil
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

type Service struct {
	nextAccountID 	int64
	accounts 		[]*types.Account
	payments 		[]*types.Payment
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

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)

	return payment, nil
}

func (s *Service)FindAccountByID(account int64) (*types.Account, error){
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

func (s *Service)FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment

	for _, pay := range s.payments {
		
		if pay.ID == paymentID {
			payment = pay
			break
		}

		if payment == nil {
			return nil, ErrPaymentNotFound
		}
	}

	return payment, nil
}

func (s *Service)Reject(paymentID string) error {
	var payCheck *types.Payment

	for _, paymentData := range s.payments {
		if paymentID == paymentData.ID {
			payCheck = paymentData
			break
		}
	}
	if payCheck == nil {
		return  ErrPaymentNotFound
	}

	account, err := s.FindAccountByID(payCheck.AccountID)

	if err == ErrAccountNotFound {
		return ErrAccountNotFound
	}

	account.Balance += payCheck.Amount
	payCheck.Amount = 0
	payCheck.Status = types.PaymentStatusFail

	return nil
}

func (s *testService) addAcount(data testAccount) (*types.Account, []*types.Payment, error) {
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

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return nil, ErrPaymentNotFound
	}

	pay, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)

	if err != nil {
		return nil, err
	}

	return pay, nil
}