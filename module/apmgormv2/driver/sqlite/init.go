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

// Package apmsqlite imports the gorm sqlite dialect package,
// and also registers the sqlite3 driver with apmsql.
package apmsqlite // import "go.elastic.co/apm/module/apmgormv2/driver/sqlite"

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go.elastic.co/apm/module/apmsql"

	_ "go.elastic.co/apm/module/apmsql/sqlite3" // register sqlite3 with apmsql
)

// Open creates a dialect with apmsql
func Open(dsn string) gorm.Dialector {

	dialect := &sqlite.Dialector{
		DSN:        dsn,
		DriverName: apmsql.DriverPrefix + sqlite.DriverName,
	}

	return dialect
}
