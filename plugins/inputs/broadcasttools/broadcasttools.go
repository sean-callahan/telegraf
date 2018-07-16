package broadcasttools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type parser func(src map[string]interface{}, index int) interface{}

var (
	regexpTemp   = regexp.MustCompile(`^T1(\d+)$`)
	regexpMeter  = regexp.MustCompile(`^M1(\d+)$`)
	regexpVC     = regexp.MustCompile(`^VCLabel(\d+)$`)
	regexpStatus = regexp.MustCompile(`^S1(\d+)$`)
	regexpRelay  = regexp.MustCompile(`^R2(\d+)$`)

	parsers = map[*regexp.Regexp]parser{
		regexpTemp: func(src map[string]interface{}, index int) interface{} {
			t := src[fmt.Sprintf("TempValue%02d", index)].(string)
			t = strings.TrimSuffix(t, " *F")
			v, _ := strconv.Atoi(t)
			return v
		},
		regexpMeter: func(src map[string]interface{}, index int) interface{} {
			return src[fmt.Sprintf("MeterValue%02d", index)]
		},
		regexpVC: func(src map[string]interface{}, index int) interface{} {
			return src[fmt.Sprintf("VCValue%02d", index)]
		},
		regexpStatus: func(src map[string]interface{}, index int) interface{} {
			return src[fmt.Sprintf("StatusIndicator%02d", index)]
		},
		regexpRelay: func(src map[string]interface{}, index int) interface{} {
			return src[fmt.Sprintf("RelayIndicator%02d", index)]
		},
	}
)

type BroadcastTools struct {
	URLs     []string
	User     string
	Password string

	devices     []Device
	initialized bool
}

type Device interface {
	Dial() error
	Close() error
	Gather(acc telegraf.Accumulator) error
}

const sampleConfig = `
  ## An array of URLs to gather stats from. i.e.,
  ##   http://example.com:3000
  urls = ["http://localhost:1776"]
  ## Username
  user = "admin"
  ## Password
  password = "password"
`

func (bt *BroadcastTools) init() error {
	if bt.initialized {
		return nil
	}

	for _, u := range bt.URLs {
		base, err := url.Parse(u)
		if err != nil {
			return err
		}

		d := &device{
			bt:   bt,
			base: base,
			c:    &http.Client{},
		}
		if err := d.Dial(); err != nil {
			return err
		}

		bt.devices = append(bt.devices, d)
	}

	bt.initialized = true
	return nil
}

func (bt *BroadcastTools) SampleConfig() string {
	return sampleConfig
}

func (bt *BroadcastTools) Description() string {
	return "Read metrics from one or many Broadcast Tools devices"
}

func (bt *BroadcastTools) Gather(acc telegraf.Accumulator) error {
	if !bt.initialized {
		if err := bt.init(); err != nil {
			return err
		}
	}

	var wg sync.WaitGroup

	for _, device := range bt.devices {
		wg.Add(1)

		go func(d Device, a telegraf.Accumulator) {
			defer wg.Done()
			if err := d.Gather(a); err != nil {
				acc.AddError(err)
			}
		}(device, acc)
	}

	wg.Wait()
	return nil
}

type device struct {
	bt   *BroadcastTools
	base *url.URL
	c    *http.Client
	ck   *http.Cookie
}

func (d device) send(method string, path string, data io.Reader, sendCookie bool) (*http.Response, error) {
	u := *d.base
	u.RawPath = path

	r, err := http.NewRequest(method, u.String(), data)
	if err != nil {
		return nil, err
	}
	if sendCookie {
		r.AddCookie(d.ck)
	}
	if method == http.MethodPost {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return d.c.Do(r)
}

func (d *device) Dial() error {
	if d.ck != nil {
		return errors.New("already logged in")
	}

	v := url.Values{}
	v.Set("AccessVal", "")
	v.Set("LoginUser", d.bt.User)
	v.Set("LoginPass", d.bt.Password)

	r, err := d.send(http.MethodPost, "/cgi-bin/postauth.cgi", strings.NewReader(v.Encode()), false)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New("authentication failed")
	}

	cks := r.Cookies()
	if len(cks) < 1 {
		return errors.New("no cookies returned")
	}
	d.ck = cks[0]

	return nil
}

func (d device) Close() error {
	v := url.Values{}
	v.Set("Logout", "1")

	_, err := d.send(http.MethodPost, "/cgi-bin/postlogout.cgi", strings.NewReader(v.Encode()), true)
	if err != nil {
		return nil // ignore logout errors
	}

	d.ck = nil
	d.c = nil
	return nil
}

func keyify(s string) string {
	s = strings.ToLower(s)
	return strings.Replace(s, " ", "_", -1)
}

func (d *device) Gather(acc telegraf.Accumulator) error {
	r, err := d.send(http.MethodGet, "/cgi-bin/getexchanger_monitor.cgi", nil, true)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status %d; got %d", http.StatusOK, r.StatusCode)
	}
	defer r.Body.Close()

	var data map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	values := data["values"].(map[string]interface{})

	fields := make(map[string]interface{})

	for key, name := range values {
		for reg, parser := range parsers {
			matches := reg.FindStringSubmatch(key)
			if len(matches) < 2 {
				continue
			}
			index, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}
			value := parser(values, index)
			fields[keyify(name.(string))] = value
		}
	}

	acc.AddFields("broadcasttools", fields, nil)

	return nil
}

func init() {
	inputs.Add("broadcasttools", func() telegraf.Input {
		bt := &BroadcastTools{}
		return bt
	})
}
