# dhcpwatch

[![build](https://github.com/HandOfGod94/dhcpwatch/actions/workflows/build.yml/badge.svg)](https://github.com/HandOfGod94/dhcpwatch/actions/workflows/build.yml)
![go-version](https://img.shields.io/badge/go-1.15-blue)

It watches dhcp database file changes and exports
prometheus metrics which can visualized as table on grafana.

Generally location for `dhcp` database on dhcp server is:
`/var/lib/dhcp/dhcpd.leases`


### Commands

```shell
# run tests locally
make test

# build binary
make build

# build binary for raspberry pi
make pi-build
```
