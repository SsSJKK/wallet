package wallet

import (
	"errors"

	"github.com/google/uuid"

	"github.com/SsSJKK/wallet/pkg/types"
)

//ErrPhoneRegistered = errors.New("phone already registred")
var ErrPhoneRegistered = errors.New("phone already registred")

//ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")

//ErrAccountNotFound = errors.New("account not found")
var ErrAccountNotFound = errors.New("account not found")

//ErrNotEnoughBalance = errors.New("You're a drifter")
var ErrNotEnoughBalance = errors.New("You're a drifter")

//ErrPaymentNotFound = errors.New("Payment not found")
var ErrPaymentNotFound = errors.New("Payment not found")

//Service s
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

//RegisterAccount R
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

//Deposit d
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

//Pay p
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		account = acc
		break
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

//FindAccountByID F
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

//FindPaymentByID f
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var Payment *types.Payment
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			Payment = pay
			break
		}
	}

	if Payment == nil {
		return nil, ErrPaymentNotFound
	}

	return Payment, nil
}

//Reject f
func (s *Service) Reject(paymentID string) error {
	pay, _ := s.FindPaymentByID(paymentID)
	if pay == nil {
		return ErrPaymentNotFound
	}
	acc, _ := s.FindAccountByID(pay.AccountID)
	pay.Status = types.PaymentStatusFail
	acc.Balance += pay.Amount
	return nil
}

//Repeat F
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(pay.AccountID, pay.Amount, pay.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
