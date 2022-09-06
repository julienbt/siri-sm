package getstopmonitoring

import (
	"bytes"
	"encoding/xml"
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

const IDENTIFIER_TIME_LAYOUT string = "20060102_150405"

const MINIMUM_STOP_VISITS_PER_LINE int = 2

type GetStopMonitoringRequest struct {
	SupplierAddress          url.URL
	RequestTimestamp         time.Time
	RequestorRef             string
	MessageIdentifier        string
	MonitoringRef            string
	MinimumStopVisitsPerLine int
}

func GetStopMonitoring(
	cfg config.ConfigCheckStatus,
	logger *logrus.Entry,
	requestTimestamp *time.Time,
	monitoringRef string,
) ([]MonitoredStopVisit, string, []byte, error) {
	var remoteErrorLoc = "GetStopMonitoring remote error"

	getStopMonitoringRequest := GetStopMonitoringRequest{}
	getStopMonitoringRequest.populate(&cfg, requestTimestamp, monitoringRef)
	httpReq, htmlReqBody, err := getStopMonitoringRequest.generateHttpSoapReq()
	if err != nil {
		return nil,
			"",
			nil,
			fmt.Errorf("error in building SOAP GetStopMonitoring request: %s", err)
	}
	// Send HTTP request and receive the response
	resp, err := siri.SoapCall(httpReq)
	if err != nil {
		return nil,
			htmlReqBody,
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("call error: %s", err)}
	}

	// Get the HTTP response body
	defer resp.Body.Close()
	htmlRespBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,
			htmlReqBody,
			htmlRespBody,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("unreadable response body: %s", err)}
	}

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil,
			htmlReqBody,
			htmlRespBody,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("bad http-response status: %s", resp.Status)}
	}

	getStopMonitoringEnv := &GetStopMonitoringEnv{}
	err = xml.Unmarshal(htmlRespBody, &getStopMonitoringEnv)
	if err != nil {
		return nil,
			htmlReqBody,
			htmlRespBody,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("unmarshallable response body: %s", err)}
	}
	monitoredStopVisits, err := checkAndExtractMonitoredStopVisit(getStopMonitoringEnv)
	if err != nil {
		return nil,
			htmlReqBody,
			htmlRespBody,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: err}
	}
	return monitoredStopVisits, htmlReqBody, htmlRespBody, nil
}

func checkAndExtractMonitoredStopVisit(envelope *GetStopMonitoringEnv) ([]MonitoredStopVisit, error) {
	const EXPECTED_NUMBER_OF_MONITORED_STOP_VISITS int = 1
	const EXPECTED_NUMBER_OF_MONITORED_STOP_VISIT_CANCELLATIONS int = 0
	stopMonitoringDelivery := envelope.StopMonitoringDelivery
	if len(stopMonitoringDelivery.MonitoredStopVisits) !=
		EXPECTED_NUMBER_OF_MONITORED_STOP_VISITS {
		err := fmt.Errorf("invalid number of MonitoredStopVisit")
		return nil, err
	}
	if len(stopMonitoringDelivery.MonitoredStopVisitCancellations) !=
		EXPECTED_NUMBER_OF_MONITORED_STOP_VISIT_CANCELLATIONS {
		err := fmt.Errorf("invalid number of MonitoredStopVisitCancellation")
		return nil, err
	}
	return stopMonitoringDelivery.MonitoredStopVisits, nil
}

func (req *GetStopMonitoringRequest) populate(
	cfg *config.ConfigCheckStatus,
	requestTimestamp *time.Time,
	monitoringRef string,
) error {
	supplierAddressUrl, err := url.Parse(cfg.SupplierAddress)
	if err != nil {
		return fmt.Errorf("error the supplier address is not a valid URL: %s", cfg.SupplierAddress)
	}
	req.RequestTimestamp = *requestTimestamp
	req.RequestorRef = cfg.SubscriberRef
	req.MessageIdentifier = cfg.SubscriberRef + ":ResponseMessage:" + requestTimestamp.Format(IDENTIFIER_TIME_LAYOUT)
	req.MonitoringRef = monitoringRef
	req.MinimumStopVisitsPerLine = MINIMUM_STOP_VISITS_PER_LINE
	req.SupplierAddress = *supplierAddressUrl
	return nil
}

func (req *GetStopMonitoringRequest) generateHttpSoapReq() (*http.Request, string, error) {
	tmpl, err := template.ParseFiles("./template/getstopmonitoring-request.tmpl")
	if err != nil {
		return nil, "", fmt.Errorf("error parsing template: %s", err)
	}
	htmlReqBodyBuffer := &bytes.Buffer{}
	err = tmpl.Execute(htmlReqBodyBuffer, req)
	if err != nil {
		return nil, "", fmt.Errorf("error building template: %s", err)
	}
	htmlReqBody := htmlReqBodyBuffer.String()
	httpReq, err := http.NewRequest(http.MethodPost, req.SupplierAddress.String(), strings.NewReader(htmlReqBody))
	headers := http.Header{
		"Content-Type": []string{"text/xml; charset=utf-8"},
		"SOAPAction":   []string{"GetStopMonitoring"},
	}
	httpReq.Header = headers // better than .Header.Set to preserve case (for "SOAPAction")
	if err != nil {
		return nil, "", fmt.Errorf("error building http-request: %s", err)
	}
	return httpReq, htmlReqBody, nil
}
