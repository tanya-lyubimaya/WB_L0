package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tanya-lyubimaya/WB_L0/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		logrus.Fatalln(err)
	}
}
