export default class ManageFilters {
  constructor(element) {

      this.taskTypeButton = element.querySelectorAll('.js-container-button');
      this.innerContainer = element.querySelector(".js-options-container");
      this.taskTypeFilterArrowUp = element.querySelector(".app-c-option-select__icon--up");
      this.taskTypeFilterArrowDown = element.querySelector(".app-c-option-select__icon--down");
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
  }

  _toggleTasktypeFilter() {
      const hiddenState = this.innerContainer.classList.contains('hide');
      this.innerContainer.classList.toggle('hide', !hiddenState)
      if (hiddenState) {
          this.taskTypeFilterArrowUp.setAttribute('aria-expanded', 'true')
          this.taskTypeFilterArrowDown.setAttribute('aria-expanded', 'false')

          this.taskTypeFilterArrowUp.classList.toggle('hide', false);
          this.taskTypeFilterArrowDown.classList.toggle('hide', true)

          window.localStorage.setItem("Open", "true")
      } else {
          this.taskTypeFilterArrowUp.setAttribute('aria-expanded', 'false')
          this.taskTypeFilterArrowDown.setAttribute('aria-expanded', 'true')
          this.taskTypeFilterArrowUp.classList.toggle('hide', true)
          this.taskTypeFilterArrowDown.classList.toggle('hide', false)

          window.localStorage.setItem("Open", "false")
      }
  }

  _retainTaskFilterMenuStateWhenReloadingPage() {
      if (window.localStorage.getItem("Open") == "true") {
          this.taskTypeFilterArrowUp.setAttribute('aria-expanded', 'true')
          this.taskTypeFilterArrowDown.setAttribute('aria-expanded', 'false')
          this.taskTypeFilterArrowUp.classList.toggle('hide', false)
          this.taskTypeFilterArrowDown.classList.toggle('hide', true)

          const hiddenState = this.innerContainer.classList.contains('hide');
          this.innerContainer.classList.toggle('hide', !hiddenState)
      } else {
          this.taskTypeFilterArrowUp.setAttribute('aria-expanded', 'false')
          this.taskTypeFilterArrowDown.setAttribute('aria-expanded', 'true')
          this.taskTypeFilterArrowUp.classList.toggle('hide', true)
          this.taskTypeFilterArrowDown.classList.toggle('hide', false)

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
          let hrefValue = this.urlForPage.split("&").filter((param) => !param.includes(assignee.value)).join("&");
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
