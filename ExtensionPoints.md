# HTTP Getaway, Part 1: HTTP extension points

Configuration options and extension points for HTTP clients.

> [Golang Leipzig](https://golangleipzig.space/), 2020-02-21, 19:00, Martin Czygan <martin.czygan@gmail.com>

# Motivation

* "the network is reliable", from [Fallacies of distributed computing](http://nighthacks.com/jag/res/Fallacies.html)
* writing data acquisition tools

# Outline

* interfaces
* configuration
* redirects
* retries
* tracing

# Interfaces

* [net] package contains 105040 LOC
* [net/http] 59154 LOC (of which 30297 are tests)

The [net/http](https://golang.org/pkg/net/http/) contains 12 interfaces (02/2020):

<!--  $ find . -type f | xargs cat | grep '^type[ ]*[A-Z].* interface {' | awk '{print $2}' | sort -->

* BufferPool
* ~~CloseNotifier~~, deprecated with [Go 1.11](https://golang.org/doc/go1.11#net/http)
* **CookieJar**
* **File**
* **FileSystem**
* Flusher
* **Handler**, `ServeHTTP(ResponseWriter, *Request)`
* **Hijacker**, expose TCP connection to response writer
* PublicSuffixList, for cookies, jars ([RFC 6265 Section 5.3, Note 5](https://tools.ietf.org/html/rfc6265#section-5.3))
* Pusher
* **ResponseWriter**
* **RoundTripper**, core interface

## Hijacker

```go
type Hijacker interface {
        Hijack() (net.Conn, *bufio.ReadWriter, error)
}
```

* can be implemented by ReponseWriters to hand over the TCP connection (and then leave it alone)
* used, e.g. by websocket libraries

![](static/hijack.gif)

## File and FileSystem

* `http.File` is a `io.ReadSeekCloser` plus `ReadDir` and `Stat`
* `http.FileSystem` is a single method interface `Open(name string) (http.File, error)`

Abstracts file system like access.

Use `http.FileServer(root FileSystem) Handler` to turn a filesystem into an http Handler.

The `http.Dir` is an `http.FileSystem` allowing access to local filesytem.

Example: [x/webshare.go](x/webshare.go).

![](static/webshare.png)

## CookieJar

```go
type CookieJar interface {
        SetCookies(u *url.URL, cookies []*Cookie)
        Cookies(u *url.URL) []*Cookie
}
```

* getting and setting cookies
* an in-memory implementation in [net/http/cookiejar](https://golang.org/pkg/net/http/cookiejar)

## ResponseWriter

>  A ResponseWriter interface is used by an HTTP handler to construct an HTTP
>  response.

```go
type ResponseWriter interface {
    Header() Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int) // First call to write will call this.
}
```

The default unexported implementation is `http.response`.

The standard library server has various implementations, e.g.

* [httptest.ResponseRecorder](https://golang.org/pkg/net/http/httptest/#ResponseRecorder), example: [x/resprec.go](x/resprec.go)
* [http.populateResponse](https://github.com/golang/go/blob/ccb95b6492ad6e7a7d1a7fda896baee4caffb3b4/src/net/http/filetransport.go#L65-L76),
  using a [io.Pipe](https://golang.org/pkg/io/#Pipe) to connect file content
and response body

## Question regarding pointer receivers?

* [In Go HTTP handlers, why is the ResponseWriter a value but the Request a
  pointer?](https://stackoverflow.com/questions/13255907/in-go-http-handlers-why-is-the-responsewriter-a-value-but-the-request-a-pointer) (SO: 76, 7y3m ago, 8k)

An interface (w) and a struct (r).

## RoundTripper

> RoundTripper is an interface representing the ability to execute a single
> HTTP transaction, obtaining the Response for a given Request.

* should not *interpret* the response (e.g. err == nil, even with HTTP errors)

```go
type RoundTripper interface {
        RoundTrip(*Request) (*Response, error)
}
```

File serving uses an internal `fileTransport` struct, that is a RoundTripper.

Example caching RoundTripper: [x/cachingrt.go](x/cachingrt.go).


# Configuration and Timeouts

Various levels:

![](static/levels.png)

Configuration can happen in the Client or on Transport level.

## Client

```go
client := &http.Client{
    Transport       RoundTripper
    CheckRedirect   func(req *Request, via []*Request) error
    Jar             CookieJar
    Timeout         time.Duration
}
```

The package defines a default client:

```go
// DefaultClient is the default Client and is used by Get, Head, and Post.
var DefaultClient = &Client{}
```

Note about timeout:

> Timeout specifies a time limit for requests made by this
> Client. The timeout includes connection time, any
> redirects, and reading the response body. The timer remains
> running after Get, Head, Post, or Do return and will
> interrupt reading of the Response.Body.
>
> **A Timeout of zero means no timeout.**

### Example Redirect Tracking

Goal: Do a request and **record all intermediate requests** between initial
request and first non-redirect request.


> **ErrUseLastResponse** can be returned by Client.CheckRedirect hooks to control
> how redirects are processed. If returned, the next request is not sent and
> the most recent response is returned with its body unclosed.

The if a *policyFunc* is given, it is called and receives the upcoming request
and the previous request.

> The arguments req and via are the upcoming request and the requests made
> already, oldest first.

```go
    ...
    Client: &http.Client{
        Timeout: 30 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    ...
```

Using special case: `http.ErrUseLastResponse` to track responses.

> As a special case, if CheckRedirect returns ErrUseLastResponse, then the
> most recent response is returned with its body unclosed, along with a nil
> error.

Example: [x/record3xx.go](x/record3xx.go) -- a client that keeps track of redirect hops.

```shell
$ go run record3xx.go http://bibpurl.oclc.org/web/6147
[1] 302 Moved Temporarily http://bibpurl.oclc.org/web/6147 <nil>
[2] 301 Moved Permanently http://www.math.washington.edu/~ejpecp/ECP/index.php <nil>
[3] 302 Found https://www.math.washington.edu/~ejpecp/ECP/index.php <nil>
[4] 301 Moved Permanently https://math.washington.edu/~ejpecp/ECP/index.php <nil>
[5] 301 Moved Permanently https://sites.math.washington.edu/~ejpecp/ECP/index.php <nil>
[6] 200 OK https://sites.math.washington.edu/~burdzy/EJPECP <nil>

$ go run record3xx.go ub.uni-leipzig.de
[1] 301 Moved Permanently http://ub.uni-leipzig.de <nil>
[2] 307 Temporary Redirect https://www.ub.uni-leipzig.de/ <nil>
[3] 200 OK https://www.ub.uni-leipzig.de/start/ <nil>
```

## Transport

Transport has a few more options.

> Transport is an implementation of RoundTripper that supports HTTP,
> HTTPS, and HTTP proxies (for either HTTP or HTTPS with CONNECT).

```go
tr := &http.Transport{
    MaxIdleConns:       10,
    IdleConnTimeout:    30 * time.Second,
    DisableCompression: true,
}
client := &http.Client{Transport: tr}
resp, err := client.Get("https://example.com")
```

The default client uses a default transport:

```go
var DefaultTransport RoundTripper = &Transport{
    Proxy: ProxyFromEnvironment,
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
        DualStack: true,
    }).DialContext,
    ForceAttemptHTTP2:     true,
    MaxIdleConns:          100,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
}
```

A few more (newer) options to control various buffer sizes.

# Alternative Client implementation with retries

* [sethgrid/pester](https://github.com/sethgrid/pester)

> pester wraps Go's standard lib http client to provide several options to
> increase resiliency in your request. If you experience poor network
> conditions or requests could experience varied delays, you can now pester the
> endpoint for data.

Very easy to swap:

```go
/* swap in replacement, just switch
   http.{Get|Post|PostForm|Head|Do} to
   pester.{Get|Post|PostForm|Head|Do}
*/
resp, err := pester.Get("http://sethammons.com")
```

Supports a variety of backoff strategies, e.g. LinearBackoff or
ExponentialJitterBackoff.

A http.Client [is
wrapped](https://github.com/sethgrid/pester/blob/68a33a018ad0ac8266f272ec669307a1829c0486/pester.go#L27-L53),
exposing additional options - a kind of decoration of a
[http.Client](https://golang.org/pkg/net/http/#Client), which itself can have
custom configuration.

Additional resiliency on the application level by supporting [429 Too Many
Requests](https://github.com/sethgrid/pester/blob/68a33a018ad0ac8266f272ec669307a1829c0486/pester.go#L52).

# Tracing

Allows interception of HTTP requests on various occasions (currently 16).

* create a usual request
* create a `http.ClientTrace` and provide callback
* decorate request
* run

```go
    // ...
    req, err := http.NewRequest("GET", "https://golangleipzig.space", nil)
    if err != nil {
        log.Fatal(err)
    }

    // callbacks
    trace := &httptrace.ClientTrace{
        GotConn: func(connInfo httptrace.GotConnInfo) {
            fmt.Printf("Got Conn: %+v\n", connInfo)
        },
        DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
            fmt.Printf("DNS Info: %+v\n", dnsInfo)
        },
    }

    // decorate
    req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

    // client "Do" or Transport "RoundTrip"
    _, err = http.DefaultTransport.RoundTrip(req)
    if err != nil {
        log.Fatal(err)
    }

    // ...
```

* basic example: [x/trace1.go](x/trace1.go)
* all callbacks: [x/trace2.go](x/trace2.go)

Tool: [httpstat](https://github.com/davecheney/httpstat)


# HTTP/2 Pop Quiz

* What percentage of sites use HTTP/2 today?

<!--

    HTTP/2 is used by 43.4% of all the websites.
    https://w3techs.com/technologies/details/ce-http2
-->

