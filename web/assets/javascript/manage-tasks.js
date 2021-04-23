export default class ManageTasks {
    constructor(element) {
        this.data = {
            selectedTasks: 0
        }
        this.checkBoxElements = element.querySelectorAll('[data-wf-module="manage-tasks_checkbox"]');
        this.selectedCountElement = element.querySelectorAll('[data-wf-module="manage-tasks_task-count"]')[0];
        this.allcheckBoxElements = element.querySelectorAll('[data-wf-module="manage-tasks_all-checkboxes"]')[0];
        
        this._bindAllCheckBox(this.allcheckBoxElements);

        this.checkBoxElements.forEach(checkbox => {
            this._bindCheckBox(checkbox);
        });
      
        this.manageTasksButton = element.querySelectorAll('[data-wf-module="manage-tasks_edit-task-btn"]')[0];
        this.cancelEditTasksButton = element.querySelectorAll('[data-wf-module="manage-tasks_cancel-button"]')[0];
        this.editPanelDiv = element.querySelectorAll('[data-wf-module="manage-tasks_edit-panel"]')[0];
      
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

    _updateSelectedState(event) {
        event.target.checked ? this.data.selectedTasks++ : this.data.selectedTasks--;
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
    this._showEditTasksPanel = this._showEditTasksPanel.bind(this);
    element.addEventListener('click', this._showEditTasksPanel);
  }

}
