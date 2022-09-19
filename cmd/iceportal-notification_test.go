package main

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func Test_shouldInformNearStop(t *testing.T) {
	t.Run("Station ahead, should notify", func(f *testing.T) {

		data := getMockData(1)

		httpmock.Activate()
		httpmock.RegisterResponder("GET", "https://iceportal.de/api1/rs/tripInfo/trip", func(*http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, string(data))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		stop, shouldInform, err := shouldInformNearStop(60)
		assert.NoError(t, err)
		assert.True(t, shouldInform)
		assert.Equal(t, "Town A", stop.Station.Name)
	},
	)

	t.Run("Station ahead, should not notify", func(f *testing.T) {

		data := getMockData(5)

		httpmock.Activate()
		httpmock.RegisterResponder("GET", "https://iceportal.de/api1/rs/tripInfo/trip", func(*http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, string(data))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		stop, shouldInform, err := shouldInformNearStop(60)
		assert.NoError(t, err)
		assert.False(t, shouldInform)
		assert.Equal(t, "Town A", stop.Station.Name)
	},
	)

	t.Run("Station in past should notify", func(f *testing.T) {

		data, err := os.ReadFile("../testdata/hef.json")
		if err != nil {
			panic(err)
		}

		httpmock.Activate()
		httpmock.RegisterResponder("GET", "https://iceportal.de/api1/rs/tripInfo/trip", func(*http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, string(data))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		stop, shouldInform, err := shouldInformNearStop(60)
		assert.NoError(t, err)
		assert.True(t, shouldInform)
		assert.Equal(t, "Bad Hersfeld", stop.Station.Name)
	},
	)

}

func getMockData(stopMinutesInTheFuture int) string {
	actualArrivalTime := time.Now().Unix()*1000 + 60*int64(stopMinutesInTheFuture)*1000
	data := fmt.Sprintf(`
{
  "trip": {
    "stops": [
      {
        "station": { "name": "Town A", "evaNr": "123" },
        "timetable": { "scheduledArrivalTime": 1663480864, "actualArrivalTime": %d },
        "info": { "passed": false }
      }
    ]
  }
}

`, actualArrivalTime)
	return data
}

func Test_run(t *testing.T) {
	t.Run("Get Notification for non notified", func(t *testing.T) {

		data := getMockData(1)

		httpmock.Activate()
		httpmock.RegisterResponder("GET", "https://iceportal.de/api1/rs/tripInfo/trip", func(*http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, string(data))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

		r := Runner{}
		r.stationsNotified = append(r.stationsNotified, "1663480000 Town B")
		err := r.run()
		assert.NoError(t, err)
		assert.Len(t, r.stationsNotified, 2)
		assert.Equal(t, "1663480864 Town A", r.stationsNotified[1])

	})
}
