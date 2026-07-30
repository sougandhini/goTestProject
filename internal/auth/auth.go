package auth

import "net/http"

// This is the only function that we write in this package
// Utility of this function is to extract the API-key from the headers of the HTTP request, otherwise return error
// eg: Authorization: ApiKey {insert apikey here} -- we are looking for something like this
func GetAPIKey(headers http.Header) (string, error) {

}
