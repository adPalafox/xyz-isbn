package processor

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"xyz-isbn/constant"
	"xyz-isbn/converter"
	"xyz-isbn/models"
	"xyz-isbn/updater"
)

type BookProcessor struct {
	wg            sync.WaitGroup
	csv           csv.Writer
	mtx           sync.Mutex
	existingISBNs map[string]struct{}
}

func NewBookProcessor(csvWriter io.Writer) *BookProcessor {
	processor := &BookProcessor{
		wg:            sync.WaitGroup{},
		csv:           *csv.NewWriter(csvWriter),
		mtx:           sync.Mutex{},
		existingISBNs: make(map[string]struct{}),
	}

	existingISBNs, err := readExistingISBNs()
	if err != nil {
		return nil
	}
	processor.existingISBNs = make(map[string]struct{}, len(existingISBNs))
	for _, isbn := range existingISBNs {
		processor.existingISBNs[isbn] = struct{}{}
	}

	return processor
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

	p.mtx.Lock()
	defer p.mtx.Unlock()

	err := p.csv.Write([]string{book.ISBN10, book.ISBN13})
	if err != nil {
		return fmt.Errorf("error writing book data to CSV: %w", err)
	}

	p.csv.Flush()

	return nil
}

func readExistingISBNs() ([]string, error) {
	file, err := os.Open(constant.CSV_FILE)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var existingISBNs []string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		isbn10 := record[0]
		isbn13 := record[1]

		existingISBNs = append(existingISBNs, isbn10, isbn13)
	}

	return existingISBNs, nil
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
