package main

/*
This API will collect Weather Data by consuming
Wunderground API and return summary for given date and city

---INPUT---
{
  "city": "birmingham",
  "country": "",
  "date": "20170101"
}
---OUTPUT---
{
    "response": {
        "version": "0.1"
    },
    "history": {
        "dailysummary": [
            {
                "fog": "0",
                "rain": "1",
                "maxtempm": "17",
                "mintempm": "12",
                "tornado": "0",
                "maxpressurem": "1014",
                "minpressurem": "1005",
                "maxwspdm": "50",
                "minwspdm": "13"
            }
        ],
        "observations": [
            {

            }
        ]
    }
}
*/

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//Default values , this can be overridden by setting ENV variables
var (
	apiURL     = "http://api.wunderground.com/api" //Wunderground  API endpoint
	apiKey     = ""                                //wundergroud API Key
	namesapce  = "default"                         //Kubernetes virtual clusters Name to read secrets
	secretName = "wunderground-secret"             // secret name
)

// Handler - this is main function, which will to prcess the incoming data
// and and query wunderground APIs
func Handler(w http.ResponseWriter, r *http.Request) {

	//process post data
	var inputData InputData
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err == io.EOF || err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	println("Query Weather Data for", inputData.City, inputData.Country, inputData.Date)

	//get API keys from Kubernetes Secrets
	getAPIKeys(w)
	//use Wundergroud AutoComplete API to get unique city link
	link, err := getCityUniqueLink(inputData.City, inputData.Country)

	//use Wundergroud API to retrieve Historical Data
	weatherDataJSON, err := getWeatherConditions(link, inputData.Date)
	if err != nil || len(weatherDataJSON) == 0 {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	//return Weather Data in JSON format
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(weatherDataJSON))

}

//[Autocomplete API](https://www.wunderground.com/weather/api/d/docs?d=autocomplete-api) :
// Its used Wundergroud Autocomplete API's  feature to list locations
// or hurricanes which match against a partial query.

//Each autocomplete result object has an l field (short for link)
// that can be used for constructing wunderground URLs or API calls
func getCityUniqueLink(city string, country string) (string, error) {

	autocompleteURL := "http://autocomplete.wunderground.com"
	//Query AutoCmplete API
	if len(country) == 0 {
		println("Country not specified to defaulting to UK...")
		country = "GB" //default country to United kingdom if blank
	}
	autocompleteURL += "/aq?query=" + url.QueryEscape(city) + "&c=" + url.QueryEscape(country)
	println("autocompleteURL : ", autocompleteURL)
	acResp, err := http.Get(autocompleteURL)
	if err != nil {
		return "", err
	}
	defer acResp.Body.Close()

	acObj := autocomplete{}
	err = json.NewDecoder(acResp.Body).Decode(&acObj)

	if err != nil {
		return "", err
	}

	if len(acObj.Results) == 0 {
		println("No result found for ", city)
		return "No results found", errors.New("No results found")
	}

	link := acObj.Results[0].Link
	println("link --- ", link)

	return link, err
}

// It uses Wunderground API History feature to query
// hostorical weather condition for the given date
// accept dateString in YYYYMMDD or YYYY/MM/DD format
func getWeatherConditions(link string, dateString string) (string, error) {

	//covert date to YYYYMMDD format
	if strings.Contains(dateString, "-") {
		dateString = strings.Replace(dateString, "-", "", 2)
	} else if strings.Contains(dateString, "/") {
		dateString = strings.Replace(dateString, "/", "", 2)
	}
	//form API URL for Historical Data
	historicalDataURL := apiURL + "/" + apiKey + "/history_" + url.QueryEscape(dateString) + link + ".json"

	println("historicalDataURL being queried : ", historicalDataURL)
	repsonse, err := http.Get(historicalDataURL)
	if err != nil {
		return "", err
	}
	defer repsonse.Body.Close()
	println("response Status for weatherAPI :", repsonse.Status)

	var historicalData HistoricalData
	err = json.NewDecoder(repsonse.Body).Decode(&historicalData)
	//Validate response
	if err != nil {
		return "" + historicalDataURL, err
	}

	//dailySummary := historicalData.History.DailySummary[0]

	//marshal to JSON
	if len(historicalData.History.DailySummary) == 0 {
		return "", errors.New("No results found")
	}
	historicalData.Status = 200
	historicalDataJSON, err := json.Marshal(historicalData)
	if err != nil {
		println(err)
		return "", err
	}

	//Historical data in JSON format
	return string(historicalDataJSON), nil
}

