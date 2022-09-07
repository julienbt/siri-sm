package subscribe

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	siri_time "github.com/julienbt/siri-sm/internal/common/time"
	"github.com/stretchr/testify/require"
)

var testDataDir string

const SECONDS_PER_HOUR int = 3_600

var EXPECTED_LOCATION *time.Location = time.FixedZone("", 2*SECONDS_PER_HOUR)

func TestMain(m *testing.M) {

	testDataDir = os.Getenv("SIRISM_TEST_DATA_DIR")
	if testDataDir == "" {
		panic("$SIRISM_TEST_DATA_DIR isn't set")
	}

	os.Exit(m.Run())
}

func TestSubscribeResponseXmlUnmarshal(t *testing.T) {
	time.Local = time.UTC
	require := require.New(t)

	htmlRespFile, err := os.Open(
		fmt.Sprintf(
			"%s/examples/SUB_RESP_000_indented.xml",
			testDataDir,
		),
	)
	require.Nil(err)
	htmlRespBody, err := ioutil.ReadAll(htmlRespFile)
	require.Nil(err)

	envelope := SubscribeEnv{}
	err = xml.Unmarshal(htmlRespBody, &envelope)
	require.Nil(err)

	// Check the number of elements
	const EXPECTED_NUMBER_OF_RESPONSE_STATUS int = 50
	responseStatusList := envelope.SubscribeResponse.ResponseStatus
	require.Len(
		responseStatusList,
		EXPECTED_NUMBER_OF_RESPONSE_STATUS,
	)

	// Check the 1st element
	{
		EXPECTED_RESPONSE_STATUS := ResponseStatus{
			XMLName: xml.Name{
				Space: "http://www.siri.org.uk/siri",
				Local: "ResponseStatus",
			},
			// 2022-08-30T04:34:46.522+02:00
			ResponseTimestamp: siri_time.Time(time.Date(
				2022, time.August, 30,
				4, 34, 46, 522_000_000,
				EXPECTED_LOCATION,
			)),
			RequestMessageRef: "SUBHOR_ILEVIA:StopPoint:BP:11N001:LOC",
			SubscriberRef:     "KISIO2",
			SubscriptionRef:   "SUBHOR_ILEVIA:StopPoint:BP:11N001:LOC",
			Status:            true,
			// 2022-08-31T02:15:00.000+02:00
			ValidUntil: siri_time.Time(time.Date(
				2022, time.August, 31,
				2, 15, 00, 000_000_000,
				EXPECTED_LOCATION,
			)),
		}
		require.Equal(
			EXPECTED_RESPONSE_STATUS,
			responseStatusList[0],
		)
	}

	// Check the last element
	{
		EXPECTED_RESPONSE_STATUS := ResponseStatus{
			XMLName: xml.Name{
				Space: "http://www.siri.org.uk/siri",
				Local: "ResponseStatus",
			},
			// 2022-08-30T04:34:46.522+02:00
			ResponseTimestamp: siri_time.Time(time.Date(
				2022, time.August, 30,
				4, 34, 46, 522_000_000,
				EXPECTED_LOCATION,
			)),
			RequestMessageRef: "SUBHOR_ILEVIA:StopPoint:BP:ACC001:LOC",
			SubscriberRef:     "KISIO2",
			SubscriptionRef:   "SUBHOR_ILEVIA:StopPoint:BP:ACC001:LOC",
			Status:            true,
			// 2022-08-31T02:15:00.000+02:00
			ValidUntil: siri_time.Time(time.Date(
				2022, time.August, 31,
				2, 15, 00, 000_000_000,
				EXPECTED_LOCATION,
			)),
		}
		require.Equal(
			EXPECTED_RESPONSE_STATUS,
			responseStatusList[EXPECTED_NUMBER_OF_RESPONSE_STATUS-1],
		)
	}
}
