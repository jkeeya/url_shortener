/* script.js
   Две вкладки:
     - create:  сгенерировать шорт по оригинальному URL
     - find:    найти шорт по оригинальному URL
   Ожидаемые эндпоинты:
     POST /create { "url": "<оригинальный URL>" } -> "https://sho.rt/abc123"
     GET  /find_short_link?url=<оригинальный URL> -> "https://sho.rt/abc123"
*/

(function () {
    "use strict";

    // Настройка путей к API
    const ENDPOINTS = {
        create: "/create",
        find:   "/find_short_link",
    };

    // Захват DOM-элементов
    const tabCreate = document.getElementById("tab-shorten");
    const tabFind   = document.getElementById("tab-expand");
    const input     = document.getElementById("input-url");
    const goBtn     = document.getElementById("go-btn");
    const output    = document.getElementById("output-box");

    if (!tabCreate || !tabFind || !input || !goBtn || !output) {
        console.error("Проверь id элементов: tab-shorten, tab-expand, input-url, go-btn, output-box");
        return;
    }

    output.readOnly = true;

    // Текущее состояние вкладки: 'create' | 'find'
    let mode = "create";

    function setMode(nextMode) {
        mode = nextMode;

        if (mode === "create") {
            tabCreate.classList.add("tab-active");
            tabCreate.classList.remove("tab-inactive");
            tabFind.classList.add("tab-inactive");
            tabFind.classList.remove("tab-active");
            input.placeholder = "Вставьте оригинальный URL (создать короткую ссылку)";
        } else {
            tabFind.classList.add("tab-active");
            tabFind.classList.remove("tab-inactive");
            tabCreate.classList.add("tab-inactive");
            tabCreate.classList.remove("tab-active");
            input.placeholder = "Вставьте оригинальный URL (найти короткую ссылку)";
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
            let resp;
            if (mode === "create") {
                // POST /create с JSON-телом
                resp = await fetch(ENDPOINTS.create, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ url: val }),
                });
            } else {
                // GET /find_short_link?url=...
                const q = new URLSearchParams({ url: val });
                resp = await fetch(`${ENDPOINTS.find}?${q.toString()}`);
            }

            if (!resp.ok) throw new Error("Сервер вернул " + resp.status);

            // Поддержим text/plain и JSON
            let text = "";
            const ct = (resp.headers.get("Content-Type") || "").toLowerCase();

            if (ct.includes("application/json")) {
                const data = await resp.json();

                if (data.Exist === false || data.exist === false) {
                    text = "Этой ссылки ещё нет в базе.";
                } else {
                    // подхватываем поле ShortLink или short, или просто JSON
                    text = data.ShortLink || data.short_link || data.result || data.url || JSON.stringify(data);
                }
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

    tabCreate.addEventListener("click", () => setMode("create"));
    tabFind.addEventListener("click", () => setMode("find"));
    goBtn.addEventListener("click", send);
    input.addEventListener("keydown", (e) => { if (e.key === "Enter") send(); });

    setMode("create");
})();
