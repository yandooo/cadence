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
	"database/sql"

	"github.com/uber/cadence/common/persistence/sql/sqlplugin"
)

const (
	createShardQry = `INSERT INTO
	cadence.shards (shard_id, range_id, data, data_encoding) VALUES (?, ?, ?, ?)`

	getShardQry = `SELECT
 shard_id, range_id, data, data_encoding
 FROM cadence.shards WHERE shard_id = ?`

	updateShardQry = `UPDATE cadence.shards 
 SET range_id = ?, data = ?, data_encoding = ? 
 WHERE shard_id = ?`

	lockShardQry     = `SELECT range_id FROM cadence.shards WHERE shard_id = ? FOR UPDATE`
	readLockShardQry = `SELECT range_id FROM cadence.shards WHERE shard_id = ? FOR READ ONLY`
)

// InsertIntoShards inserts one or more rows into shards table
func (mdb *db) InsertIntoShards(row *sqlplugin.ShardsRow) (sql.Result, error) {
	return mdb.conn.Exec(createShardQry, row.ShardID, row.RangeID, row.Data, row.DataEncoding)
}

// UpdateShards updates one or more rows into shards table
func (mdb *db) UpdateShards(row *sqlplugin.ShardsRow) (sql.Result, error) {
	return mdb.conn.Exec(updateShardQry, row.RangeID, row.Data, row.DataEncoding, row.ShardID)
}

// SelectFromShards reads one or more rows from shards table
func (mdb *db) SelectFromShards(filter *sqlplugin.ShardsFilter) (*sqlplugin.ShardsRow, error) {
	var row sqlplugin.ShardsRow
	err := mdb.Get(&row, getShardQry, filter.ShardID)
	if err != nil {
		return nil, err
	}
	return &row, err
}

// ReadLockShards acquires a read lock on a single row in shards table
func (mdb *db) ReadLockShards(filter *sqlplugin.ShardsFilter) (int, error) {
	var rangeID int
	err := mdb.Get(&rangeID, readLockShardQry, filter.ShardID)
	return rangeID, err
}

// WriteLockShards acquires a write lock on a single row in shards table
func (mdb *db) WriteLockShards(filter *sqlplugin.ShardsFilter) (int, error) {
	var rangeID int
	err := mdb.Get(&rangeID, lockShardQry, filter.ShardID)
	return rangeID, err
}
