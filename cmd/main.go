package main

import (
	"github.com/SsSJKK/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	svc.RegisterAccount("900000001")
	svc.RegisterAccount("900000002")
	svc.RegisterAccount("900000003")
	svc.RegisterAccount("900000004")
	svc.RegisterAccount("900000005")
	svc.RegisterAccount("900000006")
	svc.ExportToFile("./data/export.txt")
}
