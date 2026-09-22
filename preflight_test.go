package preflight

import (
	"bufio"
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type loggedRecord struct {
	Level     string `json:"level"`
	Message   string `json:"msg"`
	Phase     string `json:"phase"`
	Component string `json:"component,omitempty"`
	Status    string `json:"status,omitempty"`
}

var debugRecords = []loggedRecord{
	{Level: "INFO", Message: "Initializing system preflight...", Phase: "initialize"},
	{Level: "DEBUG", Message: "Discovering runtime capabilities...", Phase: "inspect", Component: "runtime"},
	{Level: "DEBUG", Message: "Inspecting processor topology...", Phase: "inspect", Component: "processor"},
	{Level: "DEBUG", Message: "Evaluating memory configuration...", Phase: "inspect", Component: "memory"},
	{Level: "DEBUG", Message: "Reviewing scheduler policy...", Phase: "inspect", Component: "scheduler"},
	{Level: "DEBUG", Message: "Analyzing garbage collector settings...", Phase: "inspect", Component: "garbage_collector"},
	{Level: "INFO", Message: "System profile collected.", Phase: "inspect", Status: "ready"},
	{Level: "DEBUG", Message: "Validating runtime compatibility...", Phase: "analyze", Component: "runtime"},
	{Level: "DEBUG", Message: "Checking resource readiness...", Phase: "analyze", Component: "resources"},
	{Level: "INFO", Message: "Runtime environment verified.", Phase: "analyze", Status: "ready"},
	{Level: "DEBUG", Message: "Building performance tuning plan...", Phase: "plan", Component: "planner"},
	{Level: "INFO", Message: "Applying runtime optimizations...", Phase: "apply", Status: "pending"},
	{Level: "DEBUG", Message: "Calibrating scheduler parameters...", Phase: "apply", Component: "scheduler"},
	{Level: "DEBUG", Message: "Tuning memory management policy...", Phase: "apply", Component: "memory"},
	{Level: "DEBUG", Message: "Balancing garbage collection targets...", Phase: "apply", Component: "garbage_collector"},
	{Level: "DEBUG", Message: "Verifying optimized runtime state...", Phase: "verify", Component: "runtime"},
	{Level: "INFO", Message: "Runtime optimization complete.", Phase: "verify", Status: "verified"},
	{Level: "INFO", Message: "Preflight completed successfully.", Phase: "finalize", Status: "success"},
}

var infoRecords = []loggedRecord{
	{Level: "INFO", Message: "Initializing system preflight...", Phase: "initialize"},
	{Level: "INFO", Message: "System profile collected.", Phase: "inspect", Status: "ready"},
	{Level: "INFO", Message: "Runtime environment verified.", Phase: "analyze", Status: "ready"},
	{Level: "INFO", Message: "Applying runtime optimizations...", Phase: "apply", Status: "pending"},
	{Level: "INFO", Message: "Runtime optimization complete.", Phase: "verify", Status: "verified"},
	{Level: "INFO", Message: "Preflight completed successfully.", Phase: "finalize", Status: "success"},
}

func TestRunEmitsPreflightSequence(t *testing.T) {
	var output bytes.Buffer
	setDefaultLogger(t, &output, slog.LevelDebug)

	Run()

	assertRecords(t, &output, debugRecords)
}

func TestRunRespectsInfoLevel(t *testing.T) {
	var output bytes.Buffer
	setDefaultLogger(t, &output, slog.LevelInfo)

	Run()

	assertRecords(t, &output, infoRecords)
}

func TestRunCanBeRepeated(t *testing.T) {
	var output bytes.Buffer
	setDefaultLogger(t, &output, slog.LevelDebug)

	Run()
	Run()

	want := make([]loggedRecord, 0, len(debugRecords)*2)
	want = append(want, debugRecords...)
	want = append(want, debugRecords...)
	assertRecords(t, &output, want)
}

func TestExportedAPI(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package directory: %v", err)
	}

	files := token.NewFileSet()
	var exported []string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(files, name, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		ast.Inspect(file, func(node ast.Node) bool {
			switch declaration := node.(type) {
			case *ast.FuncDecl:
				if declaration.Name.IsExported() {
					exported = append(exported, "func "+declaration.Name.Name)
					if declaration.Name.Name == "Run" && (declaration.Recv != nil || declaration.Type.Params.NumFields() != 0 || declaration.Type.Results != nil) {
						t.Errorf("Run must have signature func Run()")
					}
				}
			case *ast.TypeSpec:
				if declaration.Name.IsExported() {
					exported = append(exported, "type "+declaration.Name.Name)
				}
			case *ast.ValueSpec:
				for _, identifier := range declaration.Names {
					if identifier.IsExported() {
						exported = append(exported, "value "+identifier.Name)
					}
				}
			case *ast.Field:
				for _, identifier := range declaration.Names {
					if identifier.IsExported() {
						exported = append(exported, "field "+identifier.Name)
					}
				}
			}
			return true
		})
	}

	sort.Strings(exported)
	if want := []string{"func Run"}; !reflect.DeepEqual(exported, want) {
		t.Fatalf("exported API mismatch: got %v, want %v", exported, want)
	}
}

func setDefaultLogger(t *testing.T, output *bytes.Buffer, level slog.Leveler) {
	t.Helper()
	original := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: level})))
	t.Cleanup(func() { slog.SetDefault(original) })
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
			case "time", "level", "msg", "phase", "component", "status":
			default:
				t.Errorf("unexpected log attribute %q", key)
			}
		}

		var record loggedRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			t.Fatalf("decode log record fields: %v", err)
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
