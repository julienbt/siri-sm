package subscribe

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/julienbt/siri-sm/internal/config"
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

var STOP_POINT_IDS = []string{
	"CAS001",
	"CAS002",
	"CAT001",
	"CAT002",
	"CAU001",
	"CAU002",
	"CAV001",
	"CAV002",
	"CAW001",
	"CAW002",
	"CBA011",
	"CBA012",
	"CBE001",
	"CBE002",
	"CBF001",
	"CBF002",
	"CBG001",
	"CBG002",
	"CBO002",
	"CBO004",
	"CCD001",
	"CCD002",
	"CCE001",
	"CCE002",
	"CCH001",
	"CCH002",
	"CDE001",
	"CDE002",
	"CDO001",
	"CDO002",
	"CDP001",
	"CDP002",
	"CDT001",
	"CDT002",
	"CED001",
	"CED002",
	"CED002",
	"CEH001",
	"CEN001",
	"CEN001",
	"CEO001",
	"CEO002",
	"CER001",
}

func Subscribe(cfg config.ConfigSubscribe, logger *logrus.Entry, location *time.Location) (SubscribeRequestInfoResult, string, []byte, error) {
	var remoteErrorLoc = "Subscribe remote error"

	req := SubscribeRequestInfo{}
	{
		requestTimestamp := time.Now().In(location)
		req.populate(&cfg, &requestTimestamp, &requestTimestamp)
	}

	httpReq, httpReqBody, err := req.generateHttpSoapSubscribeReq()
	// Check/parse the HTTP Response
	if err != nil {
		if err != nil {
			return SubscribeRequestInfoResult{},
				"",
				nil,
				fmt.Errorf("error in building SOAP Subscribe request: %s", err)
		}
	}
	_ = httpReq
	_ = remoteErrorLoc
	return SubscribeRequestInfoResult{}, httpReqBody, nil, nil
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

func (req *SubscribeRequestInfo) generateHttpSoapSubscribeReq() (*http.Request, string, error) {
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
	numberOfSubascibeRequests := len(STOP_POINT_IDS)
	requests := make([]SubscribeRequest, 0, numberOfSubascibeRequests)
	for _, stop_point_id := range STOP_POINT_IDS {
		req := SubscribeRequest{}
		req.SubscriberRef = cfg.SubscriberRef
		req.SubscriptionIdentifier = cfg.SubscriberRef + ":Subscription"
		req.InitialTerminationTime = requestTimestamp.AddDate(0, 0, 1)
		req.PreviewInterval = "PT24H0M0.000S"
		req.RequestTimestamp = *requestTimestamp
		req.MessageIdentifier = cfg.SubscriberRef + ":Message:" + requestTimestamp.Format(IDENTIFIER_TIME_LAYOUT)
		req.MonitoringRef = cfg.ProducerRef + ":StopPoint:BP:" + stop_point_id + ":LOC"
		req.MinimumStopVisitsPerLine = MINIMUM_STOP_VISITS_PER_LINE
		req.IncrementalUpdates = INCREMENTAL_UPDATES
		req.ChangeBeforeUpdates = "PT1M"
		requests = append(requests, req)
	}
	return requests
}

type Duration time.Duration

func (d *Duration) String() string {
	return "TODO SiriSM Duration"
}
