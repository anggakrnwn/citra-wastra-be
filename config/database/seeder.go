package database

import (
	"citra-wastra-be/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	log.Println("running database seeder...")
	SeedAdmin(db)
	SeedBadges(db)
	SeedLearningPath(db)
	SeedQuestions(db)
	log.Println("seeder completed.")
}

func SeedAdmin(db *gorm.DB) {
	adminPass := os.Getenv("ADMIN_PASSWORD")
	superPass := os.Getenv("SUPER_ADMIN_PASSWORD")

	if adminPass == "" || superPass == "" {
		log.Println("warning: ADMIN_PASSWORD or SUPER_ADMIN_PASSWORD not set. skipping admin seeding.")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPass), 10)
	admin := models.User{
		ID:       "00000000-0000-0000-0000-000000000001",
		Username: "admin",
		Email:    "admin@citrawastra.com",
		Password: string(hashedPassword),
		Role:     "admin",
		IsActive: true,
	}
	if err := db.Save(&admin).Error; err != nil {
		log.Printf("failed to seed admin: %v", err)
	} else {
		log.Printf("admin user ready: %s", admin.Email)
	}

	hashedSuperPassword, _ := bcrypt.GenerateFromPassword([]byte(superPass), 10)
	superAdmin := models.User{
		ID:       "00000000-0000-0000-0000-000000000000",
		Username: "superadmin",
		Email:    "superadmin@citrawastra.com",
		Password: string(hashedSuperPassword),
		Role:     "super_admin",
		IsActive: true,
	}
	if err := db.Save(&superAdmin).Error; err != nil {
		log.Printf("failed to seed super admin: %v", err)
	} else {
		log.Printf("super admin user ready: %s", superAdmin.Email)
	}
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

func SeedLearningPath(db *gorm.DB) {
	islands := []models.Island{
		{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "Jawa", Description: "Pusat wastra dengan filosofi keraton dan pesisiran."},
		{ID: "550e8400-e29b-41d4-a716-446655440001", Name: "Bali", Description: "Wastra yang kental dengan nilai spiritual dan tradisi Hindu-Bali."},
		{ID: "550e8400-e29b-41d4-a716-446655440002", Name: "Sumatera", Description: "Perpaduan motif akulturasi budaya Melayu dan Islam."},
		{ID: "550e8400-e29b-41d4-a716-446655440003", Name: "Jakarta", Description: "Batik Betawi dengan warna cerah dan motif ikon kota."},
		{ID: "550e8400-e29b-41d4-a716-446655440004", Name: "Maluku", Description: "Wastra kepulauan yang menonjolkan kekayaan alam rempah."},
	}

	for _, is := range islands {
		db.FirstOrCreate(&is, models.Island{ID: is.ID})
	}

	modules := []models.Module{
		{ID: "a1b2c3d4-e5f6-4a5b-b6c7-d8e9f0a1b2c3", IslandID: "550e8400-e29b-41d4-a716-446655440000", Name: "Batik Gentongan", Description: "Batik khas Madura (Jawa Timur) yang dicelup dalam genting tanah liat."},
		{ID: "b2c3d4e5-f6a7-4b6c-c7d8-e9f0a1b2c3d4", IslandID: "550e8400-e29b-41d4-a716-446655440001", Name: "Batik Barong", Description: "Visualisasi makhluk mitologi pelindung kebaikan di Bali."},
		{ID: "c3d4e5f6-a7b8-4c7d-d8e9-f0a1b2c3d4e5", IslandID: "550e8400-e29b-41d4-a716-446655440002", Name: "Pintu Aceh", Description: "Simbol keramahan dan kepribadian rakyat Aceh yang kuat."},
		{ID: "d4e5f6a7-b8c9-4d8e-e9f0-a1b2c3d4e5f6", IslandID: "550e8400-e29b-41d4-a716-446655440003", Name: "Ondel-Ondel", Description: "Ikon budaya Betawi penolak bala yang ceria."},
		{ID: "e5f6a7b8-c9d0-4e9f-f0a1-b2c3d4e5f6a7", IslandID: "550e8400-e29b-41d4-a716-446655440000", Name: "Batik Kawung", Description: "Motif Yogyakarta yang melambangkan kesucian dan keadilan."},
		{ID: "f6a7b8c9-d0e1-4f0a-a1b2-c3d4e5f6a7b8", IslandID: "550e8400-e29b-41d4-a716-446655440004", Name: "Batik Pala", Description: "Representasi kejayaan Maluku sebagai kepulauan rempah dunia."},
		{ID: "a7b8c9d0-e1f2-4a1b-b2c3-d4e5f6a7b8c9", IslandID: "550e8400-e29b-41d4-a716-446655440000", Name: "Batik Parang", Description: "Motif Yogyakarta yang melambangkan jalinan yang tidak pernah putus."},
	}

	for _, m := range modules {
		db.FirstOrCreate(&m, models.Module{ID: m.ID})
	}

	levels := []models.Level{
		{ID: "11111111-1111-4111-a111-111111111111", ModuleID: "a1b2c3d4-e5f6-4a5b-b6c7-d8e9f0a1b2c3", Order: 1, Title: "Seni Celup Gentong", Content: "Keunikan Gentongan terletak pada proses pewarnaan yang merendam kain di dalam gentong tanah liat selama berminggu-minggu.", XPReward: 100},
		{ID: "22222222-2222-4222-a222-222222222222", ModuleID: "b2c3d4e5-f6a7-4b6c-c7d8-e9f0a1b2c3d4", Order: 1, Title: "Pelindung Dharma", Content: "Motif Barong melambangkan kemenangan kebajikan. Sering digunakan dalam upacara adat sebagai simbol kesucian.", XPReward: 100},
		{ID: "33333333-3333-4333-a333-333333333333", ModuleID: "c3d4e5f6-a7b8-4c7d-d8e9-f0a1b2c3d4e5", Order: 1, Title: "Simbol Keramahan", Content: "Motif ini terinspirasi dari bentuk pintu rumah adat Aceh yang rendah, melambangkan kesantunan dan rendah hati.", XPReward: 100},
		{ID: "44444444-4444-4444-a444-444444444444", ModuleID: "d4e5f6a7-b8c9-4d8e-e9f0-a1b2c3d4e5f6", Order: 1, Title: "Warna Betawi", Content: "Batik Ondel-ondel menggunakan warna-warna berani seperti merah dan pucuk rebung untuk menunjukkan semangat masyarakat Jakarta.", XPReward: 100},
		{ID: "55555555-5555-4555-a555-555555555555", ModuleID: "e5f6a7b8-c9d0-4e9f-f0a1-b2c3d4e5f6a7", Order: 1, Title: "Filosofi Empat Arah", Content: "Bentuk empat lingkaran dengan pusat di tengah melambangkan 'Sedulur Papat Limo Pancer' atau keseimbangan hidup.", XPReward: 100},
		{ID: "66666666-6666-4666-a666-666666666666", ModuleID: "f6a7b8c9-d0e1-4f0a-a1b2-c3d4e5f6a7b8", Order: 1, Title: "Emas Hitam Maluku", Content: "Buah Pala adalah identitas Maluku. Motif ini menggambarkan sejarah jalur rempah yang mengubah dunia.", XPReward: 100},
		{ID: "77777777-7777-4777-a777-777777777777", ModuleID: "a7b8c9d0-e1f2-4a1b-b2c3-d4e5f6a7b8c9", Order: 1, Title: "Semangat Pantang Menyerah", Content: "Garis miring diagonal pada Parang melambangkan ombak laut yang tidak pernah berhenti bergerak, simbol perjuangan hidup.", XPReward: 100},
	}

	for _, l := range levels {
		db.FirstOrCreate(&l, models.Level{ID: l.ID})
	}
}

