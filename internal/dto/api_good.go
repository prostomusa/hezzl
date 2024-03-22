package dto

import (
	"time"
)

type CreateGoodRequest struct {
	Name string `json:"name"`
}

type UpdateGoodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ReprioritizeGoodRequest struct {
	NewPriority int `json:"newPriority"`
}

type ReprioritizeGoodResponse struct {
	Priorities []PriorityResponse `json:"priorities"`
}

type Good struct {
	Id          int        `json:"id"`
	ProjectId   int        `json:"projectId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	Removed     bool       `json:"removed"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type DeleteGoodResponse struct {
	Id        int  `json:"id"`
	ProjectId int  `json:"projectId"`
	Removed   bool `json:"removed"`
}

type PriorityResponse struct {
	Id       int `json:"id"`
	Priority int `json:"priority"`
}

type GetGoodsResponse struct {
	Meta  MetaGoodsResponse `json:"meta"`
	Goods []Good            `json:"goods"`
}

type MetaGoodsResponse struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}
