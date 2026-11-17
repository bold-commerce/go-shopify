package goshopify

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type bodyLifecycleTransport func(*http.Request) (*http.Response, error)

func (f bodyLifecycleTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type trackedBody struct {
	io.Reader
	closes int
}

func (b *trackedBody) Close() error { b.closes++; return nil }

func TestBodyLifecycle(t *testing.T) {
	for _, status := range []int{200, 422, 503} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var bodies []*trackedBody
			calls := 0
			transport := bodyLifecycleTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				code := status
				data := `{"errors":"failure"}`
				if status == 200 || (status == 503 && calls == 2) {
					code = 200
					data = `{"id":7}`
				}
				b := &trackedBody{Reader: strings.NewReader(data)}
				bodies = append(bodies, b)
				return &http.Response{StatusCode: code, Header: make(http.Header), Body: b, Request: r}, nil
			})
			c := MustNewClient(App{}, "test", "token", WithHTTPClient(&http.Client{Transport: transport}), WithRetry(1))
			req, _ := http.NewRequestWithContext(context.Background(), "GET", "https://test.myshopify.com/admin/shop.json", nil)
			var result struct {
				Id int `json:"id"`
			}
			err := c.Do(req, &result)
			if status == 422 {
				if err == nil {
					t.Fatal("expected API error")
				}
			} else if err != nil || result.Id != 7 {
				t.Fatalf("result=%v error=%v", result, err)
			}
			for _, b := range bodies {
				if b.closes != 1 {
					t.Fatalf("original body closed %d times", b.closes)
				}
			}
			if status == 503 && calls != 2 {
				t.Fatalf("calls=%d", calls)
			}
		})
	}
}

func TestBodyLoggingPreservesClose(t *testing.T) {
	c := MustNewClient(App{}, "test", "token")
	original := &trackedBody{Reader: strings.NewReader("payload")}
	var body io.ReadCloser = original
	c.logBody(&body, "%s")
	data, err := ioutil.ReadAll(body)
	if err != nil || string(data) != "payload" {
		t.Fatalf("body=%q error=%v", data, err)
	}
	if original.closes != 0 {
		t.Fatal("closed before caller finished")
	}
	if err := body.Close(); err != nil {
		t.Fatal(err)
	}
	if original.closes != 1 {
		t.Fatal("original close was lost")
	}
}
