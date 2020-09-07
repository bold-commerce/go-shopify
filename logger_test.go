package goshopify

import (
	"bytes"
	"io"
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

	warnResExpectedStdErr := "[WARN] WARNING: body truncation may have occurred, consider increasing the value of MaxLoggedHTTPBodyBytes\n"
	reqExpected := "[DEBUG] GET: //http:%2F%2Ftest.com/foo/1\n[DEBUG] SENT: request body\n"
	resExpected := "[DEBUG] RECV 200: OK\n[DEBUG] RESP: response body\n"

	client := NewClient(app, "fooshop", "abcd", WithLogger(logger))

	client.logBody(nil, "")
	if out.String() != "" {
		t.Errorf("logBody expected empty log output but received \"%s\"", out.String())
	}

	client.logRequest(nil)
	if out.String() != "" {
		t.Errorf("logRequest expected empty log output received \"%s\"", out.String())
	}

	bdy := ioutil.NopCloser(strings.NewReader(largeReqRespBody))
	client.logBody(&bdy, "RESP: %s")
	if err.String() != warnResExpectedStdErr {
		t.Errorf("logBody expected \"%s\" but received \"%s\"", warnResExpectedStdErr, err)
	}

	if !strings.Contains(out.String(), `[DEBUG] RESP: {"orders":[{"id":123456,"email":"jon@doe.ca","closed_at":null,"created_at":"2016-05-17T04:14:36-00:00","updated_at":"2016-05-17T04:14:36-04:00",`) {
		t.Errorf("logBody expected non-empty output, but received \"%s\"", out)
	}

	err.Reset()
	out.Reset()

	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Host: "http://test.com", Path: "/foo/1"},
		Body:   ioutil.NopCloser(strings.NewReader("request body")),
	}
	client.logRequest(req)

	_, rereadErr := req.Body.Read(make([]byte, 16))
	if rereadErr != nil && rereadErr == io.EOF {
		t.Errorf("could not re-read request body, may not have been reset")
	}

	if out.String() != reqExpected {
		t.Errorf("doGetHeadersDebug expected stdout \"%s\" received \"%s\"", reqExpected, out)
	}

	err.Reset()
	out.Reset()

	client.logResponse(nil)
	if out.String() != "" {
		t.Errorf("logResponse expected empty log output received \"%s\"", out.String())
	}

	resp := &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("response body")),
	}
	client.logResponse(resp)

	_, rereadErr = resp.Body.Read(make([]byte, 16))
	if rereadErr != nil && rereadErr == io.EOF {
		t.Errorf("could not re-read response body, may not have been reset")
	}

	if out.String() != resExpected {
		t.Errorf("doGetHeadersDebug expected stdout \"%s\" received \"%s\"", resExpected, out.String())
	}
}

func TestCheckMaxBodyBytes(t *testing.T) {
	client := NewClient(app, "fooshop", "abcd")
	client.checkMaxBodyBytes()
	if client.maxBodyBytes != defaultMaxBodyBytes {
		t.Errorf("checkMaxBodyBytes failed, expected client.maxBodyBytes to be set to the the default value of %d, but was set to %d", defaultMaxBodyBytes, client.maxBodyBytes)
	}
}

var largeReqRespBody = func() string {
	b, err := ioutil.ReadFile("./fixtures/orders.json")
	if err != nil {
		panic(err)
	}

	sb := strings.Builder{}

	// the output of a single orders.json file is ~3.4KB, so we'll do that 10x to trip up the limit
	for i := 0; i < 10; i++ {
		sb.WriteString(string(b))
	}

	return sb.String()
}()
