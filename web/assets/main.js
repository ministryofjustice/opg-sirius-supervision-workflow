import './main.scss';
import GOVUKFrontend from 'govuk-frontend/govuk/all.js';
import MOJFrontend from '@ministryofjustice/frontend/moj/all.js';
import ManageTasks from './javascript/manage-tasks';

GOVUKFrontend.initAll();

const manageTasks = document.querySelectorAll('[data-module="manage-tasks"]');
  MOJFrontend.nodeListForEach(manageTasks, function (manageTask) {
    new ManageTasks(manageTask);
  });
