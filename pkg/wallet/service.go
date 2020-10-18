package wallet

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/SsSJKK/wallet/pkg/types"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, ErrAccountNotFound
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return ErrAccountNotFound
	}

	// зачисление средств пока не рассматриваем как платёж
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

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	return s.Pay(payment.AccountID, payment.Amount, payment.Category)
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Name:      name,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}
func (s *Service) ExportToFile(path string) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	text := ""
	for _, acc := range s.accounts {
		strID := strconv.FormatInt(int64(acc.ID), 10) + ";"
		strPhone := acc.Phone + ";"
		strBalance := strconv.FormatInt(int64(acc.Balance), 10)

		text += strID + string(strPhone) + strBalance + "|"
	}

	_, err = file.Write([]byte(text))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ImportFromFile(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err

	}
	defer file.Close()
	data := make([]byte, 64)
	read := ""
	for {
		n, err := file.Read(data)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
		read += string(data[:n])

	}
	//fmt.Println(a)
	importAcc := strings.Split(read, "|")
	importAcc = importAcc[:len(importAcc)-1]
	for _, acc := range importAcc {
		account := strings.Split(acc, ";")
		ID, _ := strconv.ParseInt(account[0], 10, 64)
		phone := types.Phone(account[1])
		balance, _ := strconv.ParseInt(account[2], 10, 64)

		addAcc := &types.Account{
			ID:      ID,
			Phone:   phone,
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, addAcc)
	}
	return nil
}

func (s *Service) Export(dir string) error {

	fileAccounts, err := os.Create(dir + "/accounts.dump")
	if err != nil {
		return err
	}
	defer fileAccounts.Close()
	text := ""
	for _, acc := range s.accounts {
		strID := strconv.FormatInt(int64(acc.ID), 10) + ";"
		strPhone := acc.Phone + ";"
		strBalance := strconv.FormatInt(int64(acc.Balance), 10)

		text += strID + string(strPhone) + strBalance + "\n"
	}

	_, err = fileAccounts.Write([]byte(text))
	if err != nil {
		return err
	}

	filePay, err := os.Create(dir + "/payments.dump")
	if err != nil {
		return err
	}
	defer filePay.Close()
	text = ""
	for _, pay := range s.payments {
		strID := pay.ID + ";"
		strAccountID := strconv.FormatInt(int64(pay.AccountID), 10) + ";"
		strAmount := strconv.FormatInt(int64(pay.Amount), 10) + ";"
		strCategory := string(pay.Category) + ";"
		strStatus := string(pay.Status) + ";"

		text += strID + strAccountID + strAmount + strCategory + strStatus + "\n"
	}

	_, err = filePay.Write([]byte(text))
	if err != nil {
		return err
	}

	fileFav, err := os.Create(dir + "/favorites.dump")
	if err != nil {
		return err
	}
	defer fileFav.Close()
	text = ""
	for _, fav := range s.favorites {
		strID := fav.ID + ";"
		strAccountID := strconv.FormatInt(int64(fav.AccountID), 10) + ";"
		strAmount := strconv.FormatInt(int64(fav.Amount), 10) + ";"
		strName := fav.Name + ";"
		strCategory := string(fav.Category) + ";"

		text += strID + strAccountID + strAmount + strName + strCategory + "\n"
	}

	_, err = fileFav.Write([]byte(text))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Import(dir string) error {

	fileAcc, err := os.Open(dir + "/accounts.dump")
	if err != nil {
		log.Fatal(err)
	}
	defer fileAcc.Close()
	if err == nil {
		scanner := bufio.NewScanner(fileAcc)
		for scanner.Scan() {
			line := strings.Split(scanner.Text(), ";")
			ID, _ := strconv.ParseInt(line[0], 10, 64)
			phone := types.Phone(line[1])
			balance, _ := strconv.ParseInt(line[2], 10, 64)
			acc, errAcc := s.FindAccountByID(ID)
			if errAcc == nil {
				acc.ID = ID
				acc.Phone = phone
				acc.Balance = types.Money(balance)
			}

			if errAcc != nil {
				addAcc := &types.Account{
					ID:      ID,
					Phone:   phone,
					Balance: types.Money(balance),
				}
				s.accounts = append(s.accounts, addAcc)
			}
			if s.nextAccountID < ID {
				s.nextAccountID = ID
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	filePay, err := os.Open(dir + "/payments.dump")
	if err != nil {
		log.Print(err)
	}
	defer filePay.Close()
	if err == nil {
		scanner := bufio.NewScanner(filePay)
		for scanner.Scan() {
			line := strings.Split(scanner.Text(), ";")
			ID := line[0]
			accountID, _ := strconv.ParseInt(line[1], 10, 64)
			amount, _ := strconv.ParseInt(line[2], 10, 64)
			category := types.PaymentCategory(line[3])
			status := types.PaymentStatus(line[4])
			pay, err := s.FindPaymentByID(ID)
			if err == nil {
				pay.ID = ID
				pay.AccountID = accountID
				pay.Amount = types.Money(amount)
				pay.Category = category
				pay.Status = status
			}
			if err != nil {
				addPay := &types.Payment{
					ID:        ID,
					AccountID: accountID,
					Amount:    types.Money(amount),
					Category:  category,
					Status:    status,
				}
				s.payments = append(s.payments, addPay)
			}
		}
	}

	fileFav, err := os.Open(dir + "/favorites.dump")
	if err != nil {
		log.Print(err)
	}
	defer fileFav.Close()
	if err == nil {
		scanner := bufio.NewScanner(fileFav)
		for scanner.Scan() {
			line := strings.Split(scanner.Text(), ";")
			ID := line[0]
			accountID, _ := strconv.ParseInt(line[1], 10, 64)
			amount, _ := strconv.ParseInt(line[2], 10, 64)
			name := line[3]
			category := types.PaymentCategory(line[4])

			fav, err := s.FindFavoriteByID(ID)
			if err == nil {
				fav.ID = ID
				fav.AccountID = accountID
				fav.Amount = types.Money(amount)
				fav.Name = name
				fav.Category = category
			}
			if err != nil {
				addFav := &types.Favorite{
					ID:        ID,
					AccountID: accountID,
					Amount:    types.Money(amount),
					Name:      name,
					Category:  category,
				}
				s.favorites = append(s.favorites, addFav)
			}
		}
	}

	return nil
}
