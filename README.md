# Overview

## Description
Execute collective search from multiple cloud services.  
Currently it corresponds to the following.

- Backlog
- Confluence
- Github

## Requirement
- go 1.9 or later
- dep

## Install
`go build main.go`

A configuration file is required for execution.  
Create `config.toml` in the same directory as the executable file.

##### sample

```
page_size = "100"

[backlog]
url = "https://xxxx.backlog.jp/"
api_key = "xxxxxxxxxxxxxxxxxxxx"

[github]
url = "https://api.github.com/"
api_key = "xxxxxxxxxxxxxxxxxxxxx"
optional_conditions = "+org:xxxxxx"

[confluence]
url = "https://xxxxx.atlassian.net/wiki/rest/api/content/search"
optional_conditions = "+and+space=xxxxxx+and+type=page+order+by+created+desc"
id = "hoge@fuga.co.jp"
pass = "xxxxxxxxxx"
```

## Usage
`./main search_word`