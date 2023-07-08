export class Lobby {
	view({attrs: {rooms}}) {
		console.log(rooms);

		return m("div", [
			m(".row", [
				m("div", "Raum"),
				m("div", "Spieler"),
				m("div")
			]),
			...rooms.map((room) => {
				return m(".row", [
					m("div", room.name),
					m("div", room.players.length),
					m("div", m("button", "Los!"))
				])
			}),
		])
	}
}
