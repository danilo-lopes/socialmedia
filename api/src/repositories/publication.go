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

// SearchByID return one publication from database
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
