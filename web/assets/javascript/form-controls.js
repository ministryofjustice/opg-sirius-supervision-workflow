const FormControls = () => {
  //event delegation bound to all data-module="app-select-submit"
  document.addEventListener("change", function (e) {
    if (e.target?.dataset?.module == "app-select-submit") {
      document.getElementById("Not Assigned").checked = false;
      document.getElementsByName("selected-assignee").forEach(function (checkbox) {
        checkbox.checked = false;
      });
      e.target.closest("form")?.submit();
    }
  });
};

export default FormControls;
