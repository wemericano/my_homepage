console.log("🎵 JS LOADED");
import * as Common from "./common.js";

// 공통 엘리먼트
const signupBtn = document.getElementById("signup");
const loginBtn = document.getElementById("login");
const defenseBtn = document.getElementById("defense");

let signupModal = null;
let loginModal = null;
let blogModal = null;

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

        Common.API.POST("/api/signup", { id, password, username })
        .then((res) => {
            closeSignupModal()
        })
        .catch((err) => {
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

// 로그인 버튼 클릭
if (loginBtn) {
    loginBtn.addEventListener("click", () => {
        openLoginModal();
    });
}

// 로그인 모달 생성
function createLoginModal() {
    if (loginModal) return;

    loginModal = document.createElement("div");
    loginModal.id = "loginModal";
    loginModal.className =
        "fixed inset-0 bg-black bg-opacity-50 z-50 hidden flex justify-center items-center";

    loginModal.innerHTML = `
        <div class="bg-white dark:bg-gray-800 text-black dark:text-white rounded-lg shadow-lg p-8 w-full max-w-md mx-4 relative">
            <button id="closeLogin" class="absolute top-2 right-3 text-2xl font-bold text-gray-500 hover:text-gray-700 dark:text-gray-300">&times;</button>
            <h2 class="text-2xl font-bold mb-6 text-center">로그인</h2>
            <div class="space-y-4">
                <input type="text" placeholder="아이디" id="loginId" class="w-full px-4 py-2 rounded border dark:bg-gray-700" />
                <input type="password" placeholder="패스워드" id="loginPassword" class="w-full px-4 py-2 rounded border dark:bg-gray-700" />
                <p id="loginError" class="text-red-500 text-sm mt-1 hidden"></p>
            </div>
            <div class="flex justify-end space-x-4 mt-6">
                <button id="cancelLogin" class="px-4 py-2 bg-gray-300 dark:bg-gray-600 rounded">Cancel</button>
                <button id="doLogin" class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded">Login</button>
            </div>
        </div>
    `;

    document.body.appendChild(loginModal);

    document.getElementById("cancelLogin").addEventListener("click", closeLoginModal);
    document.getElementById("closeLogin").addEventListener("click", closeLoginModal);
    const m_error = document.getElementById("loginError");

    m_error.classList.add("hidden");
    m_error.innerText = "";

    document.getElementById("doLogin").addEventListener("click", () => {
        const id = document.getElementById("loginId").value;
        const password = document.getElementById("loginPassword").value;

        Common.API.POST("/api/login", { id, password })
        .then((res) => {
            if (res.data) {
                Common.goTo("/page/main.html");
                closeLoginModal()
            }
            m_error.innerText = "아이디 또는 비밀번호가 올바르지 않습니다.";
            m_error.classList.remove("hidden");
        })
        .catch((err) => {
            closeLoginModal()
        });
    });
}

function openLoginModal() {
    if (!loginModal) {
        createLoginModal();
    }
    loginModal.classList.remove("hidden");
}

function closeLoginModal() {
    if (loginModal) {
        loginModal.classList.add("hidden");
    }
}

// Enter 키로 로그인 (index.html의 onkeyup에서 사용)
function enterkey(event) {
    if (event && event.key === "Enter") {
        login();
    }
}

// 로그인 실행 함수 (index.html의 onclick에서 사용)
function login() {
    const username = document.getElementById("username");
    const password = document.getElementById("password");
    
    if (!username || !password) {
        // 동적으로 생성된 모달 사용
        const loginId = document.getElementById("loginId");
        const loginPassword = document.getElementById("loginPassword");
        if (loginId && loginPassword) {
            const id = loginId.value;
            const pwd = loginPassword.value;
            
            Common.API.POST("/api/login", { id, password: pwd })
            .then((res) => {
                if (res.data) {
                    Common.goTo("/page/main.html");
                    closeLoginModal();
                } else {
                    const m_error = document.getElementById("loginError");
                    if (m_error) {
                        m_error.innerText = "아이디 또는 비밀번호가 올바르지 않습니다.";
                        m_error.classList.remove("hidden");
                    }
                }
            })
            .catch((err) => {
                console.error("로그인 오류:", err);
            });
        }
        return;
    }
    
    // HTML에 하드코딩된 모달 사용
    const id = username.value;
    const pwd = password.value;
    
    Common.API.POST("/api/login", { id, password: pwd })
    .then((res) => {
        if (res.data) {
            Common.goTo("/page/main.html");
            closeLoginModal();
        } else {
            alert("아이디 또는 비밀번호가 올바르지 않습니다.");
        }
    })
    .catch((err) => {
        console.error("로그인 오류:", err);
        alert("로그인 중 오류가 발생했습니다.");
    });
}

// 디펜스게임
if (defenseBtn) {
    defenseBtn.addEventListener("click", () => {
        console.log("디펜스게임");
        location.href = "./page/defense.html";
    });
}

// 페이지 이동 함수들
function defense() {
    location.href = "./page/defense.html";
}

function lotto() {
    location.href = "./page/main.html";
}

function lotto_() {
    location.href = "./page/main.html";
}

function spec() {
    alert("게임 기능은 준비 중입니다.");
}

function diff() {
    alert("아직 공개되지 않은 비밀입니다.");
}

function test() {
    alert("메뉴 기능은 준비 중입니다.");
}

function luck() {
    alert("오늘의 행운: 좋은 일이 있을 것입니다! 🍀");
}

// 회원가입 실행 함수
function signup() {
    const signupModal = document.getElementById("signupModal");
    if (!signupModal) {
        openSignupModal();
        return;
    }
    
    // HTML에 하드코딩된 모달 사용
    const inputs = signupModal.querySelectorAll("input");
    if (inputs.length >= 3) {
        const id = inputs[0].value;
        const password = inputs[1].value;
        const username = inputs[2].value;
        
        Common.API.POST("/api/signup", { id, password, username })
        .then((res) => {
            alert("회원가입이 완료되었습니다.");
            closeSignupModal();
        })
        .catch((err) => {
            console.error("회원가입 오류:", err);
            alert("회원가입 중 오류가 발생했습니다.");
        });
    }
}

// 블로그 모달 생성
function createBlogModal() {
    if (blogModal) return;

    blogModal = document.createElement("div");
    blogModal.id = "blogModal";
    blogModal.className =
        "fixed inset-0 bg-black bg-opacity-50 z-50 hidden flex justify-center items-center";

    blogModal.innerHTML = `
        <div class="bg-white dark:bg-gray-800 text-black dark:text-white rounded-lg shadow-lg p-6 w-full max-w-xs mx-4 relative">
            <h2 class="text-lg font-bold mb-4 text-center">Secret Key</h2>
            <div class="space-y-4">
                <input type="text" id="secretKey" class="w-full px-4 py-2 rounded border dark:bg-gray-700" placeholder="키를 입력하세요" />
                <p id="blogError" class="text-red-500 text-sm text-center hidden"></p>
            </div>
            <div class="flex justify-end space-x-4 mt-4">
                <button id="checkSecretKey" class="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded text-sm">확인</button>
            </div>
        </div>
    `;

    document.body.appendChild(blogModal);

    document.getElementById("checkSecretKey").addEventListener("click", () => {
        const key = document.getElementById("secretKey").value;
        const errorMsg = document.getElementById("blogError");
        
        if (key === "z") {
            location.href = "./page/blog.html";
        } else {
            errorMsg.innerText = "NOP!";
            errorMsg.classList.remove("hidden");
            setTimeout(() => {
                location.href = "./";
            }, 1000);
        }
    });

    // Enter 키로 확인
    document.getElementById("secretKey").addEventListener("keypress", (e) => {
        if (e.key === "Enter") {
            document.getElementById("checkSecretKey").click();
        }
    });
}

// 블로그 모달 열기
function openBlogModal() {
    if (!blogModal) {
        createBlogModal();
    }
    blogModal.classList.remove("hidden");
    const errorMsg = document.getElementById("blogError");
    if (errorMsg) {
        errorMsg.classList.add("hidden");
        errorMsg.innerText = "";
    }
    const input = document.getElementById("secretKey");
    if (input) {
        input.value = "";
        input.focus();
    }
}

// 전역 함수로 노출
if (typeof window !== 'undefined') {
    window.enterkey = enterkey;
    window.login = login;
    window.closeLoginModal = closeLoginModal;
    window.closeSignupModal = closeSignupModal;
    window.defense = defense;
    window.lotto = lotto;
    window.lotto_ = lotto_;
    window.spec = spec;
    window.diff = diff;
    window.test = test;
    window.luck = luck;
    window.signup = signup;
    window.openBlogModal = openBlogModal;
}