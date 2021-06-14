package sirius

type TaskDetails struct {
	ListOfPages       []int
	PreviousPage      int
	NextPage          int
	LimitedPagination []int
	FirstPage         int
	LastPage          int
	StoredTaskLimit   int
	ShowingUpperLimit int
	ShowingLowerLimit int
	LastFilter        string
	TaskTypeFilters   int
}

func (c *Client) GetTaskDetails(ctx Context, taskList TaskList, search int, displayTaskLimit int, selectedTeamMembers int, loggedInTeamId int, taskTypeSelected []string, LoadTasks []ApiTaskTypes) TaskDetails {
	var k TaskDetails

	TaskDetails := k

	for i := 1; i < taskList.Pages.PageTotal+1; i++ {
		TaskDetails.ListOfPages = append(TaskDetails.ListOfPages, i)
	}

	TaskDetails.PreviousPage = getPreviousPageNumber(search)

	TaskDetails.NextPage = getNextPageNumber(taskList, search)

	TaskDetails.StoredTaskLimit = displayTaskLimit

	TaskDetails.ShowingUpperLimit = getShowingUpperLimitNumber(taskList, displayTaskLimit)

	TaskDetails.ShowingLowerLimit = getShowingLowerLimitNumber(taskList, displayTaskLimit)

	if len(TaskDetails.ListOfPages) != 0 {
		TaskDetails.FirstPage = TaskDetails.ListOfPages[0]
		TaskDetails.LastPage = TaskDetails.ListOfPages[len(TaskDetails.ListOfPages)-1]
		TaskDetails.LimitedPagination = getPaginationLimits(taskList, TaskDetails)
	} else {
		TaskDetails.FirstPage = 0
		TaskDetails.LastPage = 0
		TaskDetails.LimitedPagination = []int{0}
	}

	return TaskDetails
}

func getPreviousPageNumber(search int) int {
	if search <= 1 {
		return 1
	} else {
		return search - 1
	}
}

func getNextPageNumber(taskList TaskList, search int) int {
	if search < taskList.Pages.PageTotal {
		if search == 0 {
			return search + 2
		} else {
			return search + 1
		}
	} else {
		return taskList.Pages.PageTotal
	}
}

func getShowingLowerLimitNumber(taskList TaskList, displayTaskLimit int) int {
	if taskList.Pages.PageCurrent == 1 && taskList.TotalTasks != 0 {
		return 1
	} else if taskList.Pages.PageCurrent == 1 && taskList.TotalTasks == 0 {
		return 0
	} else {
		previousPageNumber := taskList.Pages.PageCurrent - 1
		return previousPageNumber*displayTaskLimit + 1
	}
}

func getShowingUpperLimitNumber(taskList TaskList, displayTaskLimit int) int {
	if taskList.Pages.PageCurrent*displayTaskLimit > taskList.TotalTasks {
		return taskList.TotalTasks
	} else {
		return taskList.Pages.PageCurrent * displayTaskLimit
	}
}

func getPaginationLimits(taskList TaskList, TaskDetails TaskDetails) []int {
	var twoBeforeCurrentPage int
	var twoAfterCurrentPage int
	if taskList.Pages.PageCurrent > 2 {
		twoBeforeCurrentPage = taskList.Pages.PageCurrent - 3
	} else {
		twoBeforeCurrentPage = 0
	}
	if taskList.Pages.PageCurrent+2 <= TaskDetails.LastPage {
		twoAfterCurrentPage = taskList.Pages.PageCurrent + 2
	} else if taskList.Pages.PageCurrent+1 <= TaskDetails.LastPage {
		twoAfterCurrentPage = taskList.Pages.PageCurrent + 1
	} else {
		twoAfterCurrentPage = taskList.Pages.PageCurrent
	}
	return TaskDetails.ListOfPages[twoBeforeCurrentPage:twoAfterCurrentPage]
}
