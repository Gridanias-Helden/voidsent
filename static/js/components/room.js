export class Room {
	ready(self, player) {
		console.log("ready", self, player);

		if (self.id === player.id && self.status === "ready") {
			return m("button", "Bereit!")
		} else if (self.id === player.id) {
			return m("button", "Bereit?")
		}

		return m("span", player.status === "ready" ? "Bereit" : "")
	}
	view({attrs: {room}}) {
		console.log("room", room)

		return m("div", [
			m("h1", "Unter Gridanias Helden ..."),
			m(".row", [
				m("div"),
				m("div", "Spieler"),
				m("div", ""),
			]),
			...room.players.map((player) => {
				return m(".row", [
					m("div", m("img", { src: player.avatar})),
					m("div", player.name),
					m("div", this.ready(room.self, player)),
				])
			}),
		])
	}
}
