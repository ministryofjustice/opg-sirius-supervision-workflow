package sirius

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPreviousPageNumber(t *testing.T) {
	assert.Equal(t, getPreviousPageNumber(0), 1)
	assert.Equal(t, getPreviousPageNumber(1), 1)
	assert.Equal(t, getPreviousPageNumber(2), 1)
	assert.Equal(t, getPreviousPageNumber(3), 2)
	assert.Equal(t, getPreviousPageNumber(5), 4)
}

func setUpGetNextPageNumber(pageCurrent int, pageTotal int, totalTasks int) TaskList {
	taskList := TaskList{
		Pages: PageDetails{
			PageCurrent: pageCurrent,
			PageTotal:   pageTotal,
		},
		TotalTasks: totalTasks,
	}
	return taskList
}

func TestGetNextPageNumber(t *testing.T) {
	taskList := setUpGetNextPageNumber(1, 5, 0)

	assert.Equal(t, getNextPageNumber(taskList, 0), 2)
	assert.Equal(t, getNextPageNumber(taskList, 2), 3)
	assert.Equal(t, getNextPageNumber(taskList, 15), 5)
}

func TestGetShowingLowerLimitNumberAlwaysReturns1IfOnly1Page(t *testing.T) {
	taskList := setUpGetNextPageNumber(1, 0, 13)

	assert.Equal(t, getShowingLowerLimitNumber(taskList, 25), 1)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 50), 1)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 100), 1)
}

func TestGetShowingLowerLimitNumberAlwaysReturns0If0Tasks(t *testing.T) {
	taskList := setUpGetNextPageNumber(1, 0, 0)

	assert.Equal(t, getShowingLowerLimitNumber(taskList, 25), 0)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 50), 0)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 100), 0)
}

func TestGetShowingLowerLimitNumberCanIncrementOnPages(t *testing.T) {
	taskList := setUpGetNextPageNumber(2, 0, 100)

	assert.Equal(t, getShowingLowerLimitNumber(taskList, 25), 26)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 50), 51)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 100), 101)
}

func TestGetShowingLowerLimitNumberCanIncrementOnManyPages(t *testing.T) {
	taskList := setUpGetNextPageNumber(5, 0, 5000)

	assert.Equal(t, getShowingLowerLimitNumber(taskList, 25), 101)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 50), 201)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 100), 401)
}

func TestGetShowingUpperLimitNumberWillReturnTotalTasksIfOnFinalPage(t *testing.T) {
	taskList := setUpGetNextPageNumber(1, 0, 10)

	assert.Equal(t, getShowingUpperLimitNumber(taskList, 25), 10)
	assert.Equal(t, getShowingUpperLimitNumber(taskList, 50), 10)
	assert.Equal(t, getShowingUpperLimitNumber(taskList, 100), 10)
}

func makeListOfPagesRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
