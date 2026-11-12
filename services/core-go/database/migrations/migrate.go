package migrations

import (
	"log"

	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/app/modules/csr"
	"bank-sampah/core-go/app/modules/pickup"
	"bank-sampah/core-go/app/modules/reward"
	"bank-sampah/core-go/app/modules/tracking"
	"bank-sampah/core-go/app/modules/wallet"
	"bank-sampah/core-go/app/modules/weighbridge"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("⚡ Running database auto-migration...")
	err := db.AutoMigrate(
		&auth.User{},
		&pickup.PickupRequest{},
		&weighbridge.ScaleTicket{},
		&tracking.CourierLocation{},
		&wallet.Wallet{},
		&wallet.LedgerTransaction{},
		&wallet.RedeemRequest{},
		&reward.RewardItem{},
		&reward.RewardClaim{},
		&csr.CsrProgram{},
	)
	if err != nil {
		return err
	}
	log.Println("✓ Database auto-migration completed successfully")
	return nil
}
