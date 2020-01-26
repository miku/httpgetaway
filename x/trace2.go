package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func tlsName(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "1.0"
	case tls.VersionTLS11:
		return "1.1"
	case tls.VersionTLS12:
		return "1.2"
	case tls.VersionTLS13:
		return "1.3"
	case tls.VersionSSL30:
		return "SSL3.0 (broken)"
	default:
		return fmt.Sprintf("UNKNOWN: %x", version)
	}
}

func prependHTTP(s string) string {
	if !strings.HasPrefix(s, "http") {
		return "http://" + s
	}
	return s
}

func main() {
	site := "https://golangleipzig.space"
	flag.Parse()
	if flag.NArg() > 0 {
		site = prependHTTP(flag.Arg(0))
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	defer w.Flush()
	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			fmt.Fprintf(w, "% 16s\tGet Conn\t%s\n", time.Since(start), hostPort)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Fprintf(w, "% 16s\tGot Conn\t%+v\n", time.Since(start), connInfo)
		},
		PutIdleConn: func(err error) {
			fmt.Fprintf(w, "% 16s\tPut Idle Conn\t\n", time.Since(start))
		},
		GotFirstResponseByte: func() {
			fmt.Fprintf(w, "% 16s\tGot First Response Byte\t\n", time.Since(start))
		},
		Got100Continue: func() {
			fmt.Fprintf(w, "% 16s\tGot 100 Continue\t%d, %s\n", time.Since(start))
		},
		Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
			fmt.Fprintf(w, "% 16s\tGot 1xx\t%d, %s\n", time.Since(start), code, header)
			return nil
		},
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			fmt.Fprintf(w, "% 16s\tDNS Start\t%+v\n", time.Since(start), dnsInfo)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Fprintf(w, "% 16s\tDNS Info\t%+v\n", time.Since(start), dnsInfo)
		},
		ConnectStart: func(network, addr string) {
			fmt.Fprintf(w, "% 16s\tConn Start\t%s %s\n", time.Since(start), network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Fprintf(w, "% 16s\tConn Done\t%s %s %v\n", time.Since(start), network, addr, err)
		},
		TLSHandshakeStart: func() {
			fmt.Fprintf(w, "% 16s\tTLS Start\t\n", time.Since(start))
		},
		TLSHandshakeDone: func(connState tls.ConnectionState, err error) {
			fmt.Fprintf(w, "% 16s\tTLS Done\t%s %+v\n", time.Since(start), tlsName(connState.Version), err)
		},
		WroteHeaderField: func(key string, value []string) {
			fmt.Fprintf(w, "% 16s\tHeader\t%s: %v\n", time.Since(start), key, value)
		},
		WroteHeaders: func() {
			fmt.Fprintf(w, "% 16s\tHeader Done\t\n", time.Since(start))
		},
		Wait100Continue: func() {
			fmt.Fprintf(w, "% 16s\tWait100\t\n", time.Since(start))
		},
		WroteRequest: func(reqInfo httptrace.WroteRequestInfo) {
			fmt.Fprintf(w, "% 16s\tWrote Request\t%+v\n", time.Since(start), reqInfo)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Fprintf(w, "% 16s\tHTTP status code\t%s\n", time.Since(start), resp.Status)
	if u, err := resp.Location(); err == nil {
		fmt.Fprintf(w, "% 16s\tLocation\t%s\n", time.Since(start), u)
	}

}
