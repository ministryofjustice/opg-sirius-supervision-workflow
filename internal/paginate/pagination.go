package paginate

import "sort"

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
	GetPaginationUrl(int, ...int)
}

func (p Pagination) ShowPrevious() bool {
	return p.CurrentPage > 1
}

func (p Pagination) ShowNext() bool {
	return p.CurrentPage < p.TotalPages
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
