/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
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

	r := gin.Default()

	r.GET("/ping", ping)
	r.GET("/InitLedger", InitLedger(contract))
	r.POST("/ReadDID", ReadDID(contract))
	r.POST("/ReadMedicalData", ReadMedicalData((contract)))
	r.POST("/CreateDID", CreateDID(contract))
	r.POST("/CreateMedicalData", CreateMedicalData(contract))
	r.POST("/ValidateMedicalData", ValidateMedicalData(contract))
	r.POST("/ShareMedicalData", ShareMedicalData(contract))

	r.Run(":8085") // listen and serve on 0.0.0.0:8080
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// 	} else if funcname == "CreateCA" {
// 		if len(os.Args) != 7 {
// 			log.Fatalf("Not Expected Argument")
// 		}
// 		CreateCA(contract, os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6])
// 	} else if funcname == "CreateDID" {
// 		if len(os.Args) != 4 {
// 			log.Fatalf("Not Expected Argument")
// 		}
// 		CreateDID(contract, os.Args[2], os.Args[3])
// 	} else if funcname == "CreateMedicalData" {
// 		if len(os.Args) != 6 {
// 			log.Fatalf("Not Expected Argument")
// 		}
// 		CreateMedicalData(contract, os.Args[2], os.Args[3], os.Args[4], os.Args[5])
// 	} else if funcname == "ValidateMedicalData" {
// 		if len(os.Args) != 4 {
// 			log.Fatalf("Not Expected Argument")
// 		}
// 		ValidateMedicalData(contract, os.Args[2], os.Args[3])
// 	} else if funcname == "ShareMedicalData" {
// 		if len(os.Args) != 4 {
// 			log.Fatalf("Not Expected Argument")
// 		}
// 		ShareMedicalData(contract, os.Args[2], os.Args[3])
// 	} else if funcname == "InitLedger" {
// 		if len(os.Args) != 2 {
// 			log.Fatalf("Not Expected Argument")
// 		}
// 		InitLedger(contract)
// 	} else {
// 		log.Fatalf("Not Expected Argument")
// 	}

// 	log.Println("============ application-golang ends ============")
// }

func InitLedger(contract *gateway.Contract) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
		result, err := contract.SubmitTransaction("InitLedger")
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
		}
		c.JSON(200, gin.H{
			"message": string(result),
		})
	}
	return gin.HandlerFunc(fn)
}

func ReadDID(contract *gateway.Contract) gin.HandlerFunc {
	// params: JSON {DID: DID}
	fn := func(c *gin.Context) {
		log.Println("--> Evaluate Transaction: ReadDID, function returns DID on the ledger")

		type Data struct {
			DID string `form:"DID" json:"DID" binding:"required"`
		}

		var data Data
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   err.Error(),
			})
			return
		}

		log.Println(data.DID)
		result, err := contract.EvaluateTransaction("ReadDID", data.DID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Faild to evaluate transaction: " + err.Error(),
			})
			return
		}
		log.Println(string(result))

		c.JSON(200, gin.H{
			"message": string(result),
		})

	}
	return gin.HandlerFunc(fn)
}

func ReadMedicalData(contract *gateway.Contract) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Println("--> Evaluate Transaction: ReadMedicalData, function returns MedicalData on the ledger")

		type Data struct {
			Hash string `form:"Hash" json:"Hash" binding:"required"`
		}

		var data Data
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   err.Error(),
			})
			return
		}

		result, err := contract.EvaluateTransaction("ReadMedicalData", data.Hash)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   "Faild to evaluate transaction: " + err.Error(),
			})
			return
		}

		log.Println(string(result))
		c.JSON(200, gin.H{
			"message": string(result),
		})
	}
	return gin.HandlerFunc(fn)
}

// func CreateCA(contract *gateway.Contract) gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		log.Println("--> Submit Transaction: CreateCA, creates new CA with ID, Attrubute, Keytype, Controller, and Key arguments")

// 		type Data struct {
// 			DID        string `form:"DID" json:"DID" binding:"required"`
// 			Attribute  string `form:"Attribute" json:"Attribute" binding:"required"`
// 			Keytype    string `form:"Keytype" json:"Keytype" binding:"required"`
// 			Controller string `form:"Controller" json:"Controller" binding:"required"`
// 			Key        string `form:"Key" json:"Key" binding:"required"`
// 		}

