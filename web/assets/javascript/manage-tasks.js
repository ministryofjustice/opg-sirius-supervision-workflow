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
    }

    numberOfTasksSelected() {
        return this.data.selectedTasks;
    }

    _updateDomElements() {
        this.selectedCountElement.innerText = this.numberOfTasksSelected().toString();
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
}
