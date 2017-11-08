package main

import (
	"net/http/httptest"
	"testing"

	_ "k8s.io/client-go/kubernetes"
)

var (
	server *httptest.Server
	//Test Data TV
	userJson = ` {"tranformedData":{"ticketDetails":{"ticket":{"comment":{"html_body":"\u003cp\u003e\u003cb\u003eIf there has been any recent maintenance carried out on your home, please describe it\u003c/b\u003e : No maintenance carried out\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eIf you have any other insurance or warranties covering your home, please advise us of the company name.\u003c/b\u003e : No\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eWe have made the following assumptions about your property, you and anyone living with you\u003c/b\u003e : \u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eWhen did the incident happen?\u003c/b\u003e : 2017-01-01\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eAre you still have possession of the damage items (i.e. damaged guttering)?\u003c/b\u003e : \u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eAre you aware of anything else relevant to your claim that you would like to advise us of at this stage?\u003c/b\u003e : I would need the vendors contact for repairing the roof\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eWould you like to upload more images?\u003c/b\u003e : \u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eWhere did the incident happen? (City/town name)\u003c/b\u003e : birmingham\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003eIn as much detail as possible, please use the text box below to describe the full extent of the damage to your home and how you discovered it.\u003c/b\u003e : Roof Damaged\u003c/p\u003e\u003chr\u003e\u003cp\u003e\u003cb\u003ePlease describe the details of the condition of your home prior to discovering the damage\u003c/b\u003e : Tiles blown away\u003c/p\u003e\u003chr\u003e"},"custom_fields":null,"email":"amitkumarvarman@gmail.com","phone":"09876512345","priority":"normal","requester":{"email":"amitkumarvarman@gmail.com","locale_id":1,"name":"Amit Varman"},"status":"new","subject":"Storm surge risk data","type":"incident"}},"weatherAPIInput":{"city":"birmingham","country":"","date":"20170101"}},"weatherData":{"history":{"dailysummary":[{"fog":"0","maxpressurem":"1025","maxtempm":"7","maxwspdm":"28","minpressurem":"1014","mintempm":"0","minwspdm":"7","rain":"1","tornado":"0"}]},"response":{"version":"0.1"}},"weatherRisk":{"description":"Possible Stormy weather","riskScore":50}}`
	// ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
)

/*func TestHandler(t *testing.T) {
	//Convert string to reader and
	//Create request with JSON body
	req, err := http.NewRequest("POST", "", strings.NewReader(userJson))
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

	}
}*/

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

	emptyresponse := `{"response":{"version":""},"history":{"dailysummary":null,"observations":null}}`
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
