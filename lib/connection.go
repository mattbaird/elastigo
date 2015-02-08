// Copyright 2013 Matthew Baird
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elastigo

import (
	"fmt"
	hostpool "github.com/bitly/go-hostpool"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	Version         = "0.0.2"
	DefaultProtocol = "http"
	DefaultDomain   = "localhost"
	DefaultPort     = "9200"
	// A decay duration of zero results in the default behaviour
	DefaultDecayDuration = 0
)

type Conn struct {
	hp     hostpool.HostPool
	once   sync.Once
	client *http.Client

	mu             sync.RWMutex // protects following fields
	protocol       string
	domain         string
	clusterDomains []string
	port           string
	username       string
	password       string
	hosts          []string
	transport      *http.Transport
	requestTracer  func(method, url, body string)

	// To compute the weighting scores, we perform a weighted average of recent response times,
	// over the course of `DecayDuration`. DecayDuration may be set to 0 to use the default
	// value of 5 minutes. The EpsilonValueCalculator uses this to calculate a score
	// from the weighted average response time.
	decayDuration time.Duration
}

func NewConn() *Conn {
	// Copied from http.DefaultTransport for consistency
	//
	// See http://golang.org/pkg/net/http/#DefaultTransport
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &Conn{
		// Maintain these for backwards compatibility
		protocol:       DefaultProtocol,
		domain:         DefaultDomain,
		clusterDomains: []string{DefaultDomain},
		port:           DefaultPort,
		decayDuration:  time.Duration(DefaultDecayDuration * time.Second),

		client:    &http.Client{Transport: t},
		transport: t,
	}
}

func (c *Conn) SetProtocol(protocol string) {
	c.mu.Lock()
	c.protocol = protocol
	c.mu.Unlock()
}

func (c *Conn) SetDomain(domain string) {
	c.mu.Lock()
	c.domain = domain
	c.mu.Unlock()
}

func (c *Conn) SetClusterDomains(domains []string) {
	c.mu.Lock()
	c.clusterDomains = domains
	c.mu.Unlock()
}

func (c *Conn) SetPort(port string) {
	c.mu.Lock()
	c.port = port
	c.mu.Unlock()
}

func (c *Conn) SetUsername(username string) {
	c.mu.Lock()
	c.username = username
	c.mu.Unlock()
}

func (c *Conn) SetPassword(password string) {
	c.mu.Lock()
	c.password = password
	c.mu.Unlock()
}

func (c *Conn) SetHosts(newhosts []string) {
	c.mu.Lock()
	c.hosts = newhosts
	c.mu.Unlock()

	// Reinitialise the host pool Pretty naive as this will nuke the current
	// hostpool, and therefore reset any scoring
	c.initializeHostPool()
}

func (c *Conn) SetRequestTracer(tracer func(method, url, body string)) {
	c.mu.Lock()
	c.requestTracer = tracer
	c.mu.Unlock()
}

func (c *Conn) SetDecayDuration(duration time.Duration) {
	c.mu.Lock()
	c.decayDuration = duration
	c.mu.Unlock()

	// Reinitialise the host pool Pretty naive as this will nuke the current
	// hostpool, and therefore reset any scoring
	c.initializeHostPool()
}

func (c *Conn) SetMaxIdleConnsPerHost(n int) {
	if n < 0 {
		n = 0
	}

	c.mu.Lock()
	c.transport.MaxIdleConnsPerHost = n
	c.mu.Unlock()
}

// Set up the host pool to be used
func (c *Conn) initializeHostPool() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If no hosts are set, fallback to defaults
	if len(c.hosts) == 0 {
		c.hosts = append(c.hosts, fmt.Sprintf("%s:%s", c.domain, c.port))
	}

	// Epsilon Greedy is an algorithm that allows HostPool not only to
	// track failure state, but also to learn about "better" options in
	// terms of speed, and to pick from available hosts based on how well
	// they perform. This gives a weighted request rate to better
	// performing hosts, while still distributing requests to all hosts
	// (proportionate to their performance).  The interface is the same as
	// the standard HostPool, but be sure to mark the HostResponse
	// immediately after executing the request to the host, as that will
	// stop the implicitly running request timer.
	//
	// A good overview of Epsilon Greedy is here http://stevehanov.ca/blog/index.php?id=132
	if c.hp != nil {
		c.hp.Close()
	}
	c.hp = hostpool.NewEpsilonGreedy(
		c.hosts, c.decayDuration, &hostpool.LinearEpsilonValueCalculator{})
}

func (c *Conn) Close() {
	c.hp.Close()
}

func (c *Conn) NewRequest(method, path, query string) (*Request, error) {
	// Setup the hostpool on our first run
	c.once.Do(c.initializeHostPool)

	// Get a host from the host pool
	hr := c.hp.Get()

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Get the final host and port
	host, portNum := splitHostnamePartsFromHost(hr.Host(), c.port)

	// Build request
	var uri string
	// If query parameters are provided, the add them to the URL,
	// otherwise, leave them out
	if len(query) > 0 {
		uri = fmt.Sprintf("%s://%s:%s%s?%s", c.protocol, host, portNum, path, query)
	} else {
		uri = fmt.Sprintf("%s://%s:%s%s", c.protocol, host, portNum, path)
	}
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "elasticSearch/"+Version+" ("+runtime.GOOS+"-"+runtime.GOARCH+")")

	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	newRequest := &Request{
		Request:      req,
		hostResponse: hr,
		client:       c.client,
	}
	return newRequest, nil
}

// Split apart the hostname on colon
// Return the host and a default port if there is no separator
func splitHostnamePartsFromHost(fullHost string, defaultPortNum string) (string, string) {

	h := strings.Split(fullHost, ":")

	if len(h) == 2 {
		return h[0], h[1]
	}

	return h[0], defaultPortNum
}
