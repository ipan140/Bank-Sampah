package seeds

import (
	"log"
	"time"

	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/app/modules/drop_point"
	"bank-sampah/core-go/app/modules/inventory"
	"bank-sampah/core-go/app/modules/marketplace"
	"bank-sampah/core-go/app/modules/member_card"
	"bank-sampah/core-go/app/modules/reward"
	"bank-sampah/core-go/app/modules/wallet"
	"bank-sampah/core-go/app/modules/weighbridge"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) error {
	log.Println("🌱 Checking and seeding initial enterprise data...")

	// 1. Seed Users (Admin, Nasabah, Kurir, Mitra)
	var adminCount int64
	db.Model(&auth.User{}).Where("email = ?", "admin@banksampah.id").Count(&adminCount)

	adminID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	nasabahID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	kurirID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	mitraID := uuid.MustParse("00000000-0000-0000-0000-000000000004")

	if adminCount == 0 {
		hashedPassAdmin, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
		hashedPassUser, _ := bcrypt.GenerateFromPassword([]byte("Nasabah123!"), bcrypt.DefaultCost)
		hashedPassKurir, _ := bcrypt.GenerateFromPassword([]byte("Kurir123!"), bcrypt.DefaultCost)
		hashedPassMitra, _ := bcrypt.GenerateFromPassword([]byte("Mitra123!"), bcrypt.DefaultCost)

		users := []auth.User{
			{
				ID:        adminID,
				Name:      "Super Administrator",
				Email:     "admin@banksampah.id",
				Password:  string(hashedPassAdmin),
				Phone:     "081100000001",
				Role:      auth.RoleAdmin,
				Address:   "Kantor Pusat Bank Sampah, Jl. Daur Ulang No. 1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        nasabahID,
				Name:      "Budi Santoso",
				Email:     "nasabah@banksampah.id",
				Password:  string(hashedPassUser),
				Phone:     "081234567890",
				Role:      auth.RoleNasabah,
				Address:   "Jl. Merdeka No. 10, RT 03/RW 05, Bandung",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        kurirID,
				Name:      "Joko Supriyanto (Kurir Armada 1)",
				Email:     "kurir@banksampah.id",
				Password:  string(hashedPassKurir),
				Phone:     "081987654321",
				Role:      auth.RoleKurir,
				Address:   "Jl. Logistik Hijau No. 5, Bandung",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        mitraID,
				Name:      "PT Hijau Lestari Indonesia (CSR)",
				Email:     "mitra@banksampah.id",
				Password:  string(hashedPassMitra),
				Phone:     "081122334455",
				Role:      auth.RoleMitra,
				Address:   "Gedung Menara Hijau Lt. 12, Jakarta",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		for _, u := range users {
			if err := db.Create(&u).Error; err != nil {
				log.Printf("⚠ Warning seeding user %s: %v", u.Email, err)
			}
		}
		log.Println("✓ Seeded default users: admin@banksampah.id, nasabah@banksampah.id, kurir@banksampah.id, mitra@banksampah.id")
	}

	// 2. Seed Member Card for Nasabah
	var cardCount int64
	db.Model(&member_card.MemberCard{}).Where("user_id = ?", nasabahID).Count(&cardCount)
	if cardCount == 0 {
		card := member_card.MemberCard{
			ID:          uuid.New(),
			UserID:      nasabahID,
			CardNumber:  "BS-260101-0001",
			QRCodeToken: "BS-QR-BUDI-2026",
			IsActive:    true,
			IssuedAt:    time.Now(),
			CreatedAt:   time.Now(),
		}
		_ = db.Create(&card)
		log.Println("✓ Seeded digital member card for Budi Santoso (QR: BS-QR-BUDI-2026)")
	}

	// 3. Seed Wallet for Nasabah (Initial zero balance)
	var walletCount int64
	db.Model(&wallet.Wallet{}).Where("user_id = ?", nasabahID).Count(&walletCount)
	if walletCount == 0 {
		w := wallet.Wallet{
			ID:        uuid.New(),
			UserID:    nasabahID,
			Balance:   0,
			Points:    0,
			CreatedAt: time.Now(),
		}
		_ = db.Create(&w)
	}

	// 4. Seed Inventory Stock Categories
	var stockCount int64
	db.Model(&inventory.WasteStock{}).Count(&stockCount)
	if stockCount == 0 {
		stocks := []inventory.WasteStock{
			{ID: uuid.New(), Category: weighbridge.CategoryPlastic, CurrentKg: 150.0, AvgCostPerKg: 4000.0, LastUpdated: time.Now()},
			{ID: uuid.New(), Category: weighbridge.CategoryPaper, CurrentKg: 200.0, AvgCostPerKg: 2500.0, LastUpdated: time.Now()},
			{ID: uuid.New(), Category: weighbridge.CategoryMetal, CurrentKg: 80.0, AvgCostPerKg: 8500.0, LastUpdated: time.Now()},
			{ID: uuid.New(), Category: weighbridge.CategoryGlass, CurrentKg: 120.0, AvgCostPerKg: 1200.0, LastUpdated: time.Now()},
			{ID: uuid.New(), Category: weighbridge.CategoryOrganic, CurrentKg: 300.0, AvgCostPerKg: 500.0, LastUpdated: time.Now()},
		}
		for _, s := range stocks {
			_ = db.Create(&s)
		}
		log.Println("✓ Seeded inventory stocks (Plastic, Paper, Metal, Glass, Organic)")
	}

	// 5. Seed Drop Points (Posko RT/RW)
	var dpCount int64
	db.Model(&drop_point.DropPoint{}).Count(&dpCount)
	if dpCount == 0 {
		dps := []drop_point.DropPoint{
			{
				ID:             uuid.New(),
				Name:           "Posko Bank Sampah Unit RW 05 Coblong",
				Address:        "Balai Warga RW 05, Jl. Dago No. 120",
				RT:             "03",
				RW:             "05",
				District:       "Coblong",
				City:           "Bandung",
				Latitude:       -6.88500000,
				Longitude:      107.61400000,
				OperatingHours: "Senin - Sabtu: 08:00 - 15:00",
				ContactPerson:  "Pak RT Mulyadi",
				ContactPhone:   "081234111222",
				IsActive:       true,
			},
			{
				ID:             uuid.New(),
				Name:           "Smart Drop Box Sukajadi Mall",
				Address:        "Area Parkir Timur Sukajadi Mall, Jl. Sukajadi No. 45",
				District:       "Sukajadi",
				City:           "Bandung",
				Latitude:       -6.89000000,
				Longitude:      107.59800000,
				OperatingHours: "Setiap Hari: 24 Jam",
				ContactPerson:  "CS Smart Drop Box",
				ContactPhone:   "081234999888",
				IsActive:       true,
			},
		}
		for _, d := range dps {
			_ = db.Create(&d)
		}

		// Seed Mobile Schedule
		sched := drop_point.MobilePickupSchedule{
			ID:            uuid.New(),
			DropPointID:   dps[0].ID,
			DropPointName: dps[0].Name,
			DayOfWeek:     "Senin",
			StartTime:     "09:00",
			EndTime:       "12:00",
			CourierName:   "Joko Supriyanto",
			VehiclePlate:  "D 1234 BS",
		}
		_ = db.Create(&sched)
		log.Println("✓ Seeded drop points & mobile pickup schedules")
	}

	// 6. Seed Eco Marketplace Products
	var prodCount int64
	db.Model(&marketplace.EcoProduct{}).Count(&prodCount)
	if prodCount == 0 {
		products := []marketplace.EcoProduct{
			{
				ID:          uuid.New(),
				Name:        "Pupuk Kompos Organik Super (10 Kg)",
				Description: "Pupuk kompos alami hasil fermentasi sampah organik warga, kaya hara untuk tanaman hias dan sayur.",
				Category:    marketplace.EcoCategoryKompos,
				PriceIDR:    35000,
				Stock:       50,
				Unit:        "karung",
				ImageURL:    "/products/kompos.jpg",
				IsActive:    true,
			},
			{
				ID:          uuid.New(),
				Name:        "Tas Belanja Ramah Lingkungan (Tote Bag)",
				Description: "Anyaman tas belanja kuat & modis dibuat 100% dari daur ulang limbah plastik kresek oleh ibu-ibu PKK.",
				Category:    marketplace.EcoCategoryKerajinan,
				PriceIDR:    25000,
				Stock:       30,
				Unit:        "pcs",
				ImageURL:    "/products/totebag.jpg",
				IsActive:    true,
			},
			{
				ID:          uuid.New(),
				Name:        "Ecobrick Padat Standard (500 ml)",
				Description: "Botol plastik dipadatkan dengan limbah plastik kering siap untuk bahan bangunan ramah lingkungan / meja taman.",
				Category:    marketplace.EcoCategoryEcobrick,
				PriceIDR:    15000,
				Stock:       100,
				Unit:        "botol",
				ImageURL:    "/products/ecobrick.jpg",
				IsActive:    true,
			},
			{
				ID:          uuid.New(),
				Name:        "Pelet Biji Plastik HDPE Daur Ulang",
				Description: "Biji plastik murni hasil pencucian dan peleburan botol susu & tutup botol, siap cetak industri.",
				Category:    marketplace.EcoCategoryBijiPlastik,
				PriceIDR:    14000,
				Stock:       500,
				Unit:        "kg",
				ImageURL:    "/products/hdpe.jpg",
				IsActive:    true,
			},
		}
		for _, p := range products {
			_ = db.Create(&p)
		}
		log.Println("✓ Seeded Eco Marketplace products")
	}

	// 7. Seed Rewards Catalog
	var rewardCount int64
	db.Model(&reward.RewardItem{}).Count(&rewardCount)
	if rewardCount == 0 {
		rewards := []reward.RewardItem{
			{
				ID:             uuid.New(),
				Title:          "Voucher Belanja Minimarket Rp 25.000",
				Description:    "Tukar 2.500 poin untuk voucher belanja Indomaret/Alfamart",
				PointsRequired: 2500,
				Stock:          50,
				ImageURL:       "/rewards/voucher.jpg",
				IsActive:       true,
			},
			{
				ID:             uuid.New(),
				Title:          "Tumbler Stainless Ramah Lingkungan",
				Description:    "Tukar 1.500 poin untuk tumbler tahan panas/dingin 500ml",
				PointsRequired: 1500,
				Stock:          20,
				ImageURL:       "/rewards/tumbler.jpg",
				IsActive:       true,
			},
			{
				ID:             uuid.New(),
				Title:          "Paket Bibit Sayur Organik (Cabai, Tomat, Bayam)",
				Description:    "Tukar 500 poin untuk paket berkebun hidroponik di rumah",
				PointsRequired: 500,
				Stock:          100,
				ImageURL:       "/rewards/bibit.jpg",
				IsActive:       true,
			},
		}
		for _, r := range rewards {
			_ = db.Create(&r)
		}
		log.Println("✓ Seeded rewards catalog items")
	}

	log.Println("✨ All database seeders completed successfully!")
	return nil
}
