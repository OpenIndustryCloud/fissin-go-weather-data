package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "k8s.io/client-go/kubernetes"
)

var (
	server *httptest.Server
	//Test Data TV
	userJson = `  {"city": "Wales", "country": "", "date": "20171002" }	`
	// ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
)

func TestHandler(t *testing.T) {
	//Convert string to reader and Create request with JSON body
	req, err := http.NewRequest("POST", "", strings.NewReader(userJson))
	reqEmpty, err := http.NewRequest("POST", "", strings.NewReader(""))
	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	//TEST CASES
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test Data", args{rr, req}},
		{"Empty Data", args{rr, reqEmpty}},
	}
	for _, tt := range tests {
		// call ServeHTTP method
		// directly and pass Request and ResponseRecorder.
		handler := http.HandlerFunc(Handler)
		handler.ServeHTTP(tt.args.w, tt.args.r)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		//check content type
		if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
			t.Errorf("content type header does not match: got %v want %v",
				ctype, "application/json")
		}
		//check if weather datareturned
		//check response content
		res, err := ioutil.ReadAll(rr.Body)
		if err != nil {
			t.Error(err) //Something is wrong while read res
		}

		got := History{}
		err = json.Unmarshal(res, &got)

		if err != nil && got.DailySummary[0].Maxwspdm != "" {
			t.Errorf("%q. compute weather risk() = %v, want %v", tt.name, got, "non empty")
		}

	}
}

func Test_getCityUniqueLink(t *testing.T) {
	type args struct {
		city    string
		country string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test1", args{"wales", "GB"}, "/q/zmw:00000.123.WEGCN", false},
		{"test2", args{"london", "GB"}, "/q/zmw:00000.40.03779", false},
		{"test3", args{"r@nd0m", "GB"}, "No results found", true},
	}
	for _, tt := range tests {
		got, err := getCityUniqueLink(tt.args.city, tt.args.country)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. getCityUniqueLink() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. getCityUniqueLink() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetWeatherConditions(t *testing.T) {

	emptyresponse := `{"response":{"version":"0.1"},"history":{"dailysummary":null,"observations":null}}`
	type args struct {
		link       string
		dateString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		//{"test1", args{"/q/zmw:00000.123.WEGCN", "20170101"}, "/q/zmw:00000.123.WEGCN", false},
		//{"test2", args{"/q/zmw:00000.40.03779", "20170101"}, "/q/zmw:00000.40.03779", false},
		{"test3", args{"", ""}, emptyresponse, false},
	}
	for _, tt := range tests {
		got, err := GetWeatherConditions(tt.args.link, tt.args.dateString)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GetWeatherConditions() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. GetWeatherConditions() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
