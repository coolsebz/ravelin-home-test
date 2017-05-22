package handlers

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

// GetSession writes a json message with a random number to be used as
// a session identifier into the response stream and is by no means a
// complete session manager.
// TODO: (for becoming a full fledged network manager)
// - Keep session id unique.
// - Enforce having one session for every user.
// - Better defined session storage (db, file, in memory, etc)
// - Deal with expired sessions.
func GetSession(w http.ResponseWriter, r *http.Request) {
	newSession := rand.Int31()

	// as the session data grows, the string will become a struct
	io.WriteString(w, fmt.Sprintf("{ \"sessionId\": \"%d\" }", newSession))
	log.Println(fmt.Sprintf("emitted a new sessionId: %d", newSession))
}
