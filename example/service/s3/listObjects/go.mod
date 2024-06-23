module github.com/mvdatacenter/mvdata-sdk-go/example/service/s3/listObjects

go 1.20

require (
	github.com/mvdatacenter/mvdata-sdk-go/config v1.27.21
	github.com/mvdatacenter/mvdata-sdk-go/service/s3 v1.56.1
)

require (
	github.com/mvdatacenter/mvdata-sdk-go v1.30.0 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/mvdata/protocol/eventstream v1.6.2 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/credentials v1.17.21 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/feature/ec2/imds v1.16.8 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/internal/configsources v1.3.12 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/internal/endpoints/v2 v2.6.12 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/internal/ini v1.8.0 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/internal/v4a v1.3.12 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/service/internal/accept-encoding v1.11.2 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/service/internal/checksum v1.3.14 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/service/internal/presigned-url v1.11.14 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/service/internal/s3shared v1.17.12 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/service/sso v1.21.1 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/service/ssooidc v1.25.1 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/service/sts v1.29.1 // indirect
	github.com/aws/smithy-go v1.20.2 // indirect
)

replace github.com/mvdatacenter/mvdata-sdk-go => ../../../../

replace github.com/mvdatacenter/mvdata-sdk-go/mvdata/protocol/eventstream => ../../../../aws/protocol/eventstream/

replace github.com/mvdatacenter/mvdata-sdk-go/config => ../../../../config/

replace github.com/mvdatacenter/mvdata-sdk-go/credentials => ../../../../credentials/

replace github.com/mvdatacenter/mvdata-sdk-go/feature/ec2/imds => ../../../../feature/ec2/imds/

replace github.com/mvdatacenter/mvdata-sdk-go/internal/configsources => ../../../../internal/configsources/

replace github.com/mvdatacenter/mvdata-sdk-go/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/mvdatacenter/mvdata-sdk-go/internal/ini => ../../../../internal/ini/

replace github.com/mvdatacenter/mvdata-sdk-go/internal/v4a => ../../../../internal/v4a/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/accept-encoding => ../../../../service/internal/accept-encoding/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/checksum => ../../../../service/internal/checksum/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/presigned-url => ../../../../service/internal/presigned-url/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/s3shared => ../../../../service/internal/s3shared/

replace github.com/mvdatacenter/mvdata-sdk-go/service/s3 => ../../../../service/s3/

replace github.com/mvdatacenter/mvdata-sdk-go/service/sso => ../../../../service/sso/

replace github.com/mvdatacenter/mvdata-sdk-go/service/ssooidc => ../../../../service/ssooidc/

replace github.com/mvdatacenter/mvdata-sdk-go/service/sts => ../../../../service/sts/
