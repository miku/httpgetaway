# Ten extension points net/http

* the net/http package is both very usable out of the box and extensible
* you can customize and instrument many parts
* let's look at ten examples

# Timeouts

* How many timeouts are there?
* Conn, KeepAlive, ...

# Managing Redirects

* default behaviour
* 30X, tracking redirects

# Custom Transport

* http.RoundTripper
* rate limiting transport

# Tracing Requests

# Hijacking Connections

* http.Hijacker

# Filesystems

* http.File and http.FileSystem
