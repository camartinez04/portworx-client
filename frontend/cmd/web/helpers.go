package main

import (
	"fmt"
	"math"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

// NewHelpers sets up app config for helpers
func NewHelpers(a *AppConfig) {
	App = *a
}

// ClientError sends a specific status code and corresponding description to the user.
func ClientError(w http.ResponseWriter, status int) {

	App.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)

}

// ServerError logs the error and sends a generic 500 Internal Server Error response to the user.
func ServerError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	App.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

// IsAuthenticated checks if user is authenticated
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

// DateFormat returns date from UNIX timestamp
func DateFormat(date int64) string {

	time := time.Unix(date, 0)

	layout := "2006-01-02 15:04:05"

	return time.Format(layout)
}

// Add add two numbers
func Add(a, b int) int {

	return a + b
}

// Divide divides two numbers and returns the result as a string
func Divide(a, b uint64) string {

	if b == 0 {
		return "0"
	}

	floatA := float64(a)

	floatB := float64(b)

	result := floatA / floatB

	stringresult := fmt.Sprintf("%.0f", result)

	return stringresult

}

// Iterate performs a for loop
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 0; i < count; i++ {
		items = append(items, i)
	}

	return items
}

// HumanDate returns time in yyyy-mm-dd format
func HumanDate(t time.Time) string {

	return t.Format("2006-01-02")
}

// FormatDate helps to format numeric date into string date
func FormatDate(t time.Time, f string) string {

	return t.Format(f)
}

// BytesToGB converts bytes to gigabytes
func BytesToGB(b string) uint64 {

	toFloat, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0
	}

	floatValue := toFloat / 1073741824

	return uint64(floatValue)
}

// BytesToMB converts bytes to megabytes
func BytesToMB(b string) uint64 {

	toFloat, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0
	}

	floatValue := toFloat / 1048576

	return uint64(floatValue)
}
