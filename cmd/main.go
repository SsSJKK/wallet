package main

import (
//	"fmt"

	"github.com/SsSJKK/wallet/pkg/wallet"
)

//	"fmt"

func main() {
	svc := &wallet.Service{}
	svc.RegisterAccount("900000001")
	svc.Deposit(1, 100)
	svc.Export("./data")
	pays, _ := svc.ExportAccountHistory(1)
	svc.HistoryToFiles(pays,"./data",9)
	
}
