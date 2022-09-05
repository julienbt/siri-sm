package getstopmonitoring

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/julienbt/siri-sm/internal/config"
	"github.com/julienbt/siri-sm/internal/siri"
	"github.com/sirupsen/logrus"
)

type GetStopMonitoringRequest struct {
	SupplierAddress          string
	RequestTimestamp         string
	RequestorRef             string
	MessageIdentifier        string
	MonitoringRef            string
	MinimumStopVisitsPerLine int
}

func GetStopMonitoring(
	cfg config.ConfigCheckStatus,
	logger *logrus.Entry,
	monitoringRef string,
) ([]MonitoredStopVisit, []byte, error) {
	var remoteErrorLoc = "GetStopMonitoring remote error"
	getStopMonitoringRequest := populateGetStopMonitoringRequest(&cfg, monitoringRef)
	req, err := generateSOAPCheckStatusHttpReq(getStopMonitoringRequest)
	if err != nil {
		return nil,
			nil,
			fmt.Errorf("error in building SOAP GetStopMonitoring request: %s", err)
	}
	resp, err := siri.SoapCall(req)
	if err != nil {
		return nil,
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("call error: %s", err)}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil,
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("bad http-response status: %s", resp.Status)}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,
			body,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("unreadable response body: %s", err)}
	}
	getStopMonitoringEnv := &GetStopMonitoringEnv{}
	err = xml.Unmarshal(body, &getStopMonitoringEnv)
	if err != nil {
		return nil,
			body,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("unmarshallable response body: %s", err)}
	}
	monitoredStopVisits, err := checkAndExtractMonitoredStopVisit(getStopMonitoringEnv)
	if err != nil {
		return nil,
			body,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: err}
	}
	return monitoredStopVisits, body, nil
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

func populateGetStopMonitoringRequest(cfg *config.ConfigCheckStatus, monitoringRef string) *GetStopMonitoringRequest {
	now := time.Now()
	req := GetStopMonitoringRequest{}
	req.RequestTimestamp = now.Format(time.RFC3339)
	req.RequestorRef = cfg.SubscriberRef
	// req.MessageIdentifier = "KISIO2_ILEVIA:Message::11234:LOC"
	req.MessageIdentifier = cfg.SubscriberRef + ":ResponseMessage:" + now.Format("20060102_150405")
	req.MonitoringRef = monitoringRef
	req.MinimumStopVisitsPerLine = 2
	req.SupplierAddress = cfg.SupplierAddress
	return &req
}

func generateSOAPCheckStatusHttpReq(req *GetStopMonitoringRequest) (*http.Request, error) {
	tmpl, err := template.ParseFiles("./template/getstopmonitoring-request.tmpl")
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %s", err)
	}
	doc := &bytes.Buffer{}
	err = tmpl.Execute(doc, req)
	if err != nil {
		return nil, fmt.Errorf("error building template: %s", err)
	}
	httpReq, err := http.NewRequest(http.MethodPost, req.SupplierAddress, strings.NewReader(doc.String()))
	headers := http.Header{
		"Content-Type": []string{"text/xml; charset=utf-8"},
		"SOAPAction":   []string{"GetStopMonitoring"},
	}
	httpReq.Header = headers // better than .Header.Set to preserve case (for "SOAPAction")
	if err != nil {
		return nil, fmt.Errorf("error building http-request: %s", err)
	}
	return httpReq, nil
}
