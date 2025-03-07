package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/logger"
	"go.uber.org/zap"
)

func (s *Server) HandleUseMagicalItem(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	characterID, err := strconv.ParseInt(r.FormValue("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid item ID", zap.Error(err))
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Verify character belongs to user
	queries := db.New(s.db)
	_, err = queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Character not found or doesn't belong to user",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("user_id", user.UserID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Get item details to display appropriate message
	item, err := queries.GetChargedItem(r.Context(), db.GetChargedItemParams{
		ID:          itemID,
		CharacterID: characterID,
	})
	if err != nil {
		logger.Error("Failed to fetch charged item details",
			zap.Error(err),
			zap.Int64("item_id", itemID),
			zap.Int64("character_id", characterID))
		http.Error(w, "Item not found or is not a charged item", http.StatusNotFound)
		return
	}

	// Check if item has charges left
	if item.Charges.Valid && item.Charges.Int64 <= 0 {
		logger.Warn("Attempt to use item with no charges",
			zap.Int64("item_id", itemID),
			zap.Int64("character_id", characterID))
		renderInventoryWithMessage(w, r, characterID, "This item has no charges remaining")
		return
	}

	// Use the item by reducing charges - the parameter is a SQL query, not a Go struct
	err = queries.UseChargedItem(r.Context(), db.UseChargedItemParams{
		ID:          itemID,
		CharacterID: characterID,
	})
	if err != nil {
		logger.Error("Failed to use charged item",
			zap.Error(err),
			zap.Int64("item_id", itemID),
			zap.Int64("character_id", characterID))
		renderInventoryWithMessage(w, r, characterID, "Error using item")
		return
	}

	// If it's a one-time use item like a potion and now has 0 charges, remove it
	if item.Category == "potion" || item.Category == "scroll" {
		// Check if charges are now at 0
		updatedItem, err := queries.GetChargedItem(r.Context(), db.GetChargedItemParams{
			ID:          itemID,
			CharacterID: characterID,
		})

		if err == nil && updatedItem.Charges.Valid && updatedItem.Charges.Int64 <= 0 {
			// Remove the item as it's been consumed
			err = queries.RemoveItemFromInventory(r.Context(), db.RemoveItemFromInventoryParams{
				ID:          itemID,
				CharacterID: characterID,
			})
			if err != nil {
				logger.Error("Failed to remove consumed item",
					zap.Error(err),
					zap.Int64("item_id", itemID),
					zap.Int64("character_id", characterID))
				// Continue anyway, we'll just show the item with 0 charges
			}
		}
	}

	message := fmt.Sprintf("You used %s: %s", item.Name, item.EffectDescription)
	var remainingCharges int64 = 0
	if item.Charges.Valid {
		remainingCharges = item.Charges.Int64 - 1
		if remainingCharges < 0 {
			remainingCharges = 0
		}
	}

	logger.Info("Magical item used",
		zap.Int64("character_id", characterID),
		zap.Int64("item_id", itemID),
		zap.String("item_name", item.Name),
		zap.Int64("charges_remaining", remainingCharges))

	renderInventoryWithMessage(w, r, characterID, message)
}

func renderInventoryWithMessage(w http.ResponseWriter, r *http.Request, characterID int64, message string) {
	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=%s", characterID, message), http.StatusSeeOther)
}
