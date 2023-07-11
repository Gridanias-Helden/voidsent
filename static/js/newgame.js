import { LitElement, html } from "./libs/lit.min.js";
import "./toggle.js"

class NewGame extends LitElement {
	constructor() {
		super();

		this.selected = {
			witch: true,
			fortuneTeller: true,
			hunter: false,
			sheriff: false,
			thief: false,
			cupid: false,
			littleGirl: false,
		}

		this.enabled = {
			witch: true,
			fortuneTeller: true,
			hunter: false,
			sheriff: false,
			thief: false,
			cupid: false,
			littleGirl: false,
		}
	}

	toggleRole(role) {
		this.selected[role] = !this.selected[role]

		return false;
	}

	open() {
		this.renderRoot?.querySelector("dialog").showModal();
	}

	close() {
		this.renderRoot?.querySelector("dialog").close();
	}

	createGame() {
	}

	render() {
		return html`
			<dialog>
				<div class="dialog-content">
					<div class="dialog-header">
						<h2>Neues Spiel erstellen</h2>
					</div>
					<div class="dialog-body">
						<label for="game-name">Name des Spiels</label>
						<input id="game-name" placeholder="Name des Spiels" type="text">
					</div>
					<div class="dialog-footer">
						<button id="cancel-create-game" class="button" @click="${this.close}">Abbrechen
						</button>
						<button id="create-game" class="button" @click="${this.createGame}">Auf gehts!</button>
					</div>
				</div>

				<div>
					<void-toggle text="Die Hexe" ?enabled="${this.enabled.witch}" ?checked="${this.selected.witch}" @change="${this.toggleRole("witch")}"></void-toggle>
					<void-toggle text="Die Seherin" ?enabled="${this.enabled.fortuneTeller}" ?checked="${this.selected.fortuneTeller}" @change=" ${this.toggleRole("fortuneTeller")}
					"></void-toggle>
					<void-toggle text="Der Jäger" ?enabled="${this.enabled.hunter}" ?checked="${this.selected.hunter}" @change="${this.toggleRole("hunter")}
					"></void-toggle>
					<void-toggle text="Der Wachmann" ?enabled="${this.enabled.sheriff}" ?checked="${this.selected.sheriff}" @change="${this.toggleRole("sheriff")}
					"></void-toggle>
					<void-toggle text="Der Dieb" ?enabled="${this.enabled.thief}" ?checked="${this.selected.thief}" @change="${this.toggleRole("thief")}"></void-toggle>
					<void-toggle text="Amor" ?enabled="${this.enabled.cupid}" ?checked="${this.selected.cupid}" @change="${this.toggleRole("cupid")}
					"></void-toggle>
					<void-toggle text="Das kleine Mädchen" ?enabled="${this.enabled.littleGirl}" ?checked="${this.selected.littleGirl}" @change="${this.toggleRole("littleGirl")}"></void-toggle>
				</div>
			</dialog>
		`;
	}
}

customElements.define("void-new-game", NewGame);