package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (m *SearchModel) FindBacklog(keyword string) ([]Result, error) {
	values := url.Values{}
	values.Add("apiKey", m.Backlog.ApiKey)
	values.Add("keyword", keyword)
	resp, err := http.Get(m.Backlog.Url + "?" + values.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var j interface{}
	err = json.Unmarshal(contents, &j)
	results := []Result(nil)
	for _, ai := range j.([]interface{}) {
		result := new(Result)
		m := ai.(map[string]interface{})
		result.Title = fmt.Sprintf("%v", m["summary"])
		result.URL = fmt.Sprintf("%v", m["issueKey"])
		results = append(results, *result)
	}

	return results, nil
}
