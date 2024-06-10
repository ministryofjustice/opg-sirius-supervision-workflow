export default class CloseTab {
  constructor(element) {
    element.addEventListener("click", function () {
      window.location.href = window.close();
    });
  }
}
