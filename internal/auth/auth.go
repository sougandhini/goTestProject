package auth

import (
	"errors"
	"net/http"
	"strings"
)

// This is the only function that we write in this package
// Utility of this function is to extract the API-key from the headers of the HTTP request, otherwise return error
// eg: Authorization: ApiKey {insert apikey here} -- we are looking for something like this
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization") //Note you're searching it in header and not in authentication column
	if val == ""{
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2{
		return "", errors.New("malformed auth header") //coz we expect the value to be {apiKey, valOfApiKey}
	}
	if vals[0] != "ApiKey"{
		return "", errors.New("malformed first part of the error")
	}
	return vals[1], nil
}
