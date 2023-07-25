import { LitElement, html, css } from "./libs/lit.min.js";
import { client } from "./voidsent.js"

class Chat extends LitElement {
	static properties = {
		history: {state: true},
		msg: {state: true},
		self: {state: true},
	}

	time(t) {
		return new Date(t).toLocaleTimeString();
	}

	constructor() {
		super();

		this.history = [];

		this.self = { name: "", avatar: "" };
		client.addEventListener("session", (playerInfo) => {
			this.self = playerInfo
		})

		client.addEventListener("room:join", ({ time, from, room }) => {
			console.log(`${name} joined ${room}!`);
			this.history = [ ...this.history, html`[${this.time(time)}] <b>${from}</b> ist ${room} beigetreten.`]
		})

		client.addEventListener("room:leave", ({ time, from, room }) => {
			console.log(`${from} left ${room}!`);
			this.history = [ ...this.history, html`[${this.time(time)}] <b>${from}</b> hat ${room} verlassen.`]
			this.history.sort((a, b) => {
				if (a.time > b.time) {
					return -1;
				} else if (a.time === b.time) {
					return 0;
				} else {
					return 1
				}
			})
		})

		client.addEventListener("chat:all", ({ time, from, msg }) => {
			console.log(`${from}: ${msg}`);
			this.history = [ ...this.history, html`[${this.time(time)}] <b>${from}</b>: ${msg}`]
			this.history.sort((a, b) => {
				if (a.time > b.time) {
					return -1;
				} else if (a.time === b.time) {
					return 0;
				} else {
					return 1
				}
			})
		})

		client.addEventListener("chat:whisper", ({ time, to, from, msg }) => {
			console.log(`${name}: ${msg}`);
			let newEntry = html`[${this.time(time)}] <i><b>${from} flüstert dir zu</b>: ${msg}</i>`
			if (from === this.self.from) {
				newEntry = html`[${this.time(time)}] <i><b>Du flüsterst ${to} zu</b>: ${msg}</i>`
			}
			this.history = [ ...this.history, newEntry]
			this.history.sort((a, b) => {
				if (a.time > b.time) {
					return -1;
				} else if (a.time === b.time) {
					return 0;
				} else {
					return 1
				}
			})
		})
	}

	static styles = css`
		:host {
			background: rgba(0, 0, 0, 0.8);
			border-radius: 20px;
			padding: 20px;
			color: white;
			border: 5px solid white;
			height: 100%;
		}
		
		.root {
			display: flex;
			flex-direction: column;
			height: 100%;
		}
		
		.history {
			background-color: white;
			flex-grow: 1;
			color: black;
		}
		
		.input {
			background-color: blue;
			flex-grow: 0;
			display: flex;
		}
		
		.input > input {
			flex-grow: 3;
			border: 1x solid black;
			border-radius: 5px;
			padding: 5px;
			font-size: 1.2em;
		}
		
		.input > button {
			flex-grow: 1;
			border: 1x solid black;
			border-radius: 5px;
			padding: 5px;
			font-size: 1.2em;
			background-color: cyan;
			color: black;
		}
	`

	send() {
		if (!this.msg) {
			return;
		}

		if (this.msg.startsWith("/w")) {
			let [_, to, ...msg] = this.msg.split(" ");
			msg = msg.join(" ");
			client.ws.send(JSON.stringify({
				"type": "chat",
				"body": {
					"to": to,
					"msg": msg,
				}
			}))
			let input = this.renderRoot?.querySelector('#chat-input');
			console.log(input);
			input.value = "";
			this.msg = "";
			return;
		}

		client.ws.send(JSON.stringify({
			"type": "chat",
			"body": {
				//"to": "Da",
				"msg": this.msg,
			}
		}))

		let input = this.renderRoot?.querySelector('#chat-input');
		console.log(input);
		input.value = "";
		this.msg = "";
	}

	changeMsg(ev) {
		this.msg = ev.target.value;
	}

	render() {
		return html`
			<div class="root">
				<div class="history">
					${this.history.map((entry) => html`<div>${entry}</div>`)}
				</div>
				<div class="input">
					<input id="chat-input" type="text" placeholder="Nachricht eingeben..."  @input="${this.changeMsg}">
					<button id="chat-send" @click="${this.send}" @disabled="${!(this.msg === "")}">Senden</button>
				</div>
			</div>
		`
	}
}

customElements.define("void-chat", Chat);