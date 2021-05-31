module asset-transfer-basic

go 1.14

require (
	main.com/sdk v0.0.0
	github.com/gin-gonic/gin v1.7.2
	github.com/hyperledger/fabric-sdk-go v1.0.0
	golang.org/x/tools v0.1.2 // indirect
)

replace (
	main.com/sdk v0.0.0 => ./sdk
)