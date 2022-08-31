package checkstatus

import (
	"encoding/xml"
	"time"
)

type CheckStatusResponseEnv struct {
	XMLName                 xml.Name `xml:"Envelope"`
	CheckStatusResponseBody CheckStatusResponseBody
}

type CheckStatusResponseBody struct {
	XMLName             xml.Name `xml:"Body"`
	CheckStatusResponse CheckStatusResponse
}

type CheckStatusResponse struct {
	XMLName                   xml.Name `xml:"CheckStatusResponse"`
	CheckStatusResponseAnswer CheckStatusResponseAnswer
}

type CheckStatusResponseAnswer struct {
	XMLName            xml.Name  `xml:"Answer"`
	Status             bool      `xml:"Status"`
	ServiceStartedTime time.Time `xml:"ServiceStartedTime"`
}
