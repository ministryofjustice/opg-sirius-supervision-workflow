
export default class ManageTasks {
  constructor(element) {
      this.data = {
          selectedTasks: 0
      }
      this.teamMemberData = [];

      this.xsrfToken = element.querySelector('.js-xsrfToken');
      this.checkBoxElements = element.querySelectorAll('.js-mt-checkbox');
      this.manageTasksButton = element.querySelectorAll('.js-mt-edit-tasks-btn');
      this.allcheckBoxElements = element.querySelectorAll('.js-mt-checkbox-select-all');
      this.selectedCountElement = element.querySelectorAll('.js-mt-task-count');
      this.editPanelDiv = element.querySelectorAll('.js-mt-edit-panel');
      this.baseUrl = document.querySelector('[name=api-base-uri]').getAttribute('content')
      this.taskTypeCheckBox

      this._setupEventListeners();
    }

  _setupEventListeners() {
    document.addEventListener('click', (e) => {
        if (e.target) {
             if (e.target.classList.length > 0) {
                const hookName = Array.from(e.target.classList).filter(f => f.indexOf('js-') === 0)[0];
                switch (hookName) {
                    case "js-mt-checkbox":
                        this._updateSelectedState(e);
                        break;
                    case "js-mt-checkbox-select-all":
                        this._updateAllSelectedState(e);
                        break;
                    case "js-mt-edit-tasks-btn":
                        this._showEditTasksPanel(e);
                        break; 
                    case "js-mt-cancel":
                        this._hideEditTasksPanel(e)
                        break;
                    case "js-assign-team-select":
                        this._getCaseManagers(e);
                        break;
                    case "js-container-button":
                        this._toggleTasktypeFilter(e);
                        break;
                    default:
                        break;
                }  
            }
        }
    });
  }

  _updateDomElements() {
      this.selectedCountElement.forEach(element => {
          element.innerText = this.data.selectedTasks.toString();
      });
      this.manageTasksButton[0].classList.toggle('hide', this.data.selectedTasks === 0);
  }

  _updateSelectedRowStyles(element) {
      element.parentElement.parentElement.parentElement.classList.toggle('govuk-table__select', element.checked);
      element.parentElement.parentElement.parentElement.parentElement.classList.toggle('selected', element.checked);
  }

  _updateSelectedState(event) {
      event.target.checked ? this.data.selectedTasks++ : this.data.selectedTasks--;
      this._updateSelectedRowStyles(event.target);
      this._updateDomElements();
  }

  _updateAllSelectedState(event) {
      let isChecked = event.target.checked;

      this.checkBoxElements.forEach(checkbox => {
          checkbox.checked = isChecked;

          this._updateSelectedRowStyles(checkbox);
      });

      this.data.selectedTasks = (isChecked ? this.checkBoxElements.length : 0);
      this._updateDomElements();
  }

  _showEditTasksPanel(event) {
      this.editPanelDiv.forEach(element => {
          element.classList.toggle('hide', this.data.selectedTasks === 0);
      });
    }

  _hideEditTasksPanel(event) {
      this.editPanelDiv.forEach(element => {
          element.classList.toggle('hide', true);
      });
  }

  _getCaseManagers(event) {
      const value = event.target.value.toString();

      fetch(`${this.baseUrl}/api/v1/teams/${value}`, {
          method: "GET",
          credentials: 'include',
          headers: {
              "Content-type": "application/json",
              "X-XSRF-TOKEN": this.xsrfToken.value.toString(),
              "OPG-Bypass-Membrane": 1,
          }
      })
      .then((response) => {
          return response.json();
      })
      .then((data) => {
          let str = "<option value=''selected>Select a case manager</option>"
          let sortedAlphbetically = data.members.sort(function(a, b){
            if(a.name < b.name) { return -1; }
            if(a.name > b.name) { return 1; }
            return 0;
        });

        sortedAlphbetically.forEach( caseManager => {
             str += "<option value=" + caseManager.id + ">" + caseManager.displayName + "</option>"
          });   
      
          document.getElementById("assignCM").innerHTML = str;
      });
  }

  _toggleTasktypeFilter(event) {
      const innerContainer = event.target.parentElement.parentElement.parentElement.querySelector(".js-options-container");
      innerContainer.classList.toggle("hide")
  } 
}
