package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (*SearchModel) FindGithub() (string, error) {
	values := url.Values{}
	values.Add("apiKey", "")
	values.Add("keyword", "calinit")
	resp, err := http.Get("https://ebica.backlog.jp/api/v2/issues?" + values.Encode())
	const blank = ""
	if err != nil {
		return blank, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return blank, err
	}
	var buf bytes.Buffer
	json.Indent(&buf, []byte(contents), "", "  ")

	return buf.String(), nil
}
