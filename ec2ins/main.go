package main

import (
	"./ec2"
	"./util"
)

func main() {
	util.ParseFlags()
	session := util.NewSession()
	ec2.Run(session)
}
