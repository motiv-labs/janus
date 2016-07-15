package main

import (
	"github.com/kataras/iris"
)

type AuthAPI struct{}

// GET /auth
func (u AuthAPI) Serve(iris *iris.Context) {
	//config := &osincli.ClientConfig{
	//	ClientId:                 "1234",
	//	ClientSecret:             "aabbccdd",
	//	AuthorizeUrl:             "https://auth:3000/authentication",
	//	TokenUrl:                 "https://auth:3000/token",
	//	ErrorsInStatusCode:       true,
	//	SendClientSecretInParams: true,
	//}
	//
	//client, err := osincli.NewClient(config)
	//if err != nil {
	//	panic(err)
	//}
}
