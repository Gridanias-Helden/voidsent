import { LitElement, html, css } from "./libs/lit.min.js";
import { client } from "./voidsent.js"

class Chat extends LitElement {
	static properties = {
		history: {state: true},
	}

	constructor() {
		super();

		this.history = [];

		client.addEventListener("room:join", ({ name, room }) => {
			console.log(`${name} joined ${room}!`);
			this.history = [ ...this.history, html`<b>${name}</b> ist ${room} beigetreten.`]
		})

		client.addEventListener("room:leave", ({ name, room }) => {
			console.log(`${name} left ${room}!`);
			this.history = [ ...this.history, html`<b>${name}</b> hat ${room} verlassen.`]
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

	render() {
		return html`
			<div class="root">
				<div class="history">
					${this.history.map((entry) => html`<div>${entry}</div>`)}
				</div>
				<div class="input">
					<input type="text" placeholder="Nachricht eingeben...">
					<button>Senden</button>
				</div>
			</div>
		`
	}
}

customElements.define("void-chat", Chat);