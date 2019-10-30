package main

import (
	"encoding/hex"
	"fmt"
	"log"

	. "github.com/nknorg/nkn-sdk-go"
	"github.com/nknorg/nkn/crypto"
	"github.com/nknorg/nkn/vault"
	"github.com/nknorg/tuna"

	. "github.com/nknorg/tuna/exit"
)

func main() {
	config := Configuration{SubscriptionPrefix: tuna.DefaultSubscriptionPrefix}
	tuna.ReadJson("config.json", &config)

	Init()

	seed, _ := hex.DecodeString(config.Seed)
	privateKey := crypto.GetPrivateKeyFromSeed(seed)
	account, err := vault.NewAccountWithPrivatekey(privateKey)
	if err != nil {
		log.Panicln("Couldn't load account:", err)
	}

	wallet := NewWalletSDK(account)

	var services []Service
	tuna.ReadJson("services.json", &services)

	if config.Reverse {
		for serviceName := range config.Services {
			e := NewTunaExit(config, services, wallet)
			e.OnEntryConnected(func() {
				fmt.Printf("Service: %s, Address: %v:%v\n", serviceName, e.GetReverseIP(), e.GetReverseTCPPorts())
			})
			e.StartReverse(serviceName)
		}
	} else {
		NewTunaExit(config, services, wallet).Start()
	}

	select {}
}
