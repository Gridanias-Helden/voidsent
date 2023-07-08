export class Voidsent {
    constructor() {
        this.ws = new WebSocket(`ws://${window.location.host}/ws`);

        this.ws.onopen = this.onopen.bind(this);
        this.ws.onerror = this.onerror.bind(this);
        this.ws.onmessage = this.onmessage.bind(this);
        this.ws.onclose = this.onclose.bind(this);

        this.rooms = [];
    }

    onmessage(ev) {
        console.log("data", ev.data);
        let data = JSON.parse(ev.data);

        if (data.type === "lobby") {
            this.rooms = data.lobby;
            m.redraw();
        }
    }

    onerror(ev) {
        console.log("err", ev);
    }

    onclose(ev) {
        console.log("connection closed", ev);
    }

    onopen(ev) {
        console.log("open", ev);
    }

    newGame(name) {
        this.ws.send(JSON.stringify({
            type: "newGame",
            newGame: {
                name,
            }
        }))
    }
}