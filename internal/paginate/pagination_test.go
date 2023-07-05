package paginate

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestPagination_ShowPrevious(t *testing.T) {
	assert.False(t, Pagination{CurrentPage: 0}.ShowPrevious())
	assert.False(t, Pagination{CurrentPage: 1}.ShowPrevious())
	assert.True(t, Pagination{CurrentPage: 2}.ShowPrevious())
}

func TestPagination_ShowNext(t *testing.T) {
	assert.False(t, Pagination{CurrentPage: 0, TotalPages: 0}.ShowNext())
	assert.False(t, Pagination{CurrentPage: 1, TotalPages: 1}.ShowNext())
	assert.False(t, Pagination{CurrentPage: 2, TotalPages: 2}.ShowNext())
	assert.True(t, Pagination{CurrentPage: 1, TotalPages: 2}.ShowNext())
	assert.False(t, Pagination{CurrentPage: 2, TotalPages: 1}.ShowNext())
}

func TestPagination_GetPageNumbers(t *testing.T) {
	tests := []struct {
		currentPage int
		totalPages  int
		want        []int
	}{
		{
			currentPage: 1,
			totalPages:  1,
			want:        []int{1},
		},
		{
			currentPage: 1,
			totalPages:  3,
			want:        []int{1, 2, 3},
		},
		{
			currentPage: 1,
			totalPages:  5,
			want:        []int{1, 2, 3, 5},
		},
		{
			currentPage: 2,
			totalPages:  5,
			want:        []int{1, 2, 3, 4, 5},
		},
		{
			currentPage: 3,
			totalPages:  5,
			want:        []int{1, 2, 3, 4, 5},
		},
		{
			currentPage: 4,
			totalPages:  5,
			want:        []int{1, 2, 3, 4, 5},
		},
		{
			currentPage: 5,
			totalPages:  5,
			want:        []int{1, 3, 4, 5},
		},
		{
			currentPage: 5,
			totalPages:  10,
			want:        []int{1, 3, 4, 5, 6, 7, 10},
		},
		{
			currentPage: 0,
			totalPages:  0,
			want:        []int{1},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			pagination := Pagination{
				CurrentPage: test.currentPage,
				TotalPages:  test.totalPages,
			}
			assert.Equal(t, test.want, pagination.GetPageNumbers())
		})
	}
}

func TestPagination_GetElementsFrom(t *testing.T) {
	tests := []struct {
		currentPage     int
		elementsPerPage int
		totalElements   int
		want            int
	}{
		{
			currentPage:     1,
			elementsPerPage: 25,
			totalElements:   0,
			want:            0,
		},
		{
			currentPage:     1,
			elementsPerPage: 25,
			totalElements:   1,
			want:            1,
		},
		{
			currentPage:     2,
			elementsPerPage: 25,
			totalElements:   50,
			want:            26,
		},
		{
			currentPage:     2,
			elementsPerPage: 50,
			totalElements:   52,
			want:            51,
		},
		{
			currentPage:     0,
			elementsPerPage: 25,
			totalElements:   0,
			want:            0,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			pagination := Pagination{
				CurrentPage:     test.currentPage,
				ElementsPerPage: test.elementsPerPage,
				TotalElements:   test.totalElements,
			}
			assert.Equal(t, test.want, pagination.GetElementsFrom())
		})
	}
}

func TestPagination_GetElementsTo(t *testing.T) {
	tests := []struct {
		currentPage     int
		elementsPerPage int
		totalElements   int
		want            int
	}{
		{
			currentPage:     1,
			elementsPerPage: 25,
			totalElements:   0,
			want:            0,
		},
		{
			currentPage:     1,
			elementsPerPage: 25,
			totalElements:   1,
			want:            1,
		},
		{
			currentPage:     1,
			elementsPerPage: 25,
			totalElements:   50,
			want:            25,
		},
		{
			currentPage:     2,
			elementsPerPage: 50,
			totalElements:   52,
			want:            52,
		},
		{
			currentPage:     2,
			elementsPerPage: 50,
			totalElements:   150,
			want:            100,
		},
		{
			currentPage:     0,
			elementsPerPage: 25,
			totalElements:   0,
			want:            0,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			pagination := Pagination{
				CurrentPage:     test.currentPage,
				ElementsPerPage: test.elementsPerPage,
				TotalElements:   test.totalElements,
			}
			assert.Equal(t, test.want, pagination.GetElementsTo())
		})
	}
}

func TestGetRequestedElementsPerPage(t *testing.T) {
	assert.Equal(t, 2, GetRequestedElementsPerPage("2", []int{1, 2, 3}))
	assert.Equal(t, 25, GetRequestedElementsPerPage("2", []int{25, 50, 100}))
}

func TestGetRequestedPage(t *testing.T) {
	assert.Equal(t, 2, GetRequestedPage("2"))
	assert.Equal(t, 1, GetRequestedPage("0"))
}
