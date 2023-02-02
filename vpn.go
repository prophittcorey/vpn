package vpn

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	ErrNotFound  = fmt.Errorf(`error: does not appear to be through a known vpn`)
	ErrInvalidIP = fmt.Errorf(`error: does not appear to be a valid IP address`)
)

var (
	// A map of vpn IP range sources. All sources will be fetched concurrently
	// and merged together.
	Sources = map[string]map[string][]byte{
		// "protonvpn": {
		// 	"https://raw.githubusercontent.com/X4BNet/lists_vpn/main/input/ips/protonvpn.txt": []byte{},
		// },

		/* NOTE: This list doesn't separate providers unfortunately... */
		"x4b-merged": {
			"https://raw.githubusercontent.com/X4BNet/lists_vpn/main/ipv4.txt": []byte{},
		},
	}

	// HTTPClient is used to perform all HTTP requests. You can specify your own
	// to set a custom timeout, proxy, etc.
	HTTPClient = http.Client{
		Timeout: 3 * time.Second,
	}

	// CachePeriod specifies the amount of time an internal cache of vpn subnets are used
	// before refreshing the subnets.
	CachePeriod = 45 * time.Minute

	// UserAgent will be used in each request's user agent header field.
	UserAgent = "github.com/prophittcorey/vpn"
)

// TODO: Is there a more efficient data structure for IP -> Subnet look ups? Currently it's
// O(n), but there is likely a O(1) solution out there.

type vpns struct {
	sync.RWMutex
	subnets map[string][]*net.IPNet /* origin -> [*IPNet], ex: 'protonvpn' => [...] */
}

var (
	networks    = vpns{}
	lastFetched = time.Now()
)

func refresh() error {
	if len(networks.subnets) == 0 || time.Now().After(lastFetched.Add(CachePeriod)) {
		wg := sync.WaitGroup{}

		for origin, sources := range Sources {
			origin := origin
			sources := sources

			for url, _ := range sources {
				wg.Add(1)

				go (func(url string) {
					defer wg.Done()

					req, err := http.NewRequest(http.MethodGet, url, nil)

					if err != nil {
						return
					}

					req.Header.Set("User-Agent", UserAgent)

					res, err := HTTPClient.Do(req)

					if err != nil {
						return
					}

					if bs, err := io.ReadAll(res.Body); err == nil {
						Sources[origin][url] = bs
					}
				})(url)
			}
		}

		wg.Wait()

		/* merge / dedupe all domains */

		subnets := map[string][]*net.IPNet{}

		for origin, sources := range Sources {
			origin := origin
			sources := sources

			for _, bs := range sources {
				for _, cidr := range bytes.Fields(bs) {
					if _, subnet, err := net.ParseCIDR(string(cidr)); err == nil {
						subnets[origin] = append(subnets[origin], subnet)
					}
				}
			}
		}

		/* clear byte cache of sources */

		for origin, sources := range Sources {
			for url, _ := range sources {
				Sources[origin][url] = []byte{}
			}
		}

		/* update global networks cache */

		networks.Lock()

		networks.subnets = subnets
		lastFetched = time.Now()

		networks.Unlock()
	}

	return nil
}

// Check returns the VPN associated with the IP, or an error.
func Check(addr string) (string, error) {
	ip := net.ParseIP(addr)

	if ip == nil {
		return "", ErrInvalidIP
	}

	refresh()

	for origin, subnets := range networks.subnets {
		origin := origin
		subnets := subnets

		for _, subnet := range subnets {
			if subnet.Contains(ip) {
				return origin, nil
			}
		}
	}

	return "", ErrNotFound
}

// Subnets will return all known VPN subnets.
func Subnets() []*net.IPNet {
	refresh()

	ss := []*net.IPNet{}

	for _, subnets := range networks.subnets {
		for _, subnet := range subnets {
			ss = append(ss, subnet)
		}
	}

	return ss
}
