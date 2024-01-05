// Package api exposes a low-level interface to the MyRadio API.
//
// It exposes the Requester interface for types that represent
// connections to the API, methods for constructing Requesters, and
// functions and methods for using Requesters to make requests.
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	error2 "github.com/COMTOP1/api/handler/error"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Request represents an API request being built.
type Request struct {
	// The endpoint, as a suffix of the API root URL.
	Endpoint string
	// The set of mixins to use.
	Mixins []string
	// The map of parameters to use.
	Params map[string][]string
	// The type of request (i.e. GET/POST etc.)
	ReqType HTTPMethod
	// The body of the request
	Body bytes.Buffer
}

// HTTPMethod guards against incorrect methods being specified through strings
type HTTPMethod int

const (
	// GetReq corresponds to GET
	GetReq HTTPMethod = iota
	// PostReq corresponds to POST
	PostReq
	// PutReq corresponds to PUT
	PutReq
	// PatchReq corresponds to PATCH
	PatchReq
	// DeleteReq corresponds to DELETE
	DeleteReq
)

// String converts a HTTPMethod object into a usable request method string
func (m HTTPMethod) String() (string, error) {
	switch m {
	case GetReq:
		return "GET", nil
	case PostReq:
		return "POST", nil
	case PutReq:
		return "PUT", nil
	case PatchReq:
		return "PATCH", nil
	case DeleteReq:
		return "DELETE", nil
	default:
		return "", errors.New("invalid HTTP method specified")
	}
}

// NewRequest constructs a new request for the given endpoint.
func NewRequest(endpoint string) *Request {
	return &Request{
		Endpoint: endpoint,
		Mixins:   []string{},
		Params:   map[string][]string{},
		ReqType:  GetReq,
		Body:     bytes.Buffer{},
	}
}

// NewRequestf constructs a new request for the endpoint constructed by
// the given format string and parameters.
func NewRequestf(format string, params ...interface{}) *Request {
	return NewRequest(fmt.Sprintf(format, params...))
}

// Response represents the result of an API request.
type Response struct {
	raw *json.RawMessage
	err error
}

// IsEmpty checks whether the response payload is present, but empty.
func (r *Response) IsEmpty() bool {
	if r.err != nil {
		return false
	}

	if r.raw == nil {
		return true
	}

	// Check for 'empty' JSON payloads.
	bs, err := r.raw.MarshalJSON()
	if err != nil {
		return false
	}

	fmt.Println(string(bs))

	if len(bs) != 2 {
		return false
	}

	if bs[0] == '[' && bs[1] == ']' {
		return true
	}

	if bs[0] == '{' && bs[1] == '}' {
		return true
	}

	return false
}

// JSON returns r as raw JSON.
func (r *Response) JSON() (*json.RawMessage, error) {
	return r.raw, r.err
}

// Into unmarshals the response r into in.
func (r *Response) Into(in interface{}) error {
	if r.err != nil {
		return r.err
	}

	if r.raw == nil {
		return nil
	}
	return json.Unmarshal(*r.raw, in)
}

// Requester is the type of anything that can handle an API request.
type Requester interface {
	// Do fulfils an API request.
	Do(r *Request) *Response
	DoToken(r *Request, token string) *Response
}

// authedRequester answers API requests by making an authed API call.
type authedRequester struct {
	//apikey  string
	baseurl url.URL
}

// NewRequester creates a new 'live' requester.
func NewRequester( /*apikey string, */ url url.URL) Requester {
	return &authedRequester{
		//        apikey:  apikey,
		baseurl: url,
	}
}

