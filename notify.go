package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Lookup method
func (twilio *Twilio) Notify(body string, to []string) (notifyResponse map[string]interface{}, exception *Exception, err error) {
	twilioUrl := twilio.NotifyUrl + "/Services/" + twilio.NotifySid + "/Notifications"

	formValues := url.Values{}

	for _, phone := range to {
		formValues.Add("ToBinding", `{"binding_type": "sms", "address": "`+phone+`"}`)
	}

	formValues.Set("Body", body)

	res, err := twilio.post(formValues, twilioUrl)
	if err != nil {
		return notifyResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return notifyResponse, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return notifyResponse, exception, err
	}

	notifyResponse = make(map[string]interface{})
	err = json.Unmarshal(responseBody, &notifyResponse)
	return notifyResponse, exception, err
}
