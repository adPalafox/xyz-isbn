package models

type BookResponse struct {
	Data       []Book `json:"data"`
	TotalCount int    `json:"total_count"`
}

type Book struct {
	ID              int      `json:"id"`
	Title           string   `json:"title"`
	ISBN10          string   `json:"isbn_10"`
	ISBN13          string   `json:"isbn_13"`
	ListPrice       int      `json:"list_price"`
	PublicationYear int      `json:"publication_year"`
	Publisher       string   `json:"publisher"`
	ImageURL        string   `json:"image_url"`
	Edition         string   `json:"edition"`
	Authors         []Author `json:"authors"`
}

type Author struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}
