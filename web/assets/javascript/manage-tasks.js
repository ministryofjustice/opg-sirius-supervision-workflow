export default class ManageTasks {
    constructor(element) {
        this.data = {
            selectedTasks: 0
        }
        this.checkBoxElements = element.querySelectorAll('[data-module="manage-tasks_checkbox"]');
        this.selectedCountElement = element.querySelectorAll('[data-module="manage-tasks_task-count"]')[0];
        this.allcheckBoxElements = element.querySelectorAll('[data-module="manage-tasks_all-checkboxes"]')[0];
        
        this._bindAllCheckBox(this.allcheckBoxElements);

        this.checkBoxElements.forEach(checkbox => {
            this._bindCheckBox(checkbox);
        });
      
        this.manageTasksButton = element.querySelectorAll('[data-module="manage-tasks_edit-task-btn"]')[0];
        this.cancelEditTasksButton = element.querySelectorAll('[data-module="manage-tasks_cancel-button"]')[0];
        this.editPanelDiv = element.querySelectorAll('[data-module="manage-tasks_edit-panel"]')[0];

        this.kate = element.querySelectorAll('[data-module="manage-tasks_kate"]')[0];
      
        this._bindShowManageTasksButton(this.manageTasksButton);
        this._bindCancelTasksButton(this.cancelEditTasksButton);
        this._bindKatesFunction(this.kate);
    }
    
    numberOfTasksSelected() {
        return this.data.selectedTasks;
    }

    _updateDomElements() {
        this.selectedCountElement.innerText = this.numberOfTasksSelected().toString();
        this._showManageTasksButton();
    }

    _bindCheckBox(element) {
        this._updateSelectedState = this._updateSelectedState.bind(this);
        element.addEventListener('click', this._updateSelectedState);
    }

    _updateSelectedState(event) {
        event.target.checked ? this.data.selectedTasks++ : this.data.selectedTasks--;
        event.target.parentElement.parentElement.parentElement.classList.toggle('govuk-table__select', event.target.checked);

        event.target.parentElement.parentElement.parentElement.parentElement.classList.toggle('selected', event.target.checked);

        this._updateDomElements();   
    }

    _bindAllCheckBox(element) {
        this._updateAllSelectedState = this._updateAllSelectedState.bind(this);
        element.addEventListener('click', this._updateAllSelectedState);
    }

    _updateAllSelectedState(event) {
        let isChecked = event.target.checked; 

        this.checkBoxElements.forEach(checkbox => {
            checkbox.checked = isChecked;
        });

        this.data.selectedTasks = (isChecked ? this.checkBoxElements.length : 0);
        this._updateDomElements();
    }

    _showManageTasksButton() {
      this.manageTasksButton.classList.toggle('hide', this.data.selectedTasks === 0);
    }

    _bindShowManageTasksButton(element) {
      this._showEditTasksPanel = this._showEditTasksPanel.bind(this);
      element.addEventListener('click', this._showEditTasksPanel);
    }

   _showEditTasksPanel(event) {
    this.editPanelDiv.classList.toggle('hide', this.data.selectedTasks === 0);
   }

   _bindCancelTasksButton(element) {
    this._hideEditTasksPanel = this._hideEditTasksPanel.bind(this);
    element.addEventListener('click', this._hideEditTasksPanel);
  }

  _hideEditTasksPanel(event) {
    this.editPanelDiv.classList.toggle('hide', true);
   }

  _katesFunction() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange=function() {
      if (this.readyState == 4 && this.status == 200) {
        document.getElementById("kate").innerHTML = this.responseText;
      }
    };
    xhttp.open("GET", "/api/v1/teams/" + 14, true);
    xhttp.send();
    }

    _bindKatesFunction(element) {
      this._katesFunction = this._katesFunction(this);
      element.addEventListener('onchange', this._katesFunction);
    }
   

}
