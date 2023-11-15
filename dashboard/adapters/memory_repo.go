package adapters

import (
	"fmt"
	"strings"
	"sync"
)

var StatisticsRepoMemory = NewStatisticsRepoMemory()

type ButtonClick struct {
	Title       string
	QuestionKey string
	OptionKey   string
	Timestamp   int
	CustomerID  string
	Platform    string
}

type CouponRevel struct {
	Code       string
	CustomerID string
	Platform   string
	Timestamp  int64
}

func (c *CouponRevel) ID() string {
	return c.Platform + "_" + c.Code
}

type statisticsRepoMemory struct {
	locker       sync.RWMutex
	clicks       map[string][]ButtonClick
	couponRevels map[string][]CouponRevel
}

func NewStatisticsRepoMemory() *statisticsRepoMemory {
	return &statisticsRepoMemory{
		clicks:       make(map[string][]ButtonClick),
		couponRevels: make(map[string][]CouponRevel),
	}
}

func (repo *statisticsRepoMemory) SaveRevels(revel CouponRevel) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	revels, ok := repo.couponRevels[revel.ID()]

	if !ok {
		repo.couponRevels[revel.ID()] = []CouponRevel{revel}
		fmt.Println("save metrics revel: ", revel.ID())
		return nil
	}

	repo.couponRevels[revel.ID()] = append(revels, revel)

	return nil
}

func (repo *statisticsRepoMemory) SaveClicks(click ButtonClick) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	id := repo.makeClickID(click.Platform, click.QuestionKey, click.OptionKey)
	buttons, ok := repo.clicks[id]

	if !ok {
		repo.clicks[id] = []ButtonClick{click}
		fmt.Println("save metrics: ", id)
		return nil
	}

	repo.clicks[id] = append(buttons, click)

	return nil

}

func (repo *statisticsRepoMemory) GetRevelCount(plataform, code string) int {
	return len(repo.get_revels(plataform + "_" + code))
}

func (repo *statisticsRepoMemory) get_revels(key string) []CouponRevel {
	repo.locker.RLock()
	defer repo.locker.RUnlock()

	revels, ok := repo.couponRevels[key]

	if !ok {
		return []CouponRevel{}
	}

	return revels

}

func (repo *statisticsRepoMemory) get_clicks(key string) []ButtonClick {
	repo.locker.RLock()
	defer repo.locker.RUnlock()

	clicks, ok := repo.clicks[key]

	if !ok {
		return []ButtonClick{}
	}

	return clicks

}

func (repo *statisticsRepoMemory) makeClickID(plataform, questionKey, optionKey string) string {
	id := fmt.Sprintf("%s_%s_%s", plataform, questionKey, optionKey)
	return strings.TrimSpace(id)
}

func (repo *statisticsRepoMemory) GetClickCount(plataform, questionKey, optionKey string) int {
	id := repo.makeClickID(plataform, questionKey, optionKey)
	clicks := repo.get_clicks(id)

	distinct := make(map[string]struct{})

	for _, click := range clicks {
		distinct[click.CustomerID] = struct{}{}
	}

	return len(distinct)
}
