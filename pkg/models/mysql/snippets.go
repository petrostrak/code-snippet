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
		&s.Title,
		&s.Content,
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
	// Write the SQL statement
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	// Use the Query() on the connection pool to execute  our SQL statement.
	// This returns a sql.Rows resultset containing the  result of our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is always properly
	// closed before the Latest() returns. This defer statement should come after
	// the error check, otherwise if Query() returs an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the resultset.
	snippets := []*models.Snippet{}

	//Use rows.Next() to iterate through the rows in the resultset. This prepares
	// the first (and then each subsequent) row to be acted on by the rows.Scan()
	// method. If iteration over all the rows completes then the resultset automatically
	// closes itself and frees-up the underlying database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed snippet struct
		s := &models.Snippet{}

		if err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Content,
			&s.Created,
			&s.Expires,
		); err != nil {
			return nil, err
		}

		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any error
	// that was encountered during the iteration. It's important to call this - don't
	// assume that a successful iteration was completed over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil
}
