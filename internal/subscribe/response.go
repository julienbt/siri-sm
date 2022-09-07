package subscribe

import (
	"encoding/xml"

	siri_time "github.com/julienbt/siri-sm/internal/common/time"
)

type SubscribeEnv struct {
	XMLName           xml.Name          `xml:"Envelope"`
	SubscribeResponse SubscribeResponse `xml:"Body>SubscribeResponse"`
}

type SubscribeResponse struct {
	XMLName        xml.Name         `xml:"SubscribeResponse"`
	ResponseStatus []ResponseStatus `xml:"Answer>ResponseStatus"`
}

type ResponseStatus struct {
	XMLName           xml.Name       `xml:"ResponseStatus"`
	ResponseTimestamp siri_time.Time `xml:"ResponseTimestamp"`
	RequestMessageRef string         `xml:"RequestMessageRef"`
	SubscriberRef     string         `xml:"SubscriberRef"`
	SubscriptionRef   string         `xml:"SubscriptionRef"`
	Status            bool           `xml:"Status"`
	ValidUntil        siri_time.Time `xml:"ValidUntil"`
}
