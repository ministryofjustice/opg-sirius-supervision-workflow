export default class ManageTasks {
    constructor(element) {
        this.data = {
            selectedTasks: 0
        }
        this.teamMemberData = [];

        this.checkBoxElements = element.querySelectorAll('.js-mt-checkbox');
        this.allcheckBoxElements = element.querySelectorAll('.js-mt-checkbox-select-all');
        this.manageTasksButton = element.querySelectorAll('.js-mt-edit-tasks-btn');
        this.cancelEditTasksButton = element.querySelectorAll('.js-mt-cancel');
        this.kate = element.querySelectorAll('.manage-tasks_kate');
        // this.cmselect = element.queryselectorAll('.js-assign-cm-select');
        this.nick = element.querySelectorAll('.option-value');
        this.nickSelect = element.querySelectorAll('.js-assign-team-select');
        this.xsrfToken = element.querySelector('.js-xsrfToken');
        this.selectedCountElement = element.querySelectorAll('.js-mt-task-count');
        this.editPanelDiv = element.querySelectorAll('.js-mt-edit-panel');
        // this._bindKatesFunction(this.nick);
        this._setupEventListeners();
      }

    _setupEventListeners() {
        this.checkBoxElements.forEach(element => {
            this._updateSelectedState = this._updateSelectedState.bind(this);
            element.addEventListener('click', this._updateSelectedState);
        });

        this.allcheckBoxElements.forEach(element => {
            this._updateAllSelectedState = this._updateAllSelectedState.bind(this);
            element.addEventListener('click', this._updateAllSelectedState);
        });

        this.manageTasksButton.forEach(element => {
            this._showEditTasksPanel = this._showEditTasksPanel.bind(this);
            element.addEventListener('click', this._showEditTasksPanel);
        });

        this.cancelEditTasksButton.forEach(element => {
            this._hideEditTasksPanel = this._hideEditTasksPanel.bind(this);
            element.addEventListener('click', this._hideEditTasksPanel);
        });
        
        this.nickSelect.forEach(element => {
        console.log("nick bind func");
            this._katesFunction = this._katesFunction.bind(this);
            element.addEventListener('change', this._katesFunction);
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

    _katesFunction(event) {
      console.log("kate function")
        const value = event.target.value.toString();
        
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange=function() {
          if (this.readyState == 4 && this.status == 200) {
              console.log(this.response);
              var obj = JSON.parse(this.response);
              var items = obj.members
              console.log(obj);
          
            var str = ""
            items.forEach( item => {
              str += "<option value=" + item.id + ">" + item.displayName + "</option>"
              console.log(str)
            })
            document.getElementById("assignCM").innerHTML = str;
          }
        };
        xhttp.open("GET", `http://localhost:8080/api/v1/teams/${value}`, true);
        xhttp.withCredentials = true;
        xhttp.setRequestHeader("Content-Type", "application/json");
        xhttp.setRequestHeader("X-XSRF-TOKEN", this.xsrfToken.value.toString());
        xhttp.setRequestHeader("OPG-Bypass-Membrane", 1);
        xhttp.send();
        }
    
    _bindKatesFunction(element) {

      this._katesFunction = this._katesFunction(this);
      element.addEventListener('change', this._katesFunction());
    }
 }