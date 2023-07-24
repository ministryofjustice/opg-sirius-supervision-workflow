package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"reflect"
)

type ListPage struct {
	App            WorkflowVars
	AppliedFilters []string
	Pagination     paginate.Pagination
	PerPage        int
	UrlBuilder     urlbuilder.UrlBuilder
}

type FilterByAssignee struct {
	ListPage
	AssigneeFilterName string
	SelectedAssignees  []string
	SelectedUnassigned string
}

type FilterByTaskType struct {
	ListPage
	TaskTypes         []model.TaskType
	SelectedTaskTypes []string
}

type FilterByDueDate struct {
	ListPage
	SelectedDueDateFrom string
	SelectedDueDateTo   string
}

func (lp ListPage) HasFilterBy(page interface{}, filter string) bool {
	extends := func(parent interface{}, child interface{}) bool {
		p := reflect.TypeOf(parent)
		c := reflect.TypeOf(child)
		for i := 0; i < p.NumField(); i++ {
			if f := p.Field(i); f.Type == c && f.Anonymous {
				return true
			}
		}
		return false
	}

	switch filter {
	case "assignee":
		return extends(page, FilterByAssignee{})
	case "due-date":
		return extends(page, FilterByDueDate{})
	case "task-type":
		return extends(page, FilterByTaskType{})
	}
	return false
}
