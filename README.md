
[![Coverage Status](https://coveralls.io/repos/github/OpenIndustryCloud/fissin-go-weather-data/badge.svg?branch=master)](https://coveralls.io/github/OpenIndustryCloud/fissin-go-weather-data?branch=master)


# TYPEFORM request transformation API

`weather-data.go` uses [Wundergroud APIs](https://www.wunderground.com/weather/api/) to retrieve Weather details for given City/Country and Date.

This API works in 2 steps

- [Autocomplete API](https://www.wunderground.com/weather/api/d/docs?d=autocomplete-api)

Its used Wundergroud Autocomplete API's  feature to list all matching cities based on provided city name and country 

Each autocomplete result object has an l field (short for link) that can be used for constructing wunderground URLs or API calls:

- [Wunderground API History](https://www.wunderground.com/weather/api/d/docs?d=data/history)  feature

Then it uses Wunderground API History feature to query hostorical weather condition for the given date

PS - Country would be defaulted to United Kingdom if not provided in request payload

## API reference

- [Weather History - Wunderground](https://www.wunderground.com/weather/api/d/docs?d=data/history)

GET http://autocomplete.wunderground.com/aq?query=query

-

http://api.wunderground.com/api/Your_Key/history_YYYYMMDD/q/CA/San_Francisco.json


## Authentication

This API do not need any authentication data. 


## Error hanlding

Empty Payload or malformed JSON would result in error reponse.

- Technical error : `{"status":400,"message":"EOF"}`


## Sample Input/Output

- Request payload

```
{"city": "Wales", "country": "", "date": "20171002" }	
```

- Response

```
TBU
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