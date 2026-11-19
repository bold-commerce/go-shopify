package goshopify

import (
	"bytes"
	"io/ioutil"
	"net/http"
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
	var out bytes.Buffer
	logger := &LeveledLogger{Level: LevelDebug, stdoutOverride: &out}
	transport := bodyLifecycleTransport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Status: "OK", Header: http.Header{"X-Request-Id": {"abc"}}, Body: ioutil.NopCloser(strings.NewReader("response body"))}, nil
	})
	client := MustNewClient(app, "fooshop", "abcd", WithLogger(logger), WithHTTPClient(&http.Client{Transport: transport}))
	req, _ := http.NewRequest("POST", "https://test.com/foo", strings.NewReader("request body"))
	if err := client.Do(req, nil); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"POST: https://test.com/foo", "SENT: request body", "X-Request-Id: abc", "RECV 200: OK", "RESP: response body"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("missing %q in %q", want, out.String())
		}
	}
}
