package runn

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/tenntenn/golden"
)

func TestResultOut(t *testing.T) {
	tests := []struct {
		r *runNResult
	}{
		{newRunNResult(t, 4, []*RunResult{
			{
				Path: "testdata/book/runn_0_success.yml",
				Err:  nil,
			},
			{
				Path: "testdata/book/runn_1_fail.yml",
				Err:  ErrDummy,
			},
			{
				Path: "testdata/book/runn_2_success.yml",
				Err:  nil,
			},
			{
				Path: "testdata/book/runn_3.skip.yml",
				Err:  nil,
			},
		})},
		{newRunNResult(t, 5, []*RunResult{
			{
				Path: "testdata/book/runn_0_success.yml",
				Err:  nil,
			},
			{
				Path: "testdata/book/runn_1_fail.yml",
				Err:  ErrDummy,
			},
			{
				Path: "testdata/book/runn_2_success.yml",
				Err:  nil,
			},
			{
				Path: "testdata/book/runn_3.skip.yml",
				Err:  nil,
			},
			{
				Path: "testdata/book/always_failure.yml",
				Err:  nil,
			},
		})},
		{newRunNResult(t, 2, []*RunResult{
			{
				Path: "testdata/book/runn_0_success.yml",
				Err:  nil,
			},
			{
				Path: "testdata/book/runn_1_fail.yml",
				Err:  ErrDummy,
			},
		})},
	}
	for i, tt := range tests {
		key := fmt.Sprintf("result_out_%d", i)
		t.Run(key, func(t *testing.T) {
			got := new(bytes.Buffer)
			if err := tt.r.Out(got); err != nil {
				t.Error(err)
			}
			if os.Getenv("UPDATE_GOLDEN") != "" {
				golden.Update(t, "testdata", key, got)
				return
			}
			if diff := golden.Diff(t, "testdata", key, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestResultOutJSON(t *testing.T) {
	tests := []struct {
		r *runNResult
	}{
		{newRunNResult(t, 4, []*RunResult{
			{
				Path:        "testdata/book/runn_0_success.yml",
				Err:         nil,
				StepResults: []*StepResult{{Key: "0", Err: nil}},
			},
			{
				Path:        "testdata/book/runn_1_fail.yml",
				Err:         ErrDummy,
				StepResults: []*StepResult{{Key: "0", Err: ErrDummy}},
			},
			{
				Path:        "testdata/book/runn_2_success.yml",
				Err:         nil,
				StepResults: []*StepResult{{Key: "0", Err: nil}},
			},
			{
				Path:        "testdata/book/runn_3.skip.yml",
				Err:         nil,
				StepResults: []*StepResult{{Key: "0", Err: nil, Skipped: true}},
			},
		})},
		{newRunNResult(t, 5, []*RunResult{
			{
				Path:        "testdata/book/runn_0_success.yml",
				Err:         nil,
				StepResults: []*StepResult{{Key: "0", Err: nil}},
			},
			{
				Path:        "testdata/book/runn_1_fail.yml",
				Err:         ErrDummy,
				StepResults: []*StepResult{{Key: "0", Err: ErrDummy}},
			},
			{
				Path:        "testdata/book/runn_2_success.yml",
				Err:         nil,
				StepResults: []*StepResult{{Key: "0", Err: nil}},
			},
			{
				Path:        "testdata/book/runn_3.skip.yml",
				Err:         nil,
				StepResults: []*StepResult{{Key: "0", Err: nil, Skipped: true}},
			},
			{
				Path:        "testdata/book/always_failure.yml",
				Err:         ErrDummy,
				StepResults: []*StepResult{{Key: "0", Err: ErrDummy}},
			},
		})},
		{newRunNResult(t, 2, []*RunResult{
			{
				Path:        "testdata/book/runn_0_success.yml",
				Err:         nil,
				StepResults: []*StepResult{{Key: "0", Err: nil}},
			},
			{
				Path:        "testdata/book/runn_1_fail.yml",
				Err:         ErrDummy,
				StepResults: []*StepResult{{Key: "0", Err: ErrDummy}},
			},
		})},
	}
	for i, tt := range tests {
		key := fmt.Sprintf("result_out_json_%d", i)
		t.Run(key, func(t *testing.T) {
			got := new(bytes.Buffer)
			if err := tt.r.OutJSON(got); err != nil {
				t.Error(err)
			}
			if os.Getenv("UPDATE_GOLDEN") != "" {
				golden.Update(t, "testdata", key, got)
				return
			}
			if diff := golden.Diff(t, "testdata", key, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
