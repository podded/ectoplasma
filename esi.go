package ectoplasma

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type (
	ESIClient struct {
		HTTP         *http.Client
		UserAgent    string
		ESIRateLimit *safeCounter
		RetryLimit   int
	}
	safeCounter struct {
		cnt int
		mux sync.Mutex
	}

	ESIKillmail struct {
		Attackers     []ESIAttacker `json:"attackers" bson:"attackers"`
		KillmailID    int           `json:"killmail_id" bson:"killmail_id"`
		KillmailTime  time.Time     `json:"killmail_time" bson:"killmail_time"`
		SolarSystemID int           `json:"solar_system_id" bson:"solar_system_id"`
		Victim        ESIVictim     `json:"victim" bson:"victim"`
	}

	ESIAttacker struct {
		AllianceID     int     `json:"alliance_id,omitempty" bson:"alliance_id,omitempty"`
		CorporationID  int     `json:"corporation_id" bson:"corporation_id"`
		CharacterID    int     `json:"character_id" bson:"character_id"`
		DamageDone     int     `json:"damage_done" bson:"damage_done"`
		FinalBlow      bool    `json:"final_blow" bson:"final_blow"`
		SecurityStatus float32 `json:"security_status" bson:"security_status"`
		ShipTypeID     int     `json:"ship_type_id" bson:"ship_type_id"`
		WeaponTypeID   int     `json:"weapon_type_id" bson:"weapon_type_id"`
	}

	ESIVictim struct {
		AllianceID    int         `json:"alliance_id,omitempty" bson:"alliance_id,omitempty"`
		CorporationID int         `json:"corporation_id" bson:"corporation_id"`
		CharacterID   int         `json:"character_id" bson:"character_id"`
		DamageTaken   int         `json:"damage_taken" bson:"damage_taken"`
		Items         []ESIItem   `json:"items" bson:"items"`
		Position      ESIPosition `json:"position" bson:"position"`
		ShipTypeID    int         `json:"ship_type_id" bson:"ship_type_id"`
	}

	ESIItem struct {
		Flag              int `json:"flag" bson:"flag"`
		ItemTypeID        int `json:"item_type_id" bson:"item_type_id"`
		QuantityDropped   int `json:"quantity_dropped,omitempty" bson:"quantity_dropped,omitempty"`
		QuantityDestroyed int `json:"quantity_destroyed,omitempty" bson:"quantity_destroyed,omitempty"`
		Singleton         int `json:"singleton" bson:"singleton"`
	}

	ESIPosition struct {
		X float64 `json:"x" bson:"x"`
		Y float64 `json:"y" bson:"y"`
		Z float64 `json:"z" bson:"z"`
	}
)

func NewESIClient() *ESIClient {
	rateLimESI := &safeCounter{}

	go func() {
		for {
			time.Sleep(time.Second)
			if rateLimESI.Value() > 0 {
				rateLimESI.Dec()
			}
		}
	}()

	return &ESIClient{
		HTTP: &http.Client{
			Timeout: time.Second * 40,
			Transport: &http.Transport{
				MaxConnsPerHost:     10,
				MaxIdleConnsPerHost: 2,
			},
		},
		UserAgent:    USER_AGENT,
		ESIRateLimit: rateLimESI,
		RetryLimit:   25,
	}
}

func (c *ESIClient) makeRawHTTPGet(url string) ([]byte, int, string, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, 0, "", errors.Wrap(err, "Failed to buid http request")
	}

	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTP.Do(req)

	if err != nil {
		return nil, 0, "", errors.Wrap(err, "Failed to make request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, 0, "", errors.Wrap(err, "Failed to read response from request")
	}

	etag := res.Header.Get("Etag")

	return body, res.StatusCode, etag, nil
}

func (c *ESIClient) MakeESIGet(url string) (body []byte, status int, etag string, err error) {

	retriesRemain := c.RetryLimit
	for retriesRemain > 1 {
		retriesRemain--

		for c.ESIRateLimit.Value() > 10 {
			time.Sleep(500 * time.Millisecond)
		}

		body, status, etag, err := c.makeRawHTTPGet(url)
		if err != nil {
			if strings.Contains(err.Error(), "too many open files") {
				// This is not going to hurt to keep retrying
				retriesRemain++
			} else {
				// fmt.Printf("ESI GET ERROR - %v\n", err)
			}
			continue
		}
		if !(status == 200) {
			// Increment the counter :(
			c.ESIRateLimit.Inc()
			// fmt.Printf("ESI GET RESPONSE ERROR - %v - %v - %v\n", status, url, string(body))
			time.Sleep(250 * time.Millisecond)
			continue
		}

		return body, status, etag, err
	}

	return nil, status, etag, fmt.Errorf("Max retries exceeded for url: %s; err: %v", url, err)
}

func (client *ESIClient) RequestKillmailFromESI(hp IDHashPair) (km ESIKillmail, status int, etag string, err error) {

	url := fmt.Sprintf(ESI_KILLMAIL_URL_FMT, hp.ID, hp.Hash)
	body, status, etag, err := client.MakeESIGet(url)
	if err != nil || status != http.StatusOK {
		return
	}

	err = json.Unmarshal(body, &km)
	return

}

// Inc increments the counter.
func (c *safeCounter) Inc() {
	c.mux.Lock()
	c.cnt++
	c.mux.Unlock()
}

// Dec decrements the counter.
func (c *safeCounter) Dec() {
	c.mux.Lock()
	if c.cnt > 0 {
		c.cnt--
	}
	c.mux.Unlock()
}

// Value returns the current value of the counter.
func (c *safeCounter) Value() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.cnt
}
