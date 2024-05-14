package updater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"xyz-isbn/models"
)

// UpdateBook updates a book's ISBN via the update endpoint
func UpdateBook(book models.Book) error {
	url := fmt.Sprintf("http://localhost:8080/v1/api/book/%s", book.ISBN13)
	body := map[string]any{
		"title":            book.Title,
		"isbn_13":          book.ISBN13,
		"isbn_10":          book.ISBN10,
		"publisher":        book.Publisher,
		"publication_year": book.PublicationYear,
		"list_price":       book.ListPrice,
		"edition":          book.Edition,
		"authors":          book.Authors,
	}
	requestBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling update request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("error creating update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending update request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
