package wallet

import (
	"errors"
	"github.com/darkside1809/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater then 0")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("account balance least then amount")
var ErrFavoriteNotFound = errors.New("favorite payment not found")
var exErr = errors.New("doesn't match to expected")

type Service struct {
	nextAccountID	int64
	accounts		[]*types.Account
	payments		[]*types.Payment
	favorites	[]*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID: 		s.nextAccountID,
		Phone: 		phone,
		Balance: 	0,
	}

	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return err
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount

	paymentID := uuid.New().String()

	payment := &types.Payment{
		ID:			paymentID,
		AccountID: 	accountID,
		Amount: 	amount,
		Category: 	category,
		Status: 	types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)

	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
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

	return account, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	account.Balance += payment.Amount
	payment.Amount = 0
	payment.Status = types.PaymentStatusFail
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	newPayment, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}

	return newPayment, nil
}


func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:			uuid.New().String(),
		AccountID: 	payment.AccountID,
		Name: 		name,
		Amount: 	payment.Amount,
		Category: 	payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	var targetFavorite *types.Favorite

	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			targetFavorite = favorite
			break
		}
	}

	if targetFavorite == nil {
		return nil, ErrFavoriteNotFound
	}

	payment, err := s.Pay(targetFavorite.AccountID, targetFavorite.Amount, targetFavorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}