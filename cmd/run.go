package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/tanya-lyubimaya/WB_L0/internal/app"
	"os"
	"os/signal"
)

func Execute() error {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)

	application, err := app.New(logrus.NewEntry(logger))
	if err != nil {
		logger.Errorln(err)
		return err
	}
	go func() {
		err = application.Serve(":8080")
		logger.Traceln(err)
		logger.Traceln("application was closed!")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
	application.GracefulStop()
	return nil
}
