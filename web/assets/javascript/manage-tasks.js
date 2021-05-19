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
        this.assignTeamSelect = element.querySelectorAll('.js-assign-team-select');
        this.xsrfToken = element.querySelector('.js-xsrfToken');
        this.selectedCountElement = element.querySelectorAll('.js-mt-task-count');
        this.editPanelDiv = element.querySelectorAll('.js-mt-edit-panel');
        this.baseUrl = document.querySelector('[name=api-base-uri]').getAttribute('content')
        this.taskTypeCheckBox
        this.taskTypeButton = element.querySelectorAll('.js-container-button');
  
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
        
        this.assignTeamSelect.forEach(element => {
            this._getCaseManagers = this._getCaseManagers.bind(this);
            element.addEventListener('change', this._getCaseManagers);
        });    

        this.taskTypeButton.forEach(element => {
            this._toggleTasktypeFilter = this._toggleTasktypeFilter.bind(this);
            element.addEventListener('click', this._toggleTasktypeFilter);
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
            data.members.forEach( caseManager => {
               str += "<option value=" + caseManager.id + ">" + caseManager.displayName + "</option>"
            })
        
            document.getElementById("assignCM").innerHTML = str;
        });
    }

    _toggleTasktypeFilter(event) {
        const innerContainer = event.parentElement.querySelector(".js-options-container");
        innerContainer.classList.toggle("hide")
    }
 }