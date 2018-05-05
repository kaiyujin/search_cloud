package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/koron/go-dproxy"
)

func (model *SearchModel) FindConfluence(keyword string) ([]Result, error) {

	values := makeConfluenceValues(model, keyword)
	results := []Result(nil)
	err := findConfluence(model, model.Confluence.Url+"?"+values.Encode(), &results)
	return results, err
}

func makeConfluenceValues(model *SearchModel, keyword string) url.Values {
	values := url.Values{}
	cql := fmt.Sprintf("text~\"%v\"%v", keyword, model.Confluence.OptionalConditions)
	values.Add("cql", cql)
	values.Add("limit", model.PageSize)
	values.Add("start", "1")
	return values
}

func findConfluence(model *SearchModel, url string, results *[]Result) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(model.Confluence.Id, model.Confluence.Pass)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v %s", resp.StatusCode, contents))
	}
	var j interface{}
	if err := json.Unmarshal(contents, &j); err != nil {
		return err
	}
	dp := dproxy.New(j)
	baseUrl, err := dp.M("_links").M("base").String()
	if err != nil {
		return err
	}
	err = parseConfluence(dp, baseUrl, results)
	if err != nil {
		return err
	}
	nextUrl, err := dp.M("_links").M("next").String()
	if err != nil {
		return nil //end recursive call
	}
	return findConfluence(model, baseUrl+nextUrl, results)
}

func parseConfluence(dp dproxy.Proxy, baseUrl string, results *[]Result) error {
	ary, err := dp.M("results").Array()
	if err != nil {
		return err
	}
	for _, ai := range ary {
		result := new(Result)
		m := ai.(map[string]interface{})
		result.Title = fmt.Sprintf("%v", m["title"])
		webUrl, err := dproxy.New(m).M("_links").M("webui").String()
		if err != nil {
			return err
		}
		result.URL = fmt.Sprintf("%v%v", baseUrl, webUrl)
		*results = append(*results, *result)
	}
	return nil
}
