// Copyright (c) 2014-2020 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package timeutil

import (
	"fmt"
	"time"
)

// start-of-day
func sod(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// Human turns the time into a relative expression of time meant for human
// consumption.
// Human(t)  --> "today at 07:47"
func Human(then time.Time) string {
	return humanTimeSince(then.Local(), time.Now().Local(), 60)
}

func delta(then, now time.Time) int {
	if then.After(now) {
		return -delta(now, then)
	}

	then = sod(then)
	now = sod(now)

	n := int(then.Sub(now).Hours() / 24)
	now = now.AddDate(0, 0, n)
	for then.Before(now) {
		then = then.AddDate(0, 0, 1)
		n--
	}
	return n
}

func humanTimeSince(then, now time.Time, cutoffDays int) string {
	d := delta(then, now)
	switch {
	case d < -1 && d >= -cutoffDays:
		return fmt.Sprintf(then.Format("%d days ago, at 15:04 MST"), -d)
	case d == -1:
		return then.Format("yesterday at 15:04 MST")
	case d == 0:
		return then.Format("today at 15:04 MST")
	case d == 1:
		return then.Format("tomorrow at 15:04 MST")
	case d > 1 && d <= cutoffDays:
		return fmt.Sprintf(then.Format("in %d days, at 15:04 MST"), d)
	default:
		return then.Format("2006-01-02")
	}
}
