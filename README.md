
# Alias plugin

This plugin eliminates CNAME records from zone apex by making the subsequent resolved records look like they belong to the zone apex. This behaves similarily to [CloudFlare's Zone Flattening](https://support.cloudflare.com/hc/en-us/articles/200169056-CNAME-Flattening-RFC-compliant-support-for-CNAME-at-the-root).

Preferrably, this should not be used in favour of the RFC drafts for the new [ANAME](https://tools.ietf.org/html/draft-ietf-dnsop-aname-00) records, but the DNS library used by CoreDNS does not support ANAME records yet. 

# Usage

```
$ go get github.com/coredns/coredns
$ go get github.com/serverwentdown/alias
$ cd $GOPATH/src/github.com/coredns/coredns
$ vim plugin.cfg
# Add the line alias:github.com/serverwentdown/alias before the file middleware
$ go generate
$ go build
$ ./coredns -plugins | grep alias
```

This plugin only works with the `file` middleware with `upstream` set, or when A or AAAA records exist alongside the CNAME record.

```
example.com {
  file example.com.db {
    upstream 8.8.8.8
  }
  alias
}
```

All it does is transform records like this:

```
;; ANSWER SECTION:
example.com.	300	IN	CNAME	some.magic.example.org.
some.magic.example.org. 299 IN A	123.123.45.67
```

into

```
;; ANSWER SECTION:
example.com.	299	IN	A	123.123.45.67
```
