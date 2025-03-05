package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/logger"
	charRules "github.com/marbh56/mordezzan/internal/rules/character"
	"go.uber.org/zap"
)

// HandleXPUpdate handles updating a character's experience points
func (s *Server) HandleXPUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Error("Failed to parse form", zap.Error(err))
		renderXPError(w, "Failed to parse form")
		return
	}

	// Get form values
	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err))
		renderXPError(w, "Invalid character ID")
		return
	}

	xpChange, err := strconv.ParseInt(r.Form.Get("xp_change"), 10, 64)
	if err != nil {
		logger.Warn("Invalid XP change value", zap.Error(err))
		renderXPError(w, "Invalid XP value")
		return
	}

	calculateBonus := r.Form.Get("calculate_bonus") == "1"

	// Fetch character
	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Error fetching character", zap.Error(err))
		renderXPError(w, "Character not found")
		return
	}

	// Apply XP bonuses based on ability scores if requested
	finalXPChange := xpChange
	bonusMessage := ""
	if calculateBonus && xpChange > 0 {
		xpBonus := calculateXPBonus(character.Class, character)
		if xpBonus > 0 {
			bonusXP := (xpChange * xpBonus) / 100
			finalXPChange = xpChange + bonusXP
			bonusMessage = fmt.Sprintf(" (includes %d%% bonus: +%d XP)", xpBonus, bonusXP)
		}
	}

	// Calculate new XP total
	newXP := character.ExperiencePoints + finalXPChange
	if newXP < 0 {
		newXP = 0 // Prevent negative XP
	}

	// Create update params starting with current values
	updateParams := db.UpdateCharacterParams{
		ID:               characterID,
		UserID:           user.UserID,
		Name:             character.Name,
		Class:            character.Class,
		Level:            character.Level,
		MaxHp:            character.MaxHp,
		CurrentHp:        character.CurrentHp,
		Strength:         character.Strength,
		Dexterity:        character.Dexterity,
		Constitution:     character.Constitution,
		Intelligence:     character.Intelligence,
		Wisdom:           character.Wisdom,
		Charisma:         character.Charisma,
		ExperiencePoints: newXP,
		PlatinumPieces:   character.PlatinumPieces,
		GoldPieces:       character.GoldPieces,
		ElectrumPieces:   character.ElectrumPieces,
		SilverPieces:     character.SilverPieces,
		CopperPieces:     character.CopperPieces,
	}

	// Update character in database
	updatedChar, err := queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		logger.Error("Failed to update character XP", zap.Error(err))
		renderXPError(w, "Error updating XP")
		return
	}

	logger.Info("Character XP updated",
		zap.Int64("character_id", characterID),
		zap.Int64("old_xp", character.ExperiencePoints),
		zap.Int64("new_xp", newXP),
		zap.Int64("change", finalXPChange),
		zap.Bool("bonus_applied", calculateBonus))

	// Check if character level should change based on new XP
	progression := charRules.GetClassProgression(character.Class)
	newLevel := progression.GetLevelForXP(newXP)
	levelMessage := ""

	if newLevel != character.Level {
		// Update character level in a separate update
		levelUpdateParams := updateParams
		levelUpdateParams.Level = newLevel

		_, err := queries.UpdateCharacter(r.Context(), levelUpdateParams)
		if err != nil {
			logger.Error("Failed to update character level", zap.Error(err))
			// Continue anyway, we'll just show the XP update
		} else {
			levelMessage = fmt.Sprintf(" Level increased to %d!", newLevel)
			logger.Info("Character level increased",
				zap.Int64("character_id", characterID),
				zap.Int64("old_level", character.Level),
				zap.Int64("new_level", newLevel))
		}
	}

	// Fetch inventory for view model creation
	inventory, err := queries.GetCharacterInventoryItems(r.Context(), characterID)
	if err != nil {
		logger.Warn("Failed to fetch inventory for XP update",
			zap.Error(err),
			zap.Int64("character_id", characterID))
		// Continue anyway, we'll just show the XP update with empty inventory
		inventory = []db.GetCharacterInventoryItemsRow{}
	}

	// Create view model for template
	viewModel := NewSafeCharacterViewModel(updatedChar, inventory)

	// Add message based on XP change
	var message string
	if finalXPChange > 0 {
		message = fmt.Sprintf("Added %d XP%s%s", finalXPChange, bonusMessage, levelMessage)
	} else {
		message = fmt.Sprintf("Removed %d XP", -finalXPChange)
	}

	// Render the updated XP section
	renderXPSection(w, viewModel, message)
}

// Helper function to render XP errors
func renderXPError(w http.ResponseWriter, errMsg string) {
	w.Header().Set("HX-Retarget", "#xp-section")
	w.Header().Set("HX-Reswap", "outerHTML")
	w.WriteHeader(http.StatusBadRequest)

	data := struct {
		Error     string
		Character CharacterViewModel
		Message   string
	}{
		Error:     errMsg,
		Character: CharacterViewModel{},
		Message:   "",
	}

	RenderTemplate(w, "templates/characters/_xp_section.html", "_xp_section", data)
}

// Helper function to render the updated XP section
func renderXPSection(w http.ResponseWriter, character CharacterViewModel, message string) {
	data := struct {
		Character CharacterViewModel
		Message   string
		Error     string
	}{
		Character: character,
		Message:   message,
		Error:     "",
	}

	RenderTemplate(w, "templates/characters/_xp_section.html", "_xp_section", data)
}

// calculateXPBonus determines the XP bonus percentage based on primary ability scores
func calculateXPBonus(class string, character db.Character) int64 {
	switch class {
	case "Fighter", "Barbarian", "Ranger", "Paladin":
		// Strength-based classes
		if character.Strength >= 16 {
			return 10
		} else if character.Strength >= 13 {
			return 5
		}
	case "Thief", "Assassin":
		// Dexterity-based classes
		if character.Dexterity >= 16 {
			return 10
		} else if character.Dexterity >= 13 {
			return 5
		}
	case "Magician", "Necromancer":
		// Intelligence-based classes
		if character.Intelligence >= 16 {
			return 10
		} else if character.Intelligence >= 13 {
			return 5
		}
	case "Cleric", "Druid":
		// Wisdom-based classes
		if character.Wisdom >= 16 {
			return 10
		} else if character.Wisdom >= 13 {
			return 5
		}
	}

	return 0
}
