
[![Coverage Status](https://coveralls.io/repos/github/OpenIndustryCloud/fissin-go-weather-data/badge.svg?branch=master)](https://coveralls.io/github/OpenIndustryCloud/fissin-go-weather-data?branch=master)


# Weather Data API

`weather-data.go` uses [Wundergroud APIs](https://www.wunderground.com/weather/api/) to retrieve Weather details for given City/Country and Date.

This API works in 2 steps

- [Autocomplete API](https://www.wunderground.com/weather/api/d/docs?d=autocomplete-api) : Its used Wundergroud Autocomplete API's  feature to list locations or hurricanes which match against a partial query.

API End Point :   http://autocomplete.wunderground.com/aq?query=query

Each autocomplete result object has an l field (short for link) that can be used for constructing wunderground URLs or API calls:

- [Wunderground API History](https://www.wunderground.com/weather/api/d/docs?d=data/history)  feature
 : It uses Wunderground API History feature to query hostorical weather condition for the given date

http://api.wunderground.com/api/Your_Key/history_YYYYMMDD/q/CA/San_Francisco.json

PS - Country would be defaulted to United Kingdom if not provided in request payload

## API reference

- [Weather History - Wunderground](https://www.wunderground.com/weather/api/d/docs?d=data/history)

- [Autocomplete API - Wunderground](https://www.wunderground.com/weather/api/d/docs?d=autocomplete-api)


## Authentication


Authentication is implemented using API token, you can either configure a [Secret in Kubernetes](https://kubernetes.io/docs/concepts/configuration/secret/) 
or  have it directly configured within the API (not recommended as it exposes your secrets)


Signup and create [API Key](https://www.wunderground.com/weather/api)


## Error hanlding

Empty Payload or malformed JSON would result in error reponse.

- Technical error :  `{"status":400,"message":"EOF"}`
- No results found : `{"status":400,"message":"No results found"}`
- Malformed JSON `{"status":400,"message":"invalid character 'a' looking for beginning of object key string"}`

## Sample Input/Output

- Request payload

City :  Name of the City

Country : Country Code 

Date : Date in YYYYMMDD or YYYY/MM/DD format 


```
with Country Code

{"city": "Wales", "country": "GB", "date": "20171002" }	

without Country Code

{"city": "Wales", "date": "20171002" }	

with Country Code in YYYY/MM/DD format

{"city": "Wales", "country": "GB", "date": "2017/10/02" }	

```
- Response

```
{"status":200,"response":{"version":"0.1"},"history":{"dailysummary":[{"fog":"0","rain":"1","maxtempm":"17","mintempm":"12","tornado":"0","maxpressurem":"1014","minpressurem":"1005","maxwspdm":"50","minwspdm":"13"}],"observations":[{"tempm":"17.0","tempi":"62.6","dewptm":"13.0","dewpti":"55.4","hum":"77","wspdm":"29.6","wspdi":"18.4","wgustm":"-9999.0","wgusti":"-9999.0","wdird":"230","wdire":"SW","vism":"10.0","visi":"6.2","pressurem":"1005","pressurei":"29.68","windchillm":"-999","windchilli":"-999","heatindexm":"-9999",

//
//
dewpti":"48.2","hum":"82","wspdm":"18.5","wspdi":"11.5","wgustm":"-9999.0","wgusti":"-9999.0","wdird":"270","wdire":"West","vism":"10.0","visi":"6.2","pressurem":"1014","pressurei":"29.95","windchillm":"-999","windchilli":"-999","heatindexm":"-9999","heatindexi":"-9999","precipm":"-9999.00","precipi":"-9999.00","conds":"Scattered Clouds","icon":"partlycloudy","fog":"0","rain":"0","snow":"0","hail":"0","thunder":"0","tornado":"0","metar":"METAR EGCN 022250Z 27010KT 9999 SCT025 12/09 Q1014"}]}}
```


## Example Usage

## 1.  Deploy as Fission Functions

First, set up your fission deployment with the go environment.

```
fission env create --name go-env --image fission/go-env:1.8.1
```

To ensure that you build functions using the same version as the
runtime, fission provides a docker image and helper script for
building functions.



- Download the build helper script

```
$ curl https://raw.githubusercontent.com/fission/fission/master/environments/go/builder/go-function-build > go-function-build
$ chmod +x go-function-build
```

- Build the function as a plugin. Outputs result to 'function.so'

`$ go-function-build form-req-transformer.go`

- Upload the function to fission

`$ fission function create --name form-req-transformer --env go-env --package function.so`

- Map /form-req-transformer to the form-req-transformer function

`$ fission route create --method POST --url /form-req-transformer --function form-req-transformer`

- Run the function

```$ curl -d `{"ticket": {"subject": "My printer is on fire!", "comment": {"body": "The smoke is very colorful."}}}` -H "Content-Type: application/json" -X POST http://$FISSION_ROUTER/form-req-transformer```

## 2. Deploy as AWS Lambda

> to be updated