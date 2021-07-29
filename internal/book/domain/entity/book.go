package entity

type Book struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Year        int    `json:"yeah"`
	Description string `json:"description"`
}
