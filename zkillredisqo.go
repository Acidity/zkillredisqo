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
	Version          = "1.0.0"
	DefaultUserAgent = "zkillredisqo v" + Version + " - github.com/morpheusxaut/zkillredisqo"
	ZKillRedisQURL   = "https://redisq.zkillboard.com/listen.php"
)

type Poller struct {
	Kills  chan *Kill
	Errors chan error

	queueID     string
	timeToWait  int
	userAgent   string
	preparedURL string
	client      *http.Client
	wg          *sync.WaitGroup
	stop        bool
}

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

func (p *Poller) SetTimeToWait(ttw int) {
	p.timeToWait = ttw
	p.prepareURL()
}

func (p *Poller) SetUserAgent(ua string) {
	p.userAgent = ua
}

func (p *Poller) Stop() {
	p.stop = true
}

func (p *Poller) Wait() {
	p.wg.Wait()
}

func (p *Poller) StopAndWait() {
	p.Stop()
	p.Wait()
}

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
			continue
		}

		p.Kills <- kill
	}
}

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
