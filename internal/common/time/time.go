package time

import (
	"encoding/xml"
	"time"
)

type Time time.Time

func (ct *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const CUSTUM_TIME_LAYOUT string = "2006-01-02T15:04:05.000Z07:00"
	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}

	t, err := time.Parse(CUSTUM_TIME_LAYOUT, s)
	if err != nil {
		return err
	}
	*ct = Time(t)
	return nil
}
