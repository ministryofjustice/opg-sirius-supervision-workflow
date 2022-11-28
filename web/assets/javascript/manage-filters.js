export default class ManageFilters {
  constructor(element) {
    this.filterComponents = element.querySelectorAll(".moj-filter__options");

    this.inputElementTasktypeFilter = element.querySelectorAll(".task-type");
    this.inputElementAssigneeFilter =
      element.querySelectorAll(".assignee-type");
    this.urlForPage = window.location.search;

    this.teamSelection = element.querySelectorAll(".change-teams");
    this._clearFilters();
    this._isFiltered();
    this._isFilteredByAssignee();
    this._setupEventListeners();
    this._setupFilters();
  }

  _setupEventListeners() {
    this.filterComponents.forEach((element) => {
      this._toggleFilterVisibility = this._toggleFilterVisibility.bind(this);
      element
        .querySelectorAll(".filter-toggle-button")[0]
        .addEventListener("click", this._toggleFilterVisibility);
    });
  }

  _setupFilters() {
    this.filterComponents.forEach((element) => {
      const filterName = element.dataset.filterName;
      let isOpen = this._getFilterStatusFromLocalStorage(filterName);

      this._setFilterVisibility(element, isOpen);
    });
  }

  _setFilterStatusToLocalStorage(filterName, isOpen) {
    window.sessionStorage.setItem(
      filterName,
      JSON.stringify({ value: isOpen })
    );
  }

  _getFilterStatusFromLocalStorage(filterName) {
    let sessionStorageValue = JSON.parse(
      window.sessionStorage.getItem(filterName)
    );
    if (!sessionStorageValue) {
      sessionStorageValue = { value: false };
      this._setFilterStatusToLocalStorage(
        filterName,
        sessionStorageValue.value
      );
    }

    return sessionStorageValue.value;
  }

  _toggleFilterVisibility(e) {
    const filterElement = e.target.parentElement.parentElement.parentElement;
    const filterName = filterElement.dataset.filterName;

    let isOpen = this._getFilterStatusFromLocalStorage(filterName);
    isOpen = isOpen === true ? false : true;

    this._setFilterVisibility(filterElement, isOpen);

    this._setFilterStatusToLocalStorage(filterName, isOpen);
  }

  _setFilterVisibility(element, isOpen) {
    let filterInnerContainer = element.querySelector(".filter-inner-container");
    let filterArrowUp = element.querySelector(".filter-arrow-up");
    let filterArrowDown = element.querySelector(".filter-arrow-down");

    filterInnerContainer.classList.toggle("hide", !isOpen);

    filterArrowUp.setAttribute("aria-expanded", isOpen.toString());
    filterArrowDown.setAttribute("aria-expanded", !isOpen.toString());

    filterArrowUp.classList.toggle("hide", !isOpen);
    filterArrowDown.classList.toggle("hide", isOpen);
  }

  _isFiltered() {
    let array = [];

    this.inputElementTasktypeFilter.forEach((taskType) => {
      if (taskType.checked) {
        array.push(taskType);
      }
    });

    let append = "";

    if (array.length) {
      append +=
        '<h3 class="govuk-heading-s govuk-!-margin-bottom-0" id="Task-type-hook">Task type</h3>';
      array.forEach((taskType) => {
        let hrefValue = this.urlForPage
          .split("&")
          .filter((param) => !param.includes(taskType.value))
          .join("&");
        append +=
          `<li id=${taskType.value}><a class="moj-filter__tag" href=${hrefValue}><span class="govuk-visually-hidden">Remove this filter</span>` +
          taskType.id +
          "</li>";
      });
    }

    document.getElementById("applied-task-type-filters").innerHTML = append;
  }

  _isFilteredByAssignee() {
    let array = [];

    this.inputElementAssigneeFilter.forEach((assignee) => {
      if (assignee.checked) {
        array.push(assignee);
      }
    });

    let append = "";

    if (array.length) {
      append +=
        '<h3 class="govuk-heading-s govuk-!-margin-bottom-0" id="Assignee-hook">Assignees</h3>';
      array.forEach((assignee) => {
        let hrefValue = this.urlForPage
          .split("&")
          .filter(
            (param) => !param.includes("selected-assignee=" + assignee.value)
          )
          .join("&");
        append +=
          `<li id=${assignee.value}><a class="moj-filter__tag" href=${hrefValue}><span class="govuk-visually-hidden">Remove this filter</span>` +
          assignee.id +
          "</li>";
      });
    }

    document.getElementById("applied-assignee-filters").innerHTML = append;
  }

  _clearFilters() {
    let hrefValueWithoutSelectedTask = this.urlForPage
      .split("&")
      .filter((param) => !param.includes("selected-task-type"))
      .join("&");
    let hrefValue = hrefValueWithoutSelectedTask
      .split("&")
      .filter((param) => !param.includes("selected-assignee"))
      .join("&");
    document.getElementById("clear-filters").setAttribute("href", hrefValue);
  }
}
