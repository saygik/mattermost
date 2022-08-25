package mattermost

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const (
	HeaderRequestId                 = "X-Request-ID"
	HeaderVersionId                 = "X-Version-ID"
	HeaderClusterId                 = "X-Cluster-ID"
	HeaderEtagServer                = "ETag"
	HeaderEtagClient                = "If-None-Match"
	HeaderForwarded                 = "X-Forwarded-For"
	HeaderRealIP                    = "X-Real-IP"
	HeaderForwardedProto            = "X-Forwarded-Proto"
	HeaderToken                     = "token"
	HeaderCsrfToken                 = "X-CSRF-Token"
	HeaderBearer                    = "BEARER"
	HeaderAuth                      = "Authorization"
	HeaderCloudToken                = "X-Cloud-Token"
	HeaderRemoteclusterToken        = "X-RemoteCluster-Token"
	HeaderRemoteclusterId           = "X-RemoteCluster-Id"
	HeaderRequestedWith             = "X-Requested-With"
	HeaderRequestedWithXML          = "XMLHttpRequest"
	HeaderFirstInaccessiblePostTime = "First-Inaccessible-Post-Time"
	HeaderRange                     = "Range"
	STATUS                          = "status"
	StatusOk                        = "OK"
	StatusFail                      = "FAIL"
	StatusUnhealthy                 = "UNHEALTHY"
	StatusRemove                    = "REMOVE"

	ClientDir = "client"

	APIURLSuffixV1 = "/api/v1"
	APIURLSuffixV4 = "/api/v4"
	APIURLSuffixV5 = "/api/v5"
	APIURLSuffix   = APIURLSuffixV4
)

type Response struct {
	StatusCode    int
	RequestId     string
	Etag          string
	ServerVersion string
	Header        http.Header
}

type Client4 struct {
	URL        string       // The location of the server, for example  "http://localhost:8065"
	APIURL     string       // The api location of the server, for example "http://localhost:8065/api/v4"
	HTTPClient *http.Client // The http client
	AuthToken  string
	AuthType   string
	HTTPHeader map[string]string // Headers to be copied over for each request

	// TrueString is the string value sent to the server for true boolean query parameters.
	trueString string

	// FalseString is the string value sent to the server for false boolean query parameters.
	falseString string
}

func (c *Client4) SetToken(token string) {
	c.AuthToken = token
	c.AuthType = HeaderBearer
}
func BuildResponse(r *http.Response) *Response {
	if r == nil {
		return nil
	}

	return &Response{
		StatusCode:    r.StatusCode,
		RequestId:     r.Header.Get(HeaderRequestId),
		Etag:          r.Header.Get(HeaderEtagServer),
		ServerVersion: r.Header.Get(HeaderVersionId),
		Header:        r.Header,
	}
}
func NewAPIv4Client(url string) *Client4 {
	url = strings.TrimRight(url, "/")
	return &Client4{url, url + APIURLSuffix, &http.Client{}, "", "", map[string]string{}, "", ""}
}

func (c *Client4) CreatePost(post *Post) (*Post, *Response, error) {
	postJSON, err := json.Marshal(post)
	if err != nil {
		return nil, nil, NewAppError("CreatePost", "api.marshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	r, err := c.DoAPIPost(c.postsRoute(), string(postJSON))
	if err != nil {
		return nil, BuildResponse(r), err
	}
	defer closeBody(r)
	var p Post
	if r.StatusCode == http.StatusNotModified {
		return &p, BuildResponse(r), nil
	}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return nil, nil, NewAppError("CreatePost", "api.unmarshal_error", nil, "", http.StatusInternalServerError).Wrap(err)
	}
	return &p, BuildResponse(r), nil
}
func (c *Client4) postsRoute() string {
	return "/posts"
}

func (c *Client4) DoAPIGet(url string, etag string) (*http.Response, error) {
	return c.DoAPIRequest(http.MethodGet, c.APIURL+url, "", etag)
}

func (c *Client4) DoAPIPost(url string, data string) (*http.Response, error) {
	return c.DoAPIRequest(http.MethodPost, c.APIURL+url, data, "")
}

func (c *Client4) DoAPIRequest(method, url, data, etag string) (*http.Response, error) {
	return c.DoAPIRequestReader(method, url, strings.NewReader(data), map[string]string{HeaderEtagClient: etag})
}

func (c *Client4) DoAPIRequestWithHeaders(method, url, data string, headers map[string]string) (*http.Response, error) {
	return c.DoAPIRequestReader(method, url, strings.NewReader(data), headers)
}

func (c *Client4) DoAPIRequestBytes(method, url string, data []byte, etag string) (*http.Response, error) {
	return c.DoAPIRequestReader(method, url, bytes.NewReader(data), map[string]string{HeaderEtagClient: etag})
}

func (c *Client4) DoAPIRequestReader(method, url string, data io.Reader, headers map[string]string) (*http.Response, error) {
	rq, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		rq.Header.Set(k, v)
	}

	if c.AuthToken != "" {
		rq.Header.Set(HeaderAuth, c.AuthType+" "+c.AuthToken)
	}

	if c.HTTPHeader != nil && len(c.HTTPHeader) > 0 {
		for k, v := range c.HTTPHeader {
			rq.Header.Set(k, v)
		}
	}

	rp, err := c.HTTPClient.Do(rq)
	if err != nil {
		return rp, err
	}

	if rp.StatusCode == 304 {
		return rp, nil
	}

	if rp.StatusCode >= 300 {
		defer closeBody(rp)
		return rp, AppErrorFromJSON(rp.Body)
	}

	return rp, nil
}
func closeBody(r *http.Response) {
	if r.Body != nil {
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
	}
}
