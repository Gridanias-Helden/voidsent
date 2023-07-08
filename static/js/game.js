export class Voidsent {
    constructor() {
        this.ws = new WebSocket(`ws://${window.location.host}/ws`);

        this.ws.onopen = this.onopen.bind(this);
        this.ws.onerror = this.onerror.bind(this);
        this.ws.onmessage = this.onmessage.bind(this)
        this.ws.onclose = this.onclose.bind(this)
    }

    onmessage(ev) {
        console.log("data", ev.data);
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