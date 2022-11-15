package sirius

type PageDetails struct {
	ListOfPages        []int
	PreviousPage       int
	NextPage           int
	LimitedPagination  []int
	FirstPage          int
	LastPage           int
	StoredTaskLimit    int
	ShowingUpperLimit  int
	ShowingLowerLimit  int
	LastFilter         string
	UpperEllipsesLimit bool
}

func (c *Client) GetPageDetails(tasklist TaskList, search int, displayTaskLimit int) PageDetails {
	var k PageDetails

	PageDetails := k

	for i := 1; i < tasklist.Pages.PageTotal+1; i++ {
		PageDetails.ListOfPages = append(PageDetails.ListOfPages, i)
	}

	PageDetails.PreviousPage = GetPreviousPageNumber(search)

	PageDetails.NextPage = GetNextPageNumber(tasklist, search)

	PageDetails.StoredTaskLimit = displayTaskLimit

	PageDetails.ShowingUpperLimit = GetShowingUpperLimitNumber(tasklist, displayTaskLimit)

	PageDetails.ShowingLowerLimit = GetShowingLowerLimitNumber(tasklist, displayTaskLimit)

	PageDetails.UpperEllipsesLimit = GetUpperEllipsesLimit(tasklist.Pages, search)

	if len(PageDetails.ListOfPages) != 0 {
		PageDetails.FirstPage = PageDetails.ListOfPages[0]
		PageDetails.LastPage = PageDetails.ListOfPages[len(PageDetails.ListOfPages)-1]
		PageDetails.LimitedPagination = GetPaginationLimits(tasklist, PageDetails)
	} else {
		PageDetails.FirstPage = 0
		PageDetails.LastPage = 0
		PageDetails.LimitedPagination = []int{0}
	}

	return PageDetails
}

func GetPreviousPageNumber(search int) int {
	if search <= 1 {
		return 1
	} else {
		return search - 1
	}
}

func GetNextPageNumber(taskList TaskList, search int) int {
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

func GetShowingLowerLimitNumber(taskList TaskList, displayTaskLimit int) int {
	if taskList.Pages.PageCurrent == 1 && taskList.TotalTasks != 0 {
		return 1
	} else if taskList.Pages.PageCurrent == 1 && taskList.TotalTasks == 0 {
		return 0
	} else {
		previousPageNumber := taskList.Pages.PageCurrent - 1
		return previousPageNumber*displayTaskLimit + 1
	}
}

func GetShowingUpperLimitNumber(taskList TaskList, displayTaskLimit int) int {
	if taskList.Pages.PageCurrent*displayTaskLimit > taskList.TotalTasks {
		return taskList.TotalTasks
	} else {
		return taskList.Pages.PageCurrent * displayTaskLimit
	}
}

func GetPaginationLimits(taskList TaskList, PageDetails PageDetails) []int {
	var twoBeforeCurrentPage int
	var twoAfterCurrentPage int
	if taskList.Pages.PageCurrent > 2 {
		twoBeforeCurrentPage = taskList.Pages.PageCurrent - 3
	} else {
		twoBeforeCurrentPage = 0
	}
	if taskList.Pages.PageCurrent+2 <= PageDetails.LastPage {
		twoAfterCurrentPage = taskList.Pages.PageCurrent + 2
	} else if taskList.Pages.PageCurrent+1 <= PageDetails.LastPage {
		twoAfterCurrentPage = taskList.Pages.PageCurrent + 1
	} else {
		twoAfterCurrentPage = taskList.Pages.PageCurrent
	}
	return PageDetails.ListOfPages[twoBeforeCurrentPage:twoAfterCurrentPage]
}

func GetUpperEllipsesLimit(pages PageInformation, search int) bool {
	return (pages.PageTotal - search) >= 3
}
