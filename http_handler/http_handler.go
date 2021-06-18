package http_handler

import (
	"ltvco/c_cache"
	"ltvco/repository"
	"sync"
	"time"

	. "github.com/ahmetb/go-linq/v3"
)

type HttpInstance struct {
	Repo  *repository.Repository
	Cache *c_cache.CCache
}

func Listen() *HttpInstance {
	instance := &HttpInstance{
		Repo:  repository.New(),
		Cache: c_cache.New(),
	}

	go instance.LoadCache("2021-01-01")

	return instance
}

func (i *HttpInstance) LoadCache(from string) {
	currentTime, _ := time.Parse("2006-01-02", from)
	loadingDays := 90
	index := 1

	for {
		timeStr := currentTime.Format("2006-01-02")
		data, err := i.Repo.Daily(timeStr)

		if err == nil {
			if !i.Cache.Exist(timeStr) {
				i.Cache.Set(timeStr, data)
			}
		}

		index++
		currentTime = currentTime.AddDate(0, 0, 1)

		if index > loadingDays {
			break
		}
	}
}

func (i *HttpInstance) GetRecords(from, until, artist string) ([]repository.SongResponse, []repository.SongResponse, error) {
	var cached []repository.SongResponse
	var uncached []repository.SongResponse

	//sinle day
	if from == until {
		if i.Cache.Exist(from) {
			cachedData, _ := i.Cache.Get(from)

			cached = append(cached, cachedData.([]repository.SongResponse)...)
		} else {
			//not found in cache, retrieve data
			data, err := i.Repo.Daily(from)
			if err != nil {
				return nil, uncached, err
			}

			i.Cache.Set(from, data)

			uncached = append(uncached, data...)
		}
	} else {
		//multiple days
		untilTime, _ := time.Parse("2006-01-02", until)
		currentTime, _ := time.Parse("2006-01-02", from)

		for {
			timeStr := currentTime.Format("2006-01-02")

			if i.Cache.Exist(timeStr) {
				cachedData, _ := i.Cache.Get(timeStr)

				result := cachedData.([]repository.SongResponse)

				cached = append(cached, result...)
			} else {
				data, err := i.Repo.Daily(timeStr)
				if err != nil {
					return nil, nil, err
				}

				timeStr = currentTime.Format("2006-01-02")
				i.Cache.Set(timeStr, data)

				uncached = append(uncached, data...)
			}

			currentTime = currentTime.AddDate(0, 0, 1)

			if currentTime.After(untilTime) {
				break
			}
		}

		if artist != "" {
			var wg sync.WaitGroup
			wg.Add(2)

			go func(cached *[]repository.SongResponse, wg *sync.WaitGroup) {
				defer wg.Done()

				From(cached).Where(func(c interface{}) bool {
					return c.(repository.SongResponse).Artist == artist
				}).ToSlice(&cached)
			}(&cached, &wg)

			go func(uncached *[]repository.SongResponse, wg *sync.WaitGroup) {
				defer wg.Done()

				From(uncached).Where(func(c interface{}) bool {
					return c.(repository.SongResponse).Artist == artist
				}).ToSlice(&uncached)
			}(&uncached, &wg)

			wg.Wait()
		}
	}

	return cached, uncached, nil
}
