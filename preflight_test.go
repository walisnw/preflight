package preflight

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log/slog"
	"reflect"
	"testing"
)

type loggedRecord struct {
	Level   string `json:"level"`
	Message string `json:"msg"`
}

var debugRecords = []loggedRecord{
	{Level: "INFO", Message: "Initializing preflight..."},
	{Level: "DEBUG", Message: "Inspecting runtime configuration..."},
	{Level: "INFO", Message: "Preparing application startup..."},
	{Level: "DEBUG", Message: "Evaluating application environment..."},
	{Level: "INFO", Message: "Performing startup checks..."},
	{Level: "DEBUG", Message: "Reviewing startup conditions..."},
	{Level: "DEBUG", Message: "Analyzing runtime readiness..."},
	{Level: "DEBUG", Message: "Preparing startup recommendations..."},
	{Level: "INFO", Message: "Finalizing preflight..."},
	{Level: "INFO", Message: "Preflight completed successfully."},
}

var infoRecords = []loggedRecord{
	{Level: "INFO", Message: "Initializing preflight..."},
	{Level: "INFO", Message: "Preparing application startup..."},
	{Level: "INFO", Message: "Performing startup checks..."},
	{Level: "INFO", Message: "Finalizing preflight..."},
	{Level: "INFO", Message: "Preflight completed successfully."},
}

func TestRunUsesDefaultLogger(t *testing.T) {
	original := slog.Default()
	var output bytes.Buffer
	slog.SetDefault(newLogger(&output, slog.LevelDebug))
	t.Cleanup(func() { slog.SetDefault(original) })

	Run()

	assertRecords(t, &output, debugRecords)
}

func TestRunUsesInjectedLogger(t *testing.T) {
	original := slog.Default()
	var defaultOutput bytes.Buffer
	slog.SetDefault(newLogger(&defaultOutput, slog.LevelDebug))
	t.Cleanup(func() { slog.SetDefault(original) })

	var injectedOutput bytes.Buffer
	Run(WithLogger(newLogger(&injectedOutput, slog.LevelDebug)))

	assertRecords(t, &injectedOutput, debugRecords)
	assertRecords(t, &defaultOutput, nil)
}

func TestRunRespectsInfoLevel(t *testing.T) {
	var output bytes.Buffer

	Run(WithLogger(newLogger(&output, slog.LevelInfo)))

	assertRecords(t, &output, infoRecords)
}

func TestRunUsesLastLoggerOption(t *testing.T) {
	var firstOutput bytes.Buffer
	var lastOutput bytes.Buffer

	Run(
		WithLogger(newLogger(&firstOutput, slog.LevelDebug)),
		WithLogger(newLogger(&lastOutput, slog.LevelDebug)),
	)

	assertRecords(t, &firstOutput, nil)
	assertRecords(t, &lastOutput, debugRecords)
}

func TestRunCanBeRepeated(t *testing.T) {
	var output bytes.Buffer
	logger := newLogger(&output, slog.LevelDebug)

	Run(WithLogger(logger))
	Run(WithLogger(logger))

	want := make([]loggedRecord, 0, len(debugRecords)*2)
	want = append(want, debugRecords...)
	want = append(want, debugRecords...)
	assertRecords(t, &output, want)
}

func TestWithLoggerPanicsForNilLogger(t *testing.T) {
	assertPanic(t, "preflight: nil logger", func() {
		WithLogger(nil)
	})
}

func TestRunPanicsForNilOption(t *testing.T) {
	assertPanic(t, "preflight: nil option", func() {
		Run(nil)
	})
}

func newLogger(output *bytes.Buffer, level slog.Leveler) *slog.Logger {
	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: level}))
}

func assertRecords(t *testing.T, output *bytes.Buffer, want []loggedRecord) {
	t.Helper()

	var got []loggedRecord
	scanner := bufio.NewScanner(bytes.NewReader(output.Bytes()))
	for scanner.Scan() {
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(scanner.Bytes(), &raw); err != nil {
			t.Fatalf("decode log record: %v", err)
		}

		for key := range raw {
			switch key {
			case "time", "level", "msg":
			default:
				t.Errorf("unexpected log attribute %q", key)
			}
		}

		var record loggedRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			t.Fatalf("decode level and message: %v", err)
		}
		got = append(got, record)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan log records: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("records mismatch\ngot:  %#v\nwant: %#v", got, want)
	}
}

func assertPanic(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		if got := recover(); got != want {
			t.Fatalf("panic mismatch: got %q, want %q", got, want)
		}
	}()

	fn()
}
