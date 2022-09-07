package getstopmonitoring

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/julienbt/siri-sm/internal/common/directionname"
	siri_time "github.com/julienbt/siri-sm/internal/common/time"
)

type GetStopMonitoringEnv struct {
	XMLName                xml.Name               `xml:"Envelope"`
	StopMonitoringDelivery StopMonitoringDelivery `xml:"Body>GetStopMonitoringResponse>Answer>StopMonitoringDelivery"`
}

type StopMonitoringDelivery struct {
	XMLName                         xml.Name                         `xml:"StopMonitoringDelivery"`
	MonitoringRef                   StopPointRef                     `xml:"MonitoringRef"`
	MonitoredStopVisits             []MonitoredStopVisit             `xml:"MonitoredStopVisit"`
	MonitoredStopVisitCancellations []MonitoredStopVisitCancellation `xml:"MonitoredStopVisitCancellation"`
}

type StopPointRef string

func (mr *StopPointRef) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var innerText string
	err := d.DecodeElement(&innerText, &start)
	if err != nil {
		return err
	}

	splittedInnerText := strings.Split(innerText, ":")
	const EXPECTED_NUM_OF_PARTS int = 5
	if len(splittedInnerText) != EXPECTED_NUM_OF_PARTS {
		return fmt.Errorf("the `MonitoringRef` is not well formatted: %s", innerText)
	}
	if splittedInnerText[1] != "StopPoint" {
		return fmt.Errorf("the `MonitoringRef` is not well formatted: %s", innerText)
	}
	*mr = StopPointRef(splittedInnerText[3])
	return nil
}

type MonitoredStopVisit struct {
	XMLName                 xml.Name                `xml:"MonitoredStopVisit"`
	ItemIdentifier          string                  `xml:"ItemIdentifier"`
	MonitoringRef           StopPointRef            `xml:"MonitoringRef"`
	MonitoredVehicleJourney MonitoredVehicleJourney `xml:"MonitoredVehicleJourney"`
}

type MonitoredVehicleJourney struct {
	XMLName         xml.Name                    `xml:"MonitoredVehicleJourney"`
	LineRef         LineRef                     `xml:"LineRef"`
	DirectionName   directionname.DirectionName `xml:"DirectionName"`
	DestinationRef  StopPointRef                `xml:"DestinationRef"`
	DestinationName string                      `xml:"DestinationName"`
	MonitoredCall   MonitoredCall               `xml:"MonitoredCall"`
}

type MonitoredStopVisitCancellation struct {
	XMLName       xml.Name     `xml:"MonitoredStopVisitCancellation"`
	ItemRef       string       `xml:"ItemRef"`
	MonitoringRef StopPointRef `xml:"MonitoringRef"`
}

type LineRef string

func (lr *LineRef) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var innerText string
	err := d.DecodeElement(&innerText, &start)
	if err != nil {
		return err
	}

	splittedInnerText := strings.Split(innerText, ":")
	const EXPECTED_NUM_OF_PARTS int = 5
	if len(splittedInnerText) != EXPECTED_NUM_OF_PARTS {
		return fmt.Errorf("the `LineRef` is not well formatted: %s", innerText)
	}
	if splittedInnerText[1] != "Line" {
		return fmt.Errorf("the `LineRef` is not well formatted: %s", innerText)
	}
	*lr = LineRef(splittedInnerText[3])
	return nil
}

type MonitoredCall struct {
	XMLName               xml.Name       `xml:"MonitoredCall"`
	StopPointRef          StopPointRef   `xml:"StopPointRef"`
	AimedDepartureTime    siri_time.Time `xml:"AimedDepartureTime"`
	ExpectedDepartureTime siri_time.Time `xml:"ExpectedDepartureTime"`
}
