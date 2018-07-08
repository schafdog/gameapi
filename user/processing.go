package User

import (
	"net/http"
	"strconv"
)

// FormToUser -- fills a User struct with submitted form data
// params:
// r - request reader to fetch form data or url params (unused here)
// returns:
// User struct if successful
// array of strings of errors if any occur during processing
func FormToUser(r *http.Request) (User, []string) {
	var user User
	var errStr, hsStr string
	var errs []string
	var err error
	var highscore int

	user.Name, errStr = processFormField(r, "name")
	errs = appendError(errs, errStr)
	hsStr, errStr = processFormField(r, "highscore")
	if len(errStr) != 0 {
		user.Highscore = 0
	} else {
		highscore, err = strconv.Atoi(hsStr)
		if err != nil {
			errs = append(errs, "Parameter 'highscore' not an integer")
		} else if highscore < 0 {
			errs = append(errs, "Parameter 'highscore' is below zero")
		} else {
			user.Highscore = highscore
		}
	}
	return user, errs
}

func appendError(errs []string, errStr string) []string {
	if len(errStr) > 0 {
		errs = append(errs, errStr)
	}
	return errs
}

func processFormField(r *http.Request, field string) (string, string) {
	fieldData := r.PostFormValue(field)
	if len(fieldData) == 0 {
		return "", "Missing '" + field + "' parameter, cannot continue"
	}
	return fieldData, ""
}
