// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package apmotel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/module/apmotel"
	"go.elastic.co/apm/v2/transport/transporttest"
)

func TestSpanStartAttributes(t *testing.T) {
	tracer, apmtracer, recorder := newTestTracer()
	defer apmtracer.Close()

	tx := apmtracer.StartTransaction("root", "root")
	defer tx.End()
	ctx := context.Background()
	ctx = apm.ContextWithTransaction(ctx, tx)

	tcs := []struct {
		attrs                       []attribute.KeyValue
		spanKind                    trace.SpanKind
		spanType, subtype, resource string
	}{
		{
			attrs: []attribute.KeyValue{
				attribute.String("db.name", "myDB"),
				attribute.String("db.system", "dbSystem"),
			},
			spanKind: trace.SpanKindServer, // txType == request
		},
	}

	// TODO: Check when ctx contains span/tx
	for i, tc := range tcs {
		_, span := tracer.Start(ctx, fmt.Sprintf("tc%d", i), trace.WithAttributes(tc.attrs...), trace.WithSpanKind(tc.spanKind))
		span.End()
	}

	apmtracer.Flush(nil)
	payloads := recorder.Payloads()
	spans := payloads.Spans
	require.Len(t, spans, len(tcs))
	for i, tc := range tcs {
		assert.Equal(t, tc.spanType, spans[i].Type)
		assert.Equal(t, tc.subtype, spans[i].Subtype)
		assert.Equal(t, tc.resource, spans[i].Context.Destination.Service.Resource)
		assert.Equal(t, tc.spanKind, spans[i].Otel.SpanKind)
	}
}

func newTestTracer() (trace.Tracer, *apm.Tracer, *transporttest.RecorderTransport) {
	apmtracer, recorder := transporttest.NewRecorderTracer()
	tracer := apmotel.NewTracerProvider(apmtracer).Tracer("otel_tracer")
	return tracer, apmtracer, recorder
}
