import GOVUKFrontend from "govuk-frontend/govuk/all.js";
import ManageTasks from "./javascript/manage-tasks";
import ManageFilters from "./javascript/manage-filters";
import ManageJumpMenus from "./javascript/manage-jump-menus";

document.body.className = document.body.className
  ? document.body.className + " js-enabled"
  : "js-enabled";

GOVUKFrontend.initAll();

const mojAutoHideBanners = document.querySelectorAll('[data-module="moj-banner-auto-hide"]');
mojAutoHideBanners.forEach(function (banner) {
  setTimeout(function () { banner.classList.add("hide"); }, 5000);
});

const jumpMenus = document.querySelectorAll('[data-module="jump-menu"]');
jumpMenus.forEach(function (jumpMenu) {
  new ManageJumpMenus(jumpMenu);
});

const manageTasks = document.querySelectorAll('[data-module="manage-tasks"]');
manageTasks.forEach(function (manageTask) {
  new ManageTasks(manageTask);
});

const manageFilters = document.querySelectorAll('[data-module="moj-manage-filters"]');
manageFilters.forEach(function (manageFilter) {
  new ManageFilters(manageFilter);
});
