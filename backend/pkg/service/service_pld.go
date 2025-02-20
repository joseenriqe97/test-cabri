package service

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/joseenriqe97/test-cabri/config"
)

type PldServiceStruct struct {
	Client *resty.Client
}

type RequestPld struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type ResponsePld struct {
	IsInBlacklist bool `json:"is_in_blacklist"`
}

type PldServiceInterface interface {
	CheckPld(request RequestPld) (ResponsePld, error)
}

var PldService PldServiceInterface = &PldServiceStruct{
	Client: resty.New(),
}

func (p *PldServiceStruct) CheckPld(request RequestPld) (ResponsePld, error) {
	var responseR ResponsePld
	resp, err := p.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&responseR).
		Post(config.GetUrlPLD() + "/check-blacklist")

	if err != nil {
		return ResponsePld{}, err
	}

	response, ok := resp.Result().(*ResponsePld)
	if !ok {
		return ResponsePld{}, fmt.Errorf("failed to parse response")
	}

	return *response, nil
}
