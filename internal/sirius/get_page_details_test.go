package sirius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPreviousPageNumber(t *testing.T) {
	assert.Equal(t, GetPreviousPageNumber(0), 1)
	assert.Equal(t, GetPreviousPageNumber(1), 1)
	assert.Equal(t, GetPreviousPageNumber(2), 1)
	assert.Equal(t, GetPreviousPageNumber(3), 2)
	assert.Equal(t, GetPreviousPageNumber(5), 4)
}

func SetUpGetNextPageNumber(pageCurrent int, pageTotal int, totalTasks int) TaskList {
	taskList := TaskList{
		Pages: PageInformation{
			PageCurrent: pageCurrent,
			PageTotal:   pageTotal,
		},
		TotalTasks: totalTasks,
	}
	return taskList
}

func TestGetNextPageNumber(t *testing.T) {
	taskList := SetUpGetNextPageNumber(1, 5, 0)

	assert.Equal(t, GetNextPageNumber(taskList, 0), 2)
	assert.Equal(t, GetNextPageNumber(taskList, 2), 3)
	assert.Equal(t, GetNextPageNumber(taskList, 15), 5)
}

func TestGetShowingLowerLimitNumberAlwaysReturns1IfOnly1Page(t *testing.T) {
	taskList := SetUpGetNextPageNumber(1, 0, 13)

	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 25), 1)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 50), 1)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 100), 1)
}

func TestGetShowingLowerLimitNumberAlwaysReturns0If0Tasks(t *testing.T) {
	taskList := SetUpGetNextPageNumber(1, 0, 0)

	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 25), 0)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 50), 0)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 100), 0)
}

func TestGetShowingLowerLimitNumberCanIncrementOnPages(t *testing.T) {
	taskList := SetUpGetNextPageNumber(2, 0, 100)

	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 25), 26)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 50), 51)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 100), 101)
}

func TestGetShowingLowerLimitNumberCanIncrementOnManyPages(t *testing.T) {
	taskList := SetUpGetNextPageNumber(5, 0, 5000)

	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 25), 101)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 50), 201)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 100), 401)
}

func TestGetShowingUpperLimitNumberWillReturnTotalTasksIfOnFinalPage(t *testing.T) {
	taskList := SetUpGetNextPageNumber(1, 0, 10)

	assert.Equal(t, GetShowingUpperLimitNumber(taskList, 25), 10)
	assert.Equal(t, GetShowingUpperLimitNumber(taskList, 50), 10)
	assert.Equal(t, GetShowingUpperLimitNumber(taskList, 100), 10)
}

func MakeListOfPagesRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
