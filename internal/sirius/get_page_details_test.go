package sirius

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateTaskList(pageCurrent int, pageTotal int, totalTasks int) TaskList {
	taskList := TaskList{
		Pages: PageInformation{
			PageCurrent: pageCurrent,
			PageTotal:   pageTotal,
		},
		TotalTasks: totalTasks,
	}
	return taskList
}

func CreatePageDetails(
	ListOfPages []int,
	PreviousPage int,
	NextPage int,
	LimitedPagination []int,
	FirstPage int,
	LastPage int,
	StoredTaskLimit int,
	ShowingUpperLimit int,
	ShowingLowerLimit int,
) PageDetails {
	newPageDetails := PageDetails{
		ListOfPages,
		PreviousPage,
		NextPage,
		LimitedPagination,
		FirstPage,
		LastPage,
		StoredTaskLimit,
		ShowingUpperLimit,
		ShowingLowerLimit,
		"",
		true,
	}
	return newPageDetails
}

func TestGetPageDetails(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1, 2, 3, 4, 5}, 2, 4, []int{1, 2, 3, 4, 5}, 1, 5, 25, 75, 51)
	taskList := CreateTaskList(3, 5, 125)
	result := client.GetPageDetails(taskList, 3, 25)

	assert.Equal(t, expectedResult, result)
}

func TestGetPageDetailsPage1View25(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1}, 1, 1, []int{1}, 1, 1, 25, 10, 1)
	taskList := CreateTaskList(1, 1, 10)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 1, 25))
}

func TestGetPageDetailsPage1View50(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1}, 1, 1, []int{1}, 1, 1, 50, 10, 1)
	taskList := CreateTaskList(1, 1, 10)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 1, 50))
}

func TestGetPageDetailsPage1View100(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1}, 1, 1, []int{1}, 1, 1, 100, 99, 1)
	taskList := CreateTaskList(1, 1, 99)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 1, 100))
}

func TestGetPageDetailsPage2of2(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1, 2}, 2, 2, []int{1, 2}, 1, 2, 25, 27, 26)
	taskList := CreateTaskList(2, 2, 27)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 3, 25))
}

func TestGetPageDetailsPage2of3(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1, 2, 3}, 1, 3, []int{1, 2, 3}, 1, 3, 25, 50, 26)
	taskList := CreateTaskList(2, 3, 74)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 2, 25))
}

func TestGetPageDetailsPage4of10(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3, 5, []int{2, 3, 4, 5, 6}, 1, 10, 25, 100, 76)
	taskList := CreateTaskList(4, 10, 250)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 4, 25))
}

func TestGetPageDetailsPage4of5(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1, 2, 3, 4, 5}, 3, 5, []int{2, 3, 4, 5}, 1, 5, 50, 200, 151)
	taskList := CreateTaskList(4, 5, 1000)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 4, 50))
}

func TestGetPageDetailsPage5of5(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1, 2, 3, 4, 5}, 4, 5, []int{3, 4, 5}, 1, 5, 50, 250, 201)
	taskList := CreateTaskList(5, 5, 1000)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 5, 50))
}

func TestGetPageDetailsPage7of10(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 6, 8, []int{5, 6, 7, 8, 9}, 1, 10, 50, 350, 301)
	taskList := CreateTaskList(7, 10, 500)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 7, 50))
}

func TestGetPageDetailsPage9of10(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 8, 10, []int{7, 8, 9, 10}, 1, 10, 100, 900, 801)
	taskList := CreateTaskList(9, 10, 901)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 9, 100))
}

func TestGetPageDetailsFinalPage(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	expectedResult := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 9, 10, []int{8, 9, 10}, 1, 10, 100, 1000, 901)
	taskList := CreateTaskList(10, 10, 1000)

	assert.Equal(t, expectedResult, client.GetPageDetails(taskList, 10, 100))
}

func TestGetPreviousPageNumber(t *testing.T) {
	assert.Equal(t, GetPreviousPageNumber(0), 1)
	assert.Equal(t, GetPreviousPageNumber(1), 1)
	assert.Equal(t, GetPreviousPageNumber(2), 1)
	assert.Equal(t, GetPreviousPageNumber(3), 2)
	assert.Equal(t, GetPreviousPageNumber(5), 4)
}

func TestGetNextPageNumber(t *testing.T) {
	taskList := CreateTaskList(1, 5, 0)

	assert.Equal(t, GetNextPageNumber(taskList, 0), 2)
	assert.Equal(t, GetNextPageNumber(taskList, 2), 3)
	assert.Equal(t, GetNextPageNumber(taskList, 15), 5)
}

