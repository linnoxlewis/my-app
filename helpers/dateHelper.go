package helpers

import "time"

func DaysBetween(from time.Time, to time.Time) []time.Time {
	days := make([]time.Time, 0)

	from = toDay(from)
	to = toDay(to)

	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days
}

func Date(year, month, day int) time.Time {
	return time.Date(year,
		time.Month(month),
		day,
		0,
		0,
		0,
		0,
		time.UTC)
}

func toDay(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(),
		timestamp.Month(),
		timestamp.Day(),
		0,
		0,
		0,
		0,
		time.UTC)
}
