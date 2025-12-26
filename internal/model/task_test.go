package model

import (
	"strings"
	"testing"
)

func TestTaskCreateRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   TaskCreateRequest
		wantErr bool
	}{
		{
			name: "valid request",
			input: TaskCreateRequest{
				Title:       "My task",
				Description: "Some description",
			},
			wantErr: false,
		},
		{
			name: "missing title",
			input: TaskCreateRequest{
				Title:       "",
				Description: "Some description",
			},
			wantErr: true,
		},
		{
			name: "missing title and description (still only title error)",
			input: TaskCreateRequest{
				Title:       "",
				Description: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.input.Validate()

			if tt.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation errors, got none")
			}

			if !tt.wantErr && len(errs) > 0 {
				t.Fatalf("expected no validation errors, got %d", len(errs))
			}
		})
	}
}

func TestTaskUpdateRequest_Validate(t *testing.T) {
	tests := []struct {
		name       string
		request    TaskUpdateRequest
		wantErrLen int
		wantErrMsg map[string]string
	}{
		{
			name: "valid input",
			request: TaskUpdateRequest{
				Title:       "Test Task",
				Description: "Some description",
				Status:      "done",
			},
			wantErrLen: 0,
			wantErrMsg: nil,
		},
		{
			name: "missing title",
			request: TaskUpdateRequest{
				Title:       "",
				Description: "Some description",
				Status:      "todo",
			},
			wantErrLen: 1,
			wantErrMsg: map[string]string{
				"title": "field is required",
			},
		},
		{
			name: "invalid status",
			request: TaskUpdateRequest{
				Title:       "Task",
				Description: "Some description",
				Status:      "INVALID_STATUS",
			},
			wantErrLen: 1,
			wantErrMsg: map[string]string{
				"status": "invalid status value",
			},
		},
		{
			name: "missing title and invalid status",
			request: TaskUpdateRequest{
				Title:       "",
				Description: "Some description",
				Status:      "INVALID_STATUS",
			},
			wantErrLen: 2,
			wantErrMsg: map[string]string{
				"title":  "field is required",
				"status": "invalid status value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.request.Validate()
			if len(errs) != tt.wantErrLen {
				t.Errorf("expected %d errors, got %d", tt.wantErrLen, len(errs))
			}
			for _, e := range errs {
				if msg, ok := tt.wantErrMsg[e.Field]; ok {
					if !strings.Contains(e.Message, msg) {
						t.Errorf("expected message for field %s: %s, got %s", e.Field, msg, e.Message)
					}
				} else {
					t.Errorf("unexpected error field: %s", e.Field)
				}
			}
		})
	}
}
