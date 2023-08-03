import { LitElement, html, css } from "./libs/lit.min.js";
import "./toggle.js"
import { client } from "./voidsent.js";

class NewGame extends LitElement {
	static styles = css`
		h2 {
			margin: 0;
			padding: 0;
		}

		.dialog-content {
			display: flex;
			flex-direction: column;
			margin: 0;
			padding: 0;
			gap:10px;
		}

		.name-section, .password-section {
			display: flex;
			flex-direction: column;
		}

		.roles {
			display: flex;
			flex-direction: column;
		}

		.button-section {
			display: flex;
			gap: 10px;
		}

		.button-section button {
			flex-grow: 1;
		}
	`
	
	constructor() {
		super();

		this.selected = {
			witch: false,
			fortuneTeller: false,
			hunter: false,
			sheriff: false,
			thief: false,
			cupid: false,
			littleGirl: false,
		}

		this.enabled = {
			witch: false,
			fortuneTeller: false,
			hunter: false,
			sheriff: false,
			thief: false,
			cupid: false,
			littleGirl: false,
		}
	}

	toggleRole(role) {
		return () => {
			console.log("role", role);
			this.selected[role] = !this.selected[role]

			return false;
		}
	}

	open() {
		this.renderRoot?.querySelector("dialog").showModal();
	}

	close() {
		this.renderRoot?.querySelector("dialog").close();
	}

	createGame() {
		let roles = [];
		for (let role in this.selected) {
			if (this.selected[role]) {
				roles.push(role);
			}
		}
		client.newVoidGame({
			name: this.renderRoot?.querySelector("#game-name").value,
			password: this.renderRoot?.querySelector("#game-password").value,
			roles: roles,
		});
	}

	render() {
		return html`
			<dialog>
				<div class="dialog-content">
					<div class="dialog-header">
						<h2>Neues Spiel erstellen</h2>
					</div>

					<div class="name-section">
						<label for="game-name">Name des Spiels</label>
						<input id="game-name" placeholder="Name des Spiels" type="text">
					</div>

					<div class="password-section">
						<label for="game-password">Passwort (Optional)</label>
						<input id="game-password" placeholder="Passwort (Optional)" type="text">
					</div>

					<div class="roles">
						<void-toggle text="Die Hexe" ?enabled="${this.enabled.witch}" ?checked="${this.selected.witch}" @change="${this.toggleRole("witch")}"></void-toggle>
						<void-toggle text="Die Seherin" ?enabled="${this.enabled.fortuneTeller}" ?checked="${this.selected.fortuneTeller}" @change=" ${this.toggleRole("fortuneTeller")}
						"></void-toggle>
						<void-toggle text="Der Jäger" ?enabled="${this.enabled.hunter}" ?checked="${this.selected.hunter}" @change="${this.toggleRole("hunter")}
						"></void-toggle>
						<void-toggle text="Der Hauptmann" ?enabled="${this.enabled.sheriff}" ?checked="${this.selected.sheriff}" @change="${this.toggleRole("sheriff")}
						"></void-toggle>
						<void-toggle text="Der Dieb" ?enabled="${this.enabled.thief}" ?checked="${this.selected.thief}" @change="${this.toggleRole("thief")}"></void-toggle>
						<void-toggle text="Amor" ?enabled="${this.enabled.cupid}" ?checked="${this.selected.cupid}" @change="${this.toggleRole("cupid")}
						"></void-toggle>
						<!-- <void-toggle text="Das kleine Mädchen" ?enabled="${this.enabled.littleGirl}" ?checked="${this.selected.littleGirl}" @change="${this.toggleRole("littleGirl")}"></void-toggle> -->
					</div>

					<div class="button-section">
						<button id="cancel-create-game" class="button" @click="${this.close}">Abbrechen</button>
						<button id="create-game" class="button" @click="${this.createGame}">Auf gehts!</button>
					</div>
				</div>
			</dialog>
		`;
	}
}

customElements.define("void-new-game", NewGame);