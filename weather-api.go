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
"fog": "0",
"rain": "1",
"maxtempm": "7",
"mintempm": "0",
"tornado": "0",
"maxpressurem": "1025",
"minpressurem": "1014",
"maxwspdm": "28",
"minwspdm": "7"
}

*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//Default values , this can be overridden by setting ENV variables
var (
	apiURL          = "http://api.wunderground.com/api"
	autocompleteURL = "http://autocomplete.wunderground.com"
	apiKey          = ""
	namesapce       = "default"
	secretName      = "wunderground-secret"
)

func getAPIKeys(w http.ResponseWriter) {
	fmt.Println("[CONFIG] Reading Kubernetes secrets")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	secret, err := clientset.Core().Secrets(namesapce).Get(secretName, meta_v1.GetOptions{})
	fmt.Println("Wunderground Desk API Key : " + string(secret.Data[apiKey]))

	apiKey = string(secret.Data["apiKey"])

	if len(apiKey) == 0 {
		createErrorResponse(w, "Missing API Key", "400")
	}

}

type InputData struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Date    string `json:"date"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var inputData InputData
	err := json.NewDecoder(r.Body).Decode(&inputData)
	if err == io.EOF || err != nil {
		createErrorResponse(w, err.Error(), "400")
		return
	}
	fmt.Println("Query Weather Data for", inputData.City, inputData.Country, inputData.Date)

	//use Wundergroud AutoComplete API to get unique city link
	link, err := getCityUniqueLink(inputData.City, inputData.Country)

	//use Wundergroud API to retrieve Historical Data
	weatherDataJSON, err := GetWeatherConditions(link, inputData.Date)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(weatherDataJSON))

}

func getCityUniqueLink(city string, country string) (string, error) {

	//Query AutoCmplete API
	if len(country) == 0 {
		println("Country not specified to defaulting to UK...")
		country = "GB" //default country to United kingdom if blank
	}
	autocompleteURL = autocompleteURL + "/aq?query=" + url.QueryEscape(city) + "&c=" + url.QueryEscape(country)
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
		return "", errors.New("No results found")
	}

	link := acObj.Results[0].Link
	println("link --- ", link)

	return link, err
}

// GetLocalConditions returns weather summary for given date
func GetWeatherConditions(link string, dateString string) (string, error) {

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
	if err != nil && len(historicalData.History.DailySummary) != 0 {
		println(err)
		return "no data retrieved for " + historicalDataURL, err
	}

	//dailySummary := historicalData.History.DailySummary[0]

	//marshal to JSON
	historicalDataJSON, err := json.Marshal(historicalData)
	if err != nil {
		println(err)
		return "", err
	}

	return string(historicalDataJSON), nil
}

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

type HistoricalData struct {
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

func createErrorResponse(w http.ResponseWriter, message string, status string) {
	errorJSON, _ := json.Marshal(&Error{
		Code:    status,
		Message: message})
	//Send custom error message to caller
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(errorJSON))
}

type Error struct {
	Code    string `json:"status"`
	Message string `json:"message"`
}
