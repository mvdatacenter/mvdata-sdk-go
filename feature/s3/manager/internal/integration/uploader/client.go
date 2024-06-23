//go:build integration && perftest
// +build integration,perftest

package uploader

import (
	"net"
	"net/http"
	"time"

	"github.com/mvdatacenter/mvdata-sdk-go/mvdata"
	awshttp "github.com/mvdatacenter/mvdata-sdk-go/mvdata/transport/http"
)

func NewHTTPClient(cfg ClientConfig) aws.HTTPClient {
	return awshttp.NewBuildableClient().WithTransportOptions(func(transport *http.Transport) {
		*transport = http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   cfg.Timeouts.Connect,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			IdleConnTimeout:     90 * time.Second,

			DisableKeepAlives:     !cfg.KeepAlive,
			TLSHandshakeTimeout:   cfg.Timeouts.TLSHandshake,
			ExpectContinueTimeout: cfg.Timeouts.ExpectContinue,
			ResponseHeaderTimeout: cfg.Timeouts.ResponseHeader,
		}
	})
}
