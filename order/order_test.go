package order

import (
	"errors"
	"net/http"
	"testing"
)

type fakeContext struct {
	channel  string
	code     int
	response map[string]string
}

func (c fakeContext) Order() (Order, error) {
	return Order{
		SalesChannel: c.channel,
	}, nil
}

func (c *fakeContext) JSON(code int, v interface{}) {
	c.code = code
	c.response = v.(map[string]string)
}

func TestOnlyAcceptOnlineChannel(t *testing.T) {
	handler := &Handler{
		channel: "Online",
	}

	c := &fakeContext{channel: "Offline"}
	handler.Order(c)

	want := "Offline is not accepted"
	if want != c.response["message"] {
		t.Errorf("%q is expected but got %q\n", want, c.response["message"])
	}
}

func TestOnlyAcceptOfflineChannel(t *testing.T) {
	handler := &Handler{
		channel: "Offline",
	}

	c := &fakeContext{channel: "Online"}
	handler.Order(c)

	want := "Online is not accepted"
	if want != c.response["message"] {
		t.Errorf("%q is expected but got %q\n", want, c.response["message"])
	}
}

type fakeContextBadRequest struct {
	code     int
	response map[string]string
}

func (c fakeContextBadRequest) Order() (Order, error) {
	return Order{}, errors.New("something went wrong")
}
func (c *fakeContextBadRequest) JSON(code int, v interface{}) {
	c.code = code
	c.response = v.(map[string]string)
}

func TestBadRequestOrderWentWrong(t *testing.T) {
	handler := &Handler{}

	c := &fakeContextBadRequest{}
	handler.Order(c)

	want := http.StatusBadRequest
	if want != c.code {
		t.Errorf("%d status code expected, but got %d", want, c.code)
	}
}

type fakeContextBadRequestWithChannel struct {
	channel         string
	jsonCalledCount int
}

func (c *fakeContextBadRequestWithChannel) Order() (Order, error) {
	return Order{SalesChannel: c.channel}, errors.New("went wrong")
}

func (c *fakeContextBadRequestWithChannel) JSON(code int, v interface{}) {
	c.jsonCalledCount++
}

func TestOnlyCalledJSONOnce(t *testing.T) {
	handler := &Handler{channel: "Offline"}

	c := &fakeContextBadRequestWithChannel{channel: "Online"}
	handler.Order(c)

	want := 1
	if want != c.jsonCalledCount {
		t.Errorf("it should called one time but got %d times", c.jsonCalledCount)
	}
}

type spyStore struct {
	wasCalled bool
}

func (s *spyStore) Save(order Order) error {
	s.wasCalled = true
	return nil
}

func TestOrderWasSaved(t *testing.T) {
	spy := &spyStore{}
	handler := &Handler{
		channel: "Online",
		store:   spy,
	}
	c := &fakeContext{channel: "Online"}
	handler.Order(c)

	if !spy.wasCalled {
		t.Errorf("it should store data")
	}
}

type failStore struct{}

func (f *failStore) Save(Order) error {
	return errors.New("")
}

func TestOrderFailAtSave(t *testing.T) {
	store := &failStore{}

	handler := &Handler{
		channel: "Online",
		store:   store,
	}

	c := &fakeContext{channel: "Online"}
	handler.Order(c)

	want := http.StatusInternalServerError
	if want != c.code {
		t.Errorf("%d is expected but got %d\n", want, c.code)
	}
}
