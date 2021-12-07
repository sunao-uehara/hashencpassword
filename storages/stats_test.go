package storages

import (
	"errors"
	"reflect"
	"testing"
	// tt "github.com/sunao-uehara/hashencpassword/testutils"
)

func init() {
}

func createMockStats() {
	Stats = make(map[string]*EndpointStats)
	s := &EndpointStats{
		Total:       1,
		TotalTime:   100.0,
		AverageTime: 100.0,
	}
	Stats["endpoint1"] = s
}

func TestSaveStats(t *testing.T) {
	createMockStats()
	type in struct {
		key     string
		elapsed float64
	}
	type out struct {
		Expected *EndpointStats
	}
	type testCase struct {
		Scenario string
		In       *in
		Out      *out
	}

	testCases := []testCase{
		{
			"success case, create new stats",
			&in{
				key:     "new_endpoint1",
				elapsed: 100.0,
			},
			&out{
				Expected: &EndpointStats{
					Total:       1,
					TotalTime:   100.0,
					AverageTime: 100.0,
				},
			},
		},
		{
			"success case, update stats",
			&in{
				key:     "endpoint1",
				elapsed: 200.0,
			},
			&out{
				Expected: &EndpointStats{
					Total:       2,
					TotalTime:   300.0,
					AverageTime: 150.0,
				},
			},
		},
	}

	for _, testCase := range testCases {
		in := testCase.In
		out := testCase.Out
		SaveStats(in.key, in.elapsed)
		if !reflect.DeepEqual(out.Expected, Stats[in.key]) {
			t.Errorf("test failed, got: %v, want: %v", Stats[in.key], testCase.Out.Expected)
		}
	}
}

func TestGetStats(t *testing.T) {
	createMockStats()
	type in struct {
		key string
	}
	type out struct {
		Expected *EndpointStats
		Error    error
	}
	type testCase struct {
		Scenario string
		In       *in
		Out      *out
	}

	testCases := []testCase{
		{
			"success case, create new stats",
			&in{
				key: "endpoint1",
			},
			&out{
				Expected: &EndpointStats{
					Total:       1,
					TotalTime:   100.0,
					AverageTime: 100.0,
				},
				Error: nil,
			},
		},
		{
			"failure case, stats not found",
			&in{
				key: "badEndpoints1",
			},
			&out{
				Expected: nil,
				Error:    errors.New("stats badEndpoints1 not found"),
			},
		},
	}

	for _, testCase := range testCases {
		in := testCase.In
		out := testCase.Out
		ret, err := GetStats(in.key)
		if !reflect.DeepEqual(out.Expected, ret) {
			t.Errorf("test failed, got: %v, want: %v", ret, out.Expected)
		}
		switch {
		case err != nil && out.Error == nil:
			t.Errorf("expected non error, but some error occurred, %s", err.Error())
		case err == nil && out.Error != nil:
			t.Errorf("expected error %s, but results: no error", out.Error.Error())
		case err != nil && out.Error != nil:
		}

	}
}
