import { LitElement, html, css } from "./libs/lit.min.js";
import { client } from "./voidsent.js"

class Chat extends LitElement {
	static properties = {
		history: {state: true},
		msg: {state: true},
	}

	time(t) {
		return new Date(t).toLocaleTimeString();
	}

	constructor() {
		super();

		this.history = [];

		client.addEventListener("room:join", ({ time, name, room }) => {
			console.log(`${name} joined ${room}!`);
			this.history = [ ...this.history, html`[${this.time(time)}] <b>${name}</b> ist ${room} beigetreten.`]
		})

		client.addEventListener("room:leave", ({ time, name, room }) => {
			console.log(`${name} left ${room}!`);
			this.history = [ ...this.history, html`[${this.time(time)}] <b>${name}</b> hat ${room} verlassen.`]
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

		client.addEventListener("chat:all", ({ time, name, msg }) => {
			console.log(`${name}: ${msg}`);
			this.history = [ ...this.history, html`[${this.time(time)}] <b>${name}</b>: ${msg}`]
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

		client.addEventListener("chat:whisper", ({ time, name, msg }) => {
			console.log(`${name}: ${msg}`);
			this.history = [ ...this.history, html`[${this.time(time)}] <i><b>${name} (flüstert dir zu)</b>: ${msg}</i>`]
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