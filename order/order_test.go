package order

import (
	"testing"
)

type fakeContext struct {
	code     int
	response map[string]string
}

func (c fakeContext) Order() (Order, error) {
	return Order{
		SalesChannel: "Offline",
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

	c := &fakeContext{}
	handler.Order(c)

	want := "Offline is not accepted"
	if want != c.response["message"] {
		t.Errorf("%q is expected but got %q\n", want, c.response["message"])
	}
}
