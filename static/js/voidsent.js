class VoidClient {
	constructor() {
		let proto = window.location.protocol === "https:" ? "wss" : "ws";
		this.ws = new WebSocket(`${proto}://${window.location.host}/ws`);

		this.ws.onopen = this.onopen.bind(this);
		this.ws.onerror = this.onerror.bind(this);
		this.ws.onmessage = this.onmessage.bind(this);
		this.ws.onclose = this.onclose.bind(this);

		this.lobbyRcv = [];
		this.joinRcv = [];
		this.leaveRcv = [];
		this.chatAllRcv = [];
		this.chatWhisperRcv = [];
		this.sessionRcv = [];
	}

	onmessage(ev) {
		if (!ev.data instanceof Blob) {
			console.log("no binary frame");
			return;
		}

		// binary frame
		const reader = new FileReader();
		reader.onload = () => {
			/** @type {ArrayBuffer} */
			const res = reader.result;
			const decoder = new TextDecoder("utf-8");
			const view = new DataView(res);
			const type = decoder.decode(res.slice(0, 4));
			const timestamp = new Date(Number(view.getBigInt64(4, false)));
			let roomLen = 0;
			let room = "";
			let userLen = 0;
			let user = "";
			let msg = "";
			let avatarLen = 0;
			let avatar = "";
			let fromLen = 0;
			let from = "";
			let toLen = 0;
			let to = "";

			switch (type) {
				case "chat":
					const action = decoder.decode(res.slice(12, 14));

					switch (action) {
						case "jo":
							roomLen = view.getUint8(14);
							room = decoder.decode(res.slice(15, 15 + roomLen));
							userLen = view.getUint8(15 + roomLen);
							user = decoder.decode(res.slice(16 + roomLen, 16 + roomLen + userLen));
							for (let cb of this.joinRcv) {
								cb({
									room: room,
									user: user,
									time: timestamp,
								});
							}
							break;

						case "lv":
							roomLen = view.getUint8(14);
							room = decoder.decode(res.slice(15, 15 + roomLen));
							userLen = view.getUint8(15 + roomLen);
							user = decoder.decode(res.slice(16 + roomLen, 16 + roomLen + userLen));
							for (let cb of this.leaveRcv) {
								cb({
									room: room,
									user: user,
									time: timestamp,
								});
							}
							break;

						case "sa":
							userLen = view.getUint8(14);
							user = decoder.decode(res.slice(15, 15 + userLen));
							msg = decoder.decode(res.slice(15 + userLen));
							for (let cb of this.chatAllRcv) {
								cb({
									user: user,
									msg: msg,
									time: timestamp,
								});
							}
							break;

						case "wh":
							fromLen = view.getUint8(14);
							from = decoder.decode(res.slice(15, 15 + fromLen));
							toLen = view.getUint8(15 + fromLen);
							to = decoder.decode(res.slice(16 + fromLen, 16 + fromLen + toLen));
							msg = decoder.decode(res.slice(16 + fromLen + toLen));
							console.log("got whisper", timestamp, from, to, msg, fromLen, toLen);
							for (let cb of this.chatWhisperRcv) {
								cb({
									from: from,
									to: to,
									msg: msg,
									time: timestamp,
								});
							}
							break;
					}
					break;

				case "sess":
					userLen = view.getUint8(12);
					user = decoder.decode(res.slice(13, 13 + userLen));
					avatarLen = view.getUint8(13 + userLen);
					avatar = decoder.decode(res.slice(14 + userLen, 14 + userLen + avatarLen));
					console.log("got session", timestamp, user, avatar);

					for (let cb of this.sessionRcv) {
						cb({
							user: user,
							avatar: avatar,
						});
					}
					break;
			}
		}
		reader.readAsArrayBuffer(ev.data);

		// 	case "room:leave":
		// 		console.log("leave room");
		// 		for (let cb of this.leaveRcv) {
		// 			cb(data.body);
		// 		}
		// 		break;
		//
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

	/**
	 * Send a "new game" to the server
	 * @param {string} name
	 * @param {string} password
	 * @param {string[]} roles
	 */
	newVoidGame({
				name,
				password,
				roles
			}) {
		const encoder = new TextEncoder();
		const voidHeader = encoder.encode("void");
		const createAction = encoder.encode("cr");
		const nameLen = encoder.encode(name).length;
		const passwordLen = encoder.encode(password).length;
		const roleMap = {
			witch: 1,
			fortuneTeller: 2,
			hunter: 4,
			sheriff: 8,
			thief: 16,
			cupid: 32,
			littleGirl: 64,
		}
		let roleMask = roles.reduce((acc, role) => {
			return acc | roleMap[role];
		}, 0);

		const buffer = new Uint8Array([
				...voidHeader,
				...createAction,
				nameLen,
				...encoder.encode(name),
				passwordLen,
				...encoder.encode(password),
				roleMask,
			]
		)

		this.ws.send(buffer);
	}

	addEventListener(event, cb) {
		if (!cb) {
			return;
		}

		switch (event) {
			case "lobby":
				this.lobbyRcv.push(cb);
				break;

			case "chat:all":
				this.chatAllRcv.push(cb);
				break;

			case "chat:whisper":
				this.chatWhisperRcv.push(cb);
				break;

			case "room:join":
				this.joinRcv.push(cb);
				break;

			case "room:leave":
				this.leaveRcv.push(cb);
				break;

			case "session":
				this.sessionRcv.push(cb);
				break;
		}
	}

	chat({
			 msg,
			 to
		 }) {
		const encoder = new TextEncoder();
		const chatHeader = encoder.encode("chat");
		const msgContent = encoder.encode(msg);
		const content = [
			...chatHeader,
		];

		if (to) {
			// Whisper to a user
			content.push(...encoder.encode("wh"));
			const toContent = encoder.encode(to);
			content.push(toContent.length)
			content.push(...toContent);
		} else {
			// Say to the room
			content.push(...encoder.encode("sa"));
		}

		content.push(
			...msgContent,
		);

		const buffer = new Uint8Array(content);
		try {
			this.ws.send(buffer);
		} catch (e) {
			console.log("error sending chat", e);
		}

		console.log("sending chat", content);
	}
}

export const client = new VoidClient();