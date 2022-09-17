package main

import (
	"fmt"
	iceportalapi "github.com/craftamap/iceportal-api"
	"github.com/gen2brain/beeep"
	wifiname "github.com/yelinaung/wifi-name"
	"time"
)

type Runner struct {
	stationsNotified []string
}

func main() {
	runner := Runner{}
	for {
		fmt.Println("wifi name ", wifiname.WifiName())
		if wifiname.WifiName() == "WIFI@DB" ||
			wifiname.WifiName() == "WIFIonICE" {
			err := runner.run()
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func (r *Runner) run() error {
	stop, shouldInform, err := shouldInformNearStop(90)
	if err != nil {
		return err
	}

	if shouldInform && !contains(r.stationsNotified, stop.Station.EvaNr) {
		err := beeep.Notify(fmt.Sprintf("Nächster Halt: %s", stop.Station.Name), fmt.Sprintf("In %d Minuten", getSecondsFromNowToStop(stop)/60), "")
		if err != nil {
			return err
		}
		r.stationsNotified = append(r.stationsNotified, stop.Station.EvaNr)
	}

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func shouldInformNearStop(secondsToWarnBeforeStop int) (iceportalapi.Stop, bool, error) {
	trip, err := iceportalapi.FetchTrip()
	if err != nil {
		return iceportalapi.Stop{}, false, err
	}
	for _, stop := range trip.Trip.Stops {
		if !stop.Info.Passed {

			secondsToStop := getSecondsFromNowToStop(stop)
			if secondsToStop <= secondsToWarnBeforeStop {
				return stop, true, nil
			} else {
				return stop, false, nil
			}
		}

	}

	return iceportalapi.Stop{}, false, nil
}

func getSecondsFromNowToStop(stop iceportalapi.Stop) int {
	timeToStop := stop.Timetable.ActualArrivalTime - time.Now().Unix()*1000
	return int(timeToStop / 1000)
}
