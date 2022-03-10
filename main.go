package main

import (
	"fmt"
	"net/http"

	"github.com/darren-reddick/go-mixcloud-search/mixcloud"
	"github.com/darren-reddick/go-mixcloud-search/store"
)

func main() {

	mc := mixcloud.NewMixcloud("graeme park", mixcloud.Filter{}, &http.Client{}, store.NewStore())

	_, err := mc.Get()

	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Printf("%+v\n", rez)

}
