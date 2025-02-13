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

// Package apmpostgres imports the gorm mysql dialect package,
// and also registers the mysql driver with apmsql.
package apmpostgres // import "go.elastic.co/apm/module/apmgormv2/driver/postgres"

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	apmpgxv4 "go.elastic.co/apm/module/apmsql/pgxv4"
)

// Open creates a dialect with apmsql
func Open(dsn string) gorm.Dialector {
	dialect := &postgres.Dialector{
		Config: &postgres.Config{
			DriverName: apmpgxv4.DriverName,
			DSN:        dsn,
		},
	}

	return dialect
}
