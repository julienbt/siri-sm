package checkstatus

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

type CheckStatusRequest struct {
	SupplierAddress   url.URL
	RequestTimestamp  time.Time
	RequestorRef      string
	MessageIdentifier string
}

type CheckStatusResult struct {
	SupplierServiceStartedTime time.Time
	LastSupplierCheckStatusOk  time.Time
}

func CheckStatus(
	cfg config.ConfigCheckStatus,
	logger *logrus.Entry,
	requestTimestamp *time.Time,
) (CheckStatusResult, string, []byte, error) {
	var remoteErrorLoc = "CheckStatus remote error"
	req := CheckStatusRequest{}
	err := req.populate(&cfg, requestTimestamp)
	if err != nil {
		return CheckStatusResult{},
			"",
			nil,
			fmt.Errorf("error CheckStatus request initialization: %v", err)
	}
	httpReq, htmlReqBody, err := generateHttpSoapReq(&req)
	if err != nil {
		return CheckStatusResult{},
			"",
			nil,
			fmt.Errorf("error in building SOAP CheckStatus request: %s", err)
	}

	// Send HTTP request and receive the response
	resp, err := siri.SoapCall(httpReq)
	if err != nil {
		return CheckStatusResult{},
			htmlReqBody,
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("call error: %s", err)}
	}

	// Get the HTTP response body
	defer resp.Body.Close()
	htmlRespBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CheckStatusResult{},
			htmlReqBody,
			nil,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("unreadable response body: %s", err)}
	}

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return CheckStatusResult{},
			htmlReqBody,
			htmlRespBody,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("bad http-response status: %s", resp.Status)}
	}

	// Parse the succesfull HTTP Response
	checkStatusResponse := &CheckStatusResponseEnv{}
	err = xml.Unmarshal(htmlRespBody, &checkStatusResponse)
	if err != nil {
		return CheckStatusResult{},
			htmlReqBody,
			htmlRespBody,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("unmarshallable response body: %s", err)}
	}
	if !checkStatusResponse.CheckStatusResponseBody.CheckStatusResponse.CheckStatusResponseAnswer.Status {
		return CheckStatusResult{},
			htmlReqBody,
			htmlRespBody,
			&siri.RemoteError{Loc: remoteErrorLoc, Err: fmt.Errorf("status not true in response body")}
	}
	serviceStartedTime := checkStatusResponse.CheckStatusResponseBody.CheckStatusResponse.CheckStatusResponseAnswer.ServiceStartedTime.UTC()
	result := CheckStatusResult{
		SupplierServiceStartedTime: serviceStartedTime,
		LastSupplierCheckStatusOk:  time.Now(),
	}
	return result, htmlReqBody, htmlRespBody, nil
}

func (req *CheckStatusRequest) populate(cfg *config.ConfigCheckStatus, requestTimestamp *time.Time) error {
	supplierAddressUrl, err := url.Parse(cfg.SupplierAddress)
	if err != nil {
		return fmt.Errorf("error the supplier address is not a valid URL: %s", cfg.SupplierAddress)
	}
	req.RequestTimestamp = *requestTimestamp
	req.RequestorRef = cfg.SubscriberRef
	req.MessageIdentifier = req.RequestorRef + ":ResponseMessage:" + requestTimestamp.Format(IDENTIFIER_TIME_LAYOUT)
	req.SupplierAddress = *supplierAddressUrl
	return nil
}

func generateHttpSoapReq(req *CheckStatusRequest) (*http.Request, string, error) {
	tmpl, err := template.ParseFiles("./template/checkstatus-request.tmpl")
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
		"SOAPAction":   []string{"CheckStatus"},
	}
	httpReq.Header = headers // better than .Header.Set to preserve case (for "SOAPAction")
	if err != nil {
		return nil, "", fmt.Errorf("error building http-request: %s", err)
	}
	return httpReq, htmlReqBody, nil
}
