/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("basic")

	if len(os.Args) == 1 {
		log.Fatalf("Not Expected Argument")
	}

	funcname := os.Args[1]

	if funcname == "ReadCA" {
		if len(os.Args) != 3 {
			log.Fatalf("Not Expected Argument")
		}
		ReadCA(contract, os.Args[2])
	} else if funcname == "ReadDID" {
		if len(os.Args) != 3 {
			log.Fatalf("Not Expected Argument")
		}
		ReadDID(contract, os.Args[2])
	} else if funcname == "ReadMedicalData" {
		if len(os.Args) != 3 {
			log.Fatalf("Not Expected Argument")
		}
		ReadMedicalData(contract, os.Args[2])
	} else if funcname == "CreateCA" {
		if len(os.Args) != 7 {
			log.Fatalf("Not Expected Argument")
		}
		CreateCA(contract, os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6])
	} else if funcname == "CreateDID" {
		if len(os.Args) != 4 {
			log.Fatalf("Not Expected Argument")
		}
		CreateDID(contract, os.Args[2], os.Args[3])
	} else if funcname == "CreateMedicalData" {
		if len(os.Args) != 6 {
			log.Fatalf("Not Expected Argument")
		}
		CreateMedicalData(contract, os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	} else if funcname == "ValidateMedicalData" {
		if len(os.Args) != 4 {
			log.Fatalf("Not Expected Argument")
		}
		ValidateMedicalData(contract, os.Args[2], os.Args[3])
	} else if funcname == "ShareMedicalData" {
		if len(os.Args) != 4 {
			log.Fatalf("Not Expected Argument")
		}
		ShareMedicalData(contract, os.Args[2], os.Args[3])
	} else if funcname == "InitLedger" {
		if len(os.Args) != 2 {
			log.Fatalf("Not Expected Argument")
		}
		InitLedger(contract)
	} else {
		log.Fatalf("Not Expected Argument")
	}

	log.Println("============ application-golang ends ============")
}

func InitLedger(contract *gateway.Contract) {
	log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))
}

func ReadCA(contract *gateway.Contract, id string) {
	log.Println("--> Evaluate Transaction: ReadCA, function returns CA on the ledger")
	result, err := contract.EvaluateTransaction("ReadCA", id)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
}

func ReadDID(contract *gateway.Contract, id string) {
	log.Println("--> Evaluate Transaction: ReadDID, function returns DID on the ledger")
	result, err := contract.EvaluateTransaction("ReadDID", id)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
}

func ReadMedicalData(contract *gateway.Contract, hash string) {
	log.Println("--> Evaluate Transaction: ReadMedicalData, function returns MedicalData on the ledger")
	result, err := contract.EvaluateTransaction("ReadMedicalData", hash)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
}

func CreateCA(contract *gateway.Contract, id, attr, keytype, controller, key string) {
	log.Println("--> Submit Transaction: CreateCA, creates new CA with ID, Attrubute, Keytype, Controller, and Key arguments")
	result, err := contract.SubmitTransaction("CreateCA", id, attr, keytype, controller, key)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))
}

func CreateDID(contract *gateway.Contract, did, authid string) {
	log.Println("--> Submit Transaction: CreateDID, creates new DID with DID and AuthID arguments")
	result, err := contract.SubmitTransaction("CreateDID", did, authid)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))
}

func CreateMedicalData(contract *gateway.Contract, fileloc string, accesslevel string, metadata string, owner string) {
	log.Println("--> Submit Transaction: CreateMedicalData, creates new MedicalData with Hash, AccessLevel, Metadata and OwnerID arguments")

	dat, err := ioutil.ReadFile(fileloc)

	hash := sha256.Sum256(dat)

	hashstr := hex.EncodeToString(hash[:])

	result, err := contract.SubmitTransaction("CreateMedicalData", hashstr, accesslevel, metadata, owner)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result), `Hash: `, hashstr)
}

func ValidateMedicalData(contract *gateway.Contract, hash, did string) {
	log.Println("--> Submit Transaction: ValidateMedicalData, validates MedicalData's Owner with Hash and OwnerID arguments")
	result, err := contract.SubmitTransaction("ValidateMedicalData", hash, did)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))
}

func ShareMedicalData(contract *gateway.Contract, hash, did string) {
	log.Println("--> Submit Transaction: ValidateMedicalData, validates MedicalData's Owner with Hash and OwnerID arguments")
	result, err := contract.SubmitTransaction("ShareMedicalData", hash, did)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))
}

// log.Println("--> Evaluate Transaction: ReadAsset, function returns 'asset1' attributes")
// result, err = contract.EvaluateTransaction("ReadAsset", "asset1")
// if err != nil {
// 	log.Fatalf("Failed to evaluate transaction: %v", err)
// }
// log.Println(string(result))

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
