package box

import (
	"fmt"
	"net/http"
)

type EventService struct {
	*Client
}

type EventDetails struct {
	ServiceID   string `json:"service_id"`
	EkmID       string `json:"ekm_id"`
	VersionID   string `json:"version_id"`
	ServiceName string `json:"service_name"`
}

type EventSource struct {
	ItemID   string      `json:"item_id"`
	ItemType string      `json:"item_type"`
	Parent   *ItemParent `json:"parent"`
	ItemName string      `json:"item_name"`
}

type Event struct {
	EventType         string        `json:"event_type"`
	EventID           string        `json:"event_id"`
	Type              string        `json:"type"`
	CreatedAt         string        `json:"created_at"`
	CreatedBy         *User         `json:"created_by"`
	Source            *EventSource  `json:"source"`
	SessionID         interface{}   `json:"session_id"`
	AdditionalDetails *EventDetails `json:"additional_details"`
	IPAddress         string        `json:"ip_address"`
}

type EventsCollection struct {
	ChunkSize          int      `json:"chunk_size"`
	NextStreamPosition string   `json:"next_stream_position"`
	Entries            []*Event `json:"entries"`
}

func (e *EventService) getSeedStreamPos(startTime string) string {
	var respBoxEventsJSON EventsCollection
	req, err := http.NewRequest("GET", e.BaseUrl.String()+"/events?stream_type=admin_logs&limit=1&created_after="+startTime, nil)
	req.Header.Add("Authorization", "Bearer "+e.Token)
	_, err = e.DoWithRetries(req, &respBoxEventsJSON, 5)
	if err != nil {
		fmt.Println("this unmarshal")
		fmt.Println(err)
	}
	return respBoxEventsJSON.NextStreamPosition
}

func (e *EventService) getEvents(eventLimit string, streamPos string) *EventsCollection {
	var respBoxEventsJSON EventsCollection
	req, err := http.NewRequest("GET", e.BaseUrl.String()+"/events?stream_type=admin_logs&limit="+eventLimit+"&stream_position="+streamPos, nil)
	req.Header.Add("Authorization", "Bearer "+e.Token)
	_, err = e.DoWithRetries(req, &respBoxEventsJSON, 5)
	if err != nil {

		fmt.Println(err)
	}
	return &respBoxEventsJSON
}

func (e *EventService) streamEvents(eventLimit string, streamPos string, tunnel chan *Event) {
	for {
		if streamPos == "" {
			close(tunnel)
			return
		}
		eventsCollection := e.getEvents(eventLimit, streamPos)
		for _, event := range eventsCollection.Entries {
			tunnel <- event
		}
		streamPos = eventsCollection.NextStreamPosition
		//time.Sleep(2 * time.Second)
	}
}

func (e *EventService) Channel(eventLimit string, startTime string) chan *Event {
	streamPos := e.getSeedStreamPos(startTime)
	eventStream := make(chan *Event)
	go e.streamEvents(eventLimit, streamPos, eventStream)
	return eventStream
}
