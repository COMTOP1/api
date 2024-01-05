package handler

import (
	"bytes"
	"github.com/COMTOP1/api/handler/api"
	"net/url"
)

// Session represents an open API session.
type Session struct {
	requester api.Requester
}

// NewSession constructs a new Session with a given url.
func NewSession(baseURL string) (*Session, error) {
	url1, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &Session{requester: api.NewRequester(*url1)}, nil
}

// NewSessionForServer constructs a new Session with the given API key for a non-standard server URL.
//func NewSessionForServer(apikey, server string) (*Session, error) {
//	url1, err := url.Parse(server)
//	if err != nil {
//		return nil, err
//	}
//	return &Session{requester: api.NewRequester(apikey, *url1)}, nil
//}

//// MockSession creates a new mocked API session returning the JSON message stored in message.
//func MockSession(message []byte) (*Session, error) {
//	rm := json.RawMessage{}
//	err := rm.UnmarshalJSON(message)
//	if err != nil {
//		return nil, err
//	}
//	return &Session{requester: api.MockRequester(&rm)}, nil
//}

// do fulfils, a request for the given endpoint.
func (s *Session) do(r *api.Request) *api.Response {
	return s.requester.Do(r)
}

// doToken fulfils, a request for the given endpoint and token.
func (s *Session) doToken(r *api.Request, token string) *api.Response {
	return s.requester.DoToken(r, token)
}

// get creates, and fulfils, a GET request for the given endpoint.
func (s *Session) get(endpoint string) *api.Response {
	return s.do(api.NewRequest(endpoint))
}

// getToken creates, and fulfils, a GET request for the given endpoint.
func (s *Session) getToken(token, endpoint string) *api.Response {
	return s.doToken(api.NewRequest(endpoint), token)
}

// getfToken creates, and fulfils, a GET request for the endpoint created by
// the given format string and parameters.
func (s *Session) getfToken(token, format string, params ...interface{}) *api.Response {
	return s.doToken(api.NewRequestf(format, params...), token)
}

// getf creates, and fulfils, a GET request for the endpoint created by
// the given format string and parameters.
func (s *Session) getf(format string, params ...interface{}) *api.Response {
	return s.do(api.NewRequestf(format, params...))
}

//nolint:unused
func (s *Session) getWithQueryParams(format string, queryParams map[string][]string) *api.Response {
	r := api.NewRequest(format)
	r.Params = queryParams
	return s.do(r)
}

//// putf creates, and fulfils, a PUT request for the endpoint created by
//// the given format string and parameters.
//func (s *Session) putf(format string, body bytes.Buffer, params ...interface{}) *api.Response {
//	r := api.NewRequestf(format, params...)
//	r.ReqType = api.PutReq
//	r.Body = body
//	return s.do(r)
//}

func (s *Session) putToken(token, endpoint string, body bytes.Buffer) *api.Response {
	r := api.NewRequest(endpoint)
	r.ReqType = api.PutReq
	r.Body = body
	return s.doToken(r, token)
}

// post creates, and fulfils, a POST request for the given endpoint,
// using the given form parameters
//
//nolint:unused
func (s *Session) post(endpoint string, formParams map[string][]string) *api.Response {
	r := api.NewRequest(endpoint)
	r.ReqType = api.PostReq
	r.Params = formParams
	return s.do(r)
}

// postToken creates, and fulfils, a POST request for the given endpoint,
// using the given form parameters
//
//nolint:unused
func (s *Session) postToken(token, endpoint string, body bytes.Buffer) *api.Response {
	r := api.NewRequest(endpoint)
	r.ReqType = api.PostReq
	r.Body = body
	return s.doToken(r, token)
}

//// patch creates, and fulfils, a PATCH request for the given endpoint,
//// using the given form parameters
//func (s *Session) patch(endpoint string, formParams map[string][]string) *api.Response {
//	r := api.NewRequest(endpoint)
//	r.ReqType = api.PatchReq
//	r.Params = formParams
//	return s.do(r)
//}

// patchToken creates, and fulfils, a PATCH request for the given endpoint,
// using the given form parameters
func (s *Session) patchToken(token, endpoint string, body bytes.Buffer) *api.Response {
	r := api.NewRequest(endpoint)
	r.ReqType = api.PatchReq
	r.Body = body
	return s.doToken(r, token)
}

//// delete creates, and fulfils, a DELETE request for the given endpoint.
//func (s *Session) delete(endpoint string) *api.Response {
//	r := api.NewRequest(endpoint)
//	r.ReqType = api.DeleteReq
//	return s.do(r)
//}

//// deleteToken creates, and fulfils, a PATCH request for the given endpoint,
//// using the given form parameters
//func (s *Session) deleteToken(token, endpoint string) *api.Response {
//	r := api.NewRequest(endpoint)
//	r.ReqType = api.DeleteReq
//	return s.doToken(r, token)
//}

// deletefToken creates, and fulfils, a PATCH request for the given endpoint,
// using the given form parameters
func (s *Session) deletefToken(token, endpoint string, params ...interface{}) *api.Response {
	r := api.NewRequestf(endpoint, params...)
	r.ReqType = api.DeleteReq
	return s.doToken(r, token)
}

//// deletef creates, and fulfils, a DELETE request for the endpoint created by
//// the given format string and parameters.
//func (s *Session) deletef(format string, params ...interface{}) *api.Response {
//	r := api.NewRequestf(format, params...)
//	r.ReqType = api.DeleteReq
//	return s.do(r)
//}

// NewSessionFromKeyFile tries to open a Session with the key from an API key file.
//func NewSessionFromKeyFile() (*Session, error) {
//	apikey, err := api.GetAPIKey()
//	if err != nil {
//		return nil, err
//	}
//
//	return NewSession(apikey)
//}

// NewSessionFromKeyFileForServer tries to open a Session with the key from an API key file, with a non-standard server.
//func NewSessionFromKeyFileForServer(server string) (*Session, error) {
//	apikey, err := api.GetAPIKey()
//	if err != nil {
//		return nil, err
//	}
//
//	return NewSessionForServer(apikey, server)
//}
