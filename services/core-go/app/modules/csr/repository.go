package csr

import (
	"context"

	"bank-sampah/core-go/app/modules/auth"
	"bank-sampah/core-go/app/modules/weighbridge"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateProgram(ctx context.Context, program *CsrProgram) error
	FindProgramByID(ctx context.Context, id uuid.UUID) (*CsrProgram, error)
	ListPrograms(ctx context.Context) ([]CsrProgram, error)
	CalculateImpact(ctx context.Context) (*ImpactSummary, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateProgram(ctx context.Context, program *CsrProgram) error {
	return r.db.WithContext(ctx).Create(program).Error
}

func (r *repository) FindProgramByID(ctx context.Context, id uuid.UUID) (*CsrProgram, error) {
	var program CsrProgram
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&program).Error
	if err != nil {
		return nil, err
	}
	return &program, nil
}

func (r *repository) ListPrograms(ctx context.Context) ([]CsrProgram, error) {
	var programs []CsrProgram
	err := r.db.WithContext(ctx).Order("start_date DESC").Find(&programs).Error
	return programs, err
}

func (r *repository) CalculateImpact(ctx context.Context) (*ImpactSummary, error) {
	summary := &ImpactSummary{}

	// Aggregations from scale tickets
	type WasteStats struct {
		Category weighbridge.WasteCategory
		TotalKg  float64
		TotalRp  float64
	}
	var stats []WasteStats

	r.db.WithContext(ctx).Model(&weighbridge.ScaleTicket{}).
		Select("category, SUM(net_weight) as total_kg, SUM(total_price) as total_rp").
		Group("category").
		Scan(&stats)

	for _, s := range stats {
		summary.TotalWasteKg += s.TotalKg
		summary.TotalFundsDistributed += s.TotalRp

		switch s.Category {
		case weighbridge.CategoryPlastic:
			summary.PlasticWasteKg = s.TotalKg
		case weighbridge.CategoryPaper:
			summary.PaperWasteKg = s.TotalKg
		case weighbridge.CategoryMetal:
			summary.MetalWasteKg = s.TotalKg
		}
	}

	// ESG Environmental Formulas:
	// 1 kg plastic recycled ~ 1.5 kg CO2 avoided
	// 50 kg paper recycled ~ 1 tree saved
	summary.Co2SavedKg = (summary.PlasticWasteKg * 1.5) + (summary.MetalWasteKg * 2.0)
	summary.TreesEquivalent = summary.PaperWasteKg / 50.0

	var nasabahCount int64
	r.db.WithContext(ctx).Model(&auth.User{}).Where("role = ?", auth.RoleNasabah).Count(&nasabahCount)
	summary.TotalNasabahHelped = nasabahCount

	return summary, nil
}
