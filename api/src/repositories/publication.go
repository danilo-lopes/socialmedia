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
	"context"
	"database/sql"
)

type PublicationsRepository struct {
	db *sql.DB
}

// NewPublicationRepository creates a Publication repository
func NewPublicationRepository(db *sql.DB) *PublicationsRepository {
	return &PublicationsRepository{db}
}

// Create create post in database
func (repository PublicationsRepository) Create(post models.Publication) (uint64, error) {
	statement, erro := repository.db.Prepare(
		"INSERT INTO publications (title, content, author_id) VALUES (?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	result, erro := statement.Exec(post.Title, post.Content, post.AuthorID)
	if erro != nil {
		return 0, erro
	}

	lastInsertedID, erro := result.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(lastInsertedID), nil
}

// SearchByID return one publication in database
func (repository PublicationsRepository) SearchByID(publicationID uint64) (models.Publication, error) {
	line, erro := repository.db.Query(`
		SELECT p.*, u.nick FROM publications p JOIN users u
		ON u.id = p.author_id WHERE p.id = ?`,
		publicationID,
	)
	if erro != nil {
		return models.Publication{}, erro
	}
	defer line.Close()

	var publication models.Publication

	if line.Next() {
		if erro := line.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorID,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); erro != nil {
			return models.Publication{}, erro
		}
	}

	return publication, nil
}

// Get return all publications related by an user(self publications and his friends publications)
func (repository PublicationsRepository) Get(userID uint64) ([]models.Publication, error) {
	lines, erro := repository.db.Query(`
		SELECT DISTINCT p.*, u.nick from publications p
		JOIN users u on u.id = p.author_id
		JOIN followers f on p.author_id = f.user_id
		WHERE u.id = ? OR f.follower_id = ?
	`, userID, userID)
	if erro != nil {
		return nil, erro
	}
	defer lines.Close()

	var publications []models.Publication

	for lines.Next() {
		var publication models.Publication

		if erro = lines.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorID,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); erro != nil {
			return []models.Publication{}, erro
		}

		publications = append(publications, publication)
	}

	return publications, nil
}

// Update updates an publication in database
func (repository PublicationsRepository) Update(publicationID uint64, publication models.Publication) error {
	statement, erro := repository.db.Prepare(
		"UPDATE publications SET title = ?, content = ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(publication.Title, publication.Content, publicationID); erro != nil {
		return erro
	}

	return nil
}

// Delete deletes an publication in database
func (repository PublicationsRepository) Delete(publicationID uint64) error {
	statement, erro := repository.db.Prepare(
		"DELETE FROM publications WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(publicationID); erro != nil {
		return erro
	}

	return nil
}

// GetByUser return all publications from an user
func (repository PublicationsRepository) GetByUser(userID uint64) ([]models.Publication, error) {
	lines, erro := repository.db.Query(`
		SELECT p.*, u.nick from publications
		JOIN users u on u.id = p.author_id
		WHERE p.author_id = ?
	`, userID)
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
			&publication.AuthorNick,
		); erro != nil {
			return nil, erro
		}

		publications = append(publications, publication)
	}

	return publications, nil
}

// LikePublication likes an publication in database
func (repository PublicationsRepository) LikePublication(publicationID, likerID uint64) error {
	ctx := context.Background()
	tx, erro := repository.db.BeginTx(ctx, nil)
	if erro != nil {
		return erro
	}

	sqlQueries := []string{
		"INSERT INTO likes_of_publications(publication_id, liker_id) VALUES (?, ?)",
		"UPDATE publications SET likes = likes + 1 WHERE id = ?",
	}

	for sqlQueryIndex, query := range sqlQueries {
		if sqlQueryIndex == 0 {
			if _, erro = tx.ExecContext(
				ctx,
				query,
				publicationID, likerID,
			); erro != nil {
				tx.Rollback()
				return erro
			}
		} else if sqlQueryIndex == 1 {
			if _, erro = tx.ExecContext(
				ctx,
				query,
				publicationID,
			); erro != nil {
				tx.Rollback()
				return erro
			}
		}
	}

	if erro := tx.Commit(); erro != nil {
		return erro
	}

	return nil
}

// UnLikePublication unlikes an publication in database
func (repository PublicationsRepository) UnLikePublication(publicationID, unLikerID uint64) error {
	ctx := context.Background()
	tx, erro := repository.db.BeginTx(ctx, nil)
	if erro != nil {
		return erro
	}

	sqlQueries := []string{
		`
		UPDATE publications SET likes =
		CASE
			WHEN likes > 0 THEN likes - 1
		ELSE 0
		END
		WHERE id = ?
		`,
		`
		DELETE FROM likes_of_publications WHERE publication_id = ? AND liker_id = ?
		`,
	}

	for sqlQueryIndex, query := range sqlQueries {
		if sqlQueryIndex == 0 {
			if _, erro = tx.ExecContext(
				ctx,
				query,
				publicationID,
			); erro != nil {
				tx.Rollback()
				return erro
			}
		} else if sqlQueryIndex == 1 {
			if _, erro = tx.ExecContext(
				ctx,
				query,
				publicationID, unLikerID,
			); erro != nil {
				tx.Rollback()
				return erro
			}
		}
	}

	if erro := tx.Commit(); erro != nil {
		return erro
	}

	return nil
}

// GetLikers return all users who like an publication
func (repository PublicationsRepository) GetLikers(publicationID uint64) ([]models.User, error) {
	lines, erro := repository.db.Query(`
		SELECT u.id, u.name, u.nick, u.createdat FROM
		users u JOIN likes_of_publications p on u.id = p.liker_id WHERE p.publication_id = ?;
	`, publicationID)
	if erro != nil {
		return nil, erro
	}

	var users []models.User

	for lines.Next() {
		var user models.User

		if erro := lines.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.CreatedAt,
		); erro != nil {
			return nil, erro
		}

		users = append(users, user)
	}

	return users, nil
}
