import { LitElement, html, css } from "./libs/lit.min.js";
import "./lobby.js";
import "./chat.js";
import "./player-info.js";

class VoidGame extends LitElement {
	constructor() {
		super();
	}

	static styles = css`
		:host {
			display: block;
		}

		.right {
			display: flex;
			flex-wrap: wrap;
			width: 100%;
			gap: 20px;
		}
		
		.left {
			flex-grow: 1;
			display: flex;
			flex-direction: column;
			gap: 20px;
			/* max-width: 30%; */
		}
		
		void-player-info{
			flex-grow: 1;
		}
		
		void-chat {
			flex-grow: 2;
		}
		
		void-lobby {
			flex-grow: 2;
		}
	`

	render() {
		return html`
			<div class="right">
				<void-lobby></void-lobby>
				<div class="left">
					<void-player-info></void-player-info>
					<void-chat></void-chat>
				</div>
			</div>
		`
	}
}

customElements.define("void-game", VoidGame);