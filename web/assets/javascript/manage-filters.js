export default class ManageFilters {
  constructor(element) {
      
      this.taskTypeButton = element.querySelectorAll('.js-container-button');
      this.innerContainer = element.querySelector(".js-options-container");
      this.taskTypeFilterArrowUp = element.querySelector(".app-c-option-select__icon--up");
      this.taskTypeFilterArrowDown = element.querySelector(".app-c-option-select__icon--down");
      this.taskTypeFilterTags = element.querySelector(".task-type-filter-tags");
      this.inputElementTasktypeFilter = element.querySelectorAll(".task-type");
      this.actionFilter = element.querySelectorAll("#actionFilter");
      console.log(actionFilter)
      this._setupEventListeners();
    }

  _setupEventListeners() {
      this.taskTypeButton.forEach(element => {
          this._toggleTasktypeFilter = this._toggleTasktypeFilter.bind(this);
          element.addEventListener('click', this._toggleTasktypeFilter);
      });      
      
      this.actionFilter.forEach(element => {
          this._selectedTasktypeFilter = this._selectedTasktypeFilter.bind(this);
          element.addEventListener('click', this._selectedTasktypeFilter);
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
    // let append = ""
    // let text = JSON.parse(window.localStorage.getItem("kate"))
    // text.forEach( value => {
    //              append += "<li>" + text.id + "</li>"
    //           })
    // document.getElementById("ul-filter-tags").innerHTML = append;
  }

  _selectedTasktypeFilter() {
    let count = 0;
    let array = [];
    this.inputElementTasktypeFilter.forEach(taskType => {
      if (taskType.checked) {
        array.push(JSON.stringify(taskType.id))
      }
      window.localStorage.setItem("kate", JSON.stringify(array))
    })
    console.log(window.localStorage.getItem("kate"))

    let append = ""
    let text = window.localStorage.getItem("kate");
    console.log("text")
    console.log(text)
    let parsedText = JSON.parse(text)
    setTimeout(function(){ alert("Hello"); }, 20000);
    parsedText.forEach( value => {
                 append += "<li>" + value.id + "</li><br>"
              })
    console.log("parsedText")
    console.log(parsedText)
    document.getElementById("replaceme").innerHTML = append;
    console.log("should be appended")
    console.log(document.getElementById("replaceme"))

    
  }
 
  // let str = "<option value=''selected>Select a case manager</option>"
  //         data.members.forEach( caseManager => {
  //            str += "<option value=" + caseManager.id + ">" + caseManager.displayName + "</option>"
  //         })
      
  //         document.getElementById("assignCM").innerHTML = str;

  // <li><a class="moj-filter__tag" href="#"><span class="govuk-visually-hidden">Remove this filter</span>Report: lodge report</a></li>
}