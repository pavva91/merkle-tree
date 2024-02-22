package enums

import (
	"testing"
)

func TestTaskStatus_Itoa(t *testing.T) {
	tests := map[string]struct {
		tr      TaskStatus
		wantStr string
	}{
		"case new": {
			New,
			"new",
		},
		"case in_process": {
			InProcess,
			"in_process",
		},
		"case done": {
			Done,
			"done",
		},
		"case error": {
			Error,
			"error",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if gotStr := tt.tr.Itoa(); gotStr != tt.wantStr {
				t.Errorf("TaskStatus.ToString() = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}
