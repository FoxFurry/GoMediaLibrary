package entity

import "log"

type Book struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Year        int    `json:"year"`
	Description string `json:"description"`
}

// Equal returns true if all fields in receiver are same as in parameter
func (lhs Book) Equal(rhs Book) bool{
	return rhs==lhs
}

// EqualNoID works similar to Equal, except it ignores ID
func (lhs Book) EqualNoID(rhs Book) bool{
	log.Printf("---------------\n%v\n%v\n---------------------", lhs, rhs)
	rhs.ID = lhs.ID
	return rhs==lhs
}
