package ectoplasma

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/nulls"
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
		Attackers     []ESIAttacker `json:"attackers"`
		KillmailID    int32         `json:"killmail_id"`
		KillmailTime  time.Time     `json:"killmail_time"`
		MoonID        nulls.Int32   `json:"moon_id"`
		SolarSystemID int32         `json:"solar_system_id"`
		Victim        ESIVictim     `json:"victim"`
		WarID         nulls.Int32   `json:"war_id"`
	}

	ESIAttacker struct {
		AllianceID     nulls.Int32 `json:"alliance_id,omitempty"`
		CorporationID  nulls.Int32 `json:"corporation_id"`
		CharacterID    nulls.Int32 `json:"character_id"`
		DamageDone     int32       `json:"damage_done"`
		FinalBlow      bool        `json:"final_blow"`
		SecurityStatus float64     `json:"security_status"`
		ShipTypeID     nulls.Int32 `json:"ship_type_id"`
		WeaponTypeID   nulls.Int32 `json:"weapon_type_id"`
	}

	ESIVictim struct {
		AllianceID    nulls.Int32 `json:"alliance_id,omitempty"`
		CorporationID nulls.Int32 `json:"corporation_id"`
		CharacterID   nulls.Int32 `json:"character_id"`
		FactionID     nulls.Int32 `json:"faction_id,omitempty"`
		DamageTaken   int32       `json:"damage_taken"`
		Items         []ESIItem   `json:"items"`
		Position      ESIPosition `json:"position"`
		ShipTypeID    int32       `json:"ship_type_id"`
	}

	ESIItem struct {
		Flag              int32       `json:"flag"`
		ItemTypeID        int32       `json:"item_type_id"`
		QuantityDropped   nulls.Int64 `json:"quantity_dropped,omitempty"`
		QuantityDestroyed nulls.Int64 `json:"quantity_destroyed,omitempty"`
		Singleton         int32       `json:"singleton"`
	}

	ESIPosition struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}
)

func NewESIClient() (*ESIClient) {
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

	return nil, status, etag, fmt.Errorf("Max retries exceeded for url: ; err: %v", url, err)
}

func (client *ESIClient) RequestKillmailFromESI(hp IDHashPair) (km ESIKillmail, status int, etag string, err error) {

	url := fmt.Sprintf(ESI_KILLMAIL_URL_FMT, hp.ID, hp.Hash)
	body, status, etag, err := client.MakeESIGet(url)
	if err != nil || status != http.StatusOK{
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
