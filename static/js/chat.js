import { LitElement, html, css } from "./libs/lit.min.js";
import { client } from "./voidsent.js"

class Chat extends LitElement {
	static properties = {
		history: {state: true},
	}

	constructor() {
		super();

		this.history = [];

		client.addEventListener("join", (name) => {
			console.log(`${name} joined!`);
			this.history = [ ...this.history, html`<b>${name}</b> ist der Lobby beigetreten.`]
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
		}
	`

	render() {
		return html`
			<div class="root">
				<div class="history">
					${this.history.map((entry) => html`<div>${entry}</div>`)}
				</div>
				<div class="input">INPUT</div>
			</div>
		`
	}
}

customElements.define("void-chat", Chat);