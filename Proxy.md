# HTTP Getaway Part 2: HTTP proxy intro

> Martin Czygan, 2020-09-17, 19:00 CEST, https://golangleipzig.space

Background: I recently had the chance to learn more about proxy servers.

Proxy servers are ubiquitous, and they embody the rule:

> All problems in computer science can be solved by another level of indirection.

If you have a networking problem, a proxy might be a solution. Beside that,
proxies can be fun and useful.

# Outline

* Proxy types (forward, reverse, ...)
* What standard library support is there?
* Writing a proxy from scratch: HTTP, HTTPS.

# Proxy Types

![](static/proxytypes.png)

* forward
* reverse

## Forward

* implicit and explicit
* implicit: corporate FW; SSL decrypt
* explicit: need to set on client (mitm, archiving, caching, security, ip rotation, ...)

## Reverse

* handle incoming requests, load balancing, "HAvailability", mu-service
  configuration and routing, rate limiting, a/b and canary testing, ...

In this talk, we focus on forward proxies only.

# How does a client use a proxy?

* environment variables, [are these
  standard](https://superuser.com/questions/944958/are-http-proxy-https-proxy-and-no-proxy-environment-variables-standard)?
probably not.

Where is the code to handle the environment variables?

The `http.DefaultTransport` sets the `Proxy` field:

```go
var DefaultTransport RoundTripper = &Transport{
    Proxy: ProxyFromEnvironment,
    ...
```

And if nothing else is set, the default transport is used be the default client
(also: "make the zero value useful"):

```go
// DefaultClient is the default Client and is used by Get, Head, and Post.
var DefaultClient = &Client{}

...

func (c *Client) transport() RoundTripper {
    if c.Transport != nil {
        return c.Transport
    }
    return DefaultTransport
}
```

Which then is used for the verb functions:

```go
func Get(url string) (resp *Response, err error) {
    return DefaultClient.Get(url)
}
```

The `ProxyFromEnvironment` is actually a function type: `func(*Request)
(*url.URL, error)` - which is quite nice.

The implementation looks up the proxy from environment variables only once
(since it seems this lookup can be expensive on Windows ("This mitigates
expensive lookups on some platforms (e.g. Windows)").

Here, the standard library code branches out to:

* [x/net/http/httpproxy](https://godoc.org/golang.org/x/net/http/httpproxy)

> // Package httpproxy provides support for HTTP proxy determination
> // based on environment variables, as provided by net/http's
> // ProxyFromEnvironment function.

Fun fact, there seemingly was a CGI (catchy name) CVE in 2016 -
https://httpoxy.org/ - CGI puts the Proxy header from an incoming request into
HTTP_PROXY, there you go. (https://github.com/golang/go/issues/16405).




# Standard library support

* set proxy on transport, defaults to "ProxyFromEnvironment"
* x/net/proxy
* http/httputil/reverseproxy.go, https://gist.github.com/JalfResi/6287706

# Misc

* https://github.com/creack/goproxy
* https://github.com/jamesmoriarty/goforward
* https://github.com/smartystreets/cproxy
* https://github.com/davidfstr/nanoproxy
