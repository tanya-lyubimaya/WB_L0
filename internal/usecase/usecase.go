package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/tanya-lyubimaya/WB_L0/internal/domain"
	"github.com/tanya-lyubimaya/WB_L0/tools"
	"sync"
)

type useCase struct {
	logger         *logrus.Entry
	repoOrder      domain.OrderRepository
	cache          map[string][]byte
	natsConnection *nats.Conn
	mu             sync.Mutex
}

func New(repoOrders domain.OrderRepository, logger *logrus.Entry) (*useCase, error) {
	uc := &useCase{repoOrder: repoOrders, logger: logger}
	uc.cache = make(map[string][]byte)
	backup, err := repoOrders.Read(domain.Order{}, 1000)
	if err == nil {
		for _, v := range backup {
			temp, _ := json.Marshal(v)
			uc.cache[v.OrderUid] = temp
		}
	} else {
		uc.logger.Errorln(err)
	}
	uc.natsConnection, err = nats.Connect(tools.ConfigInstance().NATSUrl, nats.ReconnectBufSize(5*1024*1024))
	if err != nil {
		uc.logger.Errorln(err)
		return nil, err
	}

	_, err = uc.natsConnection.Subscribe("orders", uc.subscription)
	if err != nil {
		uc.logger.Errorln(err)
		return nil, err
	}

	return uc, nil
}

func (uc *useCase) subscription(m *nats.Msg) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.logger.Infoln(string(m.Data))
	var result domain.Order
	err := json.Unmarshal(m.Data, &result)
	if err != nil {
		uc.logger.Errorln(err)
		return
	}

	var validate *validator.Validate
	validate = validator.New()
	err = validate.Struct(result)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			uc.logger.Errorln(err)
			return
		}
		var hasErrors bool
		for _, err := range err.(validator.ValidationErrors) {
			if !hasErrors {
				hasErrors = true
			}
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		if hasErrors {
			return
		}
	}
	_, err = uc.repoOrder.Create(result)
	if err != nil {
		uc.logger.Errorln(err)
		return
	}
	if len(uc.cache) == 2000 {
		uc.cache = make(map[string][]byte)
		backup, err := uc.repoOrder.Read(domain.Order{}, 1000)
		if err == nil {
			for _, v := range backup {
				temp, _ := json.Marshal(v)
				uc.cache[v.OrderUid] = temp
			}
		}
	} else {
		uc.cache[result.OrderUid] = m.Data
	}
}

func (uc *useCase) ReadOrderByID(id string) (res *domain.Order, err error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	res = &domain.Order{}
	data, ok := uc.cache[id]
	if !ok {
		uc.logger.Infoln("UID not found in cache")
		res, err = uc.repoOrder.ReadByID(domain.Order{OrderUid: id})
	}
	err = json.Unmarshal(data, res)
	if err != nil {
		uc.logger.Errorln(err)
	}
	return
}

func (uc *useCase) Close() {
	uc.natsConnection.Close()
	uc.logger.Traceln("close use-case")
	uc.repoOrder.Close()
}
