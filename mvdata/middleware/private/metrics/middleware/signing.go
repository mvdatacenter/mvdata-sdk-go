package middleware

import (
	"context"

	"github.com/aws/smithy-go/middleware"
	"github.com/mvdatacenter/mvdata-sdk-go/internal/sdk"
	"github.com/mvdatacenter/mvdata-sdk-go/mvdata/middleware/private/metrics"
)

func timeSigning(stack *middleware.Stack) error {
	if err := stack.Finalize.Insert(signingStart{}, "Signing", middleware.Before); err != nil {
		return err
	}
	if err := stack.Finalize.Insert(signingEnd{}, "Signing", middleware.After); err != nil {
		return err
	}
	return nil
}

type signingStart struct{}

func (m signingStart) ID() string { return "signingStart" }

func (m signingStart) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	out middleware.FinalizeOutput, md middleware.Metadata, err error,
) {
	mctx := metrics.Context(ctx)
	attempt, err := mctx.Data().LatestAttempt()
	if err != nil {
		return out, md, err
	}

	attempt.SignStartTime = sdk.NowTime()
	return next.HandleFinalize(ctx, in)
}

type signingEnd struct{}

func (m signingEnd) ID() string { return "signingEnd" }

func (m signingEnd) HandleFinalize(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	out middleware.FinalizeOutput, md middleware.Metadata, err error,
) {
	mctx := metrics.Context(ctx)
	attempt, err := mctx.Data().LatestAttempt()
	if err != nil {
		return out, md, err
	}

	attempt.SignEndTime = sdk.NowTime()
	return next.HandleFinalize(ctx, in)
}
