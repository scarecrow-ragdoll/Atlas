// FILE: apps/api/internal/atlas/models/date.go
// VERSION: 1.0.0
// START_MODULE_CONTRACT
//   PURPOSE: Define strict calendar-date handling for Atlas GraphQL Date scalar values.
//   SCOPE: Parse, format, marshal, and unmarshal YYYY-MM-DD dates without timezone conversion; excludes timestamp handling.
//   DEPENDS: github.com/99designs/gqlgen/graphql.
//   LINKS: M-API / V-M-API / WAVE-03.
//   ROLE: TYPES
//   MAP_MODE: EXPORTS
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   Date - Calendar date wrapper for GraphQL and service input.
//   MustDate - Test/helper constructor that panics on invalid date strings.
//   ParseDate - Strict YYYY-MM-DD parser for Date values.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.0 - Added Date scalar model for WAVE-03.
// END_CHANGE_SUMMARY

package models

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

const dateLayout = "2006-01-02"

type Date struct {
	value time.Time
}

func MustDate(raw string) Date {
	d, err := ParseDate(raw)
	if err != nil {
		panic(err)
	}
	return d
}

func ParseDate(raw string) (Date, error) {
	t, err := time.Parse(dateLayout, raw)
	if err != nil || t.Format(dateLayout) != raw {
		return Date{}, fmt.Errorf("invalid date %q: must use YYYY-MM-DD", raw)
	}
	return Date{value: t}, nil
}

func (d Date) String() string {
	if d.value.IsZero() {
		return ""
	}
	return d.value.Format(dateLayout)
}

func (d Date) Time() time.Time {
	return d.value
}

func (d Date) MarshalGQL() graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, `"`+d.String()+`"`)
	})
}

func (d *Date) UnmarshalGQL(v any) error {
	raw, ok := v.(string)
	if !ok {
		return fmt.Errorf("date must be a string")
	}
	parsed, err := ParseDate(raw)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}
