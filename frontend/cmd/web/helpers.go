package main

import (
	"fmt"
	"math"
	"net/http"
	"runtime/debug"
)

// NewHelpers sets up app config for helpers
func NewHelpers(a *AppConfig) {
	app = *a
}

func ClientError(w http.ResponseWriter, status int) {

	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)

}

func ServerError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

func IsAuthenticated(r *http.Request) bool {

	exists := true

	return exists

}

// RoundFloat allows for rounding of a float to a given number of decimal places.
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// RemoveDuplicateStr removes duplicate strings from a slice of strings
func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// DeleteEmpty returns a slice of strings without empty strings.
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