// 		var data Data
// 		if err := c.ShouldBind(&data); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		result, err := contract.SubmitTransaction("CreateCA", data.DID, data.Attribute, data.Keytype, data.Controller, data.Key)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"error": "Faild to evaluate transaction: " + err.Error(),
// 			})
// 			return
// 		}

// 		log.Println(string(result))
// 		c.JSON(200, gin.H{
// 			"message": string(result),
// 		})
// 	}
// 	return gin.HandlerFunc(fn)
// }

func CreateDID(contract *gateway.Contract) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Println("--> Submit Transaction: CreateDID, creates new DID with DID and AuthID arguments")

		type Data struct {
			DID        string `form:"DID" json:"DID" binding:"required"`
			AuthID     string `form:"AuthID" json:"AuthID" binding:"required"`
			Attribute  string `form:"Attribute" json:"Attribute" binding:"required"`
			Keytype    string `form:"Keytype" json:"Keytype" binding:"required"`
			Controller string `form:"Controller" json:"Controller" binding:"required"`
			Key        string `form:"Key" json:"Key" binding:"required"`
		}

		var data Data
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   err.Error(),
			})
			return
		}

		result, err := contract.SubmitTransaction("CreateDID", data.DID, data.AuthID, data.Attribute, data.Keytype, data.Controller, data.Key)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   "Faild to evaluate transaction: " + err.Error(),
			})
			return
		}

		log.Println(string(result))
		c.JSON(200, gin.H{
			"message": data.DID + " is created.",
		})
	}
	return gin.HandlerFunc(fn)
}

func CreateMedicalData(contract *gateway.Contract) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Println("--> Submit Transaction: CreateMedicalData, creates new MedicalData with Hash, AccessLevel, Metadata and OwnerID arguments")

		type Data struct {
			// MedicalData []byte `form:"MedicalData" json:"MedicalData" binding:"required"`
			AccessLevel string `form:"AccessLevel" json:"AccessLevel" binding:"required"`
			Metadata    string `form:"Metadata" json:"Metadata" binding:"required"`
			OwnerID     string `form:"OwnerID" json:"OwnerID" binding:"required"`
		}

		file, _, err := c.Request.FormFile("MedicalData")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return
		}

		var data Data
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   err.Error(),
			})
			return
		}

		hash := sha256.Sum256(buf.Bytes())

		hashstr := hex.EncodeToString(hash[:])

		log.Println(hashstr)

		result, err := contract.SubmitTransaction("CreateMedicalData", hashstr, data.AccessLevel, data.Metadata, data.OwnerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   "Faild to evaluate transaction: " + err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": hashstr + " is created.",
		})

	}
	return gin.HandlerFunc(fn)
}

func ValidateMedicalData(contract *gateway.Contract) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Println("--> Submit Transaction: ValidateMedicalData, validates MedicalData's Owner with Hash and OwnerID arguments")

		type Data struct {
			// MedicalData []byte `form:"MedicalData" json:"MedicalData" binding:"required"`
			Hash string `form:"Hash" json:"Hash" binding:"required"`
			DID  string `form:"DID" json:"DID" binding:"required"`
		}

		var data Data
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   err.Error(),
			})
			return
		}

		result, err := contract.SubmitTransaction("ValidateMedicalData", data.Hash, data.DID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   "Faild to evaluate transaction: " + err.Error(),
			})
			return
		}

		if string(result) == "true" {
			c.JSON(200, gin.H{
				"message": "DID " + data.DID + " is validated",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "DID " + data.DID + " is not validated",
			})
		}
		log.Println(string(result))
	}
	return gin.HandlerFunc(fn)
}

func ShareMedicalData(contract *gateway.Contract) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Println("--> Submit Transaction: ValidateMedicalData, validates MedicalData's Owner with Hash and OwnerID arguments")

		type Data struct {
			// MedicalData []byte `form:"MedicalData" json:"MedicalData" binding:"required"`
			Hash string `form:"Hash" json:"Hash" binding:"required"`
			DID  string `form:"DID" json:"DID" binding:"required"`
		}

		var data Data
		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   err.Error(),
			})
			return
		}

		result, err := contract.SubmitTransaction("ShareMedicalData", data.Hash, data.DID)
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error",
				"error":   "Faild to evaluate transaction: " + err.Error(),
			})
			return
		}

		log.Println(string(result))
		c.JSON(200, gin.H{
			"message": string(result),
		})
		log.Println(string(result))
	}
	return gin.HandlerFunc(fn)
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
