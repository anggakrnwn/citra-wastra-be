package database

import (
	"citra-wastra-be/models"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	SeedBadges(db)
}

func SeedBadges(db *gorm.DB) {
	badges := []models.Badge{
		{
			ID:          "4336f332-6758-4e80-b0bc-7a33a39e803a",
			Name:        "Canting Perdana",
			Description: "Titik awal perjalananmu mengenal wastra nusantara.",
			ImageURL:    "",
			MinXP:       0,
		},
		{
			ID:          "7ba15694-8703-4c90-9515-5a7c293673f3",
			Name:        "Pencari Pola",
			Description: "Kamu mulai memahami garis dan bentuk di setiap kain.",
			ImageURL:    "",
			MinXP:       500,
		},
		{
			ID:          "12d6a592-236b-4e6b-876e-82635928f642",
			Name:        "Wastra Wiratama",
			Description: "Ksatria budaya yang aktif melestarikan kearifan lokal.",
			ImageURL:    "",
			MinXP:       1500,
		},
		{
			ID:          "9e4f2a1a-4d3b-4c2a-bd1e-5f6a7b8c9d0e",
			Name:        "Penjaga Warisan",
			Description: "Dedikasimu menjaga makna filosofi batik sangat luar biasa.",
			ImageURL:    "",
			MinXP:       3500,
		},
		{
			ID:          "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			Name:        "Empu Citra-Wastra",
			Description: "Gelar kehormatan tertinggi bagi sang ahli wastra sejati.",
			ImageURL:    "",
			MinXP:       7000,
		},
	}

	for _, b := range badges {
		db.FirstOrCreate(&b, models.Badge{
			ID: b.ID,
		})
	}
}
