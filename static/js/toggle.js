import { LitElement, html, css } from './libs/lit.min.js';

export class VoidToggle extends LitElement {
	static get properties() {
		return {
			checked: { type: Boolean },
			enabled: { type: Boolean },
			text: { type: String },
		};
	}

	static styles = css`
:host {
	display: inline-block;
	border: 1px solid black;
}

.container {
	background-color: white;
	display: flex;
	width: 300px;
	border: 1px solid white;
	padding: 5px;
}

.text {
	flex-grow: 1;
}

.options {
	width: 50px;
	overflow: hidden;
	display: flex;
}

.option-1 {
	background-color: green;
	color: white;
}

.option-2 {
	background-color: red;
	color: white;
}

.option-1, .option-2 {
	position: relative;
	min-width: 50px;
	text-align: center;
	transition: left 1s;
	left: 0;
}

.container[data-checked] .option-1,
.container[data-checked] .option-2 {
	left: -50px;
}

.container[disabled],
.container[disabled] .option-1,
.container[disabled] .option-2 {
	color: gray;
}
`;

	constructor() {
		super();
	}

	render() {
		console.log(this.checked);

		return html`
		<div class="container" ?disabled="${!this.enabled}" ?data-checked="${!this.checked}" @click="${this.click}">
			<div class="text">${this.text}</div>
			<div class="options">
				<div class="option-1">Ja</div>
				<div class="option-2">Nein</div>
			</div>
		</div>
`;
	}

	click() {
		if (!this.enabled) {
			return
		}

		this.checked = !this.checked; //ev.target.checked;
		
		let newEv = new Event("change");
		newEv.checked = this.checked;
		
		this.dispatchEvent(newEv);
	}
}

customElements.define('void-toggle', VoidToggle);
