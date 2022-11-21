import GOVUKFrontend from "govuk-frontend/govuk/all.js";
import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";
import ManageTasks from "./javascript/manage-tasks";
import ManageFilters from "./javascript/manage-filters";
import MojBannerAutoHide from "./javascript/moj-banner-auto-hide";
import FormControls from "./javascript/form-controls";

//moved from inline due to Content Security Policy
document.body.className = document.body.className
  ? document.body.className + " js-enabled"
  : "js-enabled";

GOVUKFrontend.initAll();

FormControls();

MojBannerAutoHide(document.querySelector(".app-main-class"));

const manageTasks = document.querySelectorAll('[data-module="manage-tasks"]');
MOJFrontend.nodeListForEach(manageTasks, function (manageTask) {
  new ManageTasks(manageTask);
});
const manageFilters = document.querySelectorAll(
  '[data-module="moj-manage-filters"]'
);
MOJFrontend.nodeListForEach(manageFilters, function (manageFilter) {
  new ManageFilters(manageFilter);
});
