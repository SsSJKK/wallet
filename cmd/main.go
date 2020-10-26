package main

import (
	"log"

	"github.com/SsSJKK/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	svc.Import("./data")
	ch := svc.SumPaymentsWithProgress()
	for i := range ch {
		log.Print(i)
	}
}
