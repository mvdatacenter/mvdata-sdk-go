// This package is designated as private and is intended for use only by the
// smithy client runtime. The exported API therein is not considered stable and
// is subject to breaking changes without notice.

package middleware

import (
	"context"
	"github.com/aws/smithy-go/middleware"
	"github.com/mvdatacenter/mvdata-sdk-go/internal/sdk"
	"github.com/mvdatacenter/mvdata-sdk-go/mvdata/middleware/private/metrics"
	"github.com/mvdatacenter/mvdata-sdk-go/mvdata/middleware/private/metrics/testutils"
	"testing"
	"time"
)

func TestEndpointResolutionEnd_HandleSerialize(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	ctx := metrics.InitMetricContext(context.TODO(), &metrics.SharedConnectionCounter{}, &testutils.NoopPublisher{})
	mw := GetRecordEndpointResolutionEndMiddleware()
	_, _, _ = mw.HandleSerialize(ctx, middleware.SerializeInput{}, testutils.NoopSerializeHandler{})

	actualTime := metrics.Context(ctx).Data().ResolveEndpointEndTime
	expectedTime := sdk.NowTime()
	if actualTime != expectedTime {
		t.Errorf("Unexpected ResolveEndpointEndTime, should be '%s' but was '%s'", expectedTime, actualTime)
	}
}
