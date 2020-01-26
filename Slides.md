# HTTP

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

# Thanks

Thanks, that's it! The above covers probably 80% of your needs. Do not make it
more complicated, if not necessary. Bye. And yes, it's great that so many
things are hidden behind this one `Get` method.

Any questions?

# Retry

* HTTP is a (text-based) application layer protocol ([RFC2068](https://tools.ietf.org/html/rfc2068)ff)
