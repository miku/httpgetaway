// Just follow redirect and record responses.
//
// $ go run record3xx.go http://bibpurl.oclc.org/web/6147
// [1] 302 Moved Temporarily http://bibpurl.oclc.org/web/6147 <nil>
// [2] 301 Moved Permanently http://www.math.washington.edu/~ejpecp/ECP/index.php <nil>
// [3] 302 Found https://www.math.washington.edu/~ejpecp/ECP/index.php <nil>
// [4] 301 Moved Permanently https://math.washington.edu/~ejpecp/ECP/index.php <nil>
// [5] 301 Moved Permanently https://sites.math.washington.edu/~ejpecp/ECP/index.php <nil>
// [6] 200 OK https://sites.math.washington.edu/~burdzy/EJPECP <nil>
//
// $ go run record3xx.go http://nature.com/npj-microgravity
// [1] 301 Moved Permanently http://nature.com/npj-microgravity <nil>
// [2] 303 See Other http://www.nature.com/npj-microgravity <nil>
// [3] 302 Found https://idp.nature.com/authorize?response_type=cookie&client_id=grover&redirect_uri=http://www.nature.com/npj-microgravity <nil>
// [4] 302 Found https://idp.nature.com/transit?redirect_uri=http://www.nature.com/npj-microgravity&code=e5561b78-f946-42ca-91f0-4c6f2d3c28c9 <nil>
// [5] 302 Found http://www.nature.com/npj-microgravity?error=cookies_not_supported&code=e5561b78-f946-42ca-91f0-4c6f2d3c28c9 <nil>
// [6] 303 See Other http://www.nature.com/npj-microgravity/index.html <nil>
// [7] 302 Found https://idp.nature.com/authorize?response_type=cookie&client_id=grover&redirect_uri=http://www.nature.com/npj-microgravity/index.html <nil>
// [8] 302 Found https://idp.nature.com/transit?redirect_uri=http://www.nature.com/npj-microgravity/index.html&code=35abd41b-923c-4301-a431-e3cbc988a15a <nil>
// [9] 404 Not Found http://www.nature.com/npj-microgravity/index.html?error=cookies_not_supported&code=35abd41b-923c-4301-a431-e3cbc988a15a <nil>

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	ErrMaxRedirectsExceeded = errors.New("max redirects exceeded")

	ua = flag.String("ua", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)",
		"user agent, might trigger different redirects (e.g. http://google.com)")
)

// Hop encapsulates a response and error, link serves as a shortcut.
type Hop struct {
	Link     string
	Response *http.Response
	Error    error
}

// Client allows to fetch HTTP resources and keep track of redirect responses.
type Client struct {
	Header       http.Header  // Extra headers.
	Client       *http.Client // Underlying client, we do not embed, we only use GET.
	MaxRedirects int
	Hops         []Hop
}

func prependSchema(s string) string {
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}
	return "http://" + s
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("usage: %s URL", os.Args[0])
	}
	u := prependSchema(flag.Arg(0))

	client := New()
	client.Header = http.Header{
		"User-Agent": []string{*ua},
	}
	_, err := client.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(client.DumpHops())
}

// New creates a new client with usable defaults.
func New() *Client {
	return &Client{
		Header: http.Header{},
		Client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		MaxRedirects: 100,
	}
}

// DumpHops returns a human readable list of hops.
func (c *Client) DumpHops() string {
	var buf bytes.Buffer
	for i, hop := range c.Hops {
		link, err := url.PathUnescape(hop.Link)
		if err != nil {
			link = hop.Link
		}
		if hop.Response == nil {
			fmt.Fprintf(&buf, "[%d] NA %v %v\n",
				i+1, link, hop.Error)
		}
		fmt.Fprintf(&buf, "[%d] %v %v %v\n",
			i+1, hop.Response.Status, link, hop.Error)
	}
	return buf.String()
}

// Get fetches a url and keeps track of intermediate hosts. It returns the
// response and error of the last request or the first error encountered. It is
// the responsibilty of the caller to close the response body.
func (c *Client) Get(u string) (*http.Response, error) {
	c.Hops = nil // Reset hops on new request.
	var i int

	// Redirect loop.
	for i < c.MaxRedirects {
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}
		req.Header = c.Header
		resp, err := c.Client.Do(req)
		c.Hops = append(c.Hops, Hop{
			Link:     u,
			Response: resp,
			Error:    err,
		})
		if err != nil {
			return resp, err
		}
		// Try to find location header.
		rurl, rerr := resp.Location()
		if rerr == http.ErrNoLocation {
			// We reached the final destination, return the original response
			// and error.
			return resp, err
		}
		if rerr != nil {
			return resp, rerr
		}
		if rurl == nil {
			return resp, err
		}
		u = rurl.String()
		i++
	}
	return nil, ErrMaxRedirectsExceeded
}
