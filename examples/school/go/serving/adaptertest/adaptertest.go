package adaptertest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// CheckAdapterAnswersTheSameWay Holds every Framework to one Contract: the
// same Path, Query, Token and Body Arrive, and the same Status and JSON Leave.
func CheckAdapterAnswersTheSameWay(t *testing.T, build func([]transport.Route) http.Handler) {
	t.Helper()

	var seen transport.Request
	routes := []transport.Route{
		{Method: "POST", Pattern: "/echo/{id}", Handle: func(req transport.Request) transport.Response {
			seen = req

			return transport.Response{Status: http.StatusCreated, Body: map[string]string{"id": req.Path["id"]}}
		}},
		{Method: "DELETE", Pattern: "/echo/{id}", Handle: func(transport.Request) transport.Response {
			return transport.Response{Status: http.StatusNoContent}
		}},
		{Method: "GET", Pattern: "/failing", Handle: func(transport.Request) transport.Response {
			return transport.Response{Status: http.StatusNotFound, Body: map[string]string{"error": "gone"}}
		}},
	}
	server := httptest.NewServer(build(routes))
	defer server.Close()

	request, _ := http.NewRequest("POST", server.URL+"/echo/7?page=2&size=5", strings.NewReader(`{"a":1}`))
	request.Header.Set("Authorization", "Bearer abc")
	reply := send(t, request)

	if reply.StatusCode != http.StatusCreated || reply.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("status=%d type=%q", reply.StatusCode, reply.Header.Get("Content-Type"))
	}
	if seen.Path["id"] != "7" || seen.Query["page"] != "2" || seen.Query["size"] != "5" || seen.Token != "abc" || string(seen.Body) != `{"a":1}` {
		t.Fatalf("the Request Arrived as %+v", seen)
	}

	drop, _ := http.NewRequest("DELETE", server.URL+"/echo/7", nil)
	if reply = send(t, drop); reply.StatusCode != http.StatusNoContent || reply.ContentLength > 0 {
		t.Fatalf("a 204 must Carry no Body, status=%d length=%d", reply.StatusCode, reply.ContentLength)
	}

	failing, _ := http.NewRequest("GET", server.URL+"/failing", nil)
	if reply = send(t, failing); reply.StatusCode != http.StatusNotFound {
		t.Fatalf("a Failure keeps its Status, got %d", reply.StatusCode)
	}

	unknown, _ := http.NewRequest("GET", server.URL+"/nowhere", nil)
	if reply = send(t, unknown); reply.StatusCode != http.StatusNotFound {
		t.Fatalf("an unknown Route is a 404, got %d", reply.StatusCode)
	}
}

func send(t *testing.T, request *http.Request) *http.Response {
	t.Helper()
	reply, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	reply.Body.Close()

	return reply
}
