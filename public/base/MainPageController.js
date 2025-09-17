import BasePageController from "../base/BasePageController.js";

export default class MainPageController extends BasePageController {
    constructor(router, pageHolder) {
        super(router, pageHolder);
    }

    pageInit() {
        super.pageInit();

        // 버튼 이벤트 바인딩 예시
        const backBtn = document.getElementById("goBackBtn");
        if (backBtn) {
            backBtn.addEventListener("click", () => {
                history.back();
            });
        }
    }
}
