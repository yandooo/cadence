// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package db2

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"

	"github.com/uber/cadence/common/persistence/sql"
	"github.com/uber/cadence/common/persistence/sql/sqlplugin"
	"github.com/uber/cadence/common/service/config"
)

const (
	// PluginName is the name of the plugin
	PluginName                     = "go_ibm_db"
	dataSourceNamePostgres         = "HOSTNAME=%v;DATABASE=%v;PORT=%v;UID=%v;PWD=%v"
	dataSourceNamePostgresPassword = "%v"
)

var errTLSNotImplemented = errors.New("tls for db2 has not been implemented")

type plugin struct{}

var _ sqlplugin.Plugin = (*plugin)(nil)

func init() {
	sql.RegisterPlugin(PluginName, &plugin{})
}

// CreateDB initialize the db object
func (d *plugin) CreateDB(cfg *config.SQL) (sqlplugin.DB, error) {
	conn, err := d.createDBConnection(cfg)
	if err != nil {
		return nil, err
	}
	db := newDB(conn, nil, cfg.DatabaseName)
	return db, nil
}

// CreateAdminDB initialize the adminDB object
func (d *plugin) CreateAdminDB(cfg *config.SQL) (sqlplugin.AdminDB, error) {
	conn, err := d.createDBConnection(cfg)
	if err != nil {
		return nil, err
	}
	db := newDB(conn, nil, cfg.DatabaseName)
	return db, nil
}

// CreateDBConnection creates a returns a reference to a logical connection to the
// underlying SQL database. The returned object is to tied to a single
// SQL database and the object can be used to perform CRUD operations on
// the tables in the database
func (d *plugin) createDBConnection(cfg *config.SQL) (*sqlx.DB, error) {
	err := registerTLSConfig(cfg)
	if err != nil {
		return nil, err
	}

	host, port, err := net.SplitHostPort(cfg.ConnectAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid connect address, it must be in host:port format, %v, err: %v", cfg.ConnectAddr, err)
	}

	dbName := cfg.DatabaseName
	//NOTE: db2 doesn't allow to connect with empty dbName, the admin dbName is "cdb"
	if dbName == "" {
		dbName = "cdb"
	}
	dbParts := strings.Split(dbName, "/")

	password := ""
	if cfg.Password != "" {
		password = fmt.Sprintf(dataSourceNamePostgresPassword, cfg.Password)
	}
	db, err := sqlx.Connect(PluginName, fmt.Sprintf(dataSourceNamePostgres, host, dbParts[0], port, cfg.User, password))

	if err != nil {
		return nil, err
	}

	// Maps struct names in CamelCase to snake without need for db struct tags.
	db.MapperFunc(func(s string) string { return strings.ToUpper(strcase.ToSnake(s)) })

	return db, nil
}

// TODO: implement postgres specific support for TLS
func registerTLSConfig(cfg *config.SQL) error {
	if cfg.TLS == nil || !cfg.TLS.Enabled {
		return nil
	}
	return errTLSNotImplemented
}
