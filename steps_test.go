package wizard

import (
	"errors"
	"reflect"
	"testing"
)

func TestRunStepsOrder(t *testing.T) {
	var calls []string
	err := RunSteps(
		[]Step{{Name: "A", Run: func() error { calls = append(calls, "run:A"); return nil }}, {Name: "B", Run: func() error { calls = append(calls, "run:B"); return nil }}},
		func(i, total int, name string) { calls = append(calls, "start:"+name) },
		func(i, total int) { calls = append(calls, "done") },
	)
	if err != nil {
		t.Fatalf("RunSteps error: %v", err)
	}
	want := []string{"start:A", "run:A", "done", "start:B", "run:B", "done"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
}

func TestRunStepsStopsOnError(t *testing.T) {
	boom := errors.New("boom")
	count := 0
	err := RunSteps(
		[]Step{{Name: "A", Run: func() error { count++; return boom }}, {Name: "B", Run: func() error { count++; return nil }}},
		nil,
		nil,
	)
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want boom", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
}
