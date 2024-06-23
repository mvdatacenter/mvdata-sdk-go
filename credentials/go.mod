module github.com/mvdatacenter/mvdata-sdk-go/credentials

go 1.20

require (
	github.com/mvdatacenter/mvdata-sdk-go v1.30.0
	github.com/mvdatacenter/mvdata-sdk-go/feature/ec2/imds v1.16.8
	github.com/mvdatacenter/mvdata-sdk-go/service/sso v1.21.1
	github.com/mvdatacenter/mvdata-sdk-go/service/ssooidc v1.25.1
	github.com/mvdatacenter/mvdata-sdk-go/service/sts v1.29.1
	github.com/aws/smithy-go v1.20.2
)

require (
	github.com/mvdatacenter/mvdata-sdk-go/internal/configsources v1.3.12 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/internal/endpoints/v2 v2.6.12 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/service/internal/accept-encoding v1.11.2 // indirect
	github.com/mvdatacenter/mvdata-sdk-go/service/internal/presigned-url v1.11.14 // indirect
)

replace github.com/mvdatacenter/mvdata-sdk-go => ../

replace github.com/mvdatacenter/mvdata-sdk-go/feature/ec2/imds => ../feature/ec2/imds/

replace github.com/mvdatacenter/mvdata-sdk-go/internal/configsources => ../internal/configsources/

replace github.com/mvdatacenter/mvdata-sdk-go/internal/endpoints/v2 => ../internal/endpoints/v2/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/accept-encoding => ../service/internal/accept-encoding/

replace github.com/mvdatacenter/mvdata-sdk-go/service/internal/presigned-url => ../service/internal/presigned-url/

replace github.com/mvdatacenter/mvdata-sdk-go/service/sso => ../service/sso/

replace github.com/mvdatacenter/mvdata-sdk-go/service/ssooidc => ../service/ssooidc/

replace github.com/mvdatacenter/mvdata-sdk-go/service/sts => ../service/sts/
