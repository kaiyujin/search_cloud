package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func (model *SearchModel) FindBacklog(keyword string) ([]Result, error) {
	values := makeBacklogCountValues(model, keyword)
	count, err := count(values, model.Backlog.Url)
	if err != nil {
		return nil, err
	}
	pageSize := model.PageSizeInt()
	pageMax := ((count - 1) / pageSize) + 1
	results := []Result(nil)
	for i := 0; i < pageMax; i++ {
		findValues := makeBacklogFindValues(model, keyword, i*pageSize)
		err = findBacklog(model, model.Backlog.Url+"api/v2/issues?"+findValues.Encode(), &results)
		if err != nil {
			return nil, err
		}
	}
	return results, err
}

func makeBacklogFindValues(model *SearchModel, keyword string, offset int) url.Values {
	values := makeBacklogCountValues(model, keyword)
	values.Add("offset", fmt.Sprintf("%v", offset))
	return values
}

func makeBacklogCountValues(model *SearchModel, keyword string) url.Values {
	values := url.Values{}
	values.Add("apiKey", model.Backlog.ApiKey)
	values.Add("keyword", keyword)
	values.Add("sort", "updated")
	values.Add("count", model.PageSize)
	return values
}

func findBacklog(model *SearchModel, url string, results *[]Result) error {
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
		return errors.New(fmt.Sprintf("%v %s", resp.StatusCode, contents))
	}

	var j interface{}
	if err := json.Unmarshal(contents, &j); err != nil {
		log.Printf("url:%v", url)
		log.Printf("contents:%v", contents)
		return err
	}

	for _, ai := range j.([]interface{}) {
		result := new(Result)
		m := ai.(map[string]interface{})
		result.Title = fmt.Sprintf("%v", m["summary"])
		result.URL = fmt.Sprintf("%vview/%v", model.Backlog.Url, m["issueKey"])
		*results = append(*results, *result)
	}
	return nil
}

func count(values url.Values, url string) (int, error) {
	resp, err := http.Get(url + "api/v2/issues/count?" + values.Encode())
	if err != nil {
		log.Printf("url:%v", url)
		log.Printf("response:%v", resp)
		return zero, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("url:%v", url)
		log.Printf("contents:%v", contents)
		return zero, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("url:%v", url)
		return zero, errors.New(fmt.Sprintf("%v %v", resp.StatusCode, contents))
	}
	var j interface{}
	if err := json.Unmarshal(contents, &j); err != nil {
		log.Printf("url:%v", url)
		log.Printf("contents:%v", contents)
		return zero, err
	}

	result, err := strconv.Atoi(fmt.Sprintf("%v", j.(map[string]interface{})["count"]))
	if err != nil {
		return zero, err
	}
	return result, nil
}
