package application

import (
	"time"

	"github.com/hirosato/gledger/domain"
)

type CommodityService struct {
	repository domain.CommodityRepository
}

func NewCommodityService(repository domain.CommodityRepository) *CommodityService {
	return &CommodityService{
		repository: repository,
	}
}

func (cs *CommodityService) FindOrCreateCommodity(symbol string) *domain.Commodity {
	return cs.repository.FindOrCreateCommodity(symbol)
}

func (cs *CommodityService) GetCommodity(symbol string) *domain.Commodity {
	return cs.repository.FindCommodity(symbol)
}

func (cs *CommodityService) RegisterCommodity(commodity *domain.Commodity) {
	cs.repository.RegisterCommodity(commodity)
}

func (cs *CommodityService) GetAllCommodities() []*domain.Commodity {
	return cs.repository.GetAllCommodities()
}

func (cs *CommodityService) SetDefaultCommodity(commodity *domain.Commodity) {
	cs.repository.SetDefaultCommodity(commodity)
}

func (cs *CommodityService) GetDefaultCommodity() *domain.Commodity {
	return cs.repository.GetDefaultCommodity()
}

func (cs *CommodityService) AddPrice(commoditySymbol string, date time.Time, amount *domain.Amount) {
	commodity := cs.repository.FindOrCreateCommodity(commoditySymbol)
	commodity.AddPrice(date, amount)
}

type InMemoryCommodityRepository struct {
	commodities      map[string]*domain.Commodity
	defaultCommodity *domain.Commodity
}

func NewInMemoryCommodityRepository() *InMemoryCommodityRepository {
	return &InMemoryCommodityRepository{
		commodities: make(map[string]*domain.Commodity),
	}
}

func (repo *InMemoryCommodityRepository) FindCommodity(symbol string) *domain.Commodity {
	return repo.commodities[symbol]
}

func (repo *InMemoryCommodityRepository) CreateCommodity(symbol string) *domain.Commodity {
	commodity := domain.NewCommodity(symbol)
	repo.commodities[symbol] = commodity
	return commodity
}

func (repo *InMemoryCommodityRepository) FindOrCreateCommodity(symbol string) *domain.Commodity {
	if commodity := repo.FindCommodity(symbol); commodity != nil {
		return commodity
	}
	return repo.CreateCommodity(symbol)
}

func (repo *InMemoryCommodityRepository) RegisterCommodity(commodity *domain.Commodity) {
	repo.commodities[commodity.Symbol] = commodity
}

func (repo *InMemoryCommodityRepository) GetAllCommodities() []*domain.Commodity {
	commodities := make([]*domain.Commodity, 0, len(repo.commodities))
	for _, commodity := range repo.commodities {
		commodities = append(commodities, commodity)
	}
	return commodities
}

func (repo *InMemoryCommodityRepository) SetDefaultCommodity(commodity *domain.Commodity) {
	repo.defaultCommodity = commodity
}

func (repo *InMemoryCommodityRepository) GetDefaultCommodity() *domain.Commodity {
	return repo.defaultCommodity
}