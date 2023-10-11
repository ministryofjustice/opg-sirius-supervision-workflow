package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"reflect"
)

type ListPage struct {
	App            WorkflowVars
	Error          string
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

type FilterByStatus struct {
	ListPage
	StatusOptions    []model.RefData
	SelectedStatuses []string
}

type FilterByDeputyType struct {
	ListPage
	DeputyTypes         []model.RefData
	SelectedDeputyTypes []string
}

type FilterByCaseType struct {
	ListPage
	CaseTypes         []model.RefData
	SelectedCaseTypes []string
}

func (lp ListPage) HasFilterBy(page interface{}, filter string) bool {
	filters := map[string]interface{}{
		"assignee":    FilterByAssignee{},
		"due-date":    FilterByDueDate{},
		"status":      FilterByStatus{},
		"task-type":   FilterByTaskType{},
		"deputy-type": FilterByDeputyType{},
		"case-type":   FilterByCaseType{},
	}

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

	if f, ok := filters[filter]; ok {
		return extends(page, f)
	}
	return false
}

func (ftp FilterByTaskType) ValidateSelectedTaskTypes(selectedTaskTypes []string, taskTypes []model.TaskType) []string {
	var validSelectedTaskTypes []string
	for _, selectedTaskType := range selectedTaskTypes {
		for _, taskType := range taskTypes {
			if selectedTaskType == taskType.Handle {
				validSelectedTaskTypes = append(validSelectedTaskTypes, selectedTaskType)
				break
			}
		}
	}
	return validSelectedTaskTypes
}
