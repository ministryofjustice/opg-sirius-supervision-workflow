export default class ManageJumpMenus {
  constructor(element) {
    element.addEventListener("change", function () {
      const selectedValue = this.options[this.selectedIndex].value;
      try {
        const url = new URL(selectedValue, window.location.origin);
        if (url.protocol === "http:" || url.protocol === "https:") {
          window.location.href = url.href;
        } else {
          console.error("Invalid URL protocol:", url.protocol);
        }
      } catch (e) {
        console.error("Invalid URL:", selectedValue);
      }
    });
  }
}
