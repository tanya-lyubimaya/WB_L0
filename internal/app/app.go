package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tanya-lyubimaya/WB_L0/internal/delivery/http"
	"github.com/tanya-lyubimaya/WB_L0/internal/repository/order"
	"github.com/tanya-lyubimaya/WB_L0/internal/usecase"
	"github.com/tanya-lyubimaya/WB_L0/tools"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Application struct {
	server *http.Server
	logger *logrus.Entry
}

func (a *Application) Serve(port string) error {
	return a.server.Serve(port)
}

func (a *Application) GracefulStop() {
	a.server.GracefulShutdown()
}

func New(logger *logrus.Entry) (*Application, error) {
	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
			tools.ConfigInstance().Host,
			tools.ConfigInstance().Username,
			tools.ConfigInstance().Password,
			tools.ConfigInstance().DBName,
			tools.ConfigInstance().Port,
		),
	), &gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	repoOrder, err := order.New(db, logrus.NewEntry(logger.Logger))
	if err != nil {
		return nil, err
	}
	uc, err := usecase.New(repoOrder, logger)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	server, err := http.New(uc, logger)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return &Application{server: server, logger: logger}, nil
}
