package entity

type Book struct {
	ID          uint64 `json:"id,omitempty" binding:"omitempty,validID"`
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Year        int    `json:"year" binding:"required,validYear"`
	Description string `json:"description,omitempty"`
}

// Equal returns true if all fields in receiver are same as in parameter
func (lhs Book) Equal(rhs Book) bool {
	return rhs == lhs
}

// EqualNoID works similar to Equal, except it ignores ID
func (lhs Book) EqualNoID(rhs Book) bool {
	rhs.ID = lhs.ID
	return rhs == lhs
}

// BookArrayEqualNoID compares two arrays of Book(s) using Book.EqualNoID on each element
func BookArrayEqualNoID(lhs []Book, rhs []Book) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for idx, _ := range lhs {
		if !lhs[idx].EqualNoID(rhs[idx]) {
			return false
		}
	}

	return true
}
