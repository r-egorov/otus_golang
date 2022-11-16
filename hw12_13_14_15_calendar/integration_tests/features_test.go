package main

import (
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

type GetEventsResponse struct {
	Events []Event `json:"events"`
}

type CreateEventResponse struct {
	Event Event `json:"event"`
}

type Event struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	DateTime    time.Time     `json:"datetime"`
	Duration    time.Duration `json:"duration"`
	Description string        `json:"description"`
	OwnerID     uuid.UUID     `json:"ownerId"`
}

type feature struct {
	resp          *http.Response
	gotEvent      Event
	gotListEvents []Event
}

func (f *feature) iSendARequestToWithJSONbody(method, endpoint string, jsonBody *godog.DocString) error {
	var body io.Reader
	if jsonBody != nil {
		body = strings.NewReader(jsonBody.Content)
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	f.resp, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (f *feature) otherFieldsEqualToFollowing(jsonFields *godog.DocString) error {
	var expectedEvent Event
	err := json.NewDecoder(strings.NewReader(jsonFields.Content)).Decode(&expectedEvent)
	if err != nil {
		return err
	}

	expectedEvent.ID = f.gotEvent.ID
	if expectedEvent != f.gotEvent {
		return fmt.Errorf("fields are not equal")
	}
	return nil
}

func (f *feature) theEventHasAnID() error {
	if f.gotEvent.ID == uuid.Nil {
		return fmt.Errorf("there should be an ID in the event")
	}
	return nil
}

func (f *feature) theResponseCodeIs(code int) error {
	if code != f.resp.StatusCode {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, f.resp.StatusCode)
	}
	return nil
}

func (f *feature) theResponseContainsAnEvent() error {
	var response CreateEventResponse
	err := json.NewDecoder(f.resp.Body).Decode(&response)
	defer f.resp.Body.Close()
	if err != nil {
		return err
	}

	f.gotEvent = response.Event
	return nil
}

func (f *feature) iSendARequestTo(method, endpoint string) error {

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	f.resp, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (f *feature) theResponseContainsListOfEvents() error {
	var response GetEventsResponse
	err := json.NewDecoder(f.resp.Body).Decode(&response)
	defer f.resp.Body.Close()
	if err != nil {
		return err
	}

	f.gotListEvents = response.Events
	return nil
}

func (f *feature) theEventsHaveIDs() error {
	for _, event := range f.gotListEvents {
		if event.ID == uuid.Nil {
			return fmt.Errorf("there should be an ID in the event")
		}
	}
	return nil
}

func (f *feature) eventsFieldsAreAsFollowing(jsonFields *godog.DocString) error {
	var expectedEvents []Event
	err := json.NewDecoder(strings.NewReader(jsonFields.Content)).Decode(&expectedEvents)
	if err != nil {
		return err
	}

	createCompare := func(events []Event) func(i, j int) bool {
		less := func(i, j int) bool {
			if events[i].DateTime == events[j].DateTime {
				return events[i].OwnerID.String() < events[j].OwnerID.String()
			}
			return events[i].DateTime.Before(events[j].DateTime)
		}
		return less
	}
	sort.Slice(expectedEvents, createCompare(expectedEvents))
	sort.Slice(f.gotListEvents, createCompare(f.gotListEvents))

	for i := range f.gotListEvents {
		expectedEvents[i].ID = f.gotListEvents[i].ID
		if expectedEvents[i] != f.gotListEvents[i] {
			return fmt.Errorf("events different, got: %v, expected: %v", f.gotListEvents[i], expectedEvents[i])
		}
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	f := feature{}

	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)" with JSON-body:$`, f.iSendARequestToWithJSONbody)
	ctx.Step(`^the response code is (\d+)$`, f.theResponseCodeIs)
	ctx.Step(`^the response contains an event$`, f.theResponseContainsAnEvent)
	ctx.Step(`^the event has an ID$`, f.theEventHasAnID)
	ctx.Step(`^other fields equal to following:$`, f.otherFieldsEqualToFollowing)

	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, f.iSendARequestTo)
	ctx.Step(`^the response contains list of events$`, f.theResponseContainsListOfEvents)
	ctx.Step(`^the events have IDs$`, f.theEventsHaveIDs)
	ctx.Step(`^events fields are as following:$`, f.eventsFieldsAreAsFollowing)
}
