// Copyright (c) 2022 Hevienz
// Full license can be found in the LICENSE file.

#include "bpf_helpers.h"

#define MAX_RULES   (16)

#define ETH_PROTO_IPV4 (0x0008U)

#define IP_PROTO_TCP (0x06U)
#define IP_PROTO_UDP (0x11U)

// Ethernet header
struct ethhdr {
  __u8 h_dest[6];
  __u8 h_source[6];
  __u16 h_proto;
} __attribute__((packed));

// IPv4 header
struct iphdr {
  __u8 ihl : 4;
  __u8 version : 4;
  __u8 tos;
  __u16 tot_len;
  __u16 id;
  __u16 frag_off;
  __u8 ttl;
  __u8 protocol;
  __u16 check;
  __u32 saddr;
  __u32 daddr;
} __attribute__((packed));

// TCP header
struct tcphdr {
  __u16 source;
  __u16 dest;
  __u32 seq;
  __u32 ack_seq;
  union {
    struct {
      // Field order has been converted LittleEndiand -> BigEndian
      // in order to simplify flag checking (no need to ntohs())
      __u16 ns : 1,
      reserved : 3,
      doff : 4,
      fin : 1,
      syn : 1,
      rst : 1,
      psh : 1,
      ack : 1,
      urg : 1,
      ece : 1,
      cwr : 1;
    };
  };
  __u16 window;
  __u16 check;
  __u16 urg_ptr;
};

struct acl_key {
  __u32 ip_src_addr;
  __u16 dst_port;
  __u8 ip_proto;
  __u8 reserved;
};

BPF_MAP_DEF(aclist) = {
    .map_type = BPF_MAP_TYPE_HASH,
    .key_size = sizeof(struct acl_key),
    .value_size = sizeof(__u64),
    .max_entries = MAX_RULES,
};
BPF_MAP_ADD(aclist);

// XDP program //
SEC("xdp")
int firewall(struct xdp_md *ctx) {
  void *data_end = (void *)(long)ctx->data_end;
  void *data = (void *)(long)ctx->data;

  struct ethhdr *ether = data;
  if (data + sizeof(*ether) > data_end) {
    // Malformed Ethernet header
    return XDP_ABORTED;
  }

  if (ether->h_proto == ETH_PROTO_IPV4) {
    data += sizeof(*ether);
    struct iphdr *ip = data;
    if (data + sizeof(*ip) > data_end) {
      // Malformed IPv4 header
      return XDP_ABORTED;
    }

    if (ip->protocol == IP_PROTO_TCP) {
      data += ip->ihl * 4;
      struct tcphdr *tcp = data;
      if (data + sizeof(*tcp) > data_end) {
        return XDP_ABORTED;
      }

      if (!(tcp->syn && !tcp->ack)) {
        return XDP_PASS;
      }

      struct acl_key key = {
        .ip_src_addr = ip->saddr,
        .ip_proto = ip->protocol,
        .dst_port = tcp->dest,
        .reserved = 0,
      };

      __u64 *count = bpf_map_lookup_elem(&aclist, &key);
      if (count) {
        (*count)++;

        return XDP_PASS;
      } else {
        return XDP_DROP;
      }

    } else if (ip->protocol == IP_PROTO_UDP) {
      return XDP_DROP;
    } else {
      return XDP_DROP;
    }

  } else {
    return XDP_DROP;
  }
}

char _license[] SEC("license") = "GPLv2";
