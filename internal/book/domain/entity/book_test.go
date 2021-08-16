package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBook_Equal(t *testing.T) {
	testCasesEqual := []Book{
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
	}
	testCasesDifferent := []Book{
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          2,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "2",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "2",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        2,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "2",
		},
	}

	assert.True(t, testCasesEqual[0].Equal(testCasesEqual[1]))

	for idx := 1; idx < len(testCasesDifferent); idx += 2 {
		assert.True(t, !testCasesDifferent[idx].Equal(testCasesDifferent[idx-1]))
	}
}

func TestBook_EqualNoID(t *testing.T) {
	testCasesEqual := []Book{
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          2,
			Title:       "2",
			Author:      "2",
			Year:        2,
			Description: "2",
		},
		{
			ID:          3,
			Title:       "2",
			Author:      "2",
			Year:        2,
			Description: "2",
		},
	}
	testCasesNotEqual := []Book{
		{
			ID:          1,
			Title:       "2",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "2",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        2,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			ID:          1,
			Title:       "1",
			Author:      "1",
			Year:        1,
			Description: "2",
		},
	}

	for idx := 1; idx < len(testCasesEqual); idx += 2 {
		assert.True(t, testCasesEqual[idx].EqualNoID(testCasesEqual[idx-1]))
	}

	for idx := 1; idx < len(testCasesNotEqual); idx += 2 {
		assert.True(t, !testCasesNotEqual[idx].EqualNoID(testCasesNotEqual[idx-1]))
	}
}

func TestBookArrayEqualNoID(t *testing.T) {
	testCasesEqual := [][]Book{
		{
			{
				ID:          1,
				Title:       "1",
				Author:      "1",
				Year:        1,
				Description: "1",
			},
			{
				ID:          2,
				Title:       "2",
				Author:      "2",
				Year:        2,
				Description: "2",
			},
		},
		{
			{
				ID:          3,
				Title:       "1",
				Author:      "1",
				Year:        1,
				Description: "1",
			},
			{
				ID:          4,
				Title:       "2",
				Author:      "2",
				Year:        2,
				Description: "2",
			},
		},
	}
	testCasesNotEqual := [][]Book{
		{
			{
				ID:          1,
				Title:       "1",
				Author:      "1",
				Year:        1,
				Description: "1",
			},
			{
				ID:          2,
				Title:       "2",
				Author:      "2",
				Year:        2,
				Description: "2",
			},
		},
		{
			{
				ID:          3,
				Title:       "3",
				Author:      "3",
				Year:        3,
				Description: "3",
			},
			{
				ID:          4,
				Title:       "4",
				Author:      "4",
				Year:        4,
				Description: "4",
			},
		},
	}
	testCaseNotEqualSize := [][]Book{
		{
			{
				ID:          1,
				Title:       "1",
				Author:      "1",
				Year:        1,
				Description: "1",
			},
			{
				ID:          2,
				Title:       "2",
				Author:      "2",
				Year:        2,
				Description: "2",
			},
		},
		{
			{
				ID:          3,
				Title:       "3",
				Author:      "3",
				Year:        3,
				Description: "3",
			},
			{
				ID:          4,
				Title:       "4",
				Author:      "4",
				Year:        4,
				Description: "4",
			},
			{
				ID:          5,
				Title:       "5",
				Author:      "5",
				Year:        5,
				Description: "5",
			},
		},
	}
	assert.True(t, BookArrayEqualNoID(testCasesEqual[0], testCasesEqual[1]))
	assert.True(t, !BookArrayEqualNoID(testCasesNotEqual[0], testCasesNotEqual[1]))
	assert.True(t, !BookArrayEqualNoID(testCaseNotEqualSize[0], testCaseNotEqualSize[1]))
}

func TestBook_IsValid(t *testing.T) {
	testCasesValid := []Book{
		{
			Title:       "Test Valid 1",
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			Title:       "Test valid 2",
			Author:      "2",
			Year:        2,
			Description: "2",
		},
	}
	testCasesInvalid := []Book{
		{
			Author:      "1",
			Year:        1,
			Description: "1",
		},
		{
			Title:       "2",
			Year:        2,
			Description: "2",
		},
		{
			Title:  "3",
			Author: "3",
		},
		{
			Description: "4",
		},
	}

	for _, tc := range testCasesValid {
		assert.True(t, tc.IsValid(), "Book expected to be valid, but found invalid: %v", tc)
	}
	for _, tc := range testCasesInvalid {
		assert.True(t, !tc.IsValid(), "Book expected to be invalid, but found valid: %v", tc)
	}
}
