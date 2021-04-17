package main

import (
	"github.com/darkside1809/wallet/pkg/wallet"
	"log"
)


func main() {
	svc := &wallet.Service{}

	// err := svc.ExportToFile("accounts.txt")
	// if err != nil {
	// 	log.Print(err)
	// }

	err := svc.ImportFromFile("accounts.txt")
	if err != nil {
		log.Print(err)
	}
}