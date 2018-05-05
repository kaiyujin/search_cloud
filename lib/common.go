package lib

import (
	"log"
	"strconv"
)

type SearchModel struct {
	Config
}

type Result struct {
	Title string
	URL   string
}

const emptyString = ""
const zero = 0

func (model *SearchModel) PageSizeInt() int {
	result, err := strconv.Atoi(model.PageSize)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
