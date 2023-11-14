package adapters

import (
	"fmt"
	"strings"
	"sync"
)

var ButtonStatisticsRepoMemory = NewButtonStatisticsRepoMemory()

type ButtonClick struct {
	Title      string
	Key        string
	Timestamp  int
	CustomerID string
	Platform   string
}

type buttonStatisticsRepoMemory struct {
	locker sync.RWMutex
	data   map[string][]ButtonClick
}

func NewButtonStatisticsRepoMemory() *buttonStatisticsRepoMemory {
	return &buttonStatisticsRepoMemory{
		data: make(map[string][]ButtonClick),
	}
}

func (repo *buttonStatisticsRepoMemory) Save(click ButtonClick) error {
	repo.locker.Lock()
	defer repo.locker.Unlock()

	id := repo.makeID(click.Platform, click.Key)

	buttons, ok := repo.data[id]

	if !ok {
		repo.data[id] = []ButtonClick{click}
		fmt.Println("save metrics: ", id)
		return nil
	}

	repo.data[id] = append(buttons, click)

	return nil

}

func (repo *buttonStatisticsRepoMemory) get(key string) []ButtonClick {
	repo.locker.RLock()
	defer repo.locker.RUnlock()

	clicks, ok := repo.data[key]

	if !ok {
		return []ButtonClick{}
	}

	return clicks

}

func (repo *buttonStatisticsRepoMemory) makeID(plataform, key string) string {
	id := fmt.Sprintf("%s_%s", plataform, key)
	return strings.TrimSpace(id)
}

func (repo *buttonStatisticsRepoMemory) GetClickCount(plataform, key string) int {
	id := repo.makeID(plataform, key)
	clicks := repo.get(id)

	distinct := make(map[string]struct{})

	for _, click := range clicks {
		distinct[click.CustomerID] = struct{}{}
	}

	return len(distinct)
}
