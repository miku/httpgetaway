# On HTTP Proxies

Background: I recently had the chance to learn more about proxy servers.

Proxy servers are ubiquitous, and they embody the rule:

> All problems in computer science can be solved by another level of indirection.

If you have a networking problem, a proxy might be a solution. Beside that,
proxies can be fun and useful.

# Outline

* Proxy types (forward, reverse, web, security, performance, ...)
* What standard library support is there?
* Writing a proxy from scratch: HTTP, HTTPS, websockets.
* Existing libraries, review.

# Proxy Types

* forward
* reverse (nginx, https://docs.traefik.io/getting-started/quick-start/,
* TLS offloading

# Standard library support

* set proxy on transport, defaults to "ProxyFromEnvironment"
* x/net/proxy
* http/httputil/reverseproxy.go, https://gist.github.com/JalfResi/6287706

# Misc

* https://github.com/creack/goproxy
* https://github.com/jamesmoriarty/goforward
