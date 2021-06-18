export default class ManageFilters {
  constructor(element) {

      this.taskTypeButton = element.querySelectorAll('.tasktype-button');
      this.assigneeButton = element.querySelectorAll('.assignee-button');
      console.log(this.assigneeButton)
      this.tasktypeInnerContainer = element.querySelector(".tasktype-inner-container");
      this.assigneeInnerContainer = element.querySelector(".assigned-inner-container");
      this.taskTypeFilterArrowUp = element.querySelector(".tasktype-arrow-up");
      this.taskTypeFilterArrowDown = element.querySelector(".tasktype-arrow-down");
      this.assigneeFilterArrowUp = element.querySelector(".assigned-arrow-up");
      this.assigneeFilterArrowDown = element.querySelector(".assigned-arrow-down");
      this.inputElementTasktypeFilter = element.querySelectorAll(".task-type");
      this.inputElementAssigneeFilter = element.querySelectorAll(".assignee-type");
      this.urlForPage = window.location.search;
      this.teamSelection = element.querySelectorAll('.change-teams');
      this._clearFilters();
      this._isFiltered();
      this._isFilteredByAssignee();
      this._setupEventListeners();
  }

  _setupEventListeners() {
      this.taskTypeButton.forEach(element => {
          this._toggleTasktypeFilter = this._toggleTasktypeFilter.bind(this);
          element.addEventListener('click', this._toggleTasktypeFilter);
      });

      this._retainTaskFilterMenuStateWhenReloadingPage()
      this._retainAssigneeMenuStateWhenReloadingPage()
  }

  _toggleTasktypeFilter() {
      const hiddenStateTasktype = this.tasktypeInnerContainer.classList.contains('hide');
      this.tasktypeInnerContainer.classList.toggle('hide', !hiddenStateTasktype)
      if (hiddenStateTasktype) {
        this.taskTypeFilterArrowUp.setAttribute('aria-expanded', 'true')
        this.taskTypeFilterArrowDown.setAttribute('aria-expanded', 'false')

        this.taskTypeFilterArrowUp.classList.toggle('hide', false);
        this.taskTypeFilterArrowDown.classList.toggle('hide', true)

        window.localStorage.setItem("TasktypeOpen", "true")
        console.log("set TasktypeOpen to true")
        console.log(window.localStorage.getItem("TasktypeOpen"))
    } else {
        this.taskTypeFilterArrowUp.setAttribute('aria-expanded', 'false')
        this.taskTypeFilterArrowDown.setAttribute('aria-expanded', 'true')
        this.taskTypeFilterArrowUp.classList.toggle('hide', true)
        this.taskTypeFilterArrowDown.classList.toggle('hide', false)

        window.localStorage.setItem("TasktypeOpen", "false")
        console.log("set TasktypeOpen to false")
        console.log(window.localStorage.getItem("TasktypeOpen"))
    }
}


_toggleAssigneeFilter() {
  console.log("in filter func")
  const hiddenStateAssignee = this.assigneeInnerContainer.classList.contains('hide');
  this.assigneeInnerContainer.classList.toggle('hide', !hiddenStateAssignee)
  if (hiddenStateAssignee) {
    this.assigneeFilterArrowUp.setAttribute('aria-expanded', 'true')
    this.assigneeFilterArrowDown.setAttribute('aria-expanded', 'false')

    this.assigneeFilterArrowUp.classList.toggle('hide', false);
    this.assigneeFilterArrowDown.classList.toggle('hide', true)

    window.localStorage.setItem("AssigneeOpen", "true")
    console.log("set AssigneeOpen to true")
    console.log(window.localStorage.getItem("AssigneeOpen"))
} else {
    this.assigneeFilterArrowUp.setAttribute('aria-expanded', 'false')
    this.assigneeFilterArrowDown.setAttribute('aria-expanded', 'true')
    this.assigneeFilterArrowUp.classList.toggle('hide', true)
    this.assigneeFilterArrowDown.classList.toggle('hide', false)

    window.localStorage.setItem("AssigneeOpen", "false")
    console.log("set AssigneeOpen to false")
    console.log(window.localStorage.getItem("AssigneeOpen"))
}
}

  _retainTaskFilterMenuStateWhenReloadingPage() {
    if (window.localStorage.getItem("TasktypeOpen") == "true") {
        this.taskTypeFilterArrowUp.setAttribute('aria-expanded', 'true')
        this.taskTypeFilterArrowDown.setAttribute('aria-expanded', 'false')
        this.taskTypeFilterArrowUp.classList.toggle('hide', false)
        this.taskTypeFilterArrowDown.classList.toggle('hide', true)

        const hiddenStateTasktype = this.tasktypeInnerContainer.classList.contains('hide');
        this.tasktypeInnerContainer.classList.toggle('hide', !hiddenStateTasktype)
        console.log("retained task type open is true state")
      
    } else {
        this.taskTypeFilterArrowUp.setAttribute('aria-expanded', 'false')
        this.taskTypeFilterArrowDown.setAttribute('aria-expanded', 'true')
        this.taskTypeFilterArrowUp.classList.toggle('hide', true)
        this.taskTypeFilterArrowDown.classList.toggle('hide', false)
        console.log("retained task type open is false state")
    }
  }

  _retainAssigneeMenuStateWhenReloadingPage() {
    if (window.localStorage.getItem("AssigneeOpen") == "true") {
      this.assigneeFilterArrowUp.setAttribute('aria-expanded', 'true')
      this.assigneeFilterArrowUp.setAttribute('aria-expanded', 'false')
      this.assigneeFilterArrowUp.classList.toggle('hide', false)
      this.assigneeFilterArrowUp.classList.toggle('hide', true)

      const hiddenStateAssignee = this.assigneeInnerContainer.classList.contains('hide');
      this.assigneeInnerContainer.classList.toggle('hide', !hiddenStateAssignee)
      console.log("retained assignee open is true state")
    } else {
      this.assigneeFilterArrowUp.setAttribute('aria-expanded', 'false')
      this.assigneeFilterArrowDown.setAttribute('aria-expanded', 'true')
      this.assigneeFilterArrowUp.classList.toggle('hide', true)
      this.assigneeFilterArrowDown.classList.toggle('hide', false)
      console.log("retained assignee open is false state")
    }
  }

  _isFiltered() {
      let array = [];

      this.inputElementTasktypeFilter.forEach(taskType => {
          if (taskType.checked) {
              array.push(taskType);
          }
      })

      let append = "";
    
        if (array.length) {
          append += '<h3 class="govuk-heading-s govuk-!-margin-bottom-0" id="Task-type-hook">Task type</h3>';
          array.forEach(taskType => {
              let hrefValue = this.urlForPage.split("&").filter((param) => !param.includes(taskType.value)).join("&");
              append += `<li id=${taskType.value}><a class="moj-filter__tag" href=${hrefValue}><span class="govuk-visually-hidden">Remove this filter</span>` + taskType.id + "</li>"
          }) 
        }

        document.getElementById("applied-task-type-filters").innerHTML = append
  }

  _isFilteredByAssignee() {
    let array = [];

    this.inputElementAssigneeFilter.forEach(assignee => {
      if (assignee.checked) {
        array.push(assignee);
      }
    })

    let append = "";
    
    if (array.length) {
      append += '<h3 class="govuk-heading-s govuk-!-margin-bottom-0" id="Assignee-hook">Assignees</h3>';
      array.forEach(assignee => {
          let hrefValue = this.urlForPage.split("&").filter((param) => !param.includes('selected-assignee='+assignee.value)).join("&");
          append += `<li id=${assignee.value}><a class="moj-filter__tag" href=${hrefValue}><span class="govuk-visually-hidden">Remove this filter</span>` + assignee.id + "</li>"
        })
    }

      document.getElementById("applied-assignee-filters").innerHTML = append
  }

  _clearFilters() {
    let hrefValueWithoutSelectedTask = this.urlForPage.split("&").filter((param) => !param.includes("selected-task-type")).join("&");
    let hrefValue = hrefValueWithoutSelectedTask.split("&").filter((param) => !param.includes("selected-assignee")).join("&");
    document.getElementById("clear-filters").setAttribute('href', hrefValue);
  }

}
