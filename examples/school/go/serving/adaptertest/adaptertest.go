package adaptertest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// CheckAdapterAnswersTheSameWay Holds every Framework to one Contract: the
// same Path, Query, Headers and Body Arrive, and the same Status and JSON Leave.
func CheckAdapterAnswersTheSameWay(t *testing.T, build func(transport.Libretto) http.Handler) {
	t.Helper()

	var seen transport.Request
	libretto := transport.Libretto{
		Routes: []transport.Route{
			{Method: "POST", Pattern: "/echo/{id}", Handle: func(req transport.Request) transport.Response {
				seen = req

				return transport.Response{Status: http.StatusCreated, Body: map[string]string{"id": req.Path["id"]}}
			}},
			{Method: "GET", Pattern: "/echo", Handle: func(transport.Request) transport.Response {
				return transport.Response{Status: http.StatusOK, Body: []string{}}
			}},
			{Method: "DELETE", Pattern: "/echo/{id}", Handle: func(transport.Request) transport.Response {
				return transport.Response{Status: http.StatusNoContent}
			}},
			{Method: "GET", Pattern: "/failing", Handle: func(transport.Request) transport.Response {
				return transport.Response{Status: http.StatusNotFound, Body: map[string]string{"error": "gone"}}
			}},
		},
		Unknown: func(transport.Request) transport.Response {
			return transport.Response{Status: http.StatusNotFound, Body: map[string]string{"error": "route not Found"}}
		},
	}
	server := httptest.NewServer(build(libretto))
	defer server.Close()

	request, _ := http.NewRequest("POST", server.URL+"/echo/007?page=2&size=5&page=9", strings.NewReader(`{"a":1}`))
	request.Header.Set("Authorization", "bEaReR abc")
	request.Header.Set("Content-Type", "text/plain")
	reply, body := send(t, request)

	if reply.StatusCode != http.StatusCreated || reply.Header.Get("Content-Type") != "application/json" || body != `{"id":"007"}` {
		t.Fatalf("status=%d type=%q body=%s", reply.StatusCode, reply.Header.Get("Content-Type"), body)
	}
	if seen.Path["id"] != "007" || len(seen.Query) != 2 || seen.Query["page"] != "2" || seen.Query["size"] != "5" ||
		seen.Headers["authorization"] != "bEaReR abc" || string(seen.Body) != `{"a":1}` {
		t.Fatalf("the Request Arrived as %+v", seen)
	}

	listed, _ := http.NewRequest("GET", server.URL+"/echo", nil)
	if _, body = send(t, listed); body != "[]" {
		t.Fatalf("an empty List is [], got %s", body)
	}

	drop, _ := http.NewRequest("DELETE", server.URL+"/echo/7", nil)
	if reply, body = send(t, drop); reply.StatusCode != http.StatusNoContent || body != "" {
		t.Fatalf("a 204 must Carry no Body, status=%d body=%q", reply.StatusCode, body)
	}

	failing, _ := http.NewRequest("GET", server.URL+"/failing", nil)
	if reply, _ = send(t, failing); reply.StatusCode != http.StatusNotFound {
		t.Fatalf("a Failure keeps its Status, got %d", reply.StatusCode)
	}

	for _, stranger := range [][2]string{{"GET", "/nowhere"}, {"GET", "/echo/"}, {"PUT", "/echo/7"}, {"GET", "/echo/7/more"}, {"PATCH", "/echo"}} {
		unknown, _ := http.NewRequest(stranger[0], server.URL+stranger[1], nil)
		if reply, body = send(t, unknown); reply.StatusCode != http.StatusNotFound || body != `{"error":"route not Found"}` {
			t.Errorf("%s %s: status=%d body=%s", stranger[0], stranger[1], reply.StatusCode, body)
		}
	}

	oversized, _ := http.NewRequest("POST", server.URL+"/echo/1", strings.NewReader(strings.Repeat(" ", 2*transport.MaxBodyBytes)))
	if reply, _ = send(t, oversized); reply.StatusCode != http.StatusCreated || len(seen.Body) != transport.MaxBodyBytes+1 {
		t.Fatalf("a Body past the Bound Arrives one Byte past it, status=%d length=%d", reply.StatusCode, len(seen.Body))
	}
}

func send(t *testing.T, request *http.Request) (*http.Response, string) {
	t.Helper()
	reply, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer reply.Body.Close()
	body, _ := io.ReadAll(reply.Body)

	return reply, string(body)
}
