import MainPageController from "/base/MainPageController.js";

console.log("[🎵] Main JS LOADED");

document.addEventListener("DOMContentLoaded", () => {
    const controller = new MainPageController(); // 필요시 인자 전달
    controller.pageInit();
});
