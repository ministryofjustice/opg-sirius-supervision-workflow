import './main.scss';
import GOVUKFrontend from 'govuk-frontend/govuk/all.js';
import MOJFrontend from '@ministryofjustice/frontend/moj/all.js';
import ManageTasks from './javascript/manage-tasks';
import MojBannerAutoHide from './javascript/moj-banner-auto-hide';

GOVUKFrontend.initAll();

MojBannerAutoHide(document.querySelector('.app-main-class'));

const manageTasks = document.querySelectorAll('[data-module="manage-tasks"]');
  MOJFrontend.nodeListForEach(manageTasks, function (manageTask) {
    new ManageTasks(manageTask);
  });
