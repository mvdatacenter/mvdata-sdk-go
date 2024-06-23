module github.com/mvdatacenter/mvdata-sdk-go/feature/dynamodbstreams/attributevalue

go 1.20

require (
	github.com/mvdatacenter/mvdata-sdk-go v1.30.0
	github.com/mvdatacenter/mvdata-sdk-go/service/dynamodb v1.33.2
	github.com/mvdatacenter/mvdata-sdk-go/service/dynamodbstreams v1.21.1
)

require (
	github.com/aws/smithy-go v1.20.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/mvdatacenter/mvdata-sdk-go => ../../../

replace github.com/mvdatacenter/mvdata-sdk-go/internal/configsources => ../../../internal/configsources/

replace github.com/mvdatacenter/mvdata-sdk-go/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/mvdatacenter/mvdata-sdk-go/service/dynamodb => ../../../service/dynamodb/

replace github.com/mvdatacenter/mvdata-sdk-go/service/dynamodbstreams => ../../../service/dynamodbstreams/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/endpoint-discovery => ../../../service/internal/endpoint-discovery/
