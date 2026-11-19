package goshopify

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"
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

func TestBodyReadErrorsStopRequest(t *testing.T) {
	for _, side := range []string{"request", "response"} {
		t.Run(side, func(t *testing.T) {
			broken := &trackedBody{Reader: iotest.TimeoutReader(strings.NewReader("partial"))}
			transport := bodyLifecycleTransport(func(r *http.Request) (*http.Response, error) {
				if side == "request" {
					t.Fatal("sent a request after its body failed to read")
				}
				return &http.Response{StatusCode: 200, Body: broken}, nil
			})
			c := MustNewClient(App{}, "test", "", WithHTTPClient(&http.Client{Transport: transport}))
			req, _ := http.NewRequest("POST", "https://test.myshopify.com", nil)
			if side == "request" {
				req.Body = broken
			}
			if err := c.Do(req, nil); err != iotest.ErrTimeout {
				t.Fatalf("got %v, want %v", err, iotest.ErrTimeout)
			}
			if broken.closes != 1 {
				t.Fatalf("body closed %d times", broken.closes)
			}
		})
	}
}
