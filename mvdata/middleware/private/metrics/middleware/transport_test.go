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

func TestTransportMetrics_HandleSerialize(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, &testutils.NoopPublisher{})

	data := metrics.Context(ctx).Data()

	data.NewAttempt()

	mw := GetTransportMetricsMiddleware()
	_, _, _ = mw.HandleDeserialize(ctx, middleware.DeserializeInput{}, testutils.NoopDeserializeHandler{})

	attempt, _ := data.LatestAttempt()

	actualStartTime := attempt.ServiceCallStart
	expectedStartTime := sdk.NowTime()

	if actualStartTime != expectedStartTime {
		t.Errorf("Unexpected ServiceCallStart, should be '%s' but was '%s'", expectedStartTime, expectedStartTime)
	}

	actualEndTime := attempt.ServiceCallEnd
	expectedEndTime := sdk.NowTime()

	if actualEndTime != expectedEndTime {
		t.Errorf("Unexpected ServiceCallEnd, should be '%s' but was '%s'", expectedEndTime, actualEndTime)
	}

}
