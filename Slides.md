# HTTP Getaway 1.0

Configuration options and extension points for HTTP clients.

> Golang Leipzig, 2020-02-21, 19:00, Martin Czygan <martin.czygan@gmail.com>

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
* ResponseWriter
* RoundTripper

## Hijacker

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

# Configuration and Timeouts

# Redirect Tracking

# Tracing
