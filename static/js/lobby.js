import {LitElement, html, css} from "./libs/lit.min.js";
import { client } from "./voidsent.js";
import "./newgame.js";

class VoidLobby extends LitElement {
	constructor() {
		super();

		this.lobby = [];
		client.addEventListener("lobby", (lobby) => {
			this.lobby = lobby
		})
	}

	static styles = css`
		:host {
			background: rgba(0, 0, 0, 0.8);
			border: 5px solid white;
			display: flex;
			flex-direction: column;
			height: 90vh;
			border-radius: 20px;
			padding: 20px;
			color: white;
		}
		
		h1 {
			margin: 0;
			padding: 0;
			text-align: center;
		}
		
		.lobby {
			flex-grow: 1;
		}
		
		.row {
			display: flex;
		}
		
	 	.room {
	 		flex-grow: 1;
		 }

		.players {
			width: 100px;
		}

		.join {
			width: 100px;
		}
	`

	join(ev) {
		console.log(ev);
	}

	row({name, players}) {
	return html`
		<div class="row">
			<div class="room">${name}</div>
			<div class="players">${players}</div>
			<div class="join"><button @click="${this.join}">Beitreten</button></div>
		</div>
	`}

	showNewGameDialog() {
		//this.renderRoot?.querySelector("#new-game-dialog").open();
		client.newGame("Demo")
	}

	render() {
		return html`
			<div class="lobby">
				<h1>Unter Gridanias Helden ...</h1>
				<div class="row">
					<div class="room">Raum</div>
					<div class="players">Spieler</div>
					<div class="join"></div>
				</div>
				${this.lobby.map(this.row)}
			</div>
			<button @click="${this.showNewGameDialog}">New Game!</button>
			<void-new-game id="new-game-dialog"></void-new-game>
		`
	}
}

customElements.define("void-lobby", VoidLobby)