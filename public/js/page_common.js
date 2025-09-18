console.log("🎵 PAGE COMMON LOADED");

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