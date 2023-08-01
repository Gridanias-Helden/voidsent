package websocket

import (
	"bytes"
	"crypto/rand"
	"net/http"
	"strings"
	//"strings"
	"time"

	"github.com/oklog/ulid"
	//"github.com/oklog/ulid"
	"github.com/olahol/melody"

	"github.com/gridanias-helden/voidsent/pkg/middleware"
	"github.com/gridanias-helden/voidsent/pkg/models"
	"github.com/gridanias-helden/voidsent/pkg/services"
	"github.com/gridanias-helden/voidsent/pkg/services/games/voidsent"

	//"github.com/gridanias-helden/voidsent/pkg/services/games/voidsent"
	"github.com/gridanias-helden/voidsent/pkg/storage"
	"github.com/gridanias-helden/voidsent/pkg/utils"
)

// Required for Chat Messages
var (
	ChatHeader    = []byte("chat")
	JoinAction    = []byte("jo")
	LeaveAction   = []byte("lv")
	SayAction     = []byte("sa")
	WhisperAction = []byte("wh")
)

// Required for Session Messages
var (
	SessionHeader = []byte("sess")
)

type WebSocket struct {
	sessions storage.Sessions
	broker   *services.Broker
	mel      *melody.Melody
}

func New(sessions storage.Sessions, broker *services.Broker, mel *melody.Melody) *WebSocket {
	ws := &WebSocket{
		sessions: sessions,
		broker:   broker,
		mel:      mel,
	}

	mel.HandleConnect(ws.Connect)
	mel.HandleDisconnect(ws.Disconnect)
	mel.HandleMessageBinary(ws.Message)

	return ws
}

func (ws *WebSocket) HTTPRequest(w http.ResponseWriter, r *http.Request) {
	_ = ws.mel.HandleRequest(w, r)
}

func (ws *WebSocket) Connect(s *melody.Session) {
	session, ok := s.Request.Context().Value(middleware.SessionKey).(models.Session)
	if !ok {
		_ = s.CloseWithMsg([]byte("no session found"))
		return
	}

	s.Set("session", session)
	s.Set("room", "lobby")

	time.Sleep(50 * time.Millisecond)

	// ws.broker.Send(session.ID, "lobby", "join", s)
	ws.Join(s, session, "lobby")
	ws.Session(s, session)
}

func (ws *WebSocket) Disconnect(s *melody.Session) {
	session, ok := s.MustGet("session").(models.Session)
	if !ok {
		return
	}

	room, ok := s.MustGet("room").(string)
	if !ok {
		return
	}

	ws.broker.Send(session.ID, room, "leave", s)

	ws.Leave(s, session, room)
}

func (ws *WebSocket) Message(s *melody.Session, msg []byte) {
	session, ok := s.MustGet("session").(models.Session)
	if !ok {
		return
	}

	room, ok := s.MustGet("room").(string)
	if !ok {
		return
	}

	t := string(msg[:4])
	switch t {
	case "chat":
		ws.Chat(s, session, room, msg[4:])

	case "void":
		ws.Voidsent(s, session, msg[4:])
	}
}

func (ws *WebSocket) Chat(melSess *melody.Session, voidSess models.Session, room string, msg []byte) {
	// Layout of an incoming chat message:
	// 0-1: action (sa = say, wh = whisper)
	// 2: length of receiving username ("ur") (only for whisper)
	// 2-ur: receiving username (only for whisper)
	// rest: chat message

	// Layout of an outgoing chat message:
	// 0-3: header ("chat")
	// 4-11: timestamp
	// 12-13: action (sa = say, wh = whisper)
	// 14: length of sending username ("us")
	// 15-us: sending username
	// 15+us: length of receiving username ("ur") (only for whisper)
	// 15+us-15+us+ur: receiving username (only for whisper)
	// rest: chat message

	action := string(msg[:2])
	newMsg := bytes.NewBuffer(ChatHeader)
	newMsg.Write(utils.Int64ToBytes(uint64(time.Now().UnixMilli())))

	switch action {
	case "sa":
		newMsg.WriteString("sa")
		newMsg.WriteByte(byte(len(voidSess.Username)))
		newMsg.WriteString(voidSess.Username)
		newMsg.Write(msg[2:])
		ws.Broadcast(newMsg.Bytes(), room)

	case "wh":
		toLen := msg[2]
		to := string(msg[3 : 3+toLen])
		chatMsg := string(msg[3+toLen:])

		newMsg.WriteString("wh")
		// from
		newMsg.WriteByte(byte(len(voidSess.Username)))
		newMsg.WriteString(voidSess.Username)
		// to
		newMsg.WriteByte(toLen)
		newMsg.WriteString(to)
		newMsg.Write([]byte(chatMsg))
		_ = ws.mel.BroadcastBinaryFilter(newMsg.Bytes(), ws.ToName(to))
		_ = melSess.WriteBinary(newMsg.Bytes())
	}
}

