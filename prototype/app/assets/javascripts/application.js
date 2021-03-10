/* global $ */

// Warn about using the kit in production
if (window.console && window.console.info) {
  window.console.info('Kate and Nick need to change this before it goes to production')
}

$(document).ready(function () {
  window.GOVUKFrontend.initAll()
})
