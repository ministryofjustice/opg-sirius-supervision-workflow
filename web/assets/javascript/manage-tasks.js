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
        this.nick = element.querySelectorAll('.option-value');
        this.nickSelect = element.querySelectorAll('.option-value-select');

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
        console.log("_katesFunction");
        console.log("event target attributes value");
        console.log(event.target.attributes.value);
        this.nickSelect.forEach(element => { console.log(element) })
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange=function() {
          if (this.readyState == 4 && this.status == 200) {
            document.getElementById("kate").innerHTML = "loaded"
            // console.log(xhttp.responseText);
            // console.log(xhttp.response);
          }
        };
        xhttp.open("GET", "/api/v1/teams/" + 13, true);
        xhttp.send();
        }

        // _katesFunction(event) {
        //   console.log("_katesFunction");
        //   console.log(this);
        //   var xhttp = new XMLHttpRequest();
        //   xhttp.onreadystatechange=function() {
        //     if (this.readyState == 4 && this.status == 200) {
        //       document.getElementById("kate").innerHTML = "loaded"
        //       console.log(this.response);
        //       console.log(this.responseText);
        //     }
        //   };
        //   xhttp.open("GET", "/api/v1/teams/" + 13, true);
        //   xhttp.send();
        //   }
    
    _bindKatesFunction(element) {

      this._katesFunction = this._katesFunction(this);
      element.addEventListener('change', this._katesFunction());
    }
 }