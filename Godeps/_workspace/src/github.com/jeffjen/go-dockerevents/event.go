package dockerevents

import (
	docker "github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
	"os"
)

type FilterInfo map[string]bool

type EventFilter struct {
	Image     FilterInfo
	Status    FilterInfo
	Container FilterInfo
}

func ParseEventFilter(constraints string) *EventFilter {
	// YAML encoded string describing docker events to fitler

	var filter map[interface{}][]interface{}
	if err := yaml.Unmarshal([]byte(constraints), &filter); err != nil {
		panic(err)
	}

	var Filter EventFilter

	Filter.Image = make(FilterInfo)
	Filter.Status = make(FilterInfo)
	Filter.Container = make(FilterInfo)

	for key, val := range filter {
		switch key {
		case "image":
			for _, image := range val {
				Filter.Image[image.(string)] = true
			}
			break
		case "status":
			for _, status := range val {
				Filter.Status[status.(string)] = true
			}
			break
		case "container":
			for _, container := range val {
				Filter.Container[container.(string)] = true
			}
			break
		default:
			break
		}
	}

	return &Filter
}

type EventListener struct {
	C   chan *docker.APIEvents
	Cli *docker.Client
	f   *EventFilter
}

func NewEventListener(filter *EventFilter) (event *EventListener, err error) {
	if cli, err := NewClient(); err != nil {
		return nil, err
	} else {
		if filter == nil {
			filter = new(EventFilter)
		}
		event = &EventListener{make(chan *docker.APIEvents), cli, filter}
		err = event.Cli.AddEventListener(event.C)
	}
	return
}

func (event *EventListener) Close() {
	event.Cli.RemoveEventListener(event.C)
}

func (event *EventListener) Filter(ev *docker.APIEvents) bool {
	var chkStatus, chkImage, chkID bool
	if event.f == nil {
		return true
	}
	if len(event.f.Status) == 0 {
		chkStatus = true
	} else {
		_, chkStatus = event.f.Status[ev.Status]
	}
	if len(event.f.Image) == 0 {
		chkImage = true
	} else {
		_, chkImage = event.f.Image[ev.From]
	}
	if len(event.f.Container) == 0 {
		chkID = true
	} else {
		_, chkID = event.f.Container[ev.ID]
	}
	return chkStatus && chkImage && chkID
}

type Event struct {
	Status string         `json:"status"`
	ID     string         `json:"iden"`
	From   string         `json:"image"`
	Time   int64          `json:"time"`
	Cli    *docker.Client `json:"-"`
}

func DefaultCallBack(event *Event) interface{} {
	return event
}

func EventLoop(filter *EventFilter, size uint64) (output <-chan *Event) {
	// main driver for listening docker container events

	if filter == nil {
		filter = ParseEventFilter(os.Getenv("DOCKER_EVENT_FILTER"))
	}

	sink := make(chan *Event, size)

	go func() {
		ev, _ := NewEventListener(filter)
		defer func() { ev.Close() }()

		for msg := range ev.C {
			if ev.Filter(msg) {
				sink <- &Event{msg.Status, msg.ID, msg.From, msg.Time, ev.Cli}
			}
		}
	}()

	return sink
}
