package handlers

// TODO:
// - add event type handling
// - print out the event types and their details as they roll in
// - add tests (handling each type of event, bad params, completing a struct)

import (
	"encoding/json"
	"github.com/coolsebz/ravelin-home-test/backend/storage"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int             // Seconds
}

type Dimension struct {
	Width  string
	Height string
}

type ClientEvent struct {

	// mandatory fields to have on any request
	EventType  string `json:"eventType"`
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`

	// form metadata
	TimeTaken    int             `json:"timeTaken"`
	FormId       string          `json:"formId"`
	CopiedFields map[string]bool `json:"copiedFields"`

	// resize event
	FromWidth  json.Number `json:"fromWidth,Number"`
	FromHeight json.Number `json:"fromHeight,Number"`
	ToWidth    json.Number `json:"toWidth,Number"`
	ToHeight   json.Number `json:"toHeight,Number"`

	// submitted event
	Submitted bool `json:"submitted"`
}

var behaviours = map[string]func(ClientEvent){
	"resized":      resizedWindowBehaviour,
	"copiedFields": copiedFieldsBehaviour,
	"submitted":    formSubmittedBehaviour,
}

var store = storage.New()

func ReceiveNewEvent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// unpacking the body into a byte[]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var clientEvent ClientEvent

	// unmarshalling the json into our event enum
	err = json.Unmarshal(body, &clientEvent)
	if err != nil {
		log.Panic(err.Error())
	}

	// calling the correct handler
	if behaviour, ok := behaviours[clientEvent.EventType]; ok {
		behaviour(clientEvent)
		// need to find a better json for success if we want more app code to depend on it
		io.WriteString(w, "{ \"status\": \"success\" }")
	} else {
		io.WriteString(w, "{\"error\": true, \"message\": \"An error occurred\"}")
		log.Println("could not find handler for behaviour type " + clientEvent.EventType)
	}
}

func copiedFieldsBehaviour(event ClientEvent) {
	// get Data from the store for this sessionId
	savedData, found := store.Get(event.SessionId)

	if !found || savedData == nil {
		savedData = constructInitialData(event)
	}

	// updating the fields copied
	if data, ok := savedData.(Data); ok {
		data.CopyAndPaste = event.CopiedFields
		savedData = data
	}

	// save
	store.Set(event.SessionId, savedData)

	// print the struct
	log.Printf("%+v\n", savedData)
}

func resizedWindowBehaviour(event ClientEvent) {
	// get the Data struct for this session from the store
	savedData, found := store.Get(event.SessionId)

	if !found || savedData == nil {
		savedData = constructInitialData(event)
	}

	// update the ResizedFrom and ResizedTo properties
	// TODO: need to find a nicer way to work with the store
	if data, ok := savedData.(Data); ok {
		data.ResizeFrom = Dimension{string(event.FromWidth), string(event.FromHeight)}
		data.ResizeTo = Dimension{string(event.ToWidth), string(event.ToHeight)}
		savedData = data
	}

	// save
	store.Set(event.SessionId, savedData)

	// print the struct
	log.Printf("%+v\n", savedData)
}

func formSubmittedBehaviour(event ClientEvent) {
	// get Data from the store for this sessionId
	savedData, found := store.Get(event.SessionId)

	if !found || savedData == nil {
		savedData = constructInitialData(event)
	}

	// updating by adding the total time taken
	if data, ok := savedData.(Data); ok {
		data.FormCompletionTime = event.TimeTaken
		savedData = data
	}

	// save
	store.Set(event.SessionId, savedData)

	// printing the (final) struct
	log.Printf("%+v\n", savedData)
}

func constructInitialData(event ClientEvent) Data {
	return Data{
		WebsiteUrl: event.WebsiteUrl,
		SessionId:  event.SessionId,
	}
}
