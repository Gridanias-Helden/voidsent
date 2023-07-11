export default {
	view() {
		return m("div", {},
			m("void-toggle", { text: "Hexe", checked: true }),
			m("void-toggle", { text: "Seherin", checked: true }),
			m("void-toggle", { text: "Jäger", checked: true }),
			m("void-toggle", { text: "Hauptmann", checked: true }),
			m("void-toggle", { text: "Dieb", checked: true }),
			m("void-toggle", { text: "Amor", checked: true }),
			m("void-toggle", { text: "kleinem Mädchen", checked: true }),
		)
	}
}