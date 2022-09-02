package subscribe

import (
	"net/http"
	"time"

	"github.com/julienbt/siri-sm/internal/config"
	"github.com/sirupsen/logrus"
)

type SubscribeRequest struct {
	// SupplierAddress    string
	RequestTimestamp string
	SubscriberRef    string
	// MessageIdentifier      string
	ConsumerAddress        string
	SubscriptionIdentifier string
	InitialTerminationTime string
	PreviewInterval        string
	OperatorRef            string
	ChangeBeforeUpdates    string
}

type SubscribeResponse struct{}

func Subscribe(cfg config.ConfigSubscribe, logger *logrus.Entry) (SubscribeResponse, error) {
	var remoteErrorLoc = "Subscribe remote error"
	subscribeRequest := populateSubscribeRequest(&cfg)
	_, _ = generateSOAPSubscribeHttpReq(&subscribeRequest)
	_ = remoteErrorLoc
	return SubscribeResponse{}, nil
}

func populateSubscribeRequest(cfg *config.ConfigSubscribe) SubscribeRequest {
	now := time.Now()
	req := SubscribeRequest{}
	req.RequestTimestamp = now.Format(time.RFC3339)
	req.SubscriberRef = cfg.SubscriberRef
	req.ConsumerAddress = cfg.ConsumerAddress

	return req
}

func generateSOAPSubscribeHttpReq(req *SubscribeRequest) (*http.Request, error) {
	return nil, nil
}
