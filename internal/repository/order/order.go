package order

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/tanya-lyubimaya/WB_L0/internal/domain"
	_ "github.com/tanya-lyubimaya/WB_L0/internal/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db     *gorm.DB
	logger *logrus.Entry
}

type order struct {
	gorm.Model
	OrderUid string `gorm:"notNull;unique"`
	Data     []byte `gorm:"notNull"`
}

func New(db *gorm.DB, logger *logrus.Entry) (*orderRepository, error) {
	repo := &orderRepository{db: db, logger: logger}
	if !repo.db.Migrator().HasTable(&order{}) {
		err := repo.db.Migrator().AutoMigrate(&order{})
		if err != nil {
			repo.logger.Errorln(err)
			return nil, err
		}
	}
	repo.logger.Infoln("order repo constructed")
	return repo, nil
}

func (r *orderRepository) Create(o domain.Order) (*domain.Order, error) {
	model, err := domainToRepoModel(o)
	if err != nil {
		return nil, err
	}
	err = r.db.Create(model).Error
	if err != nil {
		r.logger.Errorln(err)
		return &o, err
	}
	return &o, nil
}

func (r *orderRepository) Read(o domain.Order, limit int) (out []*domain.Order, err error) {
	var result []*order
	err = r.db.Model(&order{}).Order("created_at DESC").Limit(limit).Find(&result).Error
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}
	out, err = repoToDomainModels(result)
	return
}

func (r *orderRepository) ReadByID(o domain.Order) (out *domain.Order, err error) {
	var result order
	err = r.db.Model(&order{}).Where("order_uid = ?", o.OrderUid).Find(&result).Error
	if err != nil {
		r.logger.Errorln(err)
		return nil, err
	}
	out, err = repoToDomainModel(result)
	return
}

func (r *orderRepository) Update(o domain.Order) (*domain.Order, error) {
	panic("This method can't be implemented since the data is static")
}

func (r *orderRepository) Delete(o domain.Order) error {
	panic("This method can't be implemented since the data is static")
}

func domainToRepoModel(o domain.Order) (out *order, err error) {
	out = &order{}
	out.OrderUid = o.OrderUid
	out.Data, err = json.Marshal(o)
	return
}

func repoToDomainModel(o order) (out *domain.Order, err error) {
	out = new(domain.Order)
	err = json.Unmarshal(o.Data, out)
	return
}

func repoToDomainModels(o []*order) (out []*domain.Order, err error) {
	out = make([]*domain.Order, 0, len(o))
	var temp *domain.Order
	for _, v := range o {
		temp, err = repoToDomainModel(*v)
		if err != nil {
			return nil, err
		}
		out = append(out, temp)
	}
	return
}

func (r *orderRepository) Close() {
	sql, err := r.db.DB()
	if err != nil {
		return
	}
	_ = sql.Close()
}
