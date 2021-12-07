package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"testing"

	s "github.com/sunao-uehara/hashencpassword/storages"
	tt "github.com/sunao-uehara/hashencpassword/testutils"
)

func init() {
}

func createMockPasswordList() {
	s.PasswordList = make(map[int]*s.Password)
	p := &s.Password{
		ID:       1,
		Password: "dummy_hashed_password",
	}
	s.PasswordList[1] = p
}

func createMockStats() {
	s.Stats = make(map[string]*s.EndpointStats)
	es := &s.EndpointStats{
		Total:       3,
		TotalTime:   300.0,
		AverageTime: 100.0,
	}
	s.Stats["POST:/hash"] = es
}

func TestHelloHandler(t *testing.T) {
	testCases := []tt.HandlerTestCase{
		{
			Scenario: "success case /",
			In: tt.HandlerInput{
				Endpoint: "/",
			},
			Out: tt.HandlerOutput{
				ResponseStatusCode: 200,
				ResponseBody:       "hello",
			},
		},
	}

	wg := &sync.WaitGroup{}
	h := NewHandler(wg)
	svr := tt.NewServer()

	for _, testCase := range testCases {
		res, err := svr.GET(testCase.In.Endpoint, h.HelloHandler)
		if err != nil {
			t.Fatalf("http failed: %v", err)
		}

		if res.StatusCode != testCase.Out.ResponseStatusCode {
			t.Errorf("handler returned unexpected http stattus: got %v, want %v", res.StatusCode, testCase.Out.ResponseStatusCode)
		}

		want := testCase.Out.ResponseBody
		got, err := tt.ParseResponseBody(res)
		if err != nil {
			t.Fatalf("parse response body failed: %v", err)
		}

		if want != got {
			t.Errorf("scenario %s: handler returned unexpected body: got %v, want %v", testCase.Scenario, got, want)
		}
	}
}

func TestHashPostHandler(t *testing.T) {
	testCases := []tt.HandlerTestCase{
		{
			Scenario: "success case /hash",
			In: tt.HandlerInput{
				Endpoint:    "/hash",
				RequestBody: "angryMonkey",
			},
			Out: tt.HandlerOutput{
				ResponseStatusCode: 200,
				ResponseBody:       "1",
			},
		},
		{
			Scenario: "success case 2 /hash",
			In: tt.HandlerInput{
				Endpoint:    "/hash",
				RequestBody: "angryMonkey2",
			},
			Out: tt.HandlerOutput{
				ResponseStatusCode: 200,
				ResponseBody:       "2",
			},
		},
		{
			Scenario: "failure case /hash",
			In: tt.HandlerInput{
				Endpoint:    "/hash",
				RequestBody: "",
			},
			Out: tt.HandlerOutput{
				ResponseStatusCode: http.StatusBadRequest,
				ResponseBody:       "password field is required",
			},
		},
	}

	wg := &sync.WaitGroup{}
	h := NewHandler(wg)
	svr := tt.NewServer()
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	svr.AddHeader(headers)

	for _, testCase := range testCases {
		form := url.Values{}
		form.Add("password", testCase.In.RequestBody)

		res, err := svr.POST(testCase.In.Endpoint, h.HashPostHandler, form)
		if err != nil {
			t.Fatalf("http failed: %v", err)
		}

		if res.StatusCode != testCase.Out.ResponseStatusCode {
			t.Errorf("handler returned unexpected http stattus: got %v, want %v", res.StatusCode, testCase.Out.ResponseStatusCode)
		}

		want := testCase.Out.ResponseBody
		got, err := tt.ParseResponseBody(res)
		if err != nil {
			t.Fatalf("parse response body failed: %v", err)
		}

		if want != got {
			t.Errorf("scenario %s: handler returned unexpected body: got %v, want %v", testCase.Scenario, got, want)
		}
	}
}

func TestHashGetHandler(t *testing.T) {
	createMockPasswordList()
	testCases := []tt.HandlerTestCase{
		{
			Scenario: "success case /hash/1",
			In: tt.HandlerInput{
				Endpoint: "/hash/1",
			},
			Out: tt.HandlerOutput{
				ResponseStatusCode: 200,
				ResponseBody:       "dummy_hashed_password",
			},
		},
		{
			Scenario: "failure case /hash/invalidid",
			In: tt.HandlerInput{
				Endpoint: "/hash/invalidid",
			},
			Out: tt.HandlerOutput{
				ResponseStatusCode: http.StatusBadRequest,
				ResponseBody:       "invalid id",
			},
		},
		{
			Scenario: "failure case /hash/100",
			In: tt.HandlerInput{
				Endpoint: "/hash/100",
			},
			Out: tt.HandlerOutput{
				ResponseStatusCode: http.StatusNotFound,
				ResponseBody:       "id 100 not found",
			},
		},
	}

	wg := &sync.WaitGroup{}
	h := NewHandler(wg)
	svr := tt.NewServer()

	for _, testCase := range testCases {
		res, err := svr.GET(testCase.In.Endpoint, h.HashGetHandler)
		if err != nil {
			t.Fatalf("http failed: %v", err)
		}

		if res.StatusCode != testCase.Out.ResponseStatusCode {
			t.Errorf("handler returned unexpected http stattus: got %v, want %v", res.StatusCode, testCase.Out.ResponseStatusCode)
		}

		want := testCase.Out.ResponseBody
		got, err := tt.ParseResponseBody(res)
		if err != nil {
			t.Fatalf("parse response body failed: %v", err)
		}

		if want != got {
			t.Errorf("scenario %s: handler returned unexpected body: got %v, want %v", testCase.Scenario, got, want)
		}
	}
}

func TestStatsGetHandler(t *testing.T) {
	createMockStats()
	testCases := []tt.HandlerTestCase{
		{
			Scenario: "success case /stats",
			In: tt.HandlerInput{
				Endpoint: "/stats",
			},
			Out: tt.HandlerOutput{
				ResponseStatusCode: 200,
				ResponseBodyStruct: s.Stats["POST:/hash"],
			},
		},
		// {
		// 	Scenario: "failure case",
		// 	In: tt.HandlerInput{
		// 		Endpoint: "/stats",
		// 	},
		// 	Out: tt.HandlerOutput{
		// 		ResponseStatusCode: http.StatusBadRequest,
		// 		ResponseBody:       "stats POST:/hash not found",
		// 	},
		// },
	}

	wg := &sync.WaitGroup{}
	h := NewHandler(wg)
	svr := tt.NewServer()

	for _, testCase := range testCases {
		res, err := svr.GET(testCase.In.Endpoint, h.StatsGetHandler)
		if err != nil {
			t.Fatalf("http failed: %v", err)
		}

		if res.StatusCode != testCase.Out.ResponseStatusCode {
			t.Errorf("handler returned unexpected http stattus: got %v, want %v", res.StatusCode, testCase.Out.ResponseStatusCode)
		}

		tmp, err := json.Marshal(testCase.Out.ResponseBodyStruct)
		if err != nil {
			t.Fatalf("cannot convert struct %v to json", testCase.Out.ResponseBodyStruct)
		}
		want := string(tmp)

		got, err := tt.ParseResponseBody(res)
		if err != nil {
			t.Fatalf("parse response body failed: %v", err)
		}

		if bool, err := tt.EqualJSON(want, got); !bool || err != nil {
			t.Errorf("scenario %s: handler returned unexpected body: got %v, want %v", testCase.Scenario, got, want)
		}
	}
}
