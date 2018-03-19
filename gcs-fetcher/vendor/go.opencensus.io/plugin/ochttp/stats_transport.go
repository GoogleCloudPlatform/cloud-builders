// Copyright 2018, OpenCensus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ochttp

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

// statsTransport is an http.RoundTripper that collects stats for the outgoing requests.
type statsTransport struct {
	base    http.RoundTripper
	sampler trace.Sampler
}

// RoundTrip implements http.RoundTripper, delegating to Base and recording stats for the request.
func (t statsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, _ := tag.New(req.Context(),
		tag.Upsert(Host, req.URL.Host),
		tag.Upsert(Path, req.URL.Path),
		tag.Upsert(Method, req.Method))
	req = req.WithContext(ctx)
	track := &tracker{
		start: time.Now(),
		ctx:   ctx,
	}
	if req.Body == nil {
		// TODO: Handle cases where ContentLength is not set.
		track.reqSize = -1
	} else if req.ContentLength > 0 {
		track.reqSize = req.ContentLength
	}
	stats.Record(ctx, ClientRequestCount.M(1))

	// Perform request.
	resp, err := t.base.RoundTrip(req)

	if err != nil {
		track.statusCode = "error"
		track.end()
	} else {
		track.statusCode = strconv.Itoa(resp.StatusCode)
		if resp.Body == nil {
			track.end()
		} else {
			track.body = resp.Body
			resp.Body = track
		}
	}

	return resp, err
}

// CancelRequest cancels an in-flight request by closing its connection.
func (t statsTransport) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}
	if cr, ok := t.base.(canceler); ok {
		cr.CancelRequest(req)
	}
}

type tracker struct {
	respSize   int64
	reqSize    int64
	ctx        context.Context
	start      time.Time
	body       io.ReadCloser
	statusCode string
	endOnce    sync.Once
}

func (t *tracker) end() {
	t.endOnce.Do(func() {
		m := []stats.Measurement{
			ClientLatency.M(float64(time.Since(t.start)) / float64(time.Millisecond)),
			ClientResponseBytes.M(t.respSize),
		}
		if t.reqSize >= 0 {
			m = append(m, ClientRequestBytes.M(t.reqSize))
		}
		ctx, _ := tag.New(t.ctx, tag.Upsert(StatusCode, t.statusCode))
		stats.Record(ctx, m...)
	})
}

var _ io.ReadCloser = (*tracker)(nil)

func (t *tracker) Read(b []byte) (int, error) {
	n, err := t.body.Read(b)
	switch err {
	case nil:
		t.respSize += int64(n)
		return n, nil
	case io.EOF:
		t.end()
	}
	return n, err
}

func (t *tracker) Close() error {
	// Invoking endSpan on Close will help catch the cases
	// in which a read returned a non-nil error, we set the
	// span status but didn't end the span.
	t.end()
	return t.body.Close()
}
