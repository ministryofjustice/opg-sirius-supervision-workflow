import "govuk-frontend/dist/govuk/all.mjs";
import ManageFilters from "./javascript/manage-filters";
import ManageJumpMenus from "./javascript/manage-jump-menus";
import "opg-sirius-header/sirius-header.js";
import ManageReassign from "./javascript/manage-reassign";

document.body.className = document.body.className
  ? document.body.className + " js-enabled"
  : "js-enabled";

const mojAutoHideBanners = document.querySelectorAll('[data-module="moj-banner-auto-hide"]');
mojAutoHideBanners.forEach(function (banner) {
  setTimeout(function () { banner.classList.add("hide"); }, 5000);
});

const jumpMenus = document.querySelectorAll('[data-module="jump-menu"]');
jumpMenus.forEach(function (jumpMenu) {
  new ManageJumpMenus(jumpMenu);
});

const manageReassign = document.querySelectorAll('[data-module="manage-reassign"]');
manageReassign.forEach(function (manageReassign) {
  new ManageReassign(manageReassign);
});

const manageFilters = document.querySelectorAll('[data-module="moj-manage-filters"]');
manageFilters.forEach(function (manageFilter) {
  new ManageFilters(manageFilter);
});
