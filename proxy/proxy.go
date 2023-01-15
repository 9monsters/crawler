package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type ProxyFunc func(r *http.Request) (*url.URL, error)

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

func (r roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
	index := atomic.AddUint32(&r.index, 1) - 1

	u := r.proxyURLs[index%uint32(len(r.proxyURLs))]
	return u, nil
}

func RoundRobinProxySwitcher(proxURLs ...string) (ProxyFunc, error) {
	if len(proxURLs) < 1 {
		return nil, errors.New("Proxy URL list is empty")
	}

	urls := make([]*url.URL, len(proxURLs))

	for i, u := range proxURLs {
		parseU, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parseU
	}
	return (&roundRobinSwitcher{urls, 0}).GetProxy, nil
}
