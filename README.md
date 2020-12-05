# go-nsdp
NSDP(Netgear Switch Discovery Protocol) client library for go

# tested switches
- GS308E (V1.00.05JP with v2auth)

# usage
```
# nsdp-cli query model-name host-name ipaddr mac
{
  "host_name": "foobar",
  "model_name": "GS308E"
  "ip": "192.168.1.111",
}

# ./nsdp-cli -password <password> set tag-vlan:1000:1,2,3,4,5,6,7:
{}
```
