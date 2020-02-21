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
* Handler
* Hijacker
* PublicSuffixList
* Pusher
* ResponseWriter
* RoundTripper

