export default class ManageFilters {
  constructor(element) {
      
      this.taskTypeButton = element.querySelectorAll('.js-container-button');
      this.innerContainer = element.querySelector(".js-options-container");
      this.taskTypeFilterArrowUp = element.querySelector(".app-c-option-select__icon--up");
      this.taskTypeFilterArrowDown = element.querySelector(".app-c-option-select__icon--down");
      this.taskTypeFilterTags = element.querySelector(".task-type-filter-tags");
      this.inputElementTasktypeFilter = element.querySelectorAll(".task-type");
      this.taskTypeFilterTags = element.querySelectorAll(".task-type-filter-tags");
      this._isFiltered();
      this._setupEventListeners();
    }

  _setupEventListeners() {
      this.taskTypeButton.forEach(element => {
          this._toggleTasktypeFilter = this._toggleTasktypeFilter.bind(this);
          element.addEventListener('click', this._toggleTasktypeFilter);
      }); 
      
      this.taskTypeFilterTags.forEach(element => {
        this._toggleTasktypeFilter = this._selectedTaskTypes.bind(this);
        element.addEventListener('click', this._selectedTaskTypes);
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
        array.push(taskType.id);
      }
    })

    let append = "";
    array.forEach(value => {
      let id = value.split(" ").join("");
                 append += `<li id=${id}><a class="moj-filter__tag" href="#"><span class="govuk-visually-hidden">Remove this filter</span>` + value + "</li>"
              })
    document.getElementById("replaceme").innerHTML = append;
  }
}