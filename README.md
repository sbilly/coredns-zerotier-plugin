# coredns-zerotier-plugin

The CoreDNS plugin for ZeroTierOne. Query `https://my.zerotier.com/api` with token.

modified from [coredns-netbox-plugin](https://github.com/oz123/coredns-netbox-plugin), using [go-ztcentral](github.com/zerotier/go-ztcentral)

## Installation

Add the following to `plugin.cfg`:

```
zerotier:github.com/sbilly/coredns-zerotier-plugin
```

Build and run:

```shell
go generate
go build
./coredns -conf Corefile
```

## Example configuration

```
.:53 {
    zerotier {
        url https://my.zerotier.com/api
        token #################################
        localCacheDuration 15s
   }
}
```

## Query

```shell
dig @127.0.0.1 -p 53 <NAME>
```

returns

```
; <<>> DiG 9.10.6 <<>> @127.0.0.1 -p 53 <NAME>
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 56006
;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;example.                  IN      A

;; ANSWER SECTION:
example.           15      IN      A       <IP_ADDRESS>

;; Query time: 1230 msec
;; SERVER: 127.0.0.1#53(127.0.0.1)
;; WHEN: ...
;; MSG SIZE  rcvd: 69
```

`<NAME>`, `<IP_ADDRESS>` are from https://my.zerotier.com/central-api.html#member-member-get

## Local Development

Add the following to `go.mod`:

```
replace (
    github.com/sbilly/coredns-zerotier-plugin v0.0.0 => <PATH_TO>/coredns-zerotier-plugin
)
```
