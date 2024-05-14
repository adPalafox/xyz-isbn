package processor

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"xyz-isbn/constant"
	"xyz-isbn/converter"
	"xyz-isbn/models"
	"xyz-isbn/updater"
)

type BookProcessor struct {
	wg  sync.WaitGroup
	csv csv.Writer
}

func NewBookProcessor(csvWriter io.Writer) *BookProcessor {
	return &BookProcessor{
		wg:  sync.WaitGroup{},
		csv: *csv.NewWriter(csvWriter),
	}
}

func (p *BookProcessor) ProcessBooks(ctx context.Context, books []models.Book) error {
	semaphore := make(chan struct{}, constant.MAX_CONCURRENT_CALL)

	for _, book := range books {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled")
		case semaphore <- struct{}{}:
			p.wg.Add(1)
			go func(b models.Book) {
				defer func() {
					p.wg.Done()
					<-semaphore
				}()

				isUpdated, err := converter.IsbnConverter(&b)
				if err != nil {
					return
				}

				if isUpdated {
					if err := updater.UpdateBook(b); err != nil {
						return
					}
				}

				if err := p.writeToCSV(b); err != nil {
					return
				}
			}(book)
		}
	}

	p.wg.Wait()
	return nil
}

func (p *BookProcessor) writeToCSV(book models.Book) error {
	if err := converter.CheckIfExists(&book); err != nil {
		fmt.Printf("Error validating ISBNs for book %d: %v\n", book.ID, err)
		return err
	}

	return p.csv.Write([]string{book.ISBN10, book.ISBN13})
}

func (p *BookProcessor) FetchBookList(ctx context.Context, url string) ([]models.Book, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching book list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch book list: status code %d", resp.StatusCode)
	}

	var bookResponse models.BookResponse
	if err := json.NewDecoder(resp.Body).Decode(&bookResponse); err != nil {
		fmt.Println(resp.Body)
		return nil, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return bookResponse.Data, nil
}
