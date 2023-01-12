package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/dumacp/go-driverconsole/internal/platform"
)

const (
	urlBundle = "https://fleet.nebulae.com.co/api/emi-gateway/graphql/http"
)

func PlataformRequestItinerary(client *http.Client, url, itineraryID, groupID string) (*platform.RouteMngItinerary, int, error) {
	dataBundle, statusCode, err := requestItinerary(client, url, itineraryID, groupID)
	if err != nil {
		return nil, statusCode, err
	}

	dataType := &struct {
		Data *platform.DataItinerary `json:"data"`
	}{}

	dataType.Data = new(platform.DataItinerary)

	err = json.Unmarshal(dataBundle, dataType)
	if err != nil {
		return nil, statusCode, err
	}

	if dataType.Data == nil || dataType.Data.Data == nil {
		return nil, statusCode, fmt.Errorf("empty data, response body: %s", dataBundle)
	}

	return dataType.Data.Data, statusCode, nil
}

func PlataformRequestMetadataItinerary(client *http.Client, url, itineraryID, groupID string) (*platform.RouteMngItinerary, int, error) {
	dataBundle, statusCode, err := requestMetadataItinerary(client, url, itineraryID, groupID)
	if err != nil {
		return nil, statusCode, err
	}

	dataType := &struct {
		Data *platform.DataItinerary `json:"data"`
	}{}

	dataType.Data = new(platform.DataItinerary)

	err = json.Unmarshal(dataBundle, dataType)
	if err != nil {
		return nil, statusCode, err
	}

	if dataType.Data == nil || dataType.Data.Data == nil {
		return nil, statusCode, fmt.Errorf("empty data")
	}

	return dataType.Data.Data, statusCode, nil
}

func requestItinerary(client *http.Client, url, itineraryID, organizationID string) ([]byte, int, error) {

	var jsonStr = []byte(fmt.Sprintf(`
{
	"operationName": "RouteMngItinerary",
	"variables": {
		"id": %q,
		"organizationId": %q
	},
	"query":"query RouteMngItinerary($id: ID!, $organizationId: String!) {\n  RouteMngItinerary(id: $id, organizationId: $organizationId) {\n    id\n    name\n    active\n    direction\n    organizationId\n    routeId\n    minCompletenessPercentage\n    distanceToItineraryThreshold\n    path {\n      type\n      name\n      coords\n      radius\n      eta\n      maxSpeed\n      checkPointId\n      punishable\n      punishableValue\n      __typename\n    }\n    metadata {\n      createdBy\n      createdAt\n      updatedBy\n      updatedAt\n      __typename\n    }\n    __typename\n  }\n}\n"}`, itineraryID, organizationID))

	log.Printf("request body: %s\n", jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	log.Println("response Status: ", resp.Status)
	log.Println("response Headers: ", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	log.Println("response Body: ", string(body))
	return body, resp.StatusCode, nil
}

func requestMetadataItinerary(client *http.Client, url, itineraryID, organizationID string) ([]byte, int, error) {

	var jsonStr = []byte(fmt.Sprintf(`
{
	"operationName": "RouteMngItinerary",
	"variables": {
		"id": %q,
		"organizationId": %q
	},
	"query":"query RouteMngItinerary($id: ID!, $organizationId: String!) {\n  RouteMngItinerary(id: $id, organizationId: $organizationId) {\n    id\n    name\n    active\n    direction\n    organizationId\n    routeId\n    minCompletenessPercentage\n    distanceToItineraryThreshold\n   metadata {\n      createdBy\n      createdAt\n      updatedBy\n      updatedAt\n      __typename\n    }\n    __typename\n  }\n}\n"}`, itineraryID, organizationID))

	log.Printf("request body: %s\n", jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	log.Println("response Status: ", resp.Status)
	log.Println("response Headers: ", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	log.Println("response Body: ", string(body))
	return body, resp.StatusCode, nil
}

func requestBundle(client *http.Client, url, deviceID, organizationID string) ([]byte, error) {

	var jsonStr = []byte(fmt.Sprintf(`
{
	"operationName": "GeneralMonitorBundle",
	"variables": {
		"organizationId": %q, 
		"id": %q
	},
	"query": "query GeneralMonitorBundle($id: ID!, $organizationId: String!) {\nGeneralMonitorBundle(id: $id, organizationId: $organizationId) {\nid\nname\ndescription\nactive\norganizationId\nservices\nvehicles\n__typename\n}\n}\n"}`,
		organizationID, deviceID))

	log.Printf("request body: %s\n", jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Println("response Status: ", resp.Status)
	log.Println("response Headers: ", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	log.Println("response Body: ", string(body))
	return body, nil

}

func PlataformRequestInfo(client *http.Client, url, deviceID, groupID string) (*platform.GeneralMonitorBundle, error) {

	dataBundle, err := requestBundle(client, url, deviceID, fmt.Sprintf("%s", groupID))
	if err != nil {
		return nil, err
	}

	dataType := &struct {
		Data *platform.DataBundle `json:"data"`
	}{}

	err = json.Unmarshal(dataBundle, dataType)
	if err != nil {
		return nil, err
	}

	return dataType.Data.GeneralMonitorBundle, nil
}

func GetDataDriversInServiceBundle(svc *platform.ServiceBundle) (string, int, error) {
	if len(svc.Driver.DocumentID) <= 0 {
		return "", 0, fmt.Errorf("empty driver")
	}
	documentID := strings.Map(func(r rune) rune {
		if !unicode.IsNumber(r) || !unicode.IsDigit(r) {
			return 0
		}
		return r
	}, svc.Driver.DocumentID)
	documentIDTrim := strings.Trim(documentID, "\x00")
	docDriver, err := strconv.ParseInt(documentIDTrim, 10, 64)
	if err != nil {
		return "", 0, err
	}
	return svc.Driver.FullName, int(docDriver), nil
}
