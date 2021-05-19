package remote

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	contentType = "Content-Type"
	contentJSON = "application/json"
	contentFORM = "application/x-www-form-urlencoded"
)

type (
	ErrorCallBack    func(err error)
	RequestCallBack  func(request *http.Request)
	ResponseCallBack func(response *http.Response)

	// Client
	Remote struct {
		mux                  *sync.RWMutex
		Client               *http.Client
		errorCallBacks       []ErrorCallBack
		requestCallBacks     []RequestCallBack
		responseCallBacks    []ResponseCallBack
		errorMapCallBacks    map[string]ErrorCallBack
		requestMapCallBacks  map[string]RequestCallBack
		responseMapCallBacks map[string]ResponseCallBack
	}
)

func jsonReader(body interface{}) io.Reader {
	writer := new(bytes.Buffer)
	json.NewEncoder(writer).Encode(body)
	return writer
}

func formReader(body url.Values) *strings.Reader {
	return strings.NewReader(body.Encode())
}

func New(client *http.Client) *Remote {
	return &Remote{
		Client:               client,
		mux:                  new(sync.RWMutex),
		errorCallBacks:       make([]ErrorCallBack, 0),
		requestCallBacks:     make([]RequestCallBack, 0),
		responseCallBacks:    make([]ResponseCallBack, 0),
		errorMapCallBacks:    make(map[string]ErrorCallBack),
		requestMapCallBacks:  make(map[string]RequestCallBack),
		responseMapCallBacks: make(map[string]ResponseCallBack),
	}
}

func (remote *Remote) Clone() *Remote {
	return &Remote{
		mux:                  remote.mux,
		Client:               remote.Client,
		errorCallBacks:       remote.errorCallBacks,
		requestCallBacks:     remote.requestCallBacks,
		responseCallBacks:    remote.responseCallBacks,
		errorMapCallBacks:    remote.errorMapCallBacks,
		requestMapCallBacks:  remote.requestMapCallBacks,
		responseMapCallBacks: remote.responseMapCallBacks,
	}
}

func (remote *Remote) Request(request RequestCallBack) {
	remote.mux.Lock()
	remote.requestCallBacks = append(remote.requestCallBacks, request)
	remote.mux.Unlock()
}

func (remote *Remote) Response(response ResponseCallBack) {
	remote.mux.Lock()
	remote.responseCallBacks = append(remote.responseCallBacks, response)
	remote.mux.Unlock()
}

func (remote *Remote) RequestByURL(pattern string, request RequestCallBack) {
	remote.mux.Lock()
	remote.requestMapCallBacks[pattern] = request
	remote.mux.Unlock()
}

func (remote *Remote) ResponseByURL(pattern string, response ResponseCallBack) {
	remote.mux.Lock()
	remote.responseMapCallBacks[pattern] = response
	remote.mux.Unlock()
}

func (remote *Remote) Get(url string) (*http.Response, error) {
	return remote.do(http.MethodGet, contentFORM, url, nil)
}

func (remote *Remote) Head(url string) (*http.Response, error) {
	return remote.do(http.MethodHead, contentFORM, url, nil)
}

func (remote *Remote) Post(url string, body interface{}) (*http.Response, error) {
	return remote.do(http.MethodPost, contentJSON, url, jsonReader(body))
}

func (remote *Remote) PostForm(url string, body url.Values) (*http.Response, error) {
	return remote.do(http.MethodPost, contentFORM, url, formReader(body))
}

func (remote *Remote) Put(url string, body interface{}) (*http.Response, error) {
	return remote.do(http.MethodPut, contentJSON, url, jsonReader(body))
}

func (remote *Remote) Patch(url string, body interface{}) (*http.Response, error) {
	return remote.do(http.MethodPatch, contentJSON, url, jsonReader(body))
}

func (remote *Remote) Delete(url string, body interface{}) (*http.Response, error) {
	return remote.do(http.MethodDelete, contentJSON, url, jsonReader(body))
}

func (remote *Remote) do(method, contentTypeValue, url string, body io.Reader) (*http.Response, error) {

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	request.Header.Set(contentType, contentTypeValue)
	remote.mux.RLock()
	for _, callBack := range remote.requestCallBacks {
		callBack(request)
	}
	if callBack, ok := remote.requestMapCallBacks[url]; ok {
		callBack(request)
	}
	remote.mux.RUnlock()

	response, err := remote.Client.Do(request)
	if err != nil {
		remote.mux.RLock()
		for _, callBack := range remote.errorCallBacks {
			callBack(err)
		}
		if callBack, ok := remote.errorMapCallBacks[url]; ok {
			callBack(err)
		}
		remote.mux.RUnlock()
		return nil, err
	} else {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		remote.mux.RLock()
		for _, callBack := range remote.responseCallBacks {
			buffer := bytes.NewBuffer(bodyBytes)
			response.Body.Close()
			response.Body = ioutil.NopCloser(buffer)
			callBack(response)
		}
		if callBack, ok := remote.responseMapCallBacks[url]; ok {
			buffer := bytes.NewBuffer(bodyBytes)
			response.Body.Close()
			response.Body = ioutil.NopCloser(buffer)
			callBack(response)
		}
		remote.mux.RUnlock()

		buffer := bytes.NewBuffer(bodyBytes)
		response.Body.Close()
		response.Body = ioutil.NopCloser(buffer)
	}
	return response, nil
}
