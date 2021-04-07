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

package apmawssdkgo // import "go.elastic.co/apm/module/apmawssdkgo"

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/stacktrace"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
)

func init() {
	stacktrace.RegisterLibraryPackage(
		"github.com/aws/aws-sdk-go/aws/request",
		"github.com/aws/aws-sdk-go/aws/session",
	)
}

// WrapSession wraps the provided s3 session with handlers that hook into the
// aws sdk's request lifecycle.
func WrapSession(s *session.Session) *session.Session {
	s.Handlers.Send.PushFrontNamed(request.NamedHandler{
		Name: "go.elastic.co/apm/module/apmawssdkgo/send",
		Fn:   send,
	})
	s.Handlers.Complete.PushBackNamed(request.NamedHandler{
		Name: "go.elastic.co/apm/module/apmawssdkgo/complete",
		Fn:   complete,
	})

	return s
}

var (
	serviceTypeMap = map[string]string{
		"s3":       "storage",
		"dynamodb": "dynamodb",
	}
)

func send(req *request.Request) {
	if req.RetryCount > 0 {
		return
	}

	spanSubtype := req.ClientInfo.ServiceName
	spanType, prs := serviceTypeMap[spanSubtype]
	if !prs {
		return
	}

	var (
		err        error
		targetName string
		ctx        = req.Context()
		region     = *req.Config.Region
		values     = new(dynamoDBValues)
	)

	tx := apm.TransactionFromContext(ctx)
	if tx == nil {
		return
	}

	switch spanSubtype {
	case "s3":
		targetName = getBucketName(req)
	case "dynamodb":
		// Read the request body to parse out the TableName and
		// Statement, then supply a new io.ReadCloser with a copy of
		// the original body to the request.
		req.HTTPRequest.Body, values, err = getDynamoDBValues(req.HTTPRequest.Body)
		if err != nil {
			return
		}
		targetName = values.TableName
	default:
		// Unsupported type
		return
	}

	spanName := req.ClientInfo.ServiceID + " " + req.Operation.Name + " " + targetName
	span := tx.StartSpan(spanName, spanType, apm.SpanFromContext(ctx))
	if !span.Dropped() {
		ctx = apm.ContextWithSpan(ctx, span)
		req.HTTPRequest = apmhttp.RequestWithContext(ctx, req.HTTPRequest)
		span.Context.SetHTTPRequest(req.HTTPRequest)
	} else {
		span.End()
		span = nil
	}

	span.Subtype = spanSubtype
	span.Action = req.Operation.Name

	span.Context.SetDestinationService(apm.DestinationServiceSpanContext{
		Name:     spanSubtype,
		Resource: targetName,
		Type:     spanType,
	})
	span.Context.SetDestinationCloud(apm.DestinationCloudSpanContext{
		Region: region,
	})

	if spanType == "dynamodb" {
		dbSpanCtx := apm.DatabaseSpanContext{
			Instance: region,
			Type:     spanType,
			// TODO: What do we put here? Most users will just be
			// grabbing their environment variables.
			User: "",
		}
		if span.Action == "Query" {
			dbSpanCtx.Statement = values.KeyConditionExpression
		}
		span.Context.SetDatabase(dbSpanCtx)
	}

	req.SetContext(ctx)
}

func complete(req *request.Request) {
	ctx := req.Context()
	span := apm.SpanFromContext(ctx)
	if span.Dropped() {
		return
	}
	defer span.End()

	span.Context.SetHTTPStatusCode(req.HTTPResponse.StatusCode)

	if err := req.Error; err != nil {
		apm.CaptureError(ctx, err).Send()
	}
}

type dynamoDBValues struct {
	TableName string
	// KeyConditionExpression is only available on Query operations.
	KeyConditionExpression string
}

func getDynamoDBValues(r io.ReadCloser) (io.ReadCloser, *dynamoDBValues, error) {
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	b := dynamoDBValues{}
	json.Unmarshal(body, &b)

	return ioutil.NopCloser(bytes.NewBuffer(body)), &b, nil
}

func getBucketName(req *request.Request) string {
	host := req.HTTPRequest.URL.Host
	if strings.HasPrefix(host, req.ClientInfo.ServiceName) {
		return strings.Split(req.HTTPRequest.URL.Path[1:], "/")[0]
	}
	return strings.Split(req.HTTPRequest.URL.Host, ".")[0]
}
