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
	req.SubscribeRequests = make([]SubscribeRequest, 0)
	{
		req.SubscribeRequests = append(req.SubscribeRequests, SubscribeRequest{})
		req.SubscribeRequests[0].populate(cfg, requestTimestamp, initialTerminationTime)
	}
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
	PreviewInterval          Duration
	MonitoringRef            string
	StopVisitTypes           string
	MinimumStopVisitsPerLine int
	IncrementalUpdates       bool
	ChangeBeforeUpdates      Duration
}

func (req *SubscribeRequest) populate(
	cfg *config.ConfigSubscribe,
	requestTimestamp *time.Time,
	initialTerminationTime *time.Time,
) error {
	req.SubscriberRef = cfg.SubscriberRef
	req.SubscriptionIdentifier = "TODO SubscriptionIdentifier"
	req.InitialTerminationTime = *initialTerminationTime
	req.RequestTimestamp = *requestTimestamp
	req.MessageIdentifier = "TODO MessageIdentifier"
	req.PreviewInterval = Duration(2 * time.Hour)
	req.MonitoringRef = "TODO MonitoringRef"
	req.StopVisitTypes = "departures"
	req.MinimumStopVisitsPerLine = 2
	req.IncrementalUpdates = true
	req.ChangeBeforeUpdates = Duration(30 * time.Second)
	return nil
}

type Duration time.Duration

func (d *Duration) String() string {
	return "TODO SiriSM Duration"
}
