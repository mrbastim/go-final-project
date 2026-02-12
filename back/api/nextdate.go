package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	date, err := time.ParseInLocation(dateLayout, dstart, time.Local)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("repeat rule is empty")
	}

	baseYear, baseMonth, baseDay := date.Date()

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid daily repeat format")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", errors.New("invalid daily interval")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateLayout), nil
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid yearly repeat format")
		}
		for {
			baseYear++
			date = time.Date(baseYear, baseMonth, baseDay, 0, 0, 0, 0, time.Local)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateLayout), nil
	case "w":
		if len(parts) != 2 {
			return "", errors.New("invalid weekly repeat format")
		}
		weekdays, err := parseWeekdays(parts[1])
		if err != nil {
			return "", err
		}
		date = date.AddDate(0, 0, 1)
		for {
			if afterNow(date, now) && weekdays[weekdayNumber(date)] {
				return date.Format(dateLayout), nil
			}
			date = date.AddDate(0, 0, 1)
		}
	case "m":
		if len(parts) != 2 && len(parts) != 3 {
			return "", errors.New("invalid monthly repeat format")
		}
		daysRule, err := parseMonthDays(parts[1])
		if err != nil {
			return "", err
		}
		months := allMonths()
		if len(parts) == 3 {
			months, err = parseMonths(parts[2])
			if err != nil {
				return "", err
			}
		}
		date = date.AddDate(0, 0, 1)
		for {
			if afterNow(date, now) && months[int(date.Month())] && matchMonthDay(date, daysRule) {
				return date.Format(dateLayout), nil
			}
			date = date.AddDate(0, 0, 1)
		}
	default:
		return "", errors.New("unsupported repeat format")
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(w, map[string]string{"error": "Method not allowed"})
		return
	}

	nowStr := strings.TrimSpace(r.FormValue("now"))
	dateStr := strings.TrimSpace(r.FormValue("date"))
	repeat := strings.TrimSpace(r.FormValue("repeat"))

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		parsed, err := time.ParseInLocation(dateLayout, nowStr, time.Local)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsed
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(next))
}

type monthDayRule struct {
	days        [32]bool
	last        bool
	penultimate bool
}

func dayOnly(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func afterNow(date time.Time, now time.Time) bool {
	return dayOnly(date).After(dayOnly(now))
}

func beforeNow(date time.Time, now time.Time) bool {
	return dayOnly(date).Before(dayOnly(now))
}

func weekdayNumber(date time.Time) int {
	weekday := int(date.Weekday())
	if weekday == 0 {
		return 7
	}
	return weekday
}

func parseWeekdays(value string) ([8]bool, error) {
	var days [8]bool
	items := strings.Split(value, ",")
	if len(items) == 0 {
		return days, errors.New("weekdays list is empty")
	}
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			return days, errors.New("weekdays list is empty")
		}
		v, err := strconv.Atoi(item)
		if err != nil || v < 1 || v > 7 {
			return days, fmt.Errorf("invalid weekday: %s", item)
		}
		days[v] = true
	}
	return days, nil
}

func parseMonthDays(value string) (monthDayRule, error) {
	var rule monthDayRule
	items := strings.Split(value, ",")
	if len(items) == 0 {
		return rule, errors.New("month days list is empty")
	}
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			return rule, errors.New("month days list is empty")
		}
		v, err := strconv.Atoi(item)
		if err != nil {
			return rule, fmt.Errorf("invalid month day: %s", item)
		}
		switch {
		case v == -1:
			rule.last = true
		case v == -2:
			rule.penultimate = true
		case v >= 1 && v <= 31:
			rule.days[v] = true
		default:
			return rule, fmt.Errorf("invalid month day: %s", item)
		}
	}
	return rule, nil
}

func parseMonths(value string) ([13]bool, error) {
	var months [13]bool
	items := strings.Split(value, ",")
	if len(items) == 0 {
		return months, errors.New("months list is empty")
	}
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			return months, errors.New("months list is empty")
		}
		v, err := strconv.Atoi(item)
		if err != nil || v < 1 || v > 12 {
			return months, fmt.Errorf("invalid month: %s", item)
		}
		months[v] = true
	}
	return months, nil
}

func allMonths() [13]bool {
	var months [13]bool
	for i := 1; i <= 12; i++ {
		months[i] = true
	}
	return months
}

func lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
}

func matchMonthDay(date time.Time, rule monthDayRule) bool {
	day := date.Day()
	if rule.days[day] {
		return true
	}
	last := lastDayOfMonth(date.Year(), date.Month())
	if rule.last && day == last {
		return true
	}
	if rule.penultimate && day == last-1 {
		return true
	}
	return false
}
