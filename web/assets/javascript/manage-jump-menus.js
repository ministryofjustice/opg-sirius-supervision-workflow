export default class ManageJumpMenus {
  constructor(element) {
    element.addEventListener("change", function () {
      window.location.href = this.options[this.selectedIndex].value;
    });
  }
}
