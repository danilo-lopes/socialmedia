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
	"api/src/models"
	"database/sql"
	"fmt"
)

type users struct {
	db *sql.DB
}

// NewUsersRepositories creates a users repositories
func NewUsersRepositories(db *sql.DB) *users {
	return &users{db}
}

// Create creates a User into database
func (repository users) Create(user models.User) (uint64, error) {

	statement, erro := repository.db.Prepare(
		"INSERT INTO users (name, nick, email, pass) VALUES (?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	result, erro := statement.Exec(user.Name, user.Nick, user.Email, user.Pass)
	if erro != nil {
		return 0, erro
	}

	lastID, erro := result.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(lastID), nil
}

// Search return all Users or Nicknames matching with the filter(nameOrNick)
func (repository users) Search(nameOrNick string) ([]models.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick) // %nameOrNick%

	lines, erro := repository.db.Query(
		"SELECT id, name, nick, email, createdat FROM users WHERE name LIKE ? OR nick LIKE ?",
		nameOrNick, nameOrNick,
	)
	if erro != nil {
		return nil, erro
	}

	defer lines.Close()

	var users []models.User

	for lines.Next() {
		var user models.User

		if erro := lines.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); erro != nil {
			return nil, erro
		}

		users = append(users, user)
	}

	return users, nil
}

// SearchByID return the User matching with the ID
func (repository users) SearchByID(ID uint64) (models.User, error) {
	lines, erro := repository.db.Query(
		"SELECT id, name, nick, email, createdat FROM users WHERE id = ?",
		ID,
	)
	if erro != nil {
		return models.User{}, erro
	}
	defer lines.Close()

	var user models.User

	if lines.Next() {
		if erro := lines.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
		); erro != nil {
			return models.User{}, erro
		}
	}

	return user, nil
}

// SearchEmail search an user by email and returns the id and the password hash
func (repository users) SearchByEmail(email string) (models.User, error) {
	line, erro := repository.db.Query(
		"SELECT id, pass FROM users WHERE email = ?",
		email,
	)
	if erro != nil {
		return models.User{}, erro
	}
	defer line.Close()

	var user models.User

	if line.Next() {
		if erro := line.Scan(
			&user.ID,
			&user.Pass,
		); erro != nil {
			return models.User{}, erro
		}
	}

	return user, nil
}

// Update updates an Users attributes into database
func (repository users) Update(ID uint64, user models.User) error {

	statement, erro := repository.db.Prepare(
		"UPDATE users SET name = ?, nick = ?, email = ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(user.Name, user.Nick, user.Email, ID); erro != nil {
		return erro
	}

	return nil
}

// Delete delete an User into database
func (repository users) Delete(ID uint64) error {

	statement, erro := repository.db.Prepare(
		"DELETE FROM users WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}
