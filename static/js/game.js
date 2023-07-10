import {Voidsent} from "/js/voidsent.js";
import {Lobby} from "/js/components/lobby.js";
import {Room} from "/js/components/room.js";

let client = new Voidsent();

document.querySelector('#new-game').addEventListener('click', function () {
	document.querySelector('#new-game-dialog').showModal();
})

document.querySelector("#create-game").addEventListener("click", function () {
	client.newGame(document.querySelector("#game-name").value)
})

m.mount(document.querySelector('#lobby'), {
	view: (vnode) => {
		switch (client.page) {
			case "lobby":
				return m(Lobby, {
					rooms: client.lobby
				});

			case "room":
				return m(Room, {
					room: client.room
				});
		}
	}
})
