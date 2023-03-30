export default class ManageFilters {
  constructor(element) {
    this.applyFilters = element.querySelector('[data-module="apply-filters"]');
    this.clearFilters = element.querySelector('[data-module="clear-filters"]');
    this.filters = element.querySelectorAll('[data-module="filter"]');
    this.filterComponents = element.querySelectorAll(".moj-filter__options");

    this._setupEventListeners();
    this._setupFilters();
  }

  _setupEventListeners() {
    this._applyFilters = this._applyFilters.bind(this);
    this.applyFilters.addEventListener("click", this._applyFilters);

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
    isOpen = isOpen !== true;

    this._setFilterVisibility(filterElement, isOpen);

    this._setFilterStatusToLocalStorage(filterName, isOpen);
  }

  _setFilterVisibility(element, isOpen) {
    let filterInnerContainer = element.querySelector(".filter-inner-container");
    let filterArrowUp = element.querySelector(".filter-arrow-up");
    let filterArrowDown = element.querySelector(".filter-arrow-down");

    filterInnerContainer.classList.toggle("hide", !isOpen);

    filterArrowUp.setAttribute("aria-expanded", isOpen.toString());
    filterArrowDown.setAttribute("aria-expanded", (!isOpen).toString());

    filterArrowUp.classList.toggle("hide", !isOpen);
    filterArrowDown.classList.toggle("hide", isOpen);
  }

  _applyFilters() {
    let url = this.clearFilters.getAttribute("href");
    this.filters.forEach(function (filter) {
      if (filter.checked) {
        url += "&" + filter.name + "=" + filter.value
      }
    });
    window.location.href = url;
  }
}
