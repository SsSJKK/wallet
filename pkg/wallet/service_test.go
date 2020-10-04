package wallet

import (
	"reflect"
	"testing"

	"github.com/SsSJKK/wallet/pkg/types"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

var testSvc = newTestService()
var testAcc, _ = testSvc.RegisterAccount("992999999999")
var testPay = &types.Payment{
	ID:        "dsaafgasgsdfklgjsdf;l",
	AccountID: testAcc.ID,
	Amount:    types.Money(100),
	Category:  "TestPay",
	Status:    types.PaymentStatusInProgress,
}

func Test_RegisterAccount(t *testing.T) {
	svc := &Service{}
	_, err := svc.RegisterAccount("992000000001")

	if !reflect.DeepEqual(err, nil) {
		t.Errorf("invalid result, error: %v", err)
	}
}

func Test_RegisterAccount_2(t *testing.T) {
	svc := &Service{}
	_, err := svc.RegisterAccount("992000000001")
	_, err = svc.RegisterAccount("992000000001")
	if !reflect.DeepEqual(err, ErrPhoneRegistered) {
		t.Errorf("invalid result, error: %v", err)
	}
}

func Test_Deposit_1(t *testing.T) {
	svc := &Service{}
	account, _ := svc.RegisterAccount("992000000001")
	svc.Deposit(account.ID, 100)
	if !reflect.DeepEqual(account.Balance, types.Money(100)) {
		t.Errorf("invalid result, error: %v  need 100", account.Balance)
	}
}

func Test_Deposit_2(t *testing.T) {
	svc := &Service{}
	account, _ := svc.RegisterAccount("992000000001")
	svc.Deposit(account.ID, -100)
	if !reflect.DeepEqual(account.Balance, types.Money(0)) {
		t.Errorf("invalid result, error: %v  need 0", account.Balance)
	}
}

func Test_Deposit_3(t *testing.T) {
	svc := &Service{}
	err := svc.Deposit(-1, 100)
	if !reflect.DeepEqual(err, ErrAccountNotFound) {
		t.Errorf("invalid result, error ErrAccountFound")
	}
}

func Test_Pay_Amount(t *testing.T) {
	svc := &Service{}
	account, _ := svc.RegisterAccount("992000000001")
	_, err := svc.Pay(account.ID, -1, "A")
	if !reflect.DeepEqual(err, ErrAmountMustBePositive) {
		t.Errorf("ERROR %v", err)
	}

}

func Test_FindAccountByID_1(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("992000000001")
	account2, _ := svc.RegisterAccount("992000000002")
	svc.RegisterAccount("992000000003")

	accID, _ := svc.FindAccountByID(2)

	if !reflect.DeepEqual(accID.ID, account2.ID) {
		t.Errorf("ERROR %v, %v", accID, account2.ID)
	}
}
func Test_FindAccountByID_2(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("992000000001")
	svc.RegisterAccount("992000000002")
	svc.RegisterAccount("992000000003")

	_, err := svc.FindAccountByID(4)

	if !reflect.DeepEqual(err, ErrAccountNotFound) {
		t.Errorf("ERROR %v, %v", err, ErrAccountNotFound)
	}
}

func Test_FindPaymentById_OK(t *testing.T) {
	svc := &Service{}
	account, _ := svc.RegisterAccount("992000000001")
	svc.Deposit(1, 100)
	payment, _ := svc.Pay(1, 20, "A")
	svc.Reject(payment.ID)
	if !reflect.DeepEqual(account.Balance, types.Money(100)) {
		t.Errorf("ERROR %v", account.Balance)
	}
}

func Test_FindPaymentById_Fail(t *testing.T) {
	svc := &Service{}
	err := svc.Reject("payment.ID")
	if !reflect.DeepEqual(err, ErrPaymentNotFound) {
		t.Errorf("ERROR %v, %v", err, ErrPaymentNotFound)
	}
}

func Test_Repeat_OK(t *testing.T) {
	svc := &Service{}
	acc, _ := svc.RegisterAccount("992999999999")
	svc.Deposit(acc.ID, 100)
	pay, _ := svc.Pay(acc.ID, 20, "TestPay")
	repPay, _ := svc.Repeat(pay.ID)

	if pay.AccountID != repPay.AccountID {
		t.Errorf("ERROR: %v %v", pay.AccountID, repPay.AccountID)
	}
	if pay.Amount != repPay.Amount {
		t.Errorf("ERROR: %v %v", pay.Amount, repPay.Amount)
	}
	if pay.Category != repPay.Category {
		t.Errorf("ERROR: %v %v", pay.Category, repPay.Category)
	}
	if pay.Status != repPay.Status {
		t.Errorf("ERROR: %v %v", pay.Status, repPay.Status)
	}
}