func SeedQuestions(db *gorm.DB) {
	questions := []models.Question{
		{
			ID:       "11111111-1111-4111-b111-111111111111",
			LevelID:  "11111111-1111-4111-a111-111111111111",
			Question: "Apa wadah unik yang digunakan untuk merendam kain Batik Gentongan?",
			OptionA:  "Ember Plastik",
			OptionB:  "Gentong Tanah Liat",
			OptionC:  "Tangki Besi",
			OptionD:  "Kolam Semen",
			Correct:  "B",
		},
		{
			ID:       "22222222-2222-4222-b222-222222222222",
			LevelID:  "22222222-2222-4222-a222-222222222222",
			Question: "Motif Barong dalam budaya Bali melambangkan simbol apa?",
			OptionA:  "Kekayaan Materi",
			OptionB:  "Kekuatan Militer",
			OptionC:  "Kemenangan Kebajikan (Dharma)",
			OptionD:  "Kesedihan dan Duka",
			Correct:  "C",
		},
		{
			ID:       "33333333-3333-4333-b333-333333333333",
			LevelID:  "33333333-3333-4333-a333-333333333333",
			Question: "Motif Pintu Aceh terinspirasi dari pintu rumah adat yang rendah. Apa maknanya?",
			OptionA:  "Kesantunan dan Rendah Hati",
			OptionB:  "Kekurangan Bahan Bangunan",
			OptionC:  "Pertahanan dari Musuh",
			OptionD:  "Keindahan Arsitektur Saja",
			Correct:  "A",
		},
		{
			ID:       "55555555-5555-4555-b555-555555555555",
			LevelID:  "55555555-5555-4555-a555-555555555555",
			Question: "Filosofi 'Sedulur Papat Limo Pancer' pada motif Kawung melambangkan...",
			OptionA:  "Empat Musim di Dunia",
			OptionB:  "Keseimbangan Hidup Manusia",
			OptionC:  "Jumlah Raja di Jawa",
			OptionD:  "Arah Angin Penunjuk Jalan",
			Correct:  "B",
		},
	}

	for _, q := range questions {
		db.FirstOrCreate(&q, models.Question{ID: q.ID})
	}
}
