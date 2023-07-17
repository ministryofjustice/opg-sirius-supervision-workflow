package paginate

import (
	"sort"
	"strconv"
)

type Pagination struct {
	CurrentPage     int
	TotalPages      int
	TotalElements   int
	ElementsPerPage int
	ElementName     string
	PerPageOptions  []int
	UrlBuilder      UrlBuilder
}

type UrlBuilder interface {
	GetPaginationUrl(int, ...int) string
}

func (p Pagination) ShowPrevious() bool {
	return p.CurrentPage > 1
}

func (p Pagination) ShowNext() bool {
	return p.CurrentPage < p.TotalPages
}

func (p Pagination) GetPreviousUrl() string {
	page := p.CurrentPage - 1
	if page < 1 {
		page = 1
	}
	return p.UrlBuilder.GetPaginationUrl(page)
}

func (p Pagination) GetNextUrl() string {
	page := p.CurrentPage + 1
	if page > p.TotalPages {
		page = p.TotalPages
	}
	return p.UrlBuilder.GetPaginationUrl(page)
}

func (p Pagination) GetPageNumbers() []int {
	// always show page 1
	pageNumbers := map[int]int{1: 1}
	// show 2 pages either side of the current page
	for i := p.CurrentPage - 2; i <= p.CurrentPage+2; i++ {
		if i > 0 && i <= p.TotalPages {
			pageNumbers[i] = i
		}
	}
	// show the last page
	if p.TotalPages > 0 {
		pageNumbers[p.TotalPages] = p.TotalPages
	}
	var pages []int
	for _, pn := range pageNumbers {
		pages = append(pages, pn)
	}
	sort.Ints(pages)
	return pages
}

func (p Pagination) ShowEllipsisBetween(page1 int, page2 int) bool {
	return page2-page1 > 1
}

func (p Pagination) GetElementsFrom() int {
	if p.TotalElements == 0 {
		return 0
	}
	return (p.CurrentPage-1)*p.ElementsPerPage + 1
}

func (p Pagination) GetElementsTo() int {
	elementsTo := p.CurrentPage * p.ElementsPerPage
	if elementsTo > p.TotalElements {
		return p.TotalElements
	}
	return elementsTo
}

func GetRequestedElementsPerPage(valueFromUrl string, perPageOptions []int) int {
	elementsPerPage, _ := strconv.Atoi(valueFromUrl)
	for _, opt := range perPageOptions {
		if opt == elementsPerPage {
			return elementsPerPage
		}
	}
	return perPageOptions[0]
}

func GetRequestedPage(valueFromUrl string) int {
	page, _ := strconv.Atoi(valueFromUrl)
	if page < 1 {
		return 1
	}
	return page
}
