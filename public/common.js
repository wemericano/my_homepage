console.log("🎵 COMMON LOADED");

// 다크모드 토글
export function initDarkMode() {
  const themeToggle = document.getElementById("themeToggle");
  if (!themeToggle) return;

  const updateIcon = () => {
    themeToggle.textContent = document.documentElement.classList.contains("dark") ? "🌞" : "🌙";
  };

  themeToggle.addEventListener("click", () => {
    document.documentElement.classList.toggle("dark");
    updateIcon();
  });

  updateIcon(); // 초기 상태 반영
}

// ESC로 모달 닫기
export function initEscClose(modals = []) {
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") {
      modals.forEach((modal) => modal.classList.add("hidden"));
    }
  });
}

// API Class
export class API {
    static async GET(url) {
        try {
            const res = await fetch(url, {
                method: "GET",
                headers: { "Content-Type": "application/json" },
            });
            if (!res.ok) throw new Error(`GET ${url} 실패: ${res.status}`);
            return await res.json();
        } catch (err) {
            console.error("API.GET 오류:", err);
            throw err;
        }
    }

    static async POST(url, data = {}) {
        try {
            const res = await fetch(url, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(data),
            });
            if (!res.ok) throw new Error(`POST ${url} 실패: ${res.status}`);
            return await res.json();
        } catch (err) {
            console.error("API.POST 오류:", err);
            throw err;
        }
    }
}

// API 클래스를 전역으로 노출
if (typeof window !== 'undefined') {
    window.API = API;
}

// 페이지 이동
export function goTo(url) {
  if (!url || typeof url !== "string") {
      console.warn("유효하지 않은 URL입니다:", url);
      return;
  }
  window.location.href = url;
}