func (ws *WebSocket) Join(melSess *melody.Session, voidSess models.Session, room string) {
	lr := len(room)
	lu := len(voidSess.Username)
	if lr == 0 || lu == 0 || lr > 255 || lu > 255 {
		return
	}

	now := time.Now().UnixMilli()

	buffer := bytes.NewBuffer(ChatHeader)
	buffer.Write(utils.Int64ToBytes(uint64(now)))
	buffer.Write(JoinAction)
	buffer.WriteByte(byte(lr))
	buffer.WriteString(room)
	buffer.WriteByte(byte(lu))
	buffer.WriteString(voidSess.Username)
	msg := buffer.Bytes()

	ws.Broadcast(msg, room)
}

func (ws *WebSocket) Leave(melSess *melody.Session, voidSess models.Session, room string) {
	lr := len(room)
	lu := len(voidSess.Username)
	if lr == 0 || lu == 0 || lr > 255 || lu > 255 {
		return
	}

	now := time.Now().UnixMilli()

	buffer := bytes.NewBuffer(ChatHeader)
	buffer.Write(utils.Int64ToBytes(uint64(now)))
	buffer.Write(LeaveAction)
	buffer.WriteByte(byte(lr))
	buffer.WriteString(room)
	buffer.WriteByte(byte(lu))
	buffer.WriteString(voidSess.Username)
	msg := buffer.Bytes()

	ws.Broadcast(msg, room)
}

func (ws *WebSocket) Session(melSess *melody.Session, voidSess models.Session) {
	buffer := bytes.NewBuffer(SessionHeader)
	buffer.Write(utils.Int64ToBytes(uint64(time.Now().UnixMilli())))
	buffer.WriteByte(byte(len(voidSess.Username)))
	buffer.WriteString(voidSess.Username)
	buffer.WriteByte(byte(len(voidSess.Avatar)))
	buffer.WriteString(voidSess.Avatar)

	msg := buffer.Bytes()

	_ = melSess.WriteBinary(msg)
}

func (ws *WebSocket) ToRoom(room string) func(*melody.Session) bool {
	return func(s *melody.Session) bool {
		roomStr, ok := s.MustGet("room").(string)
		if !ok {
			return false
		}

		return roomStr == room
	}
}

func (ws *WebSocket) ToName(name string) func(*melody.Session) bool {
	return func(s *melody.Session) bool {
		session, ok := s.MustGet("session").(models.Session)
		if !ok {
			return false
		}

		return session.Username == name
	}
}

func (ws *WebSocket) Voidsent(melSess *melody.Session, voidSess models.Session, body []byte) {
	// Layout of an incoming voidsent message:
	// 0-1: action (jo = join, cr = create)
	// 2: length of game name ("gn")
	// 3-gn: game name
	// 3+gn: length of password ("pw")
	// 4+gn-4+gn+pw: password
	// 4+gn+pw: roles (bitmask)
	action := string(body[:2])

	switch action {
	case "jo":
		// Join a game
		//gameNameLen := body[2]
		//gameName := string(body[3 : 3+gameNameLen])
		//passwordLen := body[3+gameNameLen]
		//password := string(body[4+gameNameLen : 4+gameNameLen+passwordLen])
		//roles := body[4+gameNameLen+passwordLen:]
		//

	case "cr":
		// Create a game
		gameNameLen := body[2]
		gameName := string(body[3 : 3+gameNameLen])
		passwordLen := body[3+gameNameLen]
		password := string(body[4+gameNameLen : 4+gameNameLen+passwordLen])
		roles := body[4+gameNameLen+passwordLen]
		if strings.TrimSpace(gameName) == "" {
			return
		}

		id, _ := ulid.New(ulid.Now(), rand.Reader)
		ws.broker.AddService("voidsent:"+id.String(), voidsent.New(ws.broker, "voidsent:"+id.String(), password, melSess, roles))
	}

	//t, ok := body["type"]
	//if !ok {
	//	return
	//}
	//
	//switch t {
	//case "voidsent":
	//	// Create a new Voidsent game
	//	gameID := "voidsent:" + ulid.MustNew(ulid.Now(), rand.Reader).String()
	//	rolesStr, ok := body["roles"]
	//	if !ok {
	//		return
	//	}
	//	roles := strings.Split(rolesStr, ",")
	//	ws.broker.AddService(gameID, voidsent.New(ws.broker, gameID, melSess, roles...))
	//
	//	// Transfer the session to the game
	//	ws.Leave(melSess, voidSess, "lobby")
	//	melSess.Set("room", gameID)
	//}
}

func (ws *WebSocket) Broadcast(msg []byte, room string) {
	_ = ws.mel.BroadcastBinaryFilter(msg, ws.ToRoom(room))
}
