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
        this.initAnalyzeButtonV1();
        this.initAnalyzeButtonV2();
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

    initAnalyzeButtonV1() {
        const analyzeBtn = document.getElementById("analyzeBtnV1");
        if (analyzeBtn) {
            analyzeBtn.addEventListener("click", () => {
                this.analyzeV1();
            });
        }
    }

    initAnalyzeButtonV2() {
        const analyzeBtn = document.getElementById("analyzeBtnV2");
        if (analyzeBtn) {
            analyzeBtn.addEventListener("click", () => {
                this.analyzeV2();
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
                this.displayAnalyzeV1Data(res.data);
            } else {
                console.error("API 오류:", res.message);
            }
        })
        .catch((err) => {
            console.error("API 오류:", err);
        });
    }

    analyzeV2() {
        Common.API.POST("/api/analyze/v2")
        .then((res) => {
            console.log("API 응답:", res);
            if (res.code === "0000" && res.data) {
                this.displayAnalyzeV2Data(res.data);
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

    displayAnalyzeV1Data(analyzeData) {
        const tableBody = document.getElementById("lottoTableBody");
        const tableDiv = document.getElementById("lottoTable");
        
        if (!tableBody || !tableDiv) return;

        // 테이블 내용 초기화
        tableBody.innerHTML = "";

        // 분석 데이터가 있으면 테이블 표시
        if (analyzeData && analyzeData.length > 0) {
            tableDiv.style.display = "block";
            
            // 분석 결과 (1개)를 테이블 행으로 추가
            const analyzeItem = analyzeData[0];
            const row = document.createElement("tr");
            row.innerHTML = `
                <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #e3f2fd;">분석 결과</td>
                <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b; font-weight: bold;">${analyzeItem.no1}</td>
                <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b; font-weight: bold;">${analyzeItem.no2}</td>
                <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b; font-weight: bold;">${analyzeItem.no3}</td>
                <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b; font-weight: bold;">${analyzeItem.no4}</td>
                <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b; font-weight: bold;">${analyzeItem.no5}</td>
                <td style="border: 1px solid #ddd; padding: 8px; text-align: center; background: #ffeb3b; font-weight: bold;">${analyzeItem.no6}</td>
            `;
            tableBody.appendChild(row);
        } else {
            // 데이터가 없으면 메시지 표시
            const row = document.createElement("tr");
            row.innerHTML = `<td colspan="7" style="border: 1px solid #ddd; padding: 8px; text-align: center;">분석 데이터가 없습니다.</td>`;
            tableBody.appendChild(row);
            tableDiv.style.display = "block";
        }
    }

    displayAnalyzeV2Data(lottoData) {
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
                    <td style="border: 1px solid #ddd; padding: 8px; text-align: center;">>>></td>
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