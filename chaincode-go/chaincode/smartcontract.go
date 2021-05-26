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
}

// MedicalData
type MedicalData struct {
	Hash        string `json:"hash"`
	AccessLevel int    `json:"accesslevel"`
	Metadata    string `json:"metadata"`
	Owner       string `json:"owner"`
}

// // Validate Time Fucntion
// func TimeValidate(did1, did2 DID, comp string) bool {
// 	d1, _ := time.Parse(time.RFC3339, did1.Issue_date)
// 	d2, _ := time.Parse(time.RFC3339, did2.Issue_date)
// 	if comp == "equal" {
// 		return d1.Equal(d2)
// 	} else if comp == "before" {
// 		return d1.Before(d2)
// 	} else if comp == "after" {
// 		return d1.After(d2)
// 	} else {
// 		fmt.Println("third argument is expected {equal, before, after}")
// 		return false
// 	}

// }

// InitLedger adds a base set of CA to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	CA := []DID{
		{"did:CA:000001", []AuthInfo{{"did:CA:000001#keys-1", "CertificationAuthority", "Secp256k1", "did:CA:000001",
			"MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEqHA9txbYf7uFrY7toGYN/e1gkdXRlSzJdU2XBELvigoapbLIY/rObffay1/vNXSSokyuW6TdgJIg0j0nJj27qw==", time.Now().Format(time.RFC3339)}}},
		{"did:CA:000002", []AuthInfo{{"did:CA:000002#keys-1", "CertificationAuthority", "Secp256k1",
			"did:CA:000002", "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEp92Jko3ZJhd1gCy0S9gTbZI7KfJpWZu7FXC5zkyVQyzuvRJOWhA/xRuOoMHYvcMKFQRpiMEje+SkREdYQuf0Cw==", time.Now().Format(time.RFC3339)}}},
		{"did:CA:000003", []AuthInfo{{"did:CA:000003#keys-1", "CertificationAuthority", "Secp256k1",
			"did:CA:000003", "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEGIMRLmPuO8G+DAngk1nlnYaRA8bogaRqAW4Fi9MO1YddLgMiG09sJN463QPhC3p4ytW51FgjHK4Lp+phXbUyJA==", time.Now().Format(time.RFC3339)}}},
		{"did:CA:000004", []AuthInfo{{"did:CA:000004#keys-1", "CertificationAuthority", "Secp256k1",
			"did:CA:000004", "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEh4XkJVR5w1W5yXzM0M+jErjQFXKumu5DB7A1IvbPmHmTHEtKCn/uitWGfPhdk7oCsXo59chYG71uVibVwOpJSg==", time.Now().Format(time.RFC3339)}}},
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

// AuthInfo
// type CertificationAuthority struct {
// 	Id         string `json:"authid"`
// 	Attribute  string `json:"attribute"`
// 	Keytype    string `json:"keytype"`
// 	Controller string `json:"controller"`
// 	Key        string `json:"pubkey"`
// }

// CreateCA issues a new CA to the world state with given details.
// func (s *SmartContract) CreateCA(ctx contractapi.TransactionContextInterface,
// 	id, attr, keytype, controller, key string) error {
// 	exists, err := s.AssetExists(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if exists {
// 		return fmt.Errorf("the asset %s already exists", id)
// 	}

// 	asset := CertificationAuthority{
// 		id,
// 		attr,
// 		keytype,
// 		controller,
// 		key,
// 	}

// 	assetJSON, err := json.Marshal(asset)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(id, assetJSON)
// }

// CreateDID issues a new DID to the world state with given details.
func (s *SmartContract) CreateDID(ctx contractapi.TransactionContextInterface,
	did string, authid, attr, keytype, controller, key string) error {
	exists, err := s.AssetExists(ctx, did)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", did)
	}

	// aid, err := s.ReadDID(ctx, authid)

	// if err != nil {
	// 	return err
	// }

	asset := DID{
		DID:      did,
		AuthInfo: []AuthInfo{{authid, attr, keytype, controller, key, time.Now().Format(time.RFC3339)}},
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

	if int(MD.AccessLevel) < 2 {
		return false, fmt.Errorf("AccessLevel X")
	}

	// To Do: Request Share MD

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
