/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package source

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/projectriff/http-gateway/pkg/client"
	"strings"
)

type source struct {
	query        string
	update       *sql.Stmt
	streamClient *client.StreamClient
	db           *sql.DB
}

func NewSource(query string, update string, dataSourceName string,  streamClient *client.StreamClient) (*source, error) {

	var db *sql.DB
	var updateStmt *sql.Stmt
	var err error
	if db, err = sql.Open("mysql", dataSourceName); err != nil {
		return &source{}, err
	}

	if updateStmt, err = db.Prepare(update); err != nil {
		return &source{}, err
	}

	return &source{
		query:        query,
		update:       updateStmt,
		streamClient: streamClient,
		db:           db,
	}, nil
}

func (s *source) Run(ctx context.Context) (int, error) {
	i := 0
	var id interface{}
	var data string
	if rows, err := s.db.Query(s.query); err != nil {
		return 0, err
	} else {
		for rows.Next() {
			if err := rows.Scan(&id, &data); err != nil {
				return 0, err
			}
			if _, err := s.streamClient.Publish(ctx, strings.NewReader(data), nil, "text/plain", nil); err != nil {
				return 0, err
			} else {
				if _, err := s.update.Exec(id); err != nil {
					return 0, err
				}
				i = i + 1
			}
		}
	}
	return i, nil
}

func (s *source) Close() error {
	if err := s.db.Close() ; err != nil {
		return err
	} else {
		return s.streamClient.Close()
	}
}
