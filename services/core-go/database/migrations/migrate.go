package migrations

import (
	"log"

	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/app/modules/compliance"
	"bank-sampah/core-go/app/modules/csr"
	"bank-sampah/core-go/app/modules/drop_point"
	"bank-sampah/core-go/app/modules/finance_report"
	"bank-sampah/core-go/app/modules/gamification"
	"bank-sampah/core-go/app/modules/inventory"
	"bank-sampah/core-go/app/modules/marketplace"
	"bank-sampah/core-go/app/modules/member_card"
	"bank-sampah/core-go/app/modules/notification"
	"bank-sampah/core-go/app/modules/pickup"
	"bank-sampah/core-go/app/modules/ppob"
	"bank-sampah/core-go/app/modules/reward"
	"bank-sampah/core-go/app/modules/tracking"
	"bank-sampah/core-go/app/modules/wallet"
	"bank-sampah/core-go/app/modules/weighbridge"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("⚡ Running database auto-migration for all enterprise modules...")
	err := db.AutoMigrate(
		// Core Foundation Modules
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

		// Extended Enterprise Modules
		&inventory.WasteStock{},
		&inventory.StockMutation{},
		&inventory.B2BSaleOrder{},
		&member_card.MemberCard{},
		&notification.Notification{},
		&drop_point.DropPoint{},
		&drop_point.MobilePickupSchedule{},
		&finance_report.OperationalExpense{},
		&ppob.PPOBTransaction{},
		&ppob.GoldSavingsAccount{},
		&ppob.WasteDonation{},
		&marketplace.EcoProduct{},
		&marketplace.EcoOrder{},
		&marketplace.EcoOrderItem{},
		&gamification.UserBadge{},
		&compliance.EPRCertificate{},
		&compliance.B2BSubscriptionContract{},
	)
	if err != nil {
		return err
	}
	log.Println("✓ Database auto-migration completed successfully")
	return nil
}
