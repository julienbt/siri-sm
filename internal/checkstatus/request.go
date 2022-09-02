package checkstatus

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

type CheckStatusRequest struct {
	SupplierAddress   string
	RequestTimestamp  string
	RequestorRef      string
	MessageIdentifier string
}

type CheckStatusResult struct {
	SupplierServiceStartedTime time.Time
	LastSupplierCheckStatusOk  time.Time
}

func CheckStatus(cfg config.ConfigCheckStatus, logger *logrus.Entry) (CheckStatusResult, []byte, error) {
	var remoteErrorLoc = "CheckStatus remote error"
	checkStatusRequest := populateCheckStatusRequest(&cfg)
	req, err := generateSOAPCheckStatusHttpReq(checkStatusRequest)
	if err != nil {
		return CheckStatusResult{},
			nil,
			fmt.Errorf("error in building SOAP CheckStatus request: %s", err)
	}
	resp, err := siri.SoapCall(req)
	if err != nil {
		return CheckStatusResult{},
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("call error: %s", err)}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return CheckStatusResult{},
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("bad http-response status: %s", resp.Status)}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CheckStatusResult{},
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("unreadable response body: %s", err)}
	}
	checkStatusResponse := &CheckStatusResponseEnv{}
	err = xml.Unmarshal(body, &checkStatusResponse)
	if err != nil {
		return CheckStatusResult{},
			body,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("unmarshallable response body: %s", err)}
	}
	if !checkStatusResponse.CheckStatusResponseBody.CheckStatusResponse.CheckStatusResponseAnswer.Status {
		return CheckStatusResult{},
			body,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("status not true in response body")}
	}
	serviceStartedTime := checkStatusResponse.CheckStatusResponseBody.CheckStatusResponse.CheckStatusResponseAnswer.ServiceStartedTime.UTC()
	result := CheckStatusResult{
		SupplierServiceStartedTime: serviceStartedTime,
		LastSupplierCheckStatusOk:  time.Now(),
	}
	return result, body, nil
}

func populateCheckStatusRequest(cfg *config.ConfigCheckStatus) *CheckStatusRequest {
	now := time.Now()
	req := CheckStatusRequest{}
	req.RequestTimestamp = now.Format(time.RFC3339)
	req.RequestorRef = cfg.SubscriberRef
	req.MessageIdentifier = req.RequestorRef + ":ResponseMessage:" + now.Format("20060102_150405")
	req.SupplierAddress = cfg.SupplierAddress
	return &req
}

func generateSOAPCheckStatusHttpReq(req *CheckStatusRequest) (*http.Request, error) {
	tmpl, err := template.ParseFiles("./template/checkstatus-request.tmpl")
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
		"SOAPAction":   []string{"CheckStatus"},
	}
	httpReq.Header = headers // better than .Header.Set to preserve case (for "SOAPAction")
	if err != nil {
		return nil, fmt.Errorf("error building http-request: %s", err)
	}
	return httpReq, nil
}