// getAPIKeys - this funtion read kubernetes secrets for configured
// namespace and secret name
func getAPIKeys(w http.ResponseWriter) {
	println("[CONFIG] Reading Kubernetes secrets")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		createErrorResponse(w, err.Error(), http.StatusBadRequest)
	}

	//read kubernetes secrets
	secret, err := clientset.Core().Secrets(namesapce).Get(secretName, meta_v1.GetOptions{})
	println("Wunderground Desk API Key : " + string(secret.Data["apiKey"]))

	apiKey = string(secret.Data["apiKey"])

	//validate if apiKey exist
	if len(apiKey) == 0 {
		createErrorResponse(w, "Missing API Key", http.StatusBadRequest)
	}

}

// createErrorResponse - this function forms a error reposne with
// error message and http code
func createErrorResponse(w http.ResponseWriter, message string, status int) {
	errorJSON, _ := json.Marshal(&Error{
		Status:  status,
		Message: message})
	//Send custom error message to caller
	w.WriteHeader(status)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(errorJSON))
}

// Error - error object
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Input Struct
type InputData struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Date    string `json:"date"`
}

// Autocomplete API Model
type autocomplete struct {
	Results []autocompleteResult `json:"RESULTS"`
}

type autocompleteResult struct {
	Link string `json:"l"`
}

type displaylocation struct {
	Name string `json:"full"`
}

//Model for WeatherAPI

type WeatherAPIInput struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Date    string `json:"date"`
}

// Historical Data API Model
type HistoricalData struct {
	Status   int      `json:"status"`
	Response Response `json:"response"`
	History  History  `json:"history"`
}

type Response struct {
	Version string `json:"version"`
}

type History struct {
	DailySummary []DailySummary `json:"dailysummary"`
	Observations []struct {
		Tempm      string `json:"tempm"`
		Tempi      string `json:"tempi"`
		Dewptm     string `json:"dewptm"`
		Dewpti     string `json:"dewpti"`
		Hum        string `json:"hum"`
		Wspdm      string `json:"wspdm"`
		Wspdi      string `json:"wspdi"`
		Wgustm     string `json:"wgustm"`
		Wgusti     string `json:"wgusti"`
		Wdird      string `json:"wdird"`
		Wdire      string `json:"wdire"`
		Vism       string `json:"vism"`
		Visi       string `json:"visi"`
		Pressurem  string `json:"pressurem"`
		Pressurei  string `json:"pressurei"`
		Windchillm string `json:"windchillm"`
		Windchilli string `json:"windchilli"`
		Heatindexm string `json:"heatindexm"`
		Heatindexi string `json:"heatindexi"`
		Precipm    string `json:"precipm"`
		Precipi    string `json:"precipi"`
		Conds      string `json:"conds"`
		Icon       string `json:"icon"`
		Fog        string `json:"fog"`
		Rain       string `json:"rain"`
		Snow       string `json:"snow"`
		Hail       string `json:"hail"`
		Thunder    string `json:"thunder"`
		Tornado    string `json:"tornado"`
		Metar      string `json:"metar"`
	} `json:"observations"`
}

type DailySummary struct {
	Fog          string `json:"fog"`
	Rain         string `json:"rain"`
	Maxtempm     string `json:"maxtempm"`
	Mintempm     string `json:"mintempm"`
	Tornado      string `json:"tornado"`
	Maxpressurem string `json:"maxpressurem"`
	Minpressurem string `json:"minpressurem"`
	Maxwspdm     string `json:"maxwspdm"`
	Minwspdm     string `json:"minwspdm"`
}

// func main() {
// 	println("staritng app..")
// 	http.HandleFunc("/", Handler)
// 	http.ListenAndServe(":8084", nil)
// }
