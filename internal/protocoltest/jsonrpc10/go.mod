module github.com/mvdatacenter/mvdata-sdk-go/internal/protocoltest/jsonrpc10

go 1.20

require (
	github.com/mvdatacenter/mvdata-sdk-go v1.30.0
	github.com/mvdatacenter/mvdata-sdk-go/internal/configsources v1.3.12
	github.com/mvdatacenter/mvdata-sdk-go/internal/endpoints/v2 v2.6.12
	github.com/aws/smithy-go v1.20.2
)

replace github.com/mvdatacenter/mvdata-sdk-go => ../../../

replace github.com/mvdatacenter/mvdata-sdk-go/internal/configsources => ../../../internal/configsources/

replace github.com/mvdatacenter/mvdata-sdk-go/internal/endpoints/v2 => ../../../internal/endpoints/v2/
