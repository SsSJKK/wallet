package main

import (
	"fmt"

	"github.com/SsSJKK/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	/*svc.RegisterAccount("900000001")
	svc.RegisterAccount("900000002")
	svc.RegisterAccount("900000003")
	svc.RegisterAccount("900000004")
	svc.RegisterAccount("900000005")
	svc.RegisterAccount("900000006")
	svc.Deposit(1, 100)
	pay1, _ := svc.Pay(1, 50, "test")
	pay2, _ := svc.Pay(1, 10, "test2")
	svc.FavoritePayment(pay1.ID, "fav1")
	svc.FavoritePayment(pay2.ID, "fav2")
	svc.Export("./data")
	svc.ImportFromFile("./data/export.txt")
	svc.RegisterAccount("0")
	fmt.Println(svc.FindAccountByID(1))
	svc.Import("./data")
	fmt.Println(svc.FindAccountByID(1))
	fmt.Println(svc.FindAccountByID(2))
	fmt.Println(svc.FindAccountByID(3))
	fmt.Println(svc.FindAccountByID(4))
	fmt.Println(svc.FindAccountByID(5))
	fmt.Println(svc.FindAccountByID(6))*/

	svc.Import("./data")
	fmt.Println(svc.FindAccountByID(1))
	fmt.Println(svc.FindAccountByID(2))
	fmt.Println(svc.FindAccountByID(3))
	fmt.Println(svc.FindAccountByID(4))
	fmt.Println(svc.FindAccountByID(5))
	fmt.Println(svc.FindAccountByID(6))
	fmt.Println(svc.FindPaymentByID("c2925ea8-038c-4af0-8736-3c6f5fc087d6"))
	fmt.Println(svc.FindFavoriteByID("4efb2b96-2766-41cc-a2f9-821e51e37819"))

}
