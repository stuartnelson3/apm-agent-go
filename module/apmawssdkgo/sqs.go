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
	"fmt"

	"github.com/aws/aws-sdk-go/aws/request"
	"go.elastic.co/apm"
)

var operationToNames = map[string]string{
	// TODO: Get the actual key names
	deleteMessage:      "delete",
	deleteMessageBatch: "delete_batch",
	receiveMessage:     "poll",
	sendMessageBatch:   "send_batch",
	sendMessage:        "send",
	unknown:            "unknown",
}

type sqs struct {
	name string
}

// <MSG-FRAMEWORK> SEND/RECEIVE/POLL to/from <QUEUE-NAME>

func newSQS(req *request.Request) *sqs {
	queueName := "find me!"
	operationName := operationsToNames[req.Operation.Name]
	resource := req.ClientInfo.ServiceName + "/" + queueName
	name := fmt.Sprintf("%s %s %s %s",
		req.ClientInfo.ServiceID,
		req.Operation.Name,
		operationDirection(operationName),
		queueName,
	)

	return &sqs{}
}

func (s *sqs) name() string { return s.name }

func (s *sqs) resource() string { return "" }

func (s *sqs) setAdditional(*apm.Span) {}

// func operationDirection(operationName
