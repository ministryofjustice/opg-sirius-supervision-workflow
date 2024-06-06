export default class GoBack {
  constructor(element) {
    element.addEventListener("click", function () {
      window.location.href = history.back();
    });
  }
}
