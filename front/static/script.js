/* script.js
   Две вкладки:
     - shorten: сгенерировать шорт по оригинальному URL
     - lookup:  найти шорт по оригинальному URL
   Ожидаемые эндпоинты:
     POST /shorten { "url": "<оригинальный URL>" } -> "https://sho.rt/abc123" (text или {"result": "...", "short": "..."})
     POST /lookup  { "url": "<оригинальный URL>" } -> "https://sho.rt/abc123" (text или {"result": "...", "short": "..."})
*/

(function () {
    "use strict";

    // Настройка путей к API
    const ENDPOINTS = {
        shorten: "/shorten",
        lookup:  "/lookup",
    };

    // Захват DOM-элементов
    const tabShorten = document.getElementById("tab-shorten");
    const tabExpand  = document.getElementById("tab-expand"); // переиспользуем как "lookup"
    const input      = document.getElementById("input-url");
    const goBtn      = document.getElementById("go-btn");
    const output     = document.getElementById("output-box");

    if (!tabShorten || !tabExpand || !input || !goBtn || !output) {
        console.error("Проверь id элементов: tab-shorten, tab-expand, input-url, go-btn, output-box");
        return;
    }

    output.readOnly = true;

    // Текущее состояние вкладки: 'shorten' | 'lookup'
    let mode = "shorten";

    function setMode(nextMode) {
        mode = nextMode;

        if (mode === "shorten") {
            tabShorten.classList.add("tab-active");
            tabShorten.classList.remove("tab-inactive");
            tabExpand.classList.add("tab-inactive");
            tabExpand.classList.remove("tab-active");
            input.placeholder = "Вставьте оригинальный URL (сгенерировать шорт)";
        } else {
            tabExpand.classList.add("tab-active");
            tabExpand.classList.remove("tab-inactive");
            tabShorten.classList.add("tab-inactive");
            tabShorten.classList.remove("tab-active");
            input.placeholder = "Вставьте оригинальный URL (найти шорт)";
        }

        input.value = "";
        output.value = "";
    }

    function isNonEmpty(s) {
        return typeof s === "string" && /\S/.test(s);
    }

    async function send() {
        const val = input.value.trim();
        if (!isNonEmpty(val)) {
            hint("Введите URL выше.");
            return;
        }

        lock(true);

        try {
            let url = "";
            let payload = { url: val }; // в обоих режимах шлём оригинальный URL

            if (mode === "shorten") {
                url = ENDPOINTS.shorten;
            } else {
                url = ENDPOINTS.lookup;
            }

            const resp = await fetch(url, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });

            if (!resp.ok) throw new Error("Сервер вернул " + resp.status);

            // Поддержим text/plain и JSON
            let text = "";
            const ct = (resp.headers.get("Content-Type") || "").toLowerCase();

            if (ct.includes("application/json")) {
                const data = await resp.json();
                text = data.result || data.short || data.url || JSON.stringify(data);
            } else {
                text = await resp.text();
            }

            output.value = text || "";
        } catch (e) {
            output.value = "Ошибка: " + (e && e.message ? e.message : "что-то пошло не так");
        } finally {
            lock(false);
        }
    }

    function lock(state) {
        goBtn.disabled = state;
        input.disabled = state;
        goBtn.classList.toggle("is-loading", state);
    }

    function hint(msg) {
        const old = input.placeholder;
        input.placeholder = msg;
        setTimeout(() => (input.placeholder = old), 1400);
        input.focus();
    }

    tabShorten.addEventListener("click", () => setMode("shorten"));
    tabExpand .addEventListener("click", () => setMode("lookup"));
    goBtn.addEventListener("click", send);
    input.addEventListener("keydown", (e) => { if (e.key === "Enter") send(); });

    setMode("shorten");
})();
