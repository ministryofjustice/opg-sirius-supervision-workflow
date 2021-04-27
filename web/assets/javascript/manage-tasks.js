export default class ManageTasks {
    constructor(element) {
        this.data = {
            selectedTasks: 0
        }
        this.teamMemberData = [];
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
      
        this._bindShowManageTasksButton(this.manageTasksButton);
        this._bindCancelTasksButton(this.cancelEditTasksButton);
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

    _updateSelectedRowStyles(element) {
      element.parentElement.parentElement.parentElement.classList.toggle('govuk-table__select', element.checked);
      element.parentElement.parentElement.parentElement.parentElement.classList.toggle('selected', element.checked);

    }

    _updateSelectedState(event) {
        event.target.checked ? this.data.selectedTasks++ : this.data.selectedTasks--;
        this._updateSelectedRowStyles(event.target);
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

            this._updateSelectedRowStyles(checkbox);
        });

        this.data.selectedTasks = (isChecked ? this.checkBoxElements.length : 0);
        this._updateDomElements();
    }

    _showManageTasksButton() {
        console.log("_showManageTasksButton");
      this.manageTasksButton.classList.toggle('hide', this.data.selectedTasks === 0);
    }

    _bindShowManageTasksButton(element) {
      this._showEditTasksPanel = this._showEditTasksPanel.bind(this);
      element.addEventListener('click', this._showEditTasksPanel);
    }

   _showEditTasksPanel(event) {
       console.log("_showEditTasksPanel")
    this.editPanelDiv.classList.toggle('hide', false);
   }

   _bindCancelTasksButton(element) {
    this._hideEditTasksPanel = this._hideEditTasksPanel.bind(this);
    element.addEventListener('click', this._hideEditTasksPanel);
  }

  _hideEditTasksPanel(event) {
    console.log("_hideEditTasksPanel")
    this.editPanelDiv.classList.toggle('hide', true);
   }

}
