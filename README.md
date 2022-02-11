# SailFirewall

Linux firewall powered by eBPF and XDP

# Requirements
* Go 1.16+
* Linux Kernel 4.15+

# Support feature

* IPV4
* TCP

> Please contribute other protocols support

# Usage

```shell
make
# change lo if you need
sudo ./SailFirewall -iface lo
```

# API

## Add rule

POST /api/v1/rule

```json
{
    "SrcAddr": "127.0.0.1",
    "DstPort": 8000,
    "Proto": 6
}
```

> Proto 6 is TCP

## Get rule

GET /api/v1/rule

```json
{
    "SrcAddr": "127.0.0.1",
    "DstPort": 8000,
    "Proto": 6
}
```

## Delete rule

DELETE /api/v1/rule

```json
{
  "SrcAddr": "127.0.0.1",
  "DstPort": 8000,
  "Proto": 6
}
```

# Reference

[EtherType](https://zh.wikipedia.org/wiki/%E4%BB%A5%E5%A4%AA%E7%B1%BB%E5%9E%8B)

[IP protocol numbers](https://zh.wikipedia.org/wiki/IP%E5%8D%8F%E8%AE%AE%E5%8F%B7%E5%88%97%E8%A1%A8)

[TCP](https://zh.wikipedia.org/wiki/%E4%BC%A0%E8%BE%93%E6%8E%A7%E5%88%B6%E5%8D%8F%E8%AE%AE)


