package main

import "fmt"
import "tisq/conf"

func main() {
	conf.LoadConf()
	fmt.Print(conf.GlobalConf())
}
