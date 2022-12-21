const FormControls = () => {
  //event delegation bound to all data-module="app-select-submit"
  document.addEventListener("change", function (e) {
    if (e.target?.dataset?.module == "app-select-submit") {
      e.target.closest("form")?.submit();
    }
  });
};

export default FormControls;
