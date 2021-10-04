package cache

import (
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/store"
	"github.com/patrickmn/go-cache"
)

type walletCacheStore struct {
	cache       *cache.Cache
	walletStore store.WalletStore
}

func NewWalletCacheStore(cache *cache.Cache, walletStore store.WalletStore) store.WalletStore {
	return &walletCacheStore{
		cache:       cache,
		walletStore: walletStore,
	}
}

func (s *walletCacheStore) GetWallet(walletID string) (rs model.Wallet, err error) {
	wallet, found := s.cache.Get(walletID)
	if !found {
		wallet, err = s.walletStore.GetWallet(walletID)
		if err != nil {
			return rs, err
		}

		s.cache.Set(walletID, wallet, 0)
	}

	return wallet.(model.Wallet), nil
}
