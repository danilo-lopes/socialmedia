package repositories

import (
	"api/src/models"
	"database/sql"
)

type PublicationsRepository struct {
	db *sql.DB
}

// NewPublicationRepository creates a "Publication" repository
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
func (repository PublicationsRepository) LikePublication(publicationID uint64) error {

	statement, erro := repository.db.Prepare("UPDATE publications SET likes = likes + 1 WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(publicationID); erro != nil {
		return erro
	}

	return nil
}

// UnLikePublication unlikes an publication in database
func (repository PublicationsRepository) UnLikePublication(publicationID uint64) error {

	statement, erro := repository.db.Prepare(`
		UPDATE publications SET likes =
		CASE
			WHEN likes > 0 THEN likes - 1
		ELSE 0
		END
		WHERE id = ?
	`)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(publicationID); erro != nil {
		return erro
	}

	return nil
}
