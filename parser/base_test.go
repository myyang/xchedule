package parser

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/myyang/xchedule"
	"github.com/spf13/viper"
)

var singleEvent = []byte(`
title: Go to Japan
timezone: Asia/Taipei
time: 2016-01-20 10:00 ~ 2016-01-21 11:03PM
locations: Japan
members:
    - Person A
    - Slave
    - Yang man
    - Ohikiwa
schedule:
    - event1
    - event2
    - event3
    - event4
alerts:
    time:
        - 2016.01.19 10:00:00
notes:
    - "OK"
event1:
    title: "event1"
    time: 2016-01-20 10:00
event2:
    title: "event2"
    time: 2016-01-20 12:00
event3:
    title: "event3"
    time: 2016-01-20 15:00
`)

func TestNewEvent(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Can't get current path")
	}
	viper.Set("workspace", pwd)
	v := viper.New()
	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewBuffer(singleEvent))
	e := NewEvent(v, true)

	tz, err := time.LoadLocation("Asia/Taipei")
	alert, err := time.ParseInLocation("2006-01-02 15:04", "2016-01-19 10:00", tz)
	start, err := time.ParseInLocation("2006-01-02 15:04", "2016-01-20 10:00", tz)
	end, err := time.ParseInLocation("2006-01-02 15:04", "2016-01-21 23:03", tz)
	tz = time.Local
	event1Time, err := time.ParseInLocation("2006-01-02 15:04", "2016-01-20 10:00", tz)
	event2Time, err := time.ParseInLocation("2006-01-02 15:04", "2016-01-20 12:00", tz)
	event3Time, err := time.ParseInLocation("2006-01-02 15:04", "2016-01-20 15:00", tz)
	event4Time, err := time.ParseInLocation("2006-01-02 15:04", "2016-01-20 19:00", tz)
	if err != nil {
		log.Fatal(err)
	}

	expected := xchedule.Event{
		Title: "Go to Japan",
		Time:  xchedule.Time{Start: start, End: end, Period: true},
		Members: []xchedule.Member{
			xchedule.Member{Name: "Person A"},
			xchedule.Member{Name: "Slave"},
			xchedule.Member{Name: "Yang man"},
			xchedule.Member{Name: "Ohikiwa"},
		},
		Notes:     []string{"OK"},
		Locations: []xchedule.Location{},
		Schedule: []xchedule.Event{
			xchedule.Event{
				Title: "event1", Time: xchedule.Time{Start: event1Time, Period: false},
				Members: []xchedule.Member{}, Locations: []xchedule.Location{},
				Schedule: []xchedule.Event{}, Notes: []string{}},
			xchedule.Event{
				Title: "event2", Time: xchedule.Time{Start: event2Time, Period: false},
				Members: []xchedule.Member{}, Locations: []xchedule.Location{},
				Schedule: []xchedule.Event{}, Notes: []string{}},
			xchedule.Event{
				Title: "event3", Time: xchedule.Time{Start: event3Time, Period: false},
				Members: []xchedule.Member{}, Locations: []xchedule.Location{},
				Schedule: []xchedule.Event{}, Notes: []string{}},
			xchedule.Event{
				Title: "event4", Time: xchedule.Time{Start: event4Time, Period: false},
				Members: []xchedule.Member{}, Locations: []xchedule.Location{},
				Schedule: []xchedule.Event{}, Notes: []string{}},
		},
		Alert: xchedule.Alert{Times: []time.Time{alert}},
	}
	expected.SetRoot(true)

	if !reflect.DeepEqual(expected, e) {
		_, file, line, _ := runtime.Caller(0)
		t.Fatalf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, expected, e)
	}
}
