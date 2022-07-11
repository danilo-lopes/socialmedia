/*
Copyright 2022 Danilo S. Lopes.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at:

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package repositories

import (
	"context"
	"database/sql"
	"net"
)

type HealthcheckRepository struct {
	db *sql.DB
}

func NewHealthcheckRepository(db *sql.DB) *HealthcheckRepository {
	return &HealthcheckRepository{db}
}

// PingDatabase check connectivity to database
func (repository HealthcheckRepository) PingDatabase() error {
	if erro := repository.db.Ping(); erro != nil {
		return erro
	}

	return nil
}

// DNSResolver check name resolution
func (repository HealthcheckRepository) DNSResolver(address string) error {
	if _, erro := net.LookupHost(address); erro != nil {
		return erro
	}

	return nil
}

// SimulateDatabaseInsert check if our API can insert data in database
func (repository HealthcheckRepository) SimulateDatabaseInsert() error {
	ctx := context.Background()
	tx, erro := repository.db.BeginTx(ctx, nil)
	if erro != nil {
		return erro
	}

	mockSqlQuery := "INSERT INTO users (name, nick, email, pass) VALUES ('mock_user', 'mock_user_123', 'mock_user_123@gmail.com', 'password')"
	_, erro = tx.ExecContext(ctx, mockSqlQuery)
	if erro != nil {
		tx.Rollback()
		return erro
	}

	tx.Rollback()
	return nil
}
