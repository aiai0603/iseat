package models

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/gorilla/schema"
	"lib/utils"
)

type CommonQuery struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

var decoder = schema.NewDecoder()

// Set default pagination
func NewPaginate() *CommonQuery {
	return &CommonQuery{
		PageNum:  1,
		PageSize: 20,
	}
}

func ParsePages(query url.Values, q utils.Messenger) (*CommonQuery, error) {
	var (
		temp string
		err  error
	)
	p := NewPaginate()
	if temp = query.Get("pageNum"); temp != "" {
		p.PageNum, err = strconv.Atoi(temp)
		if err != nil {
			msg := q.GetMsg(utils.MsgModel, "invalid_parameter")
			return nil, fmt.Errorf(msg, "pageNum", "int", temp)
		}
		query.Del("pageNum")
	}
	if temp = query.Get("pageSize"); temp != "" {
		p.PageSize, err = strconv.Atoi(temp)
		if err != nil {
			msg := q.GetMsg(utils.MsgModel, "invalid_parameter")
			return nil, fmt.Errorf(msg, "pageSize", "int", temp)
		}
		delete(query, "pageSize")
	}

	return p, nil
}
