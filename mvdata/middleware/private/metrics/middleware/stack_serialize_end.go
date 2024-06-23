package middleware

import (
	"context"
	"github.com/aws/smithy-go/middleware"
	"github.com/mvdatacenter/mvdata-sdk-go/internal/sdk"
	"github.com/mvdatacenter/mvdata-sdk-go/mvdata/middleware/private/metrics"
)

type StackSerializeEnd struct{}

func GetRecordStackSerializeEndMiddleware() *StackSerializeEnd {
	return &StackSerializeEnd{}
}

func (m *StackSerializeEnd) ID() string {
	return "StackSerializeEnd"
}

func (m *StackSerializeEnd) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {

	mctx := metrics.Context(ctx)
	mctx.Data().SerializeEndTime = sdk.NowTime()

	out, metadata, err = next.HandleSerialize(ctx, in)

	return out, metadata, err

}
