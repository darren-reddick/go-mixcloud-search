package main

import (
	"fmt"
	"net/http"

	"github.com/darren-reddick/go-mixcloud-search/mixcloud"
)

func main() {

	mc := mixcloud.NewMixcloud("maya jane coles", mixcloud.Filter{}, &http.Client{}, mixcloud.NewStore())

	//err := mc.GetAllSync()
	err := mc.GetAllAsync()

	if err != nil {
		fmt.Println(err)
		return
	}

	mc.WriteJsonToFile()

	//fmt.Printf("%+v\n", rez)

}
