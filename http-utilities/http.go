package httputility

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	APPLICATIONJSON string = "application/json"
	TEXTXML                = "text/xml"
)

// HttpError error object
type HttpError struct {
	URL      string
	HTTPCode int
	Message  string
	Body     interface{}
	Err      error `json:"-"`
}

func (he *HttpError) Error() string {
	return fmt.Sprintf(
		"Resource error: URL: %s, status code: %v,  err: %v, body: %v",
		he.URL,
		he.HTTPCode,
		he.Err,
		he.Body,
	)
}

// RequestJSON JSON request in all types of methods
func RequestJSON(method, url string, data []byte, headers map[string]string, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	if headers == nil {
		headers = map[string]string{"Content-Type": APPLICATIONJSON}
	} else {
		headers["Content-Type"] = APPLICATIONJSON
	}

	httpStatus, responseBody, err = send(method, url, "", data, headers)
	if err != nil {
		return
	}
	if responseStruct != nil && len(responseBody) != 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}
	return
}

// AuthorizedRequestJSON Authorized JSON request in all types of methods
func AuthorizedRequestJSON(method, url, token string, data []byte, headers map[string]string, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	if headers == nil {
		headers = map[string]string{"Content-Type": APPLICATIONJSON}
	} else {
		headers["Content-Type"] = APPLICATIONJSON
	}

	httpStatus, responseBody, err = send(method, url, token, data, headers)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) != 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}
	return
}

// RequestXML XML request in all types of methods
func RequestXML(method, url string, data []byte, headers map[string]string, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	if headers == nil {
		headers = map[string]string{"Content-Type": TEXTXML}
	} else {
		headers["Content-Type"] = TEXTXML
	}

	httpStatus, responseBody, err = send(method, url, "", data, headers)
	if err != nil {
		return
	}
	if responseStruct != nil && len(responseBody) != 0 {
		err = xml.Unmarshal(responseBody, responseStruct)
	}
	return
}

// AuthorizedRequestXML Authorized XML request in all types of methods
func AuthorizedRequestXML(method, url, token string, data []byte, headers map[string]string, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	if headers == nil {
		headers = map[string]string{"Content-Type": TEXTXML}
	} else {
		headers["Content-Type"] = TEXTXML
	}

	httpStatus, responseBody, err = send(method, url, token, data, headers)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) != 0 {
		err = xml.Unmarshal(responseBody, responseStruct)
	}
	return
}

// RequestMultipart Multipart request in all types of methods
func RequestMultipart(method, url string, data []byte, headers map[string]string, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	httpStatus, responseBody, err = send(method, url, "", data, headers)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) != 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}
	return
}

// AuthorizedRequestMultipart Authorized Multipart request in all types of methods
func AuthorizedRequestMultipart(method, url, token string, data []byte, headers map[string]string, responseStruct interface{}) (httpStatus int, responseBody []byte, err error) {
	httpStatus, responseBody, err = send(method, url, token, data, headers)
	if err != nil {
		return
	}

	if responseStruct != nil && len(responseBody) != 0 {
		err = json.Unmarshal(responseBody, responseStruct)
	}
	return
}

func send(method, urlString, token string, data []byte, headers map[string]string) (httpStatus int, buf []byte, err error) {
	request, err := http.NewRequest(method, urlString, bytes.NewBuffer(data))
	if err != nil {
		return httpStatus, nil, &HttpError{URL: urlString, Err: err}
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	if token != "" {
		request.Header.Add("Authorization", token)
	}

	if strings.ContainsAny(urlString, "?") {
		urlTemp, err := url.Parse(urlString)
		if err != nil {
			return httpStatus, nil, &HttpError{URL: urlString, Err: err}
		}
		urlQuery := urlTemp.Query()
		urlTemp.RawQuery = urlQuery.Encode()
		urlString = urlTemp.String()
	}

	cfg := getInsecureTLSConfig()
	transport := &http.Transport{
		TLSClientConfig: cfg,
	}

	client := &http.Client{Transport: transport, Timeout: 60 * time.Second}

	response, err := client.Do(request)
	if err != nil {
		return httpStatus, nil, &HttpError{URL: urlString, Err: err}
	}
	defer response.Body.Close()

	buf, err = io.ReadAll(response.Body)
	if err != nil {
		return httpStatus, nil, &HttpError{URL: urlString, Err: err, HTTPCode: response.StatusCode}
	}

	httpStatus = response.StatusCode
	if response.StatusCode > 399 {
		return httpStatus, buf, &HttpError{
			URL:      urlString,
			Err:      fmt.Errorf("incorrect status code"),
			HTTPCode: response.StatusCode,
			Message:  "incorrect status code",
		}
	}

	return
}

func getInsecureTLSConfig() *tls.Config {
	skipVerifyBool := false
	if true {
		skipVerifyBool = true
	}
	return &tls.Config{
		InsecureSkipVerify: skipVerifyBool,
	}
}
