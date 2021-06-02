package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type AuthInfo struct {
	Id         string `json:"authid"`
	Attribute  string `json:"attribute"`
	Keytype    string `json:"keytype"`
	Controller string `json:"controller"`
	Key        string `json:"pubkey"`
	Issue_date string `json:"time"`
}

// DID
type DID struct {
	DID      string     `json:"DID"`
	AuthInfo []AuthInfo `json:"authinfo"`
	Sign     string     `json:"sign:`
}

// MedicalData
type MedicalData struct {
	Hash        string `json:"hash"`
	AccessLevel int    `json:"accesslevel"`
	Metadata    string `json:"metadata"`
	Owner       string `json:"owner"`
}

// InitLedger adds a base set of CA to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	CA := []DID{
		{"did:CA:000001", []AuthInfo{{"did:CA:000001#keys-1", "CertificationAuthority", "Secp256k1", "did:CA:000001",
			"MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEAiCSg/jaymN9G087hAY4vH8hn0eq1tqYuVNjJ8F+INPU7ba58YPlvbJaNLAJn22gvQ0fMi1Bc0twwKUd5NGTBQ==", time.Now().Format(time.RFC3339)}},
			"MEQCIFzf1VNdgSXIGNKdOpP80URh3Uxe6Z0apINz78aSmU6XAiAPO6BdrT4vQkwekFI1kz1pofp+XKwd5TnOPAIKC2R3Nw=="},
		{"did:CA:000002", []AuthInfo{{"did:CA:000002#keys-1", "CertificationAuthority", "Secp256k1",
			"did:CA:000002", "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAE8K51BFNKq7vWJcSaQ5yLROodJHPQdAqsjK9TCW/yfEcwEWzDGfA7hIOH1R1a0gFrYcBXSfoHczGgxUxMpMVz4A==", time.Now().Format(time.RFC3339)}},
			"MEQCIHk02CaV761fWeD6V9E9caZvaeCU9rKkyl8EIuVxR/yoAiAZ3Mk/cmXGE2rXfMjEiCw5y2ZaaaTdV2wcsjCtKH1VXQ=="},
		{"did:CA:000003", []AuthInfo{{"did:CA:000003#keys-1", "CertificationAuthority", "Secp256k1",
			"did:CA:000003", "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEmR9mRDAkmZzElNf3OCWyFuZPV2vtILIlLjj4+GL0E39f1daSBvhOiOeJdqsnlmWjf57jX5fRRvJ7lrvm6ZZpKA==", time.Now().Format(time.RFC3339)}},
			"MEYCIQD4PCNQFQ1kWQMfSjdjWI3nTLWpRJjuBCV0m04P+Bw3mQIhAOi5NUvMUpz2sttKqbWFsGtkQY3ww1IGRlGovqla3le1"},
		{"did:CA:000004", []AuthInfo{{"did:CA:000004#keys-1", "CertificationAuthority", "Secp256k1",
			"did:CA:000004", "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAELzx4l3awsI0Du1zz74w/HJe5cVFE3yORVjaZ/Hd0UXt+26bzDm8vlorT8T806ByKXA70wVesZkzgdtk8X+Yq0Q==", time.Now().Format(time.RFC3339)}},
			"MEQCIEIoeHjkn68A6jI6Fr6USoJGYA3fZF6ONBHJBR2ebEdnAiBYhEnYCrreLzYk50V+jUllsqju5Mp3LKbDHavkugi+/Q=="},
	}

	for _, CA := range CA {
		CAJSON, err := json.Marshal(CA)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(CA.DID, CAJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface,
	id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface,
	id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// CreateDID issues a new DID to the world state with given details.
func (s *SmartContract) CreateDID(ctx contractapi.TransactionContextInterface,
	did string, authid, attr, keytype, controller, key, sign string) error {
	exists, err := s.AssetExists(ctx, did)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", did)
	}

	asset := DID{
		DID:      did,
		AuthInfo: []AuthInfo{{authid, attr, keytype, controller, key, time.Now().Format(time.RFC3339)}},
		Sign:     sign,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(did, assetJSON)
}

// CreateMedicalData issues a new MedicalData to the world state with given details.
func (s *SmartContract) CreateMedicalData(ctx contractapi.TransactionContextInterface,
	hash string, accesslevel int, metadata string, owner string) error {
	exists, err := s.AssetExists(ctx, hash)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", hash)
	}

	asset := MedicalData{
		Hash:        hash,
		AccessLevel: accesslevel,
		Metadata:    metadata,
		Owner:       owner,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hash, assetJSON)
}

// ReadDID returns the DID stored in the world state with given id.
func (s *SmartContract) ReadDID(ctx contractapi.TransactionContextInterface,
	did string) (*DID, error) {
	assetJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", did)
	}

	var Ret DID
	err = json.Unmarshal(assetJSON, &Ret)
	if err != nil {
		return nil, err
	}

	return &Ret, nil
}

// ReadMedicalData returns the MedicalData stored in the world state with given id.
func (s *SmartContract) ReadMedicalData(ctx contractapi.TransactionContextInterface,
	hash string) (*MedicalData, error) {
	assetJSON, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", hash)
	}

	var Ret MedicalData
	err = json.Unmarshal(assetJSON, &Ret)
	if err != nil {
		return nil, err
	}

	return &Ret, nil
}

func (s *SmartContract) ValidateMedicalData(ctx contractapi.TransactionContextInterface,
	hash string, did string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return false, fmt.Errorf("the asset %s does not exist", hash)
	}

	var Ret MedicalData
	err = json.Unmarshal(assetJSON, &Ret)
	if err != nil {
		return false, err
	}

	exists, err := s.AssetExists(ctx, did)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("the DID %s is not exist", did)
	}

	if Ret.Owner != did {
		return false, fmt.Errorf("medicaldata's owner %s is not equal to submitted did %s", Ret.Owner, did)
	}

	return true, nil
}

func (s *SmartContract) ShareMedicalData(ctx contractapi.TransactionContextInterface,
	hash string, did string) (bool, error) {
	//
	assetJSON, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return false, fmt.Errorf("the asset %s does not exist", hash)
	}

	DIDJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return false, fmt.Errorf("the asset %s does not exist", did)
	}

	var MD MedicalData
	err = json.Unmarshal(assetJSON, &MD)
	if err != nil {
		return false, err
	}

	var DD DID
	err = json.Unmarshal(DIDJSON, &DD)
	if err != nil {
		return false, err
	}

	if DD.AuthInfo[0].Attribute == "Patient" {
		return false, fmt.Errorf("AccessLevel X")
	}

	return true, nil
}

// // GetAllAssets returns all assets found in world state
// func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
// 	// range query with empty string for startKey and endKey does an
// 	// open-ended query of all assets in the chaincode namespace.
// 	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	var assets []*Asset
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		var asset Asset
// 		err = json.Unmarshal(queryResponse.Value, &asset)
// 		if err != nil {
// 			return nil, err
// 		}
// 		assets = append(assets, &asset)
// 	}

// 	return assets, nil
// }
