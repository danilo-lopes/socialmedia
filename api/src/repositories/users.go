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

type usersRepository struct {
	db *sql.DB
}

// NewUsersRepository creates a "users" repository
func NewUsersRepository(db *sql.DB) *usersRepository {
	return &usersRepository{db}
}

// Create creates a "User" in database
func (repository usersRepository) Create(user models.User) (uint64, error) {

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
func (repository usersRepository) Search(nameOrNick string) ([]models.User, error) {
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
func (repository usersRepository) SearchByID(ID uint64) (models.User, error) {
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

// SearchByEmail search an user by email and returns the id and the password hash
func (repository usersRepository) SearchByEmail(email string) (models.User, error) {
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
func (repository usersRepository) Update(ID uint64, user models.User) error {

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
func (repository usersRepository) Delete(ID uint64) error {

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

//Follow permits an User to follow another User
func (repository usersRepository) Follow(userID, followerID uint64) error {

	statement, erro := repository.db.Prepare(
		"INSERT IGNORE INTO followers (user_id, follower_id) VALUES (?, ?)",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(userID, followerID); erro != nil {
		return erro
	}

	return nil
}

//Follow permits an User to unfollow another User
func (repository usersRepository) UnFollow(userID, followerID uint64) error {

	statement, erro := repository.db.Prepare(
		"DELETE FROM followers WHERE user_id = ? and follower_id = ? ",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(userID, followerID); erro != nil {
		return erro
	}

	return nil
}

//GetFollowers return all followers from User
func (repository usersRepository) GetFollowers(userID uint64) ([]models.User, error) {

	lines, erro := repository.db.Query(`
		SELECT u.id, u.name, u.nick, u.email, u.createdat
		FROM users u INNER JOIN followers s on u.id = s.follower_id WHERE s.user_id = ?
	`, userID,
	)
	if erro != nil {
		return nil, erro
	}
	defer lines.Close()

	var followers []models.User

	for lines.Next() {
		var follower models.User

		if erro := lines.Scan(
			&follower.ID,
			&follower.Name,
			&follower.Nick,
			&follower.Email,
			&follower.CreatedAt,
		); erro != nil {
			return nil, erro
		}

		followers = append(followers, follower)
	}

	return followers, nil
}

//GetFollowing return all users one user is following
func (repository usersRepository) GetFollowing(userID uint64) ([]models.User, error) {

	lines, erro := repository.db.Query(`
		SELECT u.id, u.name, u.nick, u.email, u.createdat
		FROM users u INNER JOIN followers s on u.id = s.user_id WHERE s.follower_id = ?
	`, userID,
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

// GetUserPass return the user password thought ID
func (repository usersRepository) GetUserPass(userID uint64) (string, error) {

	line, erro := repository.db.Query(
		"SELECT pass FROM users WHERE id = ?",
		userID,
	)
	if erro != nil {
		return "", erro
	}
	defer line.Close()

	var user models.User

	if line.Next() {
		if erro := line.Scan(
			&user.Pass,
		); erro != nil {
			return "", erro
		}
	}

	return user.Pass, nil
}

// UpadateUserPass update the user pass
func (repository usersRepository) UpadateUserPass(userID uint64, pass string) error {

	statement, erro := repository.db.Prepare(
		"UPDATE users SET pass = ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(pass, userID); erro != nil {
		return erro
	}

	return nil
}

// LikedPublication return all publications and user liked
func (repository usersRepository) LikedPublications(userID uint64) ([]models.Publication, error) {

	lines, erro := repository.db.Query(`
		SELECT DISTINCT p.* FROM publications p
		JOIN likes_of_publications l on p.id = l.publication_id
		JOIN users u on u.id = p.author_id
		WHERE u.id = ? OR l.liker_id = ?;
	`, userID, userID,
	)
	if erro != nil {
		return nil, erro
	}
	defer lines.Close()

	var publications []models.Publication

	for lines.Next() {
		var publication models.Publication

		if erro := lines.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorID,
			&publication.Likes,
			&publication.CreatedAt,
		); erro != nil {
			return nil, erro
		}

		publications = append(publications, publication)
	}

	return publications, nil
}
