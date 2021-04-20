package bot

import (
	"bytes"
	"deyforyou/dey/crawler/colorizer"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
)

const (
	contentJSON = "application/json"
	contentFORM = "application/x-www-form-urlencoded"
)

var (
	jar, _       = cookiejar.New(nil)
	defautRemote = &Remote{Client: &http.Client{Jar: jar}}
)

type (
	ErrorCallBack    func(err error)
	RequestCallBack  func(request *http.Request)
	ResponseCallBack func(response *http.Response)

	// Client
	Remote struct {
		mux                  *sync.Mutex
		Client               *http.Client
		errorCallBacks       []ErrorCallBack
		requestCallBacks     []RequestCallBack
		responseCallBacks    []ResponseCallBack
		requestMapCallBacks  map[string]RequestCallBack
		responseMapCallBacks map[string]ResponseCallBack
	}
)

func New() *Remote {
	jar, _ := cookiejar.New(nil)
	return &Remote{
		requestMapCallBacks:  make(map[string]RequestCallBack),
		responseMapCallBacks: make(map[string]ResponseCallBack),
		Client: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				MaxIdleConns:    2,
				MaxConnsPerHost: 2,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				colorizer.Warn("Redirect", req.URL)
				return nil
			},
		},
		mux: new(sync.Mutex),
	}
}

func jsonReader(body interface{}) io.Reader {
	writer := new(bytes.Buffer)
	json.NewEncoder(writer).Encode(body)
	return writer
}

func formReader(body url.Values) *strings.Reader {
	return strings.NewReader(body.Encode())
}

func Get(url string) (*http.Response, error) {
	return defautRemote.Get(url)
}

func Request(request RequestCallBack) {
	defautRemote.requestCallBacks = append(defautRemote.requestCallBacks, request)
}

func RequestByURL(pattern string, request RequestCallBack) {
	defautRemote.requestMapCallBacks[pattern] = request
}

func Response(response ResponseCallBack) {
	defautRemote.responseCallBacks = append(defautRemote.responseCallBacks, response)
}

func ResponseByURL(pattern string, response ResponseCallBack) {
	defautRemote.responseMapCallBacks[pattern] = response
}

func Post(url string, body interface{}) (*http.Response, error) {
	return defautRemote.Post(url, body)
}

func PostForm(url string, body url.Values) (*http.Response, error) {
	return defautRemote.PostForm(url, body)
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
	return remote.do(http.MethodPost, url, contentJSON, jsonReader(body))
}

func (remote *Remote) PostForm(url string, body url.Values) (*http.Response, error) {
	return remote.do(http.MethodPost, url, contentFORM, formReader(body))
}

func (remote *Remote) Put(url string, body interface{}) (*http.Response, error) {
	return remote.do(http.MethodPut, url, contentJSON, jsonReader(body))
}

func (remote *Remote) Patch(url string, body interface{}) (*http.Response, error) {
	return remote.do(http.MethodPatch, url, contentJSON, jsonReader(body))
}

func (remote *Remote) Delete(url string, body interface{}) (*http.Response, error) {
	return remote.do(http.MethodDelete, contentJSON, url, jsonReader(body))
}

func (remote *Remote) do(method, contentType, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", contentType)
	remote.mux.Lock()
	for _, callBack := range remote.requestCallBacks {
		callBack(request)
	}
	if callBack, ok := remote.requestMapCallBacks[url]; ok {
		callBack(request)
	}
	remote.mux.Unlock()

	colorizer.Info(method, url)
	response, err := remote.Client.Do(request)
	if err != nil {
		for _, callBack := range remote.errorCallBacks {
			callBack(err)
		}
	} else {
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		buffer := bytes.NewBuffer(bodyBytes)
		remote.mux.Lock()
		for _, callBack := range remote.responseCallBacks {
			response.Body.Close()
			response.Body = io.NopCloser(buffer)
			callBack(response)
		}
		response.Body.Close()
		response.Body = io.NopCloser(buffer)
		if callBack, ok := remote.responseMapCallBacks[url]; ok {
			callBack(response)
		}
		remote.mux.Unlock()

	}
	return response, nil
}

// bodyBytes, _ := ioutil.ReadAll(response.Body)
// buffer := bytes.NewBuffer(bodyBytes)

// response.Body.Close()
// response.Body = io.NopCloser(buffer)
