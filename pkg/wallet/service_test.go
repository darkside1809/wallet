package wallet

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"sort"
	"github.com/darkside1809/wallet/pkg/types"
	"github.com/google/uuid"
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

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't regist account,  error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}

	account, _ := svc.RegisterAccount("+992000000001")

	acc, e := svc.FindAccountByID(account.ID)

	if e != nil {
		t.Error(e)
	}

	if !reflect.DeepEqual(account, acc) {
		t.Error("Accounts doesn't match")
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := &Service{}

	_, e := svc.FindAccountByID(123)

	if e != ErrAccountNotFound {
		t.Error(e)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong paymen returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must retunrn ErrPaymentNotFound, returned = %v", err)
		return
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
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status can't changed, error = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(), balance didn't cahnged, account = %v", savedAccount)
	}
}

func TestService_Reject_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.Reject(uuid.New().String())
	if err == nil {
		t.Error("Reject(): must be error, returned nil")
		return
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

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	_, err = s.FavoritePayment(payment.ID, "osh")
	if err != nil {
		t.Error(err)
	}
}

func TestService_FavoritePayment_fail(t *testing.T) {
	s := newTestService()

	_, err := s.FavoritePayment(uuid.New().String(), "osh")
	if err == nil {
		t.Error("FavoritePayment(): must return error, now nil")
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()

	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error("PayFromFavorite(): can't get payments")
		return
	}

	payment := payments[0]

	favorite, err := s.FavoritePayment(payment.ID, "osh")
	if err != nil {
		t.Error("PayFromFavorite(): can't add payment to favorite")
		return
	}

	_, err = s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Error("PayFromFavorite(): can't not pay from favorite")
		return
	}
}

func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTestService()

	_, err := s.PayFromFavorite(uuid.New().String())
	if err == nil {
		t.Error("PayFromFavorite(): must be error, now returned nil")
	}
}
func TestServiceExportToFile(t *testing.T) {
	s := newTestService()

	err := s.ExportToFile("export.txt")
	if err != nil {
		t.Error("ExportToFile(): cannot export data")
	}
}
func TestService_ImportFromFile_success(t *testing.T) {
	s := newTestService()

	err := s.ImportFromFile("export.txt")

	if err != nil {
		t.Errorf("method ImportFromFile return err, err: %v", err)
	}

}

func TestService_Import_success(t *testing.T) {
	s := newTestService()
	err := s.ImportFromFile("export.txt")

	if err != nil {
		t.Errorf("method Import returned not nil error, err => %v", err)
	}

}

func TestService_Export_success(t *testing.T) {
	s := newTestService()

	s.RegisterAccount("+992000000001")


	err := s.Export("data")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

	err = s.Import("data")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}
}

func TestService_ExportHistory_success(t *testing.T) {
	s := newTestService()

	account, err := s.RegisterAccount("+992000000001")

	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = s.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	_, err = s.Pay(account.ID, 1, "car")

	if err != nil {
		t.Errorf("method Pay returned not nil error, err => %v", err)
	}

	payments, err := s.ExportAccountHistory(account.ID)

	if err != nil {
		t.Errorf("method ExportAccountHistory returned not nil error, err => %v", err)
	}
	err = s.HistoryToFiles(payments, "data", 4)

	if err != nil {
		t.Errorf("method HistoryToFiles returned not nil error, err => %v", err)
	}

} 

func BenchmarkSumPayments(b *testing.B){
	svc := Service{}
	account, err := svc.RegisterAccount("+9921283793")

	if err != nil {
		b.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		b.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	_, err = svc.Pay(account.ID, 1, "Cafe")


	if err != nil {
		b.Errorf("method Pay returned not nil error, err => %v", err)
	}
	want := types.Money(1)

	got := svc.SumPayments(2)
	if want != got{
		b.Errorf("want: %v got: %v", want, got)
	}

} 

func BenchmarkFilterPayments(b *testing.B) {
	s := newTestService()
  
	account, err := s.RegisterAccount("+992000000000")
	if err != nil {
	  b.Error(err)
	}
	for i := 0; i < 103; i++ {
	  s.payments = append(s.payments, &types.Payment{AccountID: account.ID, Amount: 1})
	}
  
	result := 103
  
	for i := 0; i < b.N; i++ {
		payments, err := s.FilterPayments(account.ID, result)
	  	if err != nil {
			b.Error(err)
		}
  
	  	if result != len(payments) {
			b.Fatalf("wrong result, got %v, want %v", len(payments), result)
	 	}
	}
}

func filter(payment types.Payment) bool {
	return payment.Amount <= 540
}

func BenchmarkService_FilterPaymentsByFn(b *testing.B) {
	s := newTestService()

	account, err := s.RegisterAccount("+99212312133")
	if err != nil {
		b.Fatal(err)
	}

	err = s.Deposit(account.ID, 100_000)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 7; i++ {
		_, err = s.Pay(account.ID, types.Money(i * 10 + 500), "car")
		if err != nil {
			b.Fatal(err)
		}
	}

	want := make([]types.Payment, 0)

	for _, payment := range s.payments {
		if filter(*payment) {
			want = append(want, *payment)
		}
	}

	for i := 0; i < b.N; i++ {
		result, err := s.FilterPaymentsByFn(filter, i)
		if err != nil {
			b.Fatal(err)
		}

		sort.Slice(want, func (i, j int) bool {
			return want[i].ID < want[j].ID
		})

		sort.Slice(result, func (i, j int) bool {
			return result[i].ID < result[j].ID
		})

		if !reflect.DeepEqual(want, result) {
			b.Fatalf("invalid result, want %v, got %v", want, result)
		}
	}
}

func BenchmarkSumPaymentsWithProgress(b *testing.B) {
	var svc Service
	account, err := svc.RegisterAccount("992918532322")
	if err != nil {
		b.Errorf("Acc already registered %v", ErrPhoneRegistered)
	}
	err = svc.Deposit(account.ID, 100_000000000000)
	if err != nil {
		b.Errorf("Acc already registered %v", ErrPhoneRegistered)
	}
	for i := types.Money(1); i <= 100_000; i++ {
		_, err := svc.Pay(account.ID, i, "house")
		if  err != nil {
			b.Errorf("something went wrong %v", err)
		}
	}
	ch := svc.SumPaymentsWithProgress()
	s, got := <-ch
	if  !got {
		b.Errorf("got => %v", got)
	}
	log.Println(s)
}