// DoToken fulfils an API request with token.
func (s *authedRequester) DoToken(r *Request, token string) *Response {
	//Validate the request method before we waste any time
	reqMethod, err := r.ReqType.String()
	if err != nil {
		return &Response{err: err}
	}

	urlParams := url.Values{ /*"api_key": []string{s.apikey}*/ }
	if len(r.Mixins) > 0 {
		urlParams.Add("mixins", strings.Join(r.Mixins, ","))
	}
	for k, vs := range r.Params {
		for _, v := range vs {
			urlParams.Add(k, v)
		}
	}

	theurl := s.baseurl
	theurl.Path += r.Endpoint
	encodedParams := urlParams.Encode()

	// POST, PUT and PATCH sends form params in the body
	if r.ReqType == PostReq || r.ReqType == PutReq || r.ReqType == PatchReq {
		r.Body.WriteString(encodedParams)
	} else {
		theurl.RawQuery = encodedParams
	}

	b := bytes.NewReader(r.Body.Bytes())
	req, err := http.NewRequest(reqMethod, theurl.String(), b)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

	// Specify content type for POST, PUT or PATCH requests, as the body format has to be specified
	if r.ReqType == PostReq || r.ReqType == PutReq || r.ReqType == PatchReq || r.ReqType == DeleteReq {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	}

	if err != nil {
		return &Response{err: err}
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &Response{err: err}
	}
	defer func() {
		err = res.Body.Close()
	}()
	if err != nil {
		return &Response{err: err}
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &Response{err: err}
	}
	if res.StatusCode != 200 && res.StatusCode != 204 {
		return &Response{err: fmt.Errorf("%s Not ok: HTTP %d\n%v", r.Endpoint, res.StatusCode, error2.New(string(data)))}
	}
	var response struct {
		Status  int
		Payload *json.RawMessage
	}
	response.Status = res.StatusCode
	raw := json.RawMessage{}
	response.Payload = &raw
	err = response.Payload.UnmarshalJSON(data)
	if err != nil {
		return &Response{err: err}
	}
	if response.Status != 200 {
		return &Response{err: fmt.Errorf(r.Endpoint + fmt.Sprintf(" Response not OK: %v", response))}
	}
	return &Response{raw: response.Payload, err: nil}
}

// Do fulfils an API request.
func (s *authedRequester) Do(r *Request) *Response {
	//Validate the request method before we waste any time
	reqMethod, err := r.ReqType.String()
	if err != nil {
		return &Response{err: err}
	}

	urlParams := url.Values{ /*"api_key": []string{s.apikey}*/ }
	if len(r.Mixins) > 0 {
		urlParams.Add("mixins", strings.Join(r.Mixins, ","))
	}
	for k, vs := range r.Params {
		for _, v := range vs {
			urlParams.Add(k, v)
		}
	}

	theurl := s.baseurl
	theurl.Path += r.Endpoint
	encodedParams := urlParams.Encode()

	// POST, PUT and PATCH sends form params in the body
	if r.ReqType == PostReq || r.ReqType == PutReq || r.ReqType == PatchReq {
		r.Body.WriteString(encodedParams)
	} else {
		theurl.RawQuery = encodedParams
	}
	req, err := http.NewRequest(reqMethod, theurl.String(), bytes.NewReader(r.Body.Bytes()))

	// Specify content type for POST, PUT or PATCH requests, as the body format has to be specified
	if r.ReqType == PostReq || r.ReqType == PutReq || r.ReqType == PatchReq {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	}

	if err != nil {
		return &Response{err: err}
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &Response{err: err}
	}
	defer func() {
		err = res.Body.Close()
	}()
	if err != nil {
		return &Response{err: err}
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &Response{err: err}
	}
	if res.StatusCode != 200 {
		return &Response{err: fmt.Errorf("%s Not ok: HTTP %d\n%v", r.Endpoint, res.StatusCode, error2.New(string(data)))}
	}
	var response struct {
		Status  int
		Payload *json.RawMessage
	}
	response.Status = res.StatusCode
	raw := json.RawMessage{}
	response.Payload = &raw
	err = response.Payload.UnmarshalJSON(data)
	if err != nil {
		return &Response{err: err}
	}
	if response.Status != 200 {
		return &Response{err: fmt.Errorf(r.Endpoint + fmt.Sprintf(" Response not OK: %v", response))}
	}
	return &Response{raw: response.Payload, err: nil}
}

// mockRequester answers API requests by returning some stock response.
//type mockRequester struct {
//	message *json.RawMessage
//}

// MockRequester creates a new mocked requester.
//func MockRequester(message *json.RawMessage) Requester {
//	return &mockRequester{message: message}
//}

// Do pretends to fulfil an API request, but actually returns the mockRequester's stock response.
//func (s *mockRequester) Do(r *Request) *Response {
//	return &Response{raw: s.message, err: nil}
//}

// DoToken pretends to fulfil an API request, but actually returns the mockRequester's stock response with token.
//func (s *mockRequester) DoToken(r *Request, token string) *Response {
//	return &Response{raw: s.message, err: nil}
//}
