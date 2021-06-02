# MedicalBlockchain

### Source Tree

MedicalBlockchain  
└ application-go  
  └─ Application.go        API code        (golang-gin)  
  └─ sdk  
    └── sdk.go               SDK code       (golang)  
└ chaincode-go  
  └─ assetTransfer.go  
  └─ chaincode  
    └── smartcontract.go    Chaindcode   (golang)  



## Quick Start

### Prerequisites

| Tech               | Version  |
| ------------------ | -------- |
| Ubuntu             | ^20.04.2 |
| Hyperledger Fabric | ^v2.3.1  |
| Golang             | ^v1.16.4 |
| Gin                | ^v1.7.2  |
| fabric-sdk-go      | ^v1.0.0  |

#### Hyperledger Fabric

Install the [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/en/latest/getting_started.html)

#### Golang

Install the [Go](https://golang.org/doc/install)

#### Gin

```shell
go get -u github.com/gin-gonic/gin
```

#### fabric-sdk-go

```shell
go get -u github.com/hyperledger/fabric-sdk-go
```

### API Start

```shell
cd application-go
go mod vendor
go run Application.go
```

## API

`localhost:8085`

### POST /CreateDID

| Key        | Description             | ValueType | Example            |
| ---------- | ----------------------- | --------- | ------------------ |
| DID        | DID                     | string    | did:Patient:000001 |
| AuthID     | Authority's DID         | string    | did:CA:000001      |
| Attribute  | DID owner's role        | string    | Patient            |
| Keytype    | Symmetric-key algorithm | string    | secp256k1          |
| Controller | DID controller          | string    | did:Patient:000001 |
| Key        | DID owner's public key  | string    | MHQCAQEEIPc1n...   |
| Sign       | Signature               | string    | MEQCIFzf1VNdgS...  |

### POST /ReadDID

| Key  | Description | ValueType | Example            |
| ---- | ----------- | --------- | ------------------ |
| DID  | DID         | string    | did:Patient:000001 |

### POST /CreateMedicalData

| Key         | Description             | ValueType | Example            |
| ----------- | ----------------------- | --------- | ------------------ |
| MedicalData | MedicalData File        | File      | example.txt        |
| AccessLevel | Access Level            | int       | 2                  |
| Metadata    | MedicalData's Metadata  | string    |                    |
| OwnerID     | MedicalData Owner's DID | string    | did:Patient:000001 |

### POST /ReadMedicalData

| Key  | Description              | ValueType | Example         |
| ---- | ------------------------ | --------- | --------------- |
| Hash | MedicalData's hash value | string    | 32d03797d413... |

### POST /ValidateMedicalData

| Key  | Description                             | ValueType | Example            |
| ---- | --------------------------------------- | --------- | ------------------ |
| Hash | MedicalData's hash value                | string    | 32d03797d413...    |
| DID  | DID to verify ownership of medical data | string    | did:Patient:000001 |

### POST /ShareMedicalData

| Key  | Description                                       | ValueType | Example            |
| ---- | ------------------------------------------------- | --------- | ------------------ |
| Hash | MedicalData's hash value                          | string    | 32d03797d413...    |
| DID  | DID to verify ownership of medical data for share | string    | did:Patient:000001 |

