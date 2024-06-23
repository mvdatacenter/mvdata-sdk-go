package middleware

import (
	"context"
	"github.com/aws/smithy-go/middleware"
	"github.com/mvdatacenter/mvdata-sdk-go/internal/sdk"
	"github.com/mvdatacenter/mvdata-sdk-go/mvdata/middleware/private/metrics"
)

type StackSerializeStart struct{}

func GetRecordStackSerializeStartMiddleware() *StackSerializeStart {
	return &StackSerializeStart{}
}

func (m *StackSerializeStart) ID() string {
	return "StackSerializeStart"
}

func (m *StackSerializeStart) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {

	mctx := metrics.Context(ctx)
	mctx.Data().SerializeStartTime = sdk.NowTime()

	out, metadata, err = next.HandleSerialize(ctx, in)

	return out, metadata, err
}
