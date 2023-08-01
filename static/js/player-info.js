import {LitElement, html, css} from "./libs/lit.min.js";
import { client } from "./voidsent.js";

class PlayerInfo extends LitElement {
	static properties = {
		playerInfo: {state: true},
	}

	constructor() {
		super();

		this.playerInfo = { name: "<Unbekannnt>", avatar: "" };
		client.addEventListener("session", (playerInfo) => {
			this.playerInfo = playerInfo
		})
	}

	static styles = css`
		:host {
			background: rgba(0, 0, 0, 0.8);
			border-radius: 20px;
			padding: 20px;
			color: white;
			border: 5px solid white;
		}
		
		div {
			display: flex;
			flex-direction: column;
			align-items: center;
		}
	`

	render() {
		return html`
			<div>
				<img src="${this.playerInfo.avatar}" alt="" />
				<span>${this.playerInfo.user}</span>
			</div>
		`
	}
}

customElements.define("void-player-info", PlayerInfo);