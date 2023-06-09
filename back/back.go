package function

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const URL_PREFIX = "http://webservice.recruit.co.jp/hotpepper/gourmet/v1/"

func init() {
	functions.HTTP("EntryPoint", EntryPoint)
}

func EntryPoint(w http.ResponseWriter, r *http.Request) {
	respBytes, code, err := entryPoint(w, r)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(code)
	w.Write(respBytes)
}

func entryPoint(w http.ResponseWriter, r *http.Request) ([]byte, int, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, http.StatusInternalServerError, fmt.Errorf("API_KEY env var is not defined")
	}
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("keyword query string is not given")
	}
	shopLocations, err := gourmetSearch(apiKey, keyword)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	respBytes, err := json.Marshal(shopLocations)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to encode JSON: %w", err)
	}
	return respBytes, http.StatusOK, nil
}

type apiResponse struct {
	Results apiResponseResults `json:"results"`
}

type apiResponseResults struct {
	Shop []apiResponseResultsShop `json:"shop"`
}

type apiResponseResultsShop struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

type shopLocation struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

func gourmetSearch(apiKey, keyword string) ([]shopLocation, error) {
	req, err := http.NewRequest(http.MethodGet, URL_PREFIX, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make new HTTP request: %w", err)
	}
	q := req.URL.Query()
	q.Add("key", apiKey)
	q.Add("keyword", keyword)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non 2xx HTTP status code: %d: %s", resp.StatusCode, resp.Status)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	fmt.Println(string(respBytes))

	apiResponse := apiResponse{}
	if err := json.Unmarshal(respBytes, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	shopLocations := []shopLocation{}
	for _, shop := range apiResponse.Results.Shop {
		shopLocations = append(shopLocations, shopLocation{Name: shop.Name, Lat: shop.Lat, Lon: shop.Lng})
	}

	return shopLocations, nil
}
