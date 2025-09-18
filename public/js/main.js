import BasePageController from "../base/BasePageController.js";
import  * as Common from "./page_common.js";

class MainPageController extends BasePageController {
    constructor(router, pageHolder) {
        super(router, pageHolder);
    }

    pageInit() {
        super.pageInit();
        console.log("🎵 Main JS LOADED");

        this.goBack();
        this.initLottoButton();
        this.initAnalyzeButton();
    }

    goBack() {
        const backBtn = document.getElementById("goBackBtn");
        if (backBtn) {
            
            backBtn.addEventListener("click", () => {
                history.back();
            });
        }
    }

    initLottoButton() {
        const lottoBtn = document.getElementById("lottoBtn");
        if (lottoBtn) {
            lottoBtn.addEventListener("click", () => {
                this.getLottoList();
            });
        }
    }

    initAnalyzeButton() {
        const analyzeBtn = document.getElementById("analyzeBtn");
        if (analyzeBtn) {
            analyzeBtn.addEventListener("click", () => {
                this.analyzeV1();
            });
        }
    }

    getLottoList() {
        Common.API.POST("/api/lotto")
        .then((res) => {
            console.log("API 응답:", res);
            if (res.code === "0000" && res.data) {
                this.displayLottoData(res.data);
            } else {
                console.error("API 오류:", res.message);
            }
        })
        .catch((err) => {
            console.error("API 오류:", err);
        });
    }

    analyzeV1() {
        Common.API.POST("/api/analyze/v1")
        .then((res) => {
            console.log("API 응답:", res);
            if (res.code === "0000" && res.data) {
            } else {
                console.error("API 오류:", res.message);
            }
        })
        .catch((err) => {
            console.error("API 오류:", err);
        });
    }

    displayLottoData(lottoData) {
        const tableBody = document.getElementById("lottoTableBody");
        const tableDiv = document.getElementById("lottoTable");
        
        if (!tableBody || !tableDiv) return;

        // 테이블 내용 초기화
        tableBody.innerHTML = "";

        // 데이터가 있으면 테이블 표시
        if (lottoData && lottoData.length > 0) {
            tableDiv.style.display = "block";
            
            // 각 로또 데이터를 테이블 행으로 추가
            lottoData.forEach((lotto) => {
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">${lotto.indexNo}</td>
                    <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b;">${lotto.no1}</td>
                    <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b;">${lotto.no2}</td>
                    <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b;">${lotto.no3}</td>
                    <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b;">${lotto.no4}</td>
                    <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b;">${lotto.no5}</td>
                    <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b;">${lotto.no6}</td>
                `;
                tableBody.appendChild(row);
            });
        } else {
            // 데이터가 없으면 메시지 표시
            const row = document.createElement("tr");
            row.innerHTML = `<td colspan="7" style="border: 1px solid #ddd; padding: 8px; text-align: center;">데이터가 없습니다.</td>`;
            tableBody.appendChild(row);
            tableDiv.style.display = "block";
        }
    }
}

document.addEventListener("DOMContentLoaded", () => {
    const controller = new MainPageController();
    controller.pageInit();
});