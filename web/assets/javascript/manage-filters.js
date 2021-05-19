export default class ManageFilters {
    constructor(element) {
        
        this.taskTypeButton = element.querySelectorAll('.js-container-button');
        this.innerContainer = element.querySelector(".js-options-container");
        this.taskTypeFilterArrowUp = element.querySelector(".app-c-option-select__icon--up");
        this.taskTypeFilterArrowDown = element.querySelector(".app-c-option-select__icon--down");
        
        this._setupEventListeners();
      }

    _setupEventListeners() {
        this.taskTypeButton.forEach(element => {
            this._toggleTasktypeFilter = this._toggleTasktypeFilter.bind(this);
            element.addEventListener('click', this._toggleTasktypeFilter);
        });
    }

    _toggleTasktypeFilter() {
        const hiddenState = this.innerContainer.classList.contains('hide');
        this.innerContainer.classList.toggle("hide", !hiddenState)
        if (hiddenState) {
            this.taskTypeFilterArrowUp.classList.toggle("hide", true)
            this.taskTypeFilterArrowDown.classList.toggle("hide", false)
        } else {
            this.taskTypeFilterArrowUp.classList.toggle("hide", false)
            this.taskTypeFilterArrowDown.classList.toggle("hide", true)
        }
    }
 }