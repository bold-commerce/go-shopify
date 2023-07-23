package goshopify

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestLeveledLogger(t *testing.T) {
	tests := []struct {
		level  int
		input  string
		stdout string
		stderr string
	}{
		{
			level:  LevelError,
			input:  "log",
			stderr: "[ERROR] error log\n",
			stdout: "",
		},
		{
			level:  LevelWarn,
			input:  "log",
			stderr: "[ERROR] error log\n[WARN] warn log\n",
			stdout: "",
		},
		{
			level:  LevelInfo,
			input:  "log",
			stderr: "[ERROR] error log\n[WARN] warn log\n",
			stdout: "[INFO] info log\n",
		},
		{
			level:  LevelDebug,
			input:  "log",
			stderr: "[ERROR] error log\n[WARN] warn log\n",
			stdout: "[INFO] info log\n[DEBUG] debug log\n",
		},
	}

	for _, test := range tests {
		err := &bytes.Buffer{}
		out := &bytes.Buffer{}
		log := &LeveledLogger{Level: test.level, stderrOverride: err, stdoutOverride: out}

		log.Errorf("error %s", test.input)
		log.Warnf("warn %s", test.input)
		log.Infof("info %s", test.input)
		log.Debugf("debug %s", test.input)

		stdout := out.String()
		stderr := err.String()

		if stdout != test.stdout {
			t.Errorf("leveled logger %d expected stdout \"%s\" received \"%s\"", test.level, test.stdout, stdout)
		}
		if stderr != test.stderr {
			t.Errorf("leveled logger %d expected stderr \"%s\" received \"%s\"", test.level, test.stderr, stderr)
		}
	}

	log := &LeveledLogger{Level: LevelDebug}
	if log.stderr() != os.Stderr {
		t.Errorf("leveled logger with no stderr override expects os.Stderr")
	}
	if log.stdout() != os.Stdout {
		t.Errorf("leveled logger with no stdout override expects os.Stdout")
	}

}

func TestDoGetHeadersDebug(t *testing.T) {
	err := &bytes.Buffer{}
	out := &bytes.Buffer{}
	logger := &LeveledLogger{Level: LevelDebug, stderrOverride: err, stdoutOverride: out}

	reqExpected := "[DEBUG] GET: //http:%2F%2Ftest.com/foo/1\n[DEBUG] SENT: request body\n"
	resExpected := "[DEBUG] Shopify X-Request-Id: 00000000-0000-0000-0000-000000000000\n[DEBUG] RECV 200: OK\n[DEBUG] RESP: response body\n"

	client := NewClient(app, "fooshop", "abcd", WithLogger(logger))

	client.logBody(nil, "")
	if out.String() != "" {
		t.Errorf("logBody expected empty log output but received \"%s\"", out.String())
	}

	client.logRequest(nil)
	if out.String() != "" {
		t.Errorf("logRequest expected empty log output received \"%s\"", out.String())
	}

	client.logRequest(&http.Request{
		Method: "GET",
		URL:    &url.URL{Host: "http://test.com", Path: "/foo/1"},
		Body:   ioutil.NopCloser(strings.NewReader("request body")),
	})

	if out.String() != reqExpected {
		t.Errorf("doGetHeadersDebug expected stdout \"%s\" received \"%s\"", reqExpected, out)
	}

	err.Reset()
	out.Reset()

	client.logResponse(nil)
	if out.String() != "" {
		t.Errorf("logResponse expected empty log output received \"%s\"", out.String())
	}

	client.logResponse(&http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Header: map[string][]string{"X-Request-Id": []string{"00000000-0000-0000-0000-000000000000"}},
		Body:       ioutil.NopCloser(strings.NewReader("response body")),
	})

	if out.String() != resExpected {
		t.Errorf("doGetHeadersDebug expected stdout \"%s\" received \"%s\"", resExpected, out.String())
	}
}
