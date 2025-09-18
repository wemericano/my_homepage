import BasePageController from "../base/BasePageController.js";

class MainPageController extends BasePageController {
    constructor(router, pageHolder) {
        super(router, pageHolder);
    }

    pageInit() {
        super.pageInit();

    }

}

document.addEventListener("DOMContentLoaded", () => {
    const controller = new MainPageController();
    controller.pageInit();
});