package main

import (
	"fmt"
	"net/http"

	"github.com/darren-reddick/go-mixcloud-search/mixcloud"
)

func main() {

	mc := mixcloud.NewMixcloud("graeme park", mixcloud.Filter{}, &http.Client{}, mixcloud.NewStore())

	err := mc.GetAll()

	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Printf("%+v\n", rez)

}
