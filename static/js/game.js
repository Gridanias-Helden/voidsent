import {Voidsent} from "/js/voidsent.js";
import {Lobby} from "/js/components/lobby.js";

let client = new Voidsent();

document.querySelector('#new-game').addEventListener('click', function () {
	document.querySelector('#new-game-dialog').showModal();
})

document.querySelector("#create-game").addEventListener("click", function () {
	client.newGame(document.querySelector("#game-name").value)
})

m.mount(document.querySelector('#lobby'), {
	view: (vnode) => {
		return m(Lobby, {
			rooms: client.rooms
		})
	}
})