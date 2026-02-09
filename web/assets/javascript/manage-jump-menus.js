export default class ManageJumpMenus {
  constructor(element) {
    element.addEventListener("change", function () {
      const url = new URL(this.options[this.selectedIndex].value, window.location.href)
      window.location.href = url.href;
    });
  }
}
