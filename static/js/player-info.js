import {LitElement, html, css} from "./libs/lit.min.js";

class PlayerInfo extends LitElement {
	constructor() {
		super();
	}

	static styles = css`
		:host {
			background: rgba(0, 0, 0, 0.8);
			border-radius: 20px;
			padding: 20px;
			color: white;
			border: 5px solid white;
		}
	`

	render() {
		return html`
			<div>PLAYER-INFO</div>
		`
	}
}

customElements.define("void-player-info", PlayerInfo);