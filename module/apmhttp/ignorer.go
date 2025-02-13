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

package apmhttp // import "go.elastic.co/apm/module/apmhttp"

import (
	"net/http"
	"regexp"
	"sync"

	"go.elastic.co/apm"
)

const (
	envIgnoreURLs           = "ELASTIC_APM_TRANSACTION_IGNORE_URLS"
	deprecatedEnvIgnoreURLs = "ELASTIC_APM_IGNORE_URLS"
)

var (
	defaultServerRequestIgnorerOnce sync.Once
	defaultServerRequestIgnorer     RequestIgnorerFunc = IgnoreNone
)

// NewDynamicServerRequestIgnorer returns the RequestIgnorer to use in
// handlers. The list of wildcard patterns comes from central config
func NewDynamicServerRequestIgnorer(t *apm.Tracer) RequestIgnorerFunc {
	return func(r *http.Request) bool {
		return t.IgnoredTransactionURL(r.URL)
	}
}

// NewRegexpRequestIgnorer returns a RequestIgnorerFunc which matches requests'
// URLs against re. Note that for server requests, typically only Path and
// possibly RawQuery will be set, so the regular expression should take this
// into account.
func NewRegexpRequestIgnorer(re *regexp.Regexp) RequestIgnorerFunc {
	if re == nil {
		panic("re == nil")
	}
	return func(r *http.Request) bool {
		return re.MatchString(r.URL.String())
	}
}

// IgnoreNone is a RequestIgnorerFunc which ignores no requests.
func IgnoreNone(*http.Request) bool {
	return false
}
