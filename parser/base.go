// Package parser provides parsers support varies type
package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/myyang/xchedule"
	"github.com/spf13/viper"
)

// Parser interface
type Parser interface{}

// NewEvent initialize a new xchedule event
func NewEvent(ev *viper.Viper, isRoot bool) xchedule.Event {
	ev.Set("isRoot", isRoot)
	e := xchedule.Event{
		Title:     getTitle(ev),
		Time:      getTime(ev),
		Locations: getLocations(ev),
		Members:   getMembers(ev),
		Alert:     getAlerts(ev),
		Notes:     getNotes(ev),
		Schedule:  getSchedule(ev),
	}
	e.SetRoot(isRoot)
	return e
}

func getTitle(ev *viper.Viper) string {
	if s := ev.GetString("title"); s != "" {
		return s
	}

	log.Fatalf("No title for event")
	return ""
}

/*
stolen/borrow from github.com/spf13/cast/caste.go
StringToDate(s string) (time.Time, error)
*/
func stringToTime(s string, ev *viper.Viper) (t time.Time, err error) {
	tz := getTimezone(ev)
	s = strings.Trim(s, " ")
	timeFormat := []string{
		"2006/01/02 15:04",
		"2006/01/02 3:04PM",
		"2006-01-02 15:04",
		"2006-01-02 3:04PM",
		"2006.01.02 15:04",
		"2006.01.02 3:04PM",
		// to seconds if needed
		"2006/01/02 15:04:05",
		"2006/01/02 3:04:05PM",
		"2006-01-02 15:04:05",
		"2006-01-02 3:04:05PM",
		"2006.01.02 15:04:05",
		"2006.01.02 3:04:05PM",
	}

	for _, format := range timeFormat {
		if t, err = time.ParseInLocation(format, s, tz); err == nil {
			return t, err
		}
	}
	return t, fmt.Errorf("Can't parse <%v>, valid format YYYY/mm/dd HH:MM:SS (digits are concerned).\n", s)
}

func getTime(ev *viper.Viper) xchedule.Time {
	period := strings.Split(ev.GetString("time"), "~")
	log.Printf("Get period: %#v(len:%v)\n", period, len(period))
	if 1 > len(period) || 2 < len(period) || (len(period) == 1 && period[0] == "") {
		log.Fatal("Wrong time present, should be single 'time.Time' or 'time.Time' ~ time.Time'\n")
	}

	start, err := stringToTime(period[0], ev)
	if err != nil {
		log.Fatalf("Error while parsing event start time, err: <%v>\n", err)
	}
	t := xchedule.Time{Start: start, Period: false}
	if len(period) == 1 {
		return t
	}
	end, err := stringToTime(period[1], ev)
	if err != nil {
		log.Fatalf("Error while parsing event end time, err: <%v>\n", err)
	}
	t.End = end
	t.Period = true
	return t
}

func getTimezone(ev *viper.Viper) *time.Location {
	// TODO: supports country/location name, e.g: Taiwan or Taipei

	var tz *time.Location
	if tzname := ev.GetString("timezone"); tzname != "" {
		ttz, err := time.LoadLocation(tzname)
		if err != nil {
			log.Fatalf("Invalid timezone name: %v (%v)\n", tzname, err)
		}
		tz = ttz
	} else {
		tz = time.Local
	}

	return tz
}

func getNotes(ev *viper.Viper) []string {
	// TODO: supports possible format ex: images, url and etc.
	notes := ev.GetStringSlice("notes")
	if len(notes) == 0 {
		notes = []string{}
	}
	return notes
}

func getMembers(ev *viper.Viper) []xchedule.Member {
	// TODO: add contact info like phone or email
	rawMembers := ev.GetStringSlice("members")
	members := []xchedule.Member{}
	for _, v := range rawMembers {
		members = append(members, xchedule.Member{Name: v})
	}
	return members
}

func getAlerts(ev *viper.Viper) xchedule.Alert {
	rawAlerts := ev.GetStringMapStringSlice("alerts")
	alerts := xchedule.Alert{}
	for _, v := range rawAlerts["time"] {
		alertTime, err := stringToTime(v, ev)
		if err != nil {
			log.Fatalf("Error while parsing alert time, err: <%v>\n", err)
		}
		alerts.Times = append(alerts.Times, alertTime)
	}
	return alerts
}

func getLocations(ev *viper.Viper) []xchedule.Location {
	return []xchedule.Location{}
}

func getSchedule(ev *viper.Viper) []xchedule.Event {
	rawEvents := ev.GetStringSlice("schedule")
	events := []xchedule.Event{}
	//ev.Debug()
	for _, e := range rawEvents {
		if raw := ev.GetStringMap(e); len(raw) >= 1 {
			log.Printf("Got raw: %v\n", raw)
			subEv := viper.New()
			subEv.SetConfigType(ev.GetString("configType"))
			for k, v := range raw {
				subEv.Set(k, v)
			}
			events = append(events, NewEvent(subEv, false))
		} else if subEv, err := findEventConfig(e); err == nil {
			events = append(events, NewEvent(subEv, false))
		} else {
			log.Fatalf("Can't find event: %v\n", e)
		}
	}
	return events
}

func findEventConfig(prefix string) (*viper.Viper, error) {
	files, err := ioutil.ReadDir(viper.GetString("workspace"))
	if err != nil {
		log.Fatalf("Fail to read event under dir: %v, error: %v\n", viper.GetString("workspace"), err)
	}
	log.Printf("Read from dir:%v, files: %v\n", viper.GetString("workspace"), files)

	for _, f := range files {
		if strings.HasPrefix(f.Name(), prefix) {
			ext := filepath.Ext(f.Name())[1:]
			ev := viper.New()
			ev.SetConfigType(ext)
			ev.Set("configType", ext)
			log.Printf("Read from file: %v with config type: %v\n", f.Name(), ext)
			ev.SetConfigFile(f.Name())
			if err := ev.ReadInConfig(); err != nil {
				log.Fatalln("Fail to read event config: err: ", err)
			}
			ev.Debug()
			return ev, nil
		}
	}

	return viper.New(), fmt.Errorf("Can't find event config file")
}
