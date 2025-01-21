package server

import (
	"time"

	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/rules"
)

type CharacterViewModel struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	MaxHp     int64  `json:"max_hp"`
	CurrentHp int64  `json:"current_hp"`

	// Ability scores with modifiers
	Strength          int64                   `json:"strength"`
	StrengthModifiers rules.StrengthModifiers `json:"strength_modifiers"`

	Dexterity          int64                    `json:"dexterity"`
	DexterityModifiers rules.DexterityModifiers `json:"dexterity_modifiers"`

	Constitution          int64                       `json:"constitution"`
	ConstitutionModifiers rules.ConstitutionModifiers `json:"constitution_modifiers"`

	Intelligence int64 `json:"intelligence"`
	Wisdom       int64 `json:"wisdom"`
	Charisma     int64 `json:"charisma"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCharacterViewModel(c db.Character) CharacterViewModel {
	return CharacterViewModel{
		ID:                    c.ID,
		UserID:                c.UserID,
		Name:                  c.Name,
		MaxHp:                 c.MaxHp,
		CurrentHp:             c.CurrentHp,
		Strength:              c.Strength,
		StrengthModifiers:     rules.CalculateStrengthModifiers(c.Strength),
		Dexterity:             c.Dexterity,
		DexterityModifiers:    rules.CalculateDexterityModifiers(c.Dexterity),
		Constitution:          c.Constitution,
		ConstitutionModifiers: rules.CalculateConstitutionModifiers(c.Constitution),
		Intelligence:          c.Intelligence,
		Wisdom:                c.Wisdom,
		Charisma:              c.Charisma,
		CreatedAt:             c.CreatedAt,
		UpdatedAt:             c.UpdatedAt,
	}
}
