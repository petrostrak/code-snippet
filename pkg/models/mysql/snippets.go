package mysql

import (
	"database/sql"

	"github.com/petrostrak/code-snippet/pkg/models"
)

// DB.Query() is used for SELECT queries which return multiple rows.
// DB.QueryRow() is used for SELECT queries which return a single row.
// DB.Exec() is used for statements which don’t return rows (like INSERT and DELETE).

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
			 VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() method on the embedded connection pool to execute the statement.
	// This method returns a sql.Result object, which contains some basic information
	// about what happend when the statement was executed.
	rs, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use the QueryRow() on the connection pool to execute our sql
	// statement. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Rpw to the
	// corresponding field in the Snippet struct. If the row returns no rows,
	// then row.Scan() will return a sql.ErrNoRows error.
	err := row.Scan(
		&s.ID,
		&s.Content,
		&s.Created,
		&s.Created,
		&s.Expires,
	)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	// If everything went OK then return the snippet object.
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
