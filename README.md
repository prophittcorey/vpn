# VPN

[![Go Reference](https://pkg.go.dev/badge/github.com/prophittcorey/vpn.svg)](https://pkg.go.dev/github.com/prophittcorey/vpn)

A golang package and command line tool for the analysis and identification of
VPN IP addresses.

## Package Usage

```golang
import "github.com/prophittcorey/vpn"

if provider, err := vpn.Check("12.34.56.78"); err == nil {
  fmt.Printf("Looks like a %s IP address.\n", provider)
}

vpn.Subnets() // => [*IPNet, ...]
```

## Tool Usage

```bash
# Install the latest tool.
go install github.com/prophittcorey/vpn/cmd/vpn@latest

# Dump all known vpn provider subnets.
vpn --subnets

# Check a specific IP address.
vpn --check 12.34.56.78
```

## License

The source code for this repository is licensed under the MIT license, which you can
find in the [LICENSE](LICENSE.md) file.