func TestGetShowingLowerLimitNumberAlwaysReturns1IfOnly1Page(t *testing.T) {
	taskList := CreateTaskList(1, 0, 13)

	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 25), 1)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 50), 1)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 100), 1)
}

func TestGetShowingLowerLimitNumberAlwaysReturns0If0Tasks(t *testing.T) {
	taskList := CreateTaskList(1, 0, 0)

	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 25), 0)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 50), 0)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 100), 0)
}

func TestGetShowingLowerLimitNumberCanIncrementOnPages(t *testing.T) {
	taskList := CreateTaskList(2, 0, 100)

	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 25), 26)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 50), 51)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 100), 101)
}

func TestGetShowingLowerLimitNumberCanIncrementOnManyPages(t *testing.T) {
	taskList := CreateTaskList(5, 0, 5000)

	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 25), 101)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 50), 201)
	assert.Equal(t, GetShowingLowerLimitNumber(taskList, 100), 401)
}

func TestGetShowingUpperLimitNumber(t *testing.T) {
	taskList := CreateTaskList(1, 0, 25)
	assert.Equal(t, GetShowingUpperLimitNumber(taskList, 25), 25)

	taskList = CreateTaskList(1, 0, 50)
	assert.Equal(t, GetShowingUpperLimitNumber(taskList, 50), 50)

	taskList = CreateTaskList(1, 0, 100)
	assert.Equal(t, GetShowingUpperLimitNumber(taskList, 100), 100)
}

func TestGetShowingUpperLimitNumberWillReturnTotalTasksIfOnFinalPage(t *testing.T) {
	taskList := CreateTaskList(1, 0, 10)

	assert.Equal(t, GetShowingUpperLimitNumber(taskList, 25), 10)
	assert.Equal(t, GetShowingUpperLimitNumber(taskList, 50), 10)
	assert.Equal(t, GetShowingUpperLimitNumber(taskList, 100), 10)
}

func TestGetPaginationLimitsPage1of10(t *testing.T) {
	taskList := CreateTaskList(1, 10, 0)
	pageDetails := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1, 2, []int{}, 1, 10, 0, 0, 0)

	assert.Equal(t, []int{1, 2, 3}, GetPaginationLimits(taskList, pageDetails))
}

func TestGetPaginationLimitsPage2of10(t *testing.T) {
	taskList := CreateTaskList(2, 10, 0)
	pageDetails := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1, 3, []int{}, 1, 10, 0, 0, 0)

	assert.Equal(t, []int{1, 2, 3, 4}, GetPaginationLimits(taskList, pageDetails))
}

func TestGetPaginationLimitsPage3of10(t *testing.T) {
	taskList := CreateTaskList(3, 10, 0)
	pageDetails := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 2, 4, []int{}, 1, 10, 0, 0, 0)

	assert.Equal(t, []int{1, 2, 3, 4, 5}, GetPaginationLimits(taskList, pageDetails))
}

func TestGetPaginationLimitsPage4of10(t *testing.T) {
	taskList := CreateTaskList(4, 10, 100)
	pageDetails := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3, 5, []int{}, 1, 10, 0, 0, 0)

	assert.Equal(t, []int{2, 3, 4, 5, 6}, GetPaginationLimits(taskList, pageDetails))
}

func TestGetPaginationLimitsPage8of10(t *testing.T) {
	taskList := CreateTaskList(8, 10, 100)
	pageDetails := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 7, 9, []int{}, 1, 10, 0, 0, 0)

	assert.Equal(t, []int{6, 7, 8, 9, 10}, GetPaginationLimits(taskList, pageDetails))
}

func TestGetPaginationLimitsPage9of10(t *testing.T) {
	taskList := CreateTaskList(9, 10, 100)
	pageDetails := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 8, 10, []int{}, 1, 10, 0, 0, 0)

	assert.Equal(t, []int{7, 8, 9, 10}, GetPaginationLimits(taskList, pageDetails))
}

func TestGetPaginationLimitsPage10of10(t *testing.T) {
	taskList := CreateTaskList(10, 10, 100)
	pageDetails := CreatePageDetails([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 9, 10, []int{}, 1, 10, 0, 0, 0)

	assert.Equal(t, []int{8, 9, 10}, GetPaginationLimits(taskList, pageDetails))
}

func MakeListOfPagesRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
