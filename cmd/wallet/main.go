package main

import (
	"log"
	"sync"
	// "github.com/darkside1809/wallet/pkg/types"
	// "github.com/darkside1809/wallet/pkg/wallet"
)

func main() {
	//svc := &wallet.Service{}

	// err := svc.Import("data")
	// if err != nil {
	// 	log.Print(err)
	// }
	// err := svc.Export("data")
	// if err != nil {
	// 	log.Print(err)
	// }
	// err = svc.Import("data")
	// if err != nil {
	// 	log.Print(err)
	// }
	// err = svc.ExportToFile("data")
	// if err != nil {
	// 	log.Print(err)
	// }
	// err = svc.ImportFromFile("data")
	// if err != nil {
	// 	log.Print(err)
	// }

	// account, err := svc.ExportAccountHistory(125)
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Print(account)

	// payments := []types.Payment{}
	// err = svc.HistoryToFiles(payments, "data", 125)
	
	// err = svc.ImportFromFile("data")
	// if err != nil {
	// 	log.Print(err)
	// }

	log.Print("main started")
	wg := sync.WaitGroup{}
	wg.Add(2)
	sum := 0
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			sum++
		}
	}()	
	go func() {
		wg.Done()
		for i := 0; i < 1000; i++ {
			sum++
		}
	}()	
	
	wg.Wait()
	log.Print("main finished")
	log.Print(sum)
}