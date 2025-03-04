package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/marbh56/mordezzan/internal/currency"
	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/logger"
	"go.uber.org/zap"
)

// HandleCurrencyUpdate handles updating a character's currency
func (s *Server) HandleCurrencyUpdate(w http.ResponseWriter, r *http.Request) {
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
		renderCurrencyError(w, "Failed to parse form")
		return
	}

	// Get form values
	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err))
		renderCurrencyError(w, "Invalid character ID")
		return
	}

	amount, err := strconv.ParseInt(r.Form.Get("amount"), 10, 64)
	if err != nil {
		logger.Warn("Invalid currency amount", zap.Error(err))
		renderCurrencyError(w, "Invalid amount")
		return
	}

	denomination := r.Form.Get("denomination")
	if !isValidDenomination(denomination) {
		logger.Warn("Invalid currency denomination", zap.String("denomination", denomination))
		renderCurrencyError(w, "Invalid denomination")
		return
	}

	// Fetch character
	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Error fetching character", zap.Error(err))
		renderCurrencyError(w, "Character not found")
		return
	}

	// Create a purse from the character's current currency
	purse := currency.Purse{
		PlatinumPieces: character.PlatinumPieces,
		GoldPieces:     character.GoldPieces,
		ElectrumPieces: character.ElectrumPieces,
		SilverPieces:   character.SilverPieces,
		CopperPieces:   character.CopperPieces,
	}

	// Record old values for logging
	oldPurse := purse

	// Update the currency based on denomination
	var success bool
	if amount >= 0 {
		// Adding currency
		currency.AddToPurse(&purse, amount, currency.Denomination(denomination))
		success = true
	} else {
		// Removing currency (convert to positive for removal)
		success = currency.RemoveFromPurse(&purse, -amount, currency.Denomination(denomination))
		if !success {
			renderCurrencyError(w, fmt.Sprintf("Not enough %s", getDenominationName(denomination)))
			return
		}
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
		ExperiencePoints: character.ExperiencePoints,
		PlatinumPieces:   purse.PlatinumPieces,
		GoldPieces:       purse.GoldPieces,
		ElectrumPieces:   purse.ElectrumPieces,
		SilverPieces:     purse.SilverPieces,
		CopperPieces:     purse.CopperPieces,
	}

	// Update character in database
	updatedChar, err := queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		logger.Error("Failed to update character currency", zap.Error(err))
		renderCurrencyError(w, "Error updating currency")
		return
	}

	logger.Info("Character currency updated",
		zap.Int64("character_id", characterID),
		zap.String("denomination", denomination),
		zap.Int64("old_platinum", oldPurse.PlatinumPieces),
		zap.Int64("new_platinum", purse.PlatinumPieces),
		zap.Int64("old_gold", oldPurse.GoldPieces),
		zap.Int64("new_gold", purse.GoldPieces),
		zap.Int64("old_electrum", oldPurse.ElectrumPieces),
		zap.Int64("new_electrum", purse.ElectrumPieces),
		zap.Int64("old_silver", oldPurse.SilverPieces),
		zap.Int64("new_silver", purse.SilverPieces),
		zap.Int64("old_copper", oldPurse.CopperPieces),
		zap.Int64("new_copper", purse.CopperPieces),
		zap.Int64("change", amount))

	// Fetch inventory for coin weight calculation
	inventory, err := queries.GetCharacterInventoryItems(r.Context(), characterID)
	if err != nil {
		logger.Warn("Failed to fetch inventory for coin weight",
			zap.Error(err),
			zap.Int64("character_id", characterID))
		// Continue anyway, we'll just show currency without weight
	}

	// Create view model for template
	viewModel := NewSafeCharacterViewModel(updatedChar, inventory)

	// Add message with proper currency name
	var message string
	if amount > 0 {
		message = fmt.Sprintf("Added %d %s", amount, getDenominationName(denomination))
	} else {
		message = fmt.Sprintf("Removed %d %s", -amount, getDenominationName(denomination))
	}

	// Render the updated currency section
	renderCurrencySectionUpdate(w, viewModel, message)
}

// Helper to render currency error response
func renderCurrencyError(w http.ResponseWriter, errMsg string) {
	w.Header().Set("HX-Retarget", "#currency-section")
	w.Header().Set("HX-Reswap", "outerHTML")
	w.WriteHeader(http.StatusBadRequest)

	// We need to include both Error and Character fields to match the template structure
	data := struct {
		Error     string
		Character CharacterViewModel // Empty struct for template compatibility
		Message   string             // Empty string for template compatibility
	}{
		Error:     errMsg,
		Character: CharacterViewModel{}, // Empty character model
		Message:   "",
	}

	RenderTemplate(w, "templates/characters/_currency_section.html", "_currency_section", data)
}

// Helper to render updated currency section
func renderCurrencySectionUpdate(w http.ResponseWriter, character CharacterViewModel, message string) {
	data := struct {
		Character CharacterViewModel
		Message   string
		Error     string
	}{
		Character: character,
		Message:   message,
		Error:     "",
	}

	RenderTemplate(w, "templates/characters/_currency_section.html", "_currency_section", data)
}

// Helper to get denomination full name
func getDenominationName(denom string) string {
	switch denom {
	case "pp":
		return "platinum pieces"
	case "gp":
		return "gold pieces"
	case "ep":
		return "electrum pieces"
	case "sp":
		return "silver pieces"
	case "cp":
		return "copper pieces"
	default:
		return "coins"
	}
}

// isValidDenomination checks if the denomination is valid
func isValidDenomination(denom string) bool {
	validDenoms := map[string]bool{
		"pp": true,
		"gp": true,
		"ep": true,
		"sp": true,
		"cp": true,
	}
	return validDenoms[denom]
}

// Helper function to calculate coin weight
func calculateCoinWeight(pp, gp, ep, sp, cp int64) int {
	// Count total coins
	totalCoins := pp + gp + ep + sp + cp

	// Using 50 coins per pound which is a common RPG standard
	coinsPerPound := 50.0

	// Calculate weight with proper rounding
	weight := float64(totalCoins) / coinsPerPound

	// Round to nearest pound, but ensure at least 1 pound if there are any coins
	if totalCoins > 0 && weight < 0.5 {
		return 1
	}

	return int(weight + 0.5) // Round to nearest integer
}
