package db

import (
	"github.com/matbarofex/mtz-crypto/pkg/model"
	"github.com/matbarofex/mtz-crypto/pkg/store"
	"gorm.io/gorm"
)

type walletStore struct {
	db *gorm.DB
}

func NewWalletStore(db *gorm.DB) store.WalletStore {
	return &walletStore{db: db}
}

func (s *walletStore) GetWallet(walletID string) (rs model.Wallet, err error) {
	items := []model.WalletItem{}

	err = s.db.Find(&items, "wallet_id = ?", walletID).Error
	if err != nil {
		return rs, err
	}

	rs.ID = walletID
	rs.Items = items

	return rs, err
}
