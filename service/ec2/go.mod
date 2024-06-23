module github.com/mvdatacenter/mvdata-sdk-go/service/ec2

go 1.20

require (
	github.com/mvdatacenter/mvdata-sdk-go v1.30.0
	github.com/mvdatacenter/mvdata-sdk-go/internal/configsources v1.3.12
	github.com/mvdatacenter/mvdata-sdk-go/internal/endpoints/v2 v2.6.12
	github.com/mvdatacenter/mvdata-sdk-go/service/internal/accept-encoding v1.11.2
	github.com/mvdatacenter/mvdata-sdk-go/service/internal/presigned-url v1.11.14
	github.com/aws/smithy-go v1.20.2
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/mvdatacenter/mvdata-sdk-go => ../../

replace github.com/mvdatacenter/mvdata-sdk-go/internal/configsources => ../../internal/configsources/

replace github.com/mvdatacenter/mvdata-sdk-go/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/presigned-url => ../../service/internal/presigned-url/
