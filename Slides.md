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
* timeouts
* retries
* tracing

# Interfaces

* [net] package contains 105040 LOC
* [net/http] 59154 LOC (of which 30297 are tests)

The [net/http](https://golang.org/pkg/net/http/) contains 12 interfaces (02/2020):

<!--  $ find . -type f | xargs cat | grep '^type[ ]*[A-Z].* interface {' | awk '{print $2}' | sort -->

* BufferPool
* CloseNotifier
* CookieJar
* File
* FileSystem
* Flusher
* **Handler**, `ServeHTTP(ResponseWriter, *Request)`
* Hijacker
* PublicSuffixList
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

Example: [x/wbshare.go](x/webshare.go).



