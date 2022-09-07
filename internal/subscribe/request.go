package subscribe

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/julienbt/siri-sm/internal/config"
	"github.com/julienbt/siri-sm/internal/siri"
	"github.com/sirupsen/logrus"
)

var LOCATION_NAME = "Europ/Paris"
var LOCATION, _ = time.LoadLocation(LOCATION_NAME)

var INITIAL_TERMINATION_TIME time.Time
var REQUEST_TIMESTAMP time.Time

const SUBSCRIBER_REF string = "KISIO2"
const STOP_VISIT_TYPES string = "departures"
const MINIMUM_STOP_VISITS_PER_LINE int = 2
const PREVIEW_INTERVAL_DURATION time.Duration = 2 * time.Hour
const INCREMENTAL_UPDATES bool = true
const CHANGE_BEFORE_UPDATES_DURATION time.Duration = 2 * time.Second

const IDENTIFIER_TIME_LAYOUT string = "20060102_150405"

var STOP_POINT_IDS_LILLE_BUS = []string{
	"CAS001",
	// "CAS002",
	// "CAT001",
	// "CAT002",
	// "CAU001",
	// "CAU002",
	// "CAV001",
	// "CAV002",
	// "CAW001",
	// "CAW002",
	// "CBA011",
	// "CBA012",
	// "CBE001",
	// "CBE002",
	// "CBF001",
	// "CBF002",
	// "CBG001",
	// "CBG002",
	// "CBO002",
	// "CBO004",
	// "CCD001",
	// "CCD002",
	// "CCE001",
	// "CCE002",
	// "CCH001",
	// "CCH002",
	// "CDE001",
	// "CDE002",
	// "CDO001",
	// "CDO002",
	// "CDP001",
	// "CDP002",
	// "CDT001",
	// "CDT002",
	// "CED001",
	// "CED002",
	// "CED002",
	// "CEH001",
	// "CEN001",
	// "CEN001",
	// "CEO001",
	// "CEO002",
	// "CER001",
}

var STOP_POINT_IDS_AMIENS = []string{
	"RAMPO1",
}

func Subscribe(cfg config.ConfigSubscribe, logger *logrus.Entry, requestTimestamp *time.Time) (SubscribeRequestInfoResult, string, []byte, error) {
	var remoteErrorLoc = "Subscribe remote error"
	req := SubscribeRequestInfo{}
	err := req.populate(&cfg, requestTimestamp, requestTimestamp)
	if err != nil {
		return SubscribeRequestInfoResult{},
			"",
			nil,
			fmt.Errorf("error Subscibe request initialization: %v", err)
	}

	httpReq, htmlReqBody, err := req.generateHttpSoapReq()
	// Check/parse the HTTP Response
	if err != nil {
		if err != nil {
			return SubscribeRequestInfoResult{},
				"",
				nil,
				fmt.Errorf("error in building SOAP Subscribe request: %s", err)
		}
	}

	// Send HTTP request and receive the response
	resp, err := siri.SoapCall(httpReq)
	if err != nil {
		return SubscribeRequestInfoResult{},
			htmlReqBody,
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("call error: %s", err)}
	}

	// Get the HTTP response body
	defer resp.Body.Close()
	htmlRespBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SubscribeRequestInfoResult{},
			htmlReqBody,
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("unreadable response body: %s", err)}
	}

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return SubscribeRequestInfoResult{},
			htmlReqBody,
			htmlRespBody,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("bad http-response status: %s", resp.Status)}
	}

	return SubscribeRequestInfoResult{}, htmlReqBody, htmlRespBody, nil
}

type SubscribeRequestInfo struct {
	SupplierAddress   url.URL
	RequestTimestamp  time.Time
	SubscriberRef     string
	ConsumerAddress   string
	SubscribeRequests []SubscribeRequest
}

type SubscribeRequestInfoResult struct{}

func (req *SubscribeRequestInfo) populate(
	cfg *config.ConfigSubscribe,
	requestTimestamp *time.Time,
	initialTerminationTime *time.Time) error {
	supplierAddressUrl, err := url.Parse(cfg.SupplierAddress)
	if err != nil {
		return fmt.Errorf("error the supplier address is not a valid URL: %s", cfg.SupplierAddress)
	}
	req.SupplierAddress = *supplierAddressUrl
	req.RequestTimestamp = *requestTimestamp
	req.SubscriberRef = cfg.SubscriberRef
	req.ConsumerAddress = cfg.ConsumerAddress
	req.SubscribeRequests = initSubscribeRequests(cfg, requestTimestamp, initialTerminationTime)
	return nil
}

func (req *SubscribeRequestInfo) generateHttpSoapReq() (*http.Request, string, error) {
	tmpl, err := template.ParseFiles("./template/subscription-request.tmpl")
	if err != nil {
		return nil, "", fmt.Errorf("error parsing template: %s", err)
	}

	httpReqBodyBuffer := &bytes.Buffer{}
	err = tmpl.Execute(httpReqBodyBuffer, req)
	if err != nil {
		return nil, "", fmt.Errorf("error building template: %s", err)
	}
	httpReqBody := httpReqBodyBuffer.String()
	httpReq, err := http.NewRequest(http.MethodPost, req.SupplierAddress.String(), strings.NewReader(httpReqBody))
	headers := http.Header{
		"Content-Type": []string{"text/xml; charset=utf-8"},
		"SOAPAction":   []string{"Subscribe"},
	}
	httpReq.Header = headers // better than .Header.Set to preserve case (for "SOAPAction")
	if err != nil {
		return nil, "", fmt.Errorf("error building http-request: %s", err)
	}
	return httpReq, httpReqBody, nil
}

type SubscribeRequest struct {
	SubscriberRef            string
	SubscriptionIdentifier   string
	InitialTerminationTime   time.Time
	RequestTimestamp         time.Time
	MessageIdentifier        string
	PreviewInterval          string
	MonitoringRef            string
	StopVisitTypes           string
	MinimumStopVisitsPerLine int
	IncrementalUpdates       bool
	ChangeBeforeUpdates      string
}

func initSubscribeRequests(
	cfg *config.ConfigSubscribe,
	requestTimestamp *time.Time,
	initialTerminationTime *time.Time,
) []SubscribeRequest {
	numberOfSubascibeRequests := len(STOP_POINT_IDS_LILLE_BUS)
	requests := make([]SubscribeRequest, 0, numberOfSubascibeRequests)
	for _, stop_point_id := range STOP_POINT_IDS_LILLE_BUS {
		req := SubscribeRequest{}
		req.SubscriberRef = cfg.SubscriberRef
		req.SubscriptionIdentifier = cfg.SubscriberRef + ":Subscription:" + "arret_" + stop_point_id + ":LOC"
		req.InitialTerminationTime = requestTimestamp.AddDate(0, 0, 1)
		req.PreviewInterval = "PT2H0M0.000S"
		req.RequestTimestamp = *requestTimestamp
		req.MessageIdentifier = cfg.SubscriberRef + ":Message:" + requestTimestamp.Format(IDENTIFIER_TIME_LAYOUT)
		req.MonitoringRef = cfg.ProducerRef + ":StopPoint:BP:" + stop_point_id + ":LOC"
		req.StopVisitTypes = STOP_VISIT_TYPES
		req.MinimumStopVisitsPerLine = MINIMUM_STOP_VISITS_PER_LINE
		req.IncrementalUpdates = INCREMENTAL_UPDATES
		req.ChangeBeforeUpdates = "PT0M30.000S"
		requests = append(requests, req)
	}
	return requests
}

type Duration time.Duration

func (d *Duration) String() string {
	return "TODO SiriSM Duration"
}
