package main

import (
	"fmt"
	"net/http"

	"github.com/darren-reddick/go-mixcloud-search/mixcloud"
)

func main() {

	filter, err := mixcloud.NewFilter(
		[]string{"ibiza"},
		[]string{""},
	)

	if err != nil {
		panic(err)
	}

	mc := mixcloud.NewSearch("maya jane coles", filter, &http.Client{}, mixcloud.NewStore())

	//err := mc.GetAllSync()
	err = mc.GetAllAsync()

	if err != nil {
		fmt.Println(err)
		return
	}

	mc.WriteJsonToFile()

	//fmt.Printf("%+v\n", rez)

}
