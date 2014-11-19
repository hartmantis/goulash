package cookbook

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RoboticCheese/goulash/apiinstance"
	"github.com/RoboticCheese/goulash/common"
)

func cdata() (data Cookbook) {
	data = Cookbook{
		Component:         common.Component{Endpoint: "https://example1.com"},
		Name:              "test1",
		Maintainer:        "someuser",
		Description:       "A cookbook",
		Category:          "Other",
		LatestVersion:     "1.2.3",
		ExternalURL:       "https://extexample1.com",
		AverageRating:     0,
		CreatedAt:         "2014-09-01T01:01:01.123Z",
		UpdatedAt:         "2014-09-02T01:01:01.123Z",
		Deprecated:        false,
		FoodcriticFailure: false,
		Versions:          []string{"1.2.3", "1.2.0", "1.1.0"},
		Metrics: Metrics{
			Downloads: Downloads{
				Total: 99,
				Versions: map[string]int{
					"1.2.3": 32,
					"1.2.0": 33,
					"1.1.0": 34,
				},
			},
			Followers: 123,
		},
	}
	return
}

var jsonData = map[string]string{
	"name":               "chef-dk",
	"maintainer":         "roboticcheese",
	"description":        "Installs/configures the Chef-DK",
	"category":           "Other",
	"latest_version":     "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1",
	"external_url":       "https://github.com/RoboticCheese/chef-dk-chef",
	"average_rating":     "null",
	"created_at":         "2014-06-24T01:14:49.000Z",
	"updated_at":         "2014-09-20T04:46:00.780Z",
	"deprecated":         "false",
	"foodcritic_failure": "false",
	"versions": `
		[ "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1",
		 "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0"]`,
	"metrics": `{
		"downloads": {
			"total": 100,
			"versions": {
				"2.0.0": 50,
				"2.0.1": 50
			}
		},
		"followers": 20
	}`,
}

func jsonified() (res string) {
	res = `{"name": "` + jsonData["name"] + `",` +
		`"maintainer": "` + jsonData["maintainer"] + `",` +
		`"description": "` + jsonData["description"] + `",` +
		`"category": "` + jsonData["category"] + `",` +
		`"latest_version": "` + jsonData["latest_version"] + `",` +
		`"external_url": "` + jsonData["external_url"] + `",` +
		`"average_rating": ` + jsonData["average_rating"] + `,` +
		`"created_at": "` + jsonData["created_at"] + `",` +
		`"updated_at": "` + jsonData["updated_at"] + `",` +
		`"deprecated": ` + jsonData["deprecated"] + `,` +
		`"foodcritic_failure": ` + jsonData["foodcritic_failure"] + `,` +
		`"versions": ` + jsonData["versions"] + `,` +
		`"metrics": ` + jsonData["metrics"] + `}`
	return
}

func startHTTP() (ts *httptest.Server) {
	ts = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, jsonified())
			},
		),
	)
	return
}

func Test_Equals_1_Equal(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	res := data1.Equals(&data2)
	if res != true {
		t.Fatalf("Expected: true, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != true {
		t.Fatalf("Expected: true, got: %v", res)
	}
}

func Test_Equals_2_DifferentEndpoints(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Endpoint = "https://somewherelse.com"
	res := data1.Equals(&data2)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func Test_Equals_3_DifferentName(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Name = "ansible"
	res := data1.Equals(&data2)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func Test_Equals_4_DifferentLatestVersion(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.LatestVersion = "9.9.9"
	res := data1.Equals(&data2)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func Test_Equals_5_DifferentVersions(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions = append(data2.Versions, "9.9.9")
	res := data1.Equals(&data2)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func Test_Equals_6_DifferentMetrics(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Metrics.Downloads.Versions["1.2.3"] = 999
	res := data1.Equals(&data2)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func Test_New_1_NoError(t *testing.T) {
	ts := startHTTP()
	defer ts.Close()

	i := new(apiinstance.APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	c, err := New(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for k, v := range map[string]string{
		c.Endpoint:      ts.URL + "/api/v1/cookbooks/chef-dk",
		c.Name:          jsonData["name"],
		c.Maintainer:    jsonData["maintainer"],
		c.Description:   jsonData["description"],
		c.Category:      jsonData["category"],
		c.LatestVersion: jsonData["latest_version"],
		c.ExternalURL:   jsonData["external_url"],
		c.CreatedAt:     jsonData["created_at"],
		c.UpdatedAt:     jsonData["updated_at"],
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
	if c.Deprecated != false {
		t.Fatalf("Expected: false, got: %v", c.Deprecated)
	}
	if c.FoodcriticFailure != false {
		t.Fatalf("Expected: false, got: %v", c.FoodcriticFailure)
	}
	if c.AverageRating != 0 {
		t.Fatalf("Expected: 0, got: %v", c.AverageRating)
	}
	if len(c.Versions) != 2 {
		t.Fatalf("Expected: 2 versions, got: %v", len(c.Versions))
	}
	ver := "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1"
	if c.Versions[0] != ver {
		t.Fatalf("Expected: %v, got: %v", ver, c.Versions[0])
	}
	ver = "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0"
	if c.Versions[1] != ver {
		t.Fatalf("Expected: %v, got: %v", ver, c.Versions[1])
	}
	if c.Metrics.Downloads.Total != 100 {
		t.Fatalf("Expected: 100, got: %v", c.Metrics.Downloads.Total)
	}
	if c.Metrics.Downloads.Versions["2.0.0"] != 50 {
		t.Fatalf("Expected: 50, got: %v", c.Metrics.Downloads.Versions["2.0.0"])
	}
	if c.Metrics.Downloads.Versions["2.0.1"] != 50 {
		t.Fatalf("Expected: 50, got: %v", c.Metrics.Downloads.Versions["2.0.1"])
	}
	if c.Metrics.Followers != 20 {
		t.Fatalf("Expected: 20, got: %v", c.Metrics.Followers)
	}
}

func Test_New_2_NilFoodcriticFailure(t *testing.T) {
	jsonData["foodcritic_failure"] = "null"
	ts := startHTTP()
	defer ts.Close()

	i := new(apiinstance.APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	c, err := New(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if c.FoodcriticFailure != false {
		t.Fatalf("Expected: nil, got: %v", c.FoodcriticFailure)
	}
}

func Test_New_3_AverageRating(t *testing.T) {
	jsonData["average_rating"] = "20"
	ts := startHTTP()
	defer ts.Close()

	i := new(apiinstance.APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	c, err := New(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if c.AverageRating != 20 {
		t.Fatalf("Expected: 20, got: %v", c.AverageRating)
	}
}

func Test_New_4_ConnError(t *testing.T) {
	ts := startHTTP()
	ts.Close()

	i := new(apiinstance.APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	_, err := New(i, "chef-dk")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_5_404Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer ts.Close()

	i := new(apiinstance.APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	_, err := New(i, "chef-dk")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_6_RealData(t *testing.T) {
	i := new(apiinstance.APIInstance)
	i.Endpoint = "https://supermarket.getchef.com/api/v1"
	c, err := New(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for k, v := range map[string]string{
		c.Name:        "chef-dk",
		c.Maintainer:  "roboticcheese",
		c.Category:    "Other",
		c.ExternalURL: "https://github.com/RoboticCheese/chef-dk-chef",
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
}
