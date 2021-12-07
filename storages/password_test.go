package storages

import (
	"errors"
	"reflect"
	"testing"
	// tt "github.com/sunao-uehara/hashencpassword/testutils"
)

func init() {
}

func createMockPasswordList() {
	PasswordList = make(map[int]*Password)
	p := &Password{
		ID:       1,
		Password: "",
	}
	PasswordList[1] = p
}

func TestCreatePasswordID(t *testing.T) {
	createMockPasswordList()
	type out struct {
		Expected int
	}
	type testCase struct {
		Scenario string
		Out      *out
	}

	testCases := []testCase{
		{
			"success case 1",
			&out{
				Expected: 2,
			},
		},
		{
			"success case 2",
			&out{
				Expected: 3,
			},
		},
	}

	for _, testCase := range testCases {
		id := CreatePasswordID()
		if testCase.Out.Expected != id {
			t.Errorf("test failed, got: %v, want: %v", id, testCase.Out.Expected)
		}
	}
}

func TestUpdatePassword(t *testing.T) {
	createMockPasswordList()
	type in struct {
		id       int
		password string
	}
	type out struct {
		// Expected nil
		Error error
	}
	type testCase struct {
		Scenario string
		In       *in
		Out      *out
	}

	testCases := []testCase{
		{
			"success case",
			&in{
				id:       1,
				password: "angryMonkey",
			},
			&out{
				Error: nil,
			},
		},
		{
			"failure case, id not found",
			&in{
				id:       100,
				password: "angryMonkey",
			},
			&out{
				Error: errors.New("id 100 not found"),
			},
		},
		{
			"failure case, password is empty",
			&in{
				id:       1,
				password: "",
			},
			&out{
				Error: errors.New("password is empty"),
			},
		},
	}

	for _, testCase := range testCases {
		in := testCase.In
		out := testCase.Out
		err := UpdatePassword(in.id, in.password)
		switch {
		case err != nil && out.Error == nil:
			t.Errorf("expected non error, but some error occurred, %s", err.Error())
		case err == nil && out.Error != nil:
			t.Errorf("expected error %s, but results: no error", out.Error.Error())
		case err != nil && out.Error != nil:
		}
	}
}

func TestGetPassword(t *testing.T) {
	createMockPasswordList()
	type in struct {
		id int
	}
	type out struct {
		Expected *Password
		Error    error
	}
	type testCase struct {
		Scenario string
		In       *in
		Out      *out
	}

	testCases := []testCase{
		{
			"success case",
			&in{
				id: 1,
			},
			&out{
				Expected: PasswordList[1],
				Error:    nil,
			},
		},
		{
			"failure case, id not found",
			&in{
				id: 100,
			},
			&out{
				Expected: nil,
				Error:    errors.New("id 100 not found"),
			},
		},
	}

	for _, testCase := range testCases {
		in := testCase.In
		out := testCase.Out
		ret, err := GetPassword(in.id)
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

func TestHashPassword(t *testing.T) {
	createMockPasswordList()
	type in struct {
		password string
	}
	type out struct {
		Expected string
		Error    error
	}
	type testCase struct {
		Scenario string
		In       *in
		Out      *out
	}

	testCases := []testCase{
		{
			"success case",
			&in{
				password: "angryMonkey",
			},
			&out{
				Expected: "this_is_dummy_password",
				Error:    nil,
			},
		},
		{
			"failure case, password is empty",
			&in{
				password: "",
			},
			&out{
				Expected: "",
				Error:    errors.New("password is empty"),
			},
		},
	}

	for _, testCase := range testCases {
		in := testCase.In
		out := testCase.Out
		ret, err := hashPassword(in.password)
		if out.Expected == "" && ret != "" {
			t.Errorf("test failed, got: %v, want: %v", ret, out.Expected)
		}
		if len(out.Expected) > 0 && len(ret) == 0 {
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
