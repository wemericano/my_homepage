console.log("JS LOADED");
import { API } from "./common.js";

// 공통 엘리먼트
const signupBtn = document.getElementById("signup");

let signupModal = null;

// 회원가입 버튼 클릭 시
if (signupBtn) {
    signupBtn.addEventListener("click", () => {
        openSignupModal();
    });
}

// 모달 생성 (최초 1번만)
function createSignupModal() {
    if (signupModal) return;

    signupModal = document.createElement("div");
    signupModal.id = "signupModal";
    signupModal.className =
        "fixed inset-0 bg-black bg-opacity-50 z-50 hidden flex justify-center items-center";

    signupModal.innerHTML = `
        <div class="bg-white dark:bg-gray-800 text-black dark:text-white rounded-lg shadow-lg p-8 w-full max-w-md mx-4 relative">
            <button id="closeSignup" class="absolute top-2 right-3 text-2xl font-bold text-gray-500 hover:text-gray-700 dark:text-gray-300">&times;</button>
            <h2 class="text-2xl font-bold mb-6 text-center">회원가입</h2>
            <div class="space-y-4">
                <input type="text" placeholder="아이디" id="signupId" class="w-full px-4 py-2 rounded border dark:bg-gray-700" />
                <input type="password" placeholder="패스워드" id="signupPassword" class="w-full px-4 py-2 rounded border dark:bg-gray-700" />
                <input type="text" placeholder="유저네임" id="signupUsername" class="w-full px-4 py-2 rounded border dark:bg-gray-700" />
            </div>
            <div class="flex justify-end space-x-4 mt-6">
                <button id="cancelSignup" class="px-4 py-2 bg-gray-300 dark:bg-gray-600 rounded">Cancel</button>
                <button id="doSignup" class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded">Sign Up</button>
            </div>
        </div>
    `;

    document.body.appendChild(signupModal);

    document.getElementById("cancelSignup").addEventListener("click", closeSignupModal);
    document.getElementById("closeSignup").addEventListener("click", closeSignupModal);

    document.getElementById("doSignup").addEventListener("click", () => {
        const id = document.getElementById("signupId").value;
        const password = document.getElementById("signupPassword").value;
        const username = document.getElementById("signupUsername").value;

        console.log("회원가입 입력값:", { id, password, username });
        API.POST("/api/signup", { id, password, username })
        .then((res) => {
            console.log("회원가입 성공:", res);
            closeSignupModal()
        })
        .catch((err) => {
            console.error("회원가입 실패:", err);
            closeSignupModal()
        });
    });
}

// 모달 열기
function openSignupModal() {
    if (!signupModal) {
        createSignupModal();
    }
    signupModal.classList.remove("hidden");
}

// 모달 닫기
function closeSignupModal() {
    if (signupModal) {
        signupModal.classList.add("hidden");
    }
}




// 로그인