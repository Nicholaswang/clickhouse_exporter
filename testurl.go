package main

import (
	"fmt"
	"net/url"
)

func main() {
	u := "http://11:11"
	uri, _ := url.Parse(u)
	fmt.Printf("%v", uri)
	u = "http://22:11/;http://1323sdf:dfs"
	uri, _ = url.Parse(u)
	fmt.Printf("%v", uri)
}