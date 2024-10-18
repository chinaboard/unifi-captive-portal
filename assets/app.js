const config = {
    firstPrompt: {
        role: "system",
        content:
            "**Location**: Home manager, responsible for interactive dialogue when visitors connect to WiFi hotspots.\n" +
            "\n" +
            "**Ability**:\n" +
            "- Use " + navigator.language + "language for the conversation.\n" +
            "- Ability to interact on nonsensical topics.\n" +
            "- After chatting about a few topics, ask what the cat in the house is called.\n" +
            "- Each speech should not exceed 50 words.\n" +
            "\n" +
            "**Behavior**:\n" +
            "- You say hello.\n"
    },
    maxChatLength: 15,
    matchKeywordList: ["miaomiao","meow"],
}

const elements = {
    line: document.getElementById("line"),
    box: document.getElementById("box"),
    loader: document.getElementById("loader")
};

const messages = [];

window.onSend = onSend;
elements.line.addEventListener("keydown", (e) => {
    if (e.key.toLowerCase() === "enter") {
        e.preventDefault();
        onSend();
    }
});

function onSend() {
    const value = (elements.line.value || elements.line.innerText).trim();
    if (!value) return;
    elements.line.value = "";
    elements.line.innerText = "";
    addItem("user", value);
    postLine(value);
}

function addItem(type, content) {
    const request = document.createElement("div");
    request.className = type;
    request.innerText = content;
    elements.box.appendChild(request);
    window.scrollTo({top: document.body.scrollHeight, behavior: "auto"});
    elements.line.focus();
    return request;
}

function postLine(text) {
    for (const item of config.matchKeywordList) {
        if (text.includes(item)) {
            text = "You should say Bingo!";
            break;
        }
    }
    messages.push({role: "user", content: text});
    if (messages.length > config.maxChatLength) {
        messages.splice(2, messages.length - 3);
    }
    chat(messages);
}

function chat(msg) {
    const assistantElem = addItem('', '');
    const baseUrl = window.location.origin;
    send(`${baseUrl}/chat`, {
        model: config.model,
        messages: msg,
        stream: true,
        temperature: config.temperature,
    }, (data) => {
        const msg = data.choices[0].delta || data.choices[0].message || {};
        assistantElem.className = 'assistant';
        assistantElem.innerText += msg.content || "";
    }, () => onSuccessed(assistantElem));
}

function onSuccessed(assistantElem) {
    const msg = assistantElem.innerText;
    messages.push({role: "assistant", content: msg});
    if (msg.includes("Bingo")) {
        const params = new URLSearchParams(window.location.search);
        const id = params.get('id');
        const ap = params.get('ap');
        const url = params.get('url');
        const baseUrl = window.location.origin;
        window.location.href = `${baseUrl}/auth?&id=${id}&ap=${ap}&url=${url}`;
    }
}

function send(reqUrl, body, onMessage, scussionCall) {
    elements.loader.hidden = false;
    const onError = (data) => {
        elements.loader.hidden = true;
        addItem("system", `Network error: ${data}`);
    };
    const source = new SSE(reqUrl, {
        headers: {
            "Content-Type": "application/json",
        },
        method: "POST",
        payload: JSON.stringify(body),
    });

    source.addEventListener("message", (e) => {
        if (e.data === "[DONE]") {
            elements.loader.hidden = true;
            scussionCall();
        } else {
            try {
                onMessage(JSON.parse(e.data));
            } catch (error) {
                onError(error);
            }
        }
    });

    source.addEventListener("error", (e) => {
        onError(e.data);
    });

    source.stream();
}

function init() {
    window.scrollTo(0, document.body.clientHeight);
    elements.box.innerHTML = '';
    messages.length = 0;
    messages.push(config.firstPrompt);
    chat(messages);
}

init();