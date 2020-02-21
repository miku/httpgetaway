# HTTP Getaway

HTTP extension points and alternative implementations in Go.

* X-Location, X-Date
* Martin Czygan <martin.czygan@gmail.com> -- Open Data Engineer at [Internet
  Archive](https://archive.org/).

# All you need to know

```go
r, err := http.Get("https://golangleipzig.space")
if err != nil {
    return err
}
defer r.Body.Close()
```

* [x/hello.go](x/hello.go)

# Thanks

Thanks, that's it! The above covers probably 80% of your needs. Do not make it
more complicated, if not necessary. Bye. And yes, it's great that so many
things are hidden behind this one `Get` method.

Any questions?

# Retry

* HTTP is a (text-based) application layer protocol
* HTTP/1.1, [RFC2616](https://tools.ietf.org/html/rfc2616)
  ([RFC2068](https://tools.ietf.org/html/rfc2068) ...)
* HTTP/2 provides an optimized transport, but: "HTTP's existing semantics
  remain unchanged." ([RFC7540](https://tools.ietf.org/html/rfc7540))

> HTTP/2 is used by 42.8% of all the websites. --
> [https://w3techs.com/technologies/details/ce-http2](https://w3techs.com/technologies/details/ce-http2)

* HTTP/3 is in the pipeline (but [not usable](https://caniuse.com/#feat=http3) today)

# Core elements

* Resources, Representations, limited set of operations (verbs)
* Status codes
* Redirection
* Transport
* Security

# Today

* Tracing
* Alternative clients
* Utilities

# Go net and net/http packages

Go comes with solid networking support in the standard library. Especially, the
`net/http` suite is both usable and extendable.

First, we want to look at the extension points.

# Tracing

We saw that a HTTP round trip takes some time. WHW?

* [Introducing HTTP Tracing](https://blog.golang.org/http-tracing) (2016)

> In Go 1.7 we introduced HTTP tracing, a facility to gather fine-grained
> information throughout the lifecycle of an HTTP client request.

Facilities in `net/http/httptrace` ([https://golang.org/pkg/net/http/httptrace](https://golang.org/pkg/net/http/httptrace)).

# A struct full of callbacks

```go
trace := &httptrace.ClientTrace{
        GotConn: func(connInfo httptrace.GotConnInfo) {
            fmt.Printf("Got Conn: %+v\n", connInfo)
        },
        DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
            fmt.Printf("DNS Info: %+v\n", dnsInfo)
        },
    }
```

Wrap a `http.Request` with client trace (using [WithContext](https://golang.org/pkg/net/http/#Request.WithContext), 1.7).

```go
req, _ := http.NewRequest("GET", "http://example.com", nil)
req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
```

# Implementations may use the hooks

Implementations of `http.RoundTripper` can choose, whether they call back to the tracer.

# Example

* [x/trace1.go](x/trace1.go)

```
$ go run x/trace1.go
DNS Info: {Addrs:[{IP:185.199.109.153 Zone:} {IP:185.199.110.153 Zone:}] Err:<nil> Coalesced:false}
Got Conn: {Conn:0xc000108000 Reused:false WasIdle:false IdleTime:0s}
```

# Tracing hooks

The `httptrace.ClientTrace` contains 16 hook methods.

Example:

* [x/trace2.go](x/trace2.go)

```
$ go run x/trace2.go
        91.878µs    |Get Conn                   |golangleipzig.space:443
       390.407µs    |DNS Start                  |{Host:golangleipzig.space}
      2.515361ms    |DNS Info                   |{Addrs:[{IP:185.199.109.153 Zone:} ...
      2.582607ms    |Conn Start                 |tcp 185.199.109.153:443
     20.334243ms    |Conn Done                  |tcp 185.199.109.153:443 <nil>
             ...
    187.396936ms    |Wrote Request              |{Err:<nil>}
    207.149618ms    |Got First Response Byte    |
    207.441984ms    |HTTP status code           |200 OK
```

# Tracing across redirects

* possibly by using a custom `http.RoundTripper` which keeps track of the current URL
* functions are first class values in Go, the ClientTrace hook may be a method
  on a struct (with state, e.g. the current URL)

Example:

* [x/trace3.go](x/trace3.go)

# Command line tool

* [httpstat](https://github.com/davecheney/httpstat)

![](static/httpstat.png)

# What is a RoundTripper?

> RoundTripper is an interface representing the ability to execute a single
> HTTP transaction, obtaining the Response for a given Request.

* single HTTP transaction
* responsible to retrieve a response (only)
* the `http.Request` should be left alone by the caller until `resp.Body` is closed
* must close request body (maybe after return)

```go
type RoundTripper interface {
    RoundTrip(*Request) (*Response, error)
}
```

# Use cases

* implement custom round tripper for testing request, response scenarios
* a RT that adds auth headers
* intercept request for caching
* intercept for auditing
* rate limiting
* not concerned with HTTP status codes

# Example caching transport

Keep response body in memory, very basic.

* [x/cachingrt.go](x/cachingrt.go)

# Let's talk about timeouts

* often first contact with alternatives - the default client has no timeouts