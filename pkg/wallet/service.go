package wallet

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

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
	accounts			[]*types.Account
	payments			[]*types.Payment
	favorites		[]*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {

	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID: 			s.nextAccountID,
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
		Amount: 		amount,
		Category: 	category,
		Status: 		types.PaymentStatusInProgress,
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
		Amount: 		payment.Amount,
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

func (s *Service) ExportToFile(path string) error {
	content := make([]byte, 0)
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	} ()

	for _, account := range s.accounts {
		content = append(content, []byte(strconv.FormatInt(account.ID, 10))...)
		content = append(content, []byte(";")...)
		content = append(content, []byte(account.Phone)...)
		content = append(content, []byte(";")...)
		content = append(content, []byte(strconv.FormatInt(int64(account.Balance), 10))...)
		content = append(content, []byte("|")...)
	}

	_, err = file.Write(content)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil

}
func (s *Service) ImportFromFile(path string) error {
	content := make([]byte, 0)
	buf := make([]byte, 4)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	} ()
	
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}

		content = append(content, buf[:read]...)
	}

	log.Print(string(content))

	for _, rows := range strings.Split(string(content), "|") {
		columns := strings.Split(rows, ";")
		if len(columns) == 3 {
			s.RegisterAccount(types.Phone(columns[1]))
		}
	}
	for _, account := range s.accounts {
		log.Print(account)
	}

	return nil
}

func (s *Service) Export(dir string) error {
	accountFile := ""
	for _, account := range s.accounts {
		accounts := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" + strconv.FormatInt(int64(account.Balance), 10) + "\r\n"
		accountFile += accounts
	}
	if len(accountFile) > 0 {
		accountPath := dir + "/accounts.dump"
		accounts, err1 := os.Create(accountPath)
		if err1 != nil {
			log.Print(err1)
		return err1
		}
		_, accErr := accounts.Write([]byte(accountFile))
		if accErr != nil {
			log.Print(accErr)
		}
		defer func ()  {
			err := accounts.Close()
			if err != nil {
				log.Print(err)
				return
			}
		}()
	}	

	paymentFile := ""
	for _, payment := range s.payments {
		payments := string(payment.ID) + ";" + strconv.FormatInt(payment.AccountID, 10) + ";" + strconv.FormatInt(int64(payment.Amount),10) + ";" +string(payment.Category) + ";" +string(payment.Status) + "\r\n"
		paymentFile += payments
	}
	if len(paymentFile) > 0 {
		payPath := dir +"/payments.dump"
		paymentsFile, err2 := os.Create(payPath)
		if err2 != nil {
			log.Print(err2)
		return err2
		}
		_, payErr := paymentsFile.Write([]byte(paymentFile))
		if payErr != nil {
			log.Print(payErr)
		}
		defer func ()  {
			err := paymentsFile.Close()
			if err != nil {
				log.Print(err)
				return
			}
		}()
	}

	favoriteFile := ""
	for _, favorite := range s.favorites {
		favorite := string(favorite.ID) + ";" + strconv.FormatInt(favorite.AccountID,10) + ";" + string(favorite.Name) + ";" +strconv.FormatInt(int64(favorite.Amount),10) + ";" + string(favorite.Category) + "\r\n"
		favoriteFile += favorite
	}
	if len(favoriteFile) > 0 {
		favoritePath := dir + "/favorites.dump"
		favFile, err3 := os.Create(favoritePath)
		if err3 != nil {
			log.Print(err3)
		return err3
		}
		_, favorites_error := favFile.Write([]byte(favoriteFile))
		if favorites_error != nil {
			log.Print(favorites_error)
		}
		defer func ()  {
			err := favFile.Close()
			if err != nil {
				log.Print(err)
				return
			}
		}()
	}
	return nil	
}

func (s *Service) Import(dir string) error {
	// edit accounts file
	accountPath := dir + "/accounts.dump"
	accFile, err := os.Open(accountPath)
	if err != nil{
		log.Print(err)
		return err
	}

	defer func(){
		err := accFile.Close()
		if err != nil {
			log.Print(err)
			return
		}
	} ()
	accountContent := make([]byte,0)
	buf := make([]byte,4)

	for {
		read_accounts, err := accFile.Read(buf)
		if err == io.EOF {
			accountContent = append(accountContent, buf[:read_accounts]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		accountContent = append(accountContent, buf[:read_accounts]...)
	}

	data := strings.Split(string(accountContent), "\r\n")
	for _, accounts := range data {
		if len(accounts)>1{
		account := strings.Split(accounts, ";")
		id, err := strconv.ParseInt(account[0],10,64)
		if err != nil {
			log.Print(err)
		}
		balance,err := strconv.ParseInt(account[2],10,64)
		if err != nil {
			log.Print(err)
		}
		accountt := &types.Account{
			ID: id,
			Phone: types.Phone(account[1]),
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, accountt)
		}
	}

	//payments
	paymentPath := dir + "/payments.dump"
	payFile, err := os.Open(paymentPath)
	if err != nil{
		log.Print(err)
		return err
	}
	defer func(){
		err := payFile.Close()
		if err != nil {
			log.Print(err)
			return
		}
	} ()

	paymentContent := make([]byte,0)
	buff := make([]byte,4)

	for {
		readPayment, err := payFile.Read(buff)
		if err == io.EOF {
			paymentContent = append(paymentContent, buff[:readPayment]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		paymentContent = append(paymentContent, buff[:readPayment]...)
	}
	
	dataa := strings.Split(string(paymentContent), "\r\n")
	for _, payments := range dataa {
		if len(payments) > 1{
		payment := strings.Split(payments, ";")
		id_account, err := strconv.ParseInt(payment[1],10,64)
		if err != nil {
			log.Print(err)
		}
		amount,err := strconv.ParseInt(payment[2],10,64)
		if err != nil {
			log.Print(err)
		}
		paymentt := &types.Payment{
			ID: 			payment[0],
			AccountID: 	id_account,
			Amount: 		types.Money(amount),
			Category: 	types.PaymentCategory(payment[3]),
			Status: 		types.PaymentStatus(payment[4]),
		}
		s.payments = append(s.payments, paymentt)
		}
	}


	//favorites
	favPath := dir + "/favorites.dump"
	favoriteFile, err := os.Open(favPath)
	if err == nil{

	defer func(){
		err := favoriteFile.Close()
		if err != nil {
			log.Print(err)
			return
		}
	} ()
	favoriteContent := make([]byte,0)
	bufff := make([]byte,4)
	for {
		read_favorite, err := favoriteFile.Read(bufff)
		if err == io.EOF {
			favoriteContent = append(favoriteContent, bufff[:read_favorite]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		favoriteContent = append(favoriteContent, bufff[:read_favorite]...)
	}
	dataaa := strings.Split(string(favoriteContent), "\r\n")
	for _, favorites := range dataaa {
		if len(favorites)>1{
		favorite := strings.Split(favorites, ";")
		id_account, err := strconv.ParseInt(favorite[1],10,64)
		if err != nil {
			log.Print(err)
		}
		amount,err := strconv.ParseInt(favorite[3],10,64)
		if err != nil {
			log.Print(err)
		}
		favoritee := &types.Favorite{
			ID: favorite[0],
			AccountID: id_account,
			Name: favorite[2],
			Amount: types.Money(amount),
			Category: types.PaymentCategory(favorite[4]),
		}
		s.favorites = append(s.favorites,favoritee)
		}
	}
	}
	return nil
}
