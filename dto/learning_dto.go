package dto

import "citra-wastra-be/models"

type LevelWithProgress struct {
	models.Level
	IsCompleted bool `json:"is_completed"`
	IsLocked    bool `json:"is_locked"`
}

type ModuleWithProgress struct {
	models.Module
	LevelsWithProgress []LevelWithProgress `json:"levels_with_progress"`
}
