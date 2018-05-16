package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"log"

	"github.com/koron/go-dproxy"
)

func (model *SearchModel) FindGithub(keyword string) ([]Result, error) {
	values := makeGithubValues(model, keyword)
	url := model.Github.Url + "search/issues?" + values.Encode()
	results := []Result(nil)
	err := findGithub(url, keyword, model, &results)
	return results, err
}

func makeGithubValues(model *SearchModel, keyword string) url.Values {
	values := url.Values{}
	values.Add("access_token", model.Github.ApiKey)
	values.Add("q", keyword+model.Github.OptionalConditions)
	values.Add("sort", "updated")
	values.Add("order", "desc")
	values.Add("per_page", model.PageSize)
	return values
}

func findGithub(url, keyword string, model *SearchModel, results *[]Result) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("url:%v", url)
		log.Printf("response:%v", resp)
		return err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("url:%v", url)
		log.Printf("contents:%v", contents)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("url:%v", url)
		return errors.New(fmt.Sprintf("%v %v", resp.StatusCode, contents))
	}

	var j interface{}
	if err := json.Unmarshal(contents, &j); err != nil {
		log.Printf("url:%v", url)
		log.Printf("contents:%v", contents)
		return err
	}

	results, err = parseGithub(dproxy.New(j), keyword, results)
	if err != nil {
		log.Printf("url:%v", url)
		log.Printf("results:%v", results)
		return err
	}
	nextUrl := relNext(resp.Header.Get("Link"))
	if nextUrl != emptyString {
		return findGithub(nextUrl, keyword, model, results)
	}
	return nil
}

func parseGithub(dp dproxy.Proxy, keyword string, results *[]Result) (*[]Result, error) {
	ary, err := dp.M("items").Array()
	if err != nil {
		return nil, err
	}
	for _, ai := range ary {
		m := ai.(map[string]interface{})
		title := fmt.Sprintf("%v", m["title"])
		result := new(Result)
		result.Title = fmt.Sprintf("%v", title)
		result.URL = fmt.Sprintf("%v", m["html_url"])
		*results = append(*results, *result)
	}
	return results, nil
}

var replacer = strings.NewReplacer("<", "", ">", "")

func relNext(header string) string {
	a := strings.Split(header, ",")
	result := emptyString
	for _, v := range a {
		if strings.Contains(v, "rel=\"next\"") {
			result = replacer.Replace(v)
			result = strings.TrimSpace(strings.Split(result, ";")[0])
			return result
		}
	}
	return emptyString
}
