package zkillredisqo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	// Version of zkillredisqo library
	Version = "1.0.1"
	// DefaultUserAgent defines a user agent to use for HTTP requests if separate one is specified by the user
	DefaultUserAgent = "zkillredisqo v" + Version + " - github.com/morpheusxaut/zkillredisqo"
	// ZKillRedisQURL defines the URL of zKillboard's RedisQ service
	ZKillRedisQURL = "https://redisq.zkillboard.com/listen.php"
)

// Poller allows for polling of Kills from zKillboard's RedisQ service
type Poller struct {
	// Kills will receive all parsed kills as they are received from zKillboard
	Kills chan *Kill
	// Errors will receive any error encountered while retrieving or parsing kills
	Errors chan error

	// queueID can be used to identify an app to the RedisQ service (default empty)
	queueID string
	// timeToWait indicates how long RedisQ should wait before returning "null" kills if no new data is available
	timeToWait int
	// userAgent defines the user agent used for HTTP requests
	userAgent string
	// preparedURL stores the RedisQ URL and all optional parameters pre-parsed
	preparedURL string
	// client represents the HTTP client to use for requests
	client *http.Client
	// wg is used to allow for proper shutdown, allowing to wait for all requests to finish before exiting
	wg *sync.WaitGroup
	// stop indicates whether the polling loop should be stopped
	stop bool
}

// NewPoller creates a new poller and starts the polling loop
func NewPoller(client *http.Client) *Poller {
	if client == nil {
		client = &http.Client{
			// zKillboard's RedisQ returns a null package after 10 seconds, allow two extra seconds to account for latency
			Timeout: time.Second * 12,
		}
	}

	p := &Poller{
		Kills:       make(chan *Kill),
		Errors:      make(chan error),
		queueID:     "",
		timeToWait:  10,
		userAgent:   DefaultUserAgent,
		preparedURL: "",
		client:      client,
		wg:          &sync.WaitGroup{},
		stop:        false,
	}

	p.prepareURL()

	p.wg.Add(1)
	go p.poll()

	return p
}

// SetTimeToWait adjusts the timeToWait value passed to RedisQ
func (p *Poller) SetTimeToWait(ttw int) {
	p.timeToWait = ttw
	p.prepareURL()
}

// SetUserAgent allows for a custom user agent to be configured
func (p *Poller) SetUserAgent(ua string) {
	p.userAgent = ua
}

// Stop notifies the poller to stop its loop after the next iteration
func (p *Poller) Stop() {
	p.stop = true
}

// Wait blocks until all requests have been finished
func (p *Poller) Wait() {
	p.wg.Wait()
}

// StopAndWait notifies the poller to stop and waits for all requests to finish
func (p *Poller) StopAndWait() {
	p.Stop()
	p.Wait()
}

// poll retrieves kills from zKillboard's RedisQ until a stop is requested
func (p *Poller) poll() {
	defer p.wg.Done()

	for !p.stop {
		kill, err := p.retrieveKill()
		if err != nil {
			p.Errors <- err
			continue
		}
		if kill == nil {
			p.Errors <- errors.New("nil kill")
			continue
		}
		if kill.IsNullKill() {
			// Silently ignore "null" kills and wait for next proper kill package
			continue
		}

		p.Kills <- kill
	}
}

// retrieveKill tries fetching a kill from RedisQ and parsing the returned JSON
func (p *Poller) retrieveKill() (*Kill, error) {
	req, err := http.NewRequest("GET", p.preparedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", p.userAgent)
	req.Header.Add("Content-Type", "application/json")

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	var kill *Kill
	if err = json.NewDecoder(res.Body).Decode(&kill); err != nil {
		return nil, err
	}

	return kill, nil
}

// prepareURL prepares the URL for requests to RedisQ and adds optional parameters as needed
func (p *Poller) prepareURL() {
	u, err := url.Parse(ZKillRedisQURL)
	if err != nil {
		p.Errors <- fmt.Errorf("failed to parse URL: %v", err)
		p.preparedURL = ZKillRedisQURL
		return
	}

	q := u.Query()
	q.Set("ttw", fmt.Sprintf("%d", p.timeToWait))
	if len(p.queueID) > 0 {
		q.Set("queueID", p.queueID)
	}

	u.RawQuery = q.Encode()
	p.preparedURL = u.String()
}
