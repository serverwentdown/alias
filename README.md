
# alias

The *alias* plugin eliminates CNAME records from zone apex by making the subsequent resolved records look like they belong to the zone apex. This behaves similarily to [CloudFlare's Zone Flattening](https://support.cloudflare.com/hc/en-us/articles/200169056-CNAME-Flattening-RFC-compliant-support-for-CNAME-at-the-root).

This plugin works only with plugins that produce A or AAAA records alongside the CNAME record. Examples include `auto` and `file`. However, you might need to adjust the order of this plugin to use it with other plugins. 

> Preferrably, this should not be used in favour of the RFC drafts for the new [ANAME](https://tools.ietf.org/html/draft-ietf-dnsop-aname-00) records, but the DNS library used by CoreDNS does not support ANAME records yet. 

Release builds can be found [here](https://github.com/serverwentdown/alias/releases)

## Syntax

```
alias
```

## Examples

```
example.com {
	file db.example.com
	alias
}
# This is used to resolve CNAME records by the `file` plugin. Modify accordingly
. {
	forward . 1.1.1.1 1.0.0.1
}
```

This will transform responses like this:

```
;; ANSWER SECTION:
example.com.		3600	IN	CNAME	two.example.org.
two.example.org.	3600	IN	CNAME	one.example.net.
one.example.net.	3600	IN	A	127.0.0.1
```

into

```
;; ANSWER SECTION:
example.com.		3600	IN	A	127.0.0.1
```

See [`example/`](example/) for a more extensive example. 

## Installation

As per [CoreDNS docs](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/), there are two ways.

### Build with compile-time configuration file

```
$ git clone https://github.com/coredns/coredns
$ cd coredns
$ vim plugin.cfg
# Add the line alias:github.com/serverwentdown/alias before the file middleware
$ go generate
$ go build
$ ./coredns -plugins | grep alias
```

### Build with external golang source code

```
$ git clone https://github.com/serverwentdown/alias
$ cd alias/coredns
$ go build
$ ./coredns -plugins | grep alias
```
