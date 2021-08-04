package db

import (
	"github.com/foxfurry/simple-rest/internal/book/domain/entity"
	"testing"
)

var mock1 = entity.Book{
	Author: "TestAuthor_1",
	Title: "TestTitle_1",
	Year: 1,
	Description: "TestDescription_1",
}

var mock2 = entity.Book{
	Author: "TestAuthor_2",
	Title: "TestTitle_2",
	Year: 2,
	Description: "TestDescription_2",
}

var mock3 = entity.Book{
	Author: "TestAuthor_3",
	Title: "TestTitle_3",
	Year: 3,
	Description: "TestDescription_3",
}

var testBookDBRepo *BookDBRepository

func TestBookDBRepository_GetBook(t *testing.T) {
	//obj, err := testBookDBRepo.SaveBook()
}