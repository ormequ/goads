package filters

import (
	"github.com/stretchr/testify/assert"
	"goads/internal/entities/ads"
	"testing"
	"time"
)

func FuzzAdByAuthor(f *testing.F) {
	f.Fuzz(func(t *testing.T, filterV int64, realV int64) {
		filter := AdByAuthor(AdsOptions{AuthorID: filterV})
		assert.Equal(t, filter(ads.Ad{ID: realV}), filterV == realV)
	})
}

func FuzzAdByCreateDate(f *testing.F) {
	f.Fuzz(func(t *testing.T, checkV uint64, realV uint64) {
		checkT := time.UnixMilli(int64(checkV)).Truncate(time.Second)
		realT := time.UnixMilli(int64(realV)).Truncate(time.Second)
		if checkT.IsZero() || checkT.UnixMilli() == 0 {
			return
		}
		filter := AdByAuthor(AdsOptions{Date: checkT})
		assert.Equal(t, filter(ads.Ad{CreateDate: realT}), realT == checkT)
	})
}

func FuzzAdPublished(f *testing.F) {
	f.Fuzz(func(t *testing.T, checkV bool, realV bool) {
		filter := AdByAuthor(AdsOptions{All: !checkV})
		assert.Equal(t, filter(ads.Ad{Published: realV}), realV == checkV)
	})
}
