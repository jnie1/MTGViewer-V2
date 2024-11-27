package books

import (
	"database/sql"

	"github.com/jnie1/MTGViewer-V2/database"
)

func GetBooks() ([]BookData, error) {
	books := []BookData{}

	db := database.Instance()
	row, err := db.Query(`
		SELECT books.name, author.name
		FROM books
		JOIN author ON books.author_id = author.id`)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		book := BookData{}

		if err = row.Scan(&book.Name, &book.AuthorName); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func GetBook(bookId int) (BookData, error) {
	book := BookData{}

	db := database.Instance()
	row := db.QueryRow(`
		SELECT books.name, author.name
		FROM books
		JOIN author ON books.author_id = author.id
		WHERE id = $1`, bookId)

	err := row.Scan(&book.Name, &book.AuthorName)

	return book, err
}

func AddBook(book BookData) error {
	db := database.Instance()
	authorId, err := getAuthorId(db, book.AuthorName)

	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO books (name, author_id) VALUES ($1, $2)", book.Name, authorId)
	return err
}

func getAuthorId(db *sql.DB, authorName string) (int, error) {
	var authorId int

	row := db.QueryRow("SELECT id FROM author WHERE name = $1", authorName)

	if err := row.Scan(&authorId); err == nil {
		return authorId, nil
	}

	row = db.QueryRow("INSERT INTO author (name) VALUES ($1) RETURNING id", authorName)
	err := row.Scan(&authorId)

	return authorId, err
}
