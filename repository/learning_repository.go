package repository

import (
	"citra-wastra-be/models"

	"gorm.io/gorm"
)

type LearningRepository interface {
	GetIslands() ([]models.Island, error)
	GetModulesByIsland(islandID string) ([]models.Module, error)
	CreateModule(module *models.Module) error
	GetLevelsByModule(moduleID string) ([]models.Level, error)
	GetLevelByID(levelID string) (models.Level, error)
	CreateLevel(level *models.Level) error
	UpdateLevel(level *models.Level) error
	DeleteLevel(levelID string) error
	SaveProgress(progress *models.UserLevelProgress) error
	GetProgress(userID string, levelID string) (*models.UserLevelProgress, error)
	GetQuestionsByLevel(levelID string) ([]models.Question, error)
	CreateQuestion(q *models.Question) error
	UpdateQuestion(q *models.Question) error
	DeleteQuestion(id string) error
}

type learningRepository struct {
	db *gorm.DB
}

func NewLearningRepository(db *gorm.DB) LearningRepository {
	return &learningRepository{db}
}

func (r *learningRepository) GetIslands() ([]models.Island, error) {
	var islands []models.Island
	err := r.db.Preload("Modules").Find(&islands).Error
	return islands, err
}

func (r *learningRepository) GetModulesByIsland(islandID string) ([]models.Module, error) {
	var modules []models.Module
	err := r.db.Where("island_id = ?", islandID).Preload("Levels").Find(&modules).Error
	return modules, err
}

func (r *learningRepository) CreateModule(module *models.Module) error {
	return r.db.Create(module).Error
}

func (r *learningRepository) GetLevelsByModule(moduleID string) ([]models.Level, error) {
	var levels []models.Level
	err := r.db.Where("module_id = ?", moduleID).Order("\"order\" ASC").Find(&levels).Error
	return levels, err
}

func (r *learningRepository) GetLevelByID(levelID string) (models.Level, error) {
	var level models.Level
	err := r.db.First(&level, "id = ?", levelID).Error
	return level, err
}

func (r *learningRepository) CreateLevel(level *models.Level) error {
	return r.db.Create(level).Error
}

func (r *learningRepository) UpdateLevel(level *models.Level) error {
	return r.db.Save(level).Error
}

func (r *learningRepository) DeleteLevel(levelID string) error {
	return r.db.Delete(&models.Level{}, "id = ?", levelID).Error
}

func (r *learningRepository) SaveProgress(progress *models.UserLevelProgress) error {
	return r.db.FirstOrCreate(progress, models.UserLevelProgress{
		UserID:  progress.UserID,
		LevelID: progress.LevelID,
	}).Error
}

func (r *learningRepository) GetProgress(userID string, levelID string) (*models.UserLevelProgress, error) {
	var progress models.UserLevelProgress
	err := r.db.Where("user_id = ? AND level_id = ?", userID, levelID).First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *learningRepository) GetQuestionsByLevel(levelID string) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Where("level_id = ?", levelID).Find(&questions).Error
	return questions, err
}

func (r *learningRepository) CreateQuestion(q *models.Question) error {
	return r.db.Create(q).Error
}

func (r *learningRepository) UpdateQuestion(q *models.Question) error {
	return r.db.Save(q).Error
}

func (r *learningRepository) DeleteQuestion(id string) error {
	return r.db.Delete(&models.Question{}, "id = ?", id).Error
}
