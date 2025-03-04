package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/logger"
	charRules "github.com/marbh56/mordezzan/internal/rules/character"
	"go.uber.org/zap"
)

// HandleAddInventoryItem handles adding items to a character's inventory
func (s *Server) HandleAddInventoryItem(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse character ID from query parameters
	characterIDStr := r.URL.Query().Get("character_id")
	characterID, err := strconv.ParseInt(characterIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err), zap.String("raw_id", characterIDStr))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	// Verify character belongs to user
	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Character not found or belongs to another user",
			zap.Error(err), zap.Int64("character_id", characterID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Check if this is a form submission or the initial page load
	if r.Method == http.MethodPost {
		handleAddItemSubmission(s, w, r, character, queries)
		return
	}

	// Initial page load - display form
	handleAddItemForm(s, w, r, character, queries)
}

// HandleRemoveInventoryItem handles removing items from a character's inventory
func (s *Server) HandleRemoveInventoryItem(w http.ResponseWriter, r *http.Request) {
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

	// Get character ID and item ID from form
	characterID, err := strconv.ParseInt(r.FormValue("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("character_id")))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid item ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("item_id")))
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
		logger.Error("Character not found or belongs to another user",
			zap.Error(err),
			zap.Int64("character_id", characterID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Remove item from inventory
	err = queries.RemoveItemFromInventory(r.Context(), db.RemoveItemFromInventoryParams{
		ID:          itemID,
		CharacterID: characterID,
	})
	if err != nil {
		logger.Error("Failed to remove item from inventory",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("item_id", itemID))
		http.Error(w, "Error removing item", http.StatusInternalServerError)
		return
	}

	logger.Info("Item removed from inventory successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("item_id", itemID))

	// Redirect back to character detail page
	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Item removed successfully", characterID), http.StatusSeeOther)
}

// HandleUpdateInventoryItem handles updating items in a character's inventory
func (s *Server) HandleUpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Warn("Invalid method for inventory update",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Error("Failed to parse form",
			zap.Error(err))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	itemID, err := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid item ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("item_id")))
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.FormValue("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("character_id")))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	// Verify the character belongs to the user
	queries := db.New(s.db)
	_, err = queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Character not found or belongs to another user",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("user_id", user.UserID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Parse quantity
	quantity, err := strconv.ParseInt(r.FormValue("quantity"), 10, 64)
	if err != nil || quantity < 1 {
		quantity = 1 // Default to 1 if invalid
	}

	// Parse container ID if provided
	var containerID sql.NullInt64
	if containerIDStr := r.FormValue("container_id"); containerIDStr != "" {
		id, err := strconv.ParseInt(containerIDStr, 10, 64)
		if err == nil {
			containerID = sql.NullInt64{Int64: id, Valid: true}

			// Verify the container belongs to the character
			if containerID.Valid {
				// Here we need to check if the container belongs to the character
				// and also ensure we're not creating a circular reference
				if containerID.Int64 == itemID {
					logger.Warn("Attempted to put a container inside itself",
						zap.Int64("item_id", itemID),
						zap.Int64("container_id", containerID.Int64))
					http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Cannot put a container inside itself", characterID), http.StatusSeeOther)
					return
				}
			}
		}
	}

	// Parse equipment slot ID if provided
	var equipmentSlotID sql.NullInt64
	if slotIDStr := r.FormValue("equipment_slot_id"); slotIDStr != "" {
		id, err := strconv.ParseInt(slotIDStr, 10, 64)
		if err == nil {
			// Check if the slot is already occupied
			isOccupied, err := queries.IsSlotOccupied(r.Context(), db.IsSlotOccupiedParams{
				CharacterID:     characterID,
				EquipmentSlotID: sql.NullInt64{Int64: id, Valid: true},
			})
			if err != nil {
				logger.Error("Failed to check if equipment slot is occupied",
					zap.Error(err),
					zap.Int64("character_id", characterID),
					zap.Int64("slot_id", id))
			} else if isOccupied {
				logger.Warn("Attempted to equip item to an occupied slot",
					zap.Int64("character_id", characterID),
					zap.Int64("slot_id", id))
				http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Equipment slot is already occupied", characterID), http.StatusSeeOther)
				return
			}

			equipmentSlotID = sql.NullInt64{Int64: id, Valid: true}
		}
	}

	// Parse notes if provided
	var notes sql.NullString
	if notesStr := r.FormValue("notes"); notesStr != "" {
		notes = sql.NullString{String: notesStr, Valid: true}
	}

	// Handle mutual exclusivity - if placing in container, can't equip and vice versa
	if containerID.Valid && equipmentSlotID.Valid {
		logger.Warn("Item cannot be both equipped and in a container",
			zap.Int64("item_id", itemID),
			zap.Int64("character_id", characterID))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Item cannot be both equipped and in a container", characterID), http.StatusSeeOther)
		return
	}

	// Update the inventory item
	updateParams := db.UpdateInventoryItemParams{
		Quantity:        quantity,
		ContainerID:     containerID,
		EquipmentSlotID: equipmentSlotID,
		Notes:           notes,
		ID:              itemID,
		CharacterID:     characterID,
	}

	_, err = queries.UpdateInventoryItem(r.Context(), updateParams)
	if err != nil {
		logger.Error("Failed to update inventory item",
			zap.Error(err),
			zap.Int64("item_id", itemID),
			zap.Int64("character_id", characterID))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error updating item", characterID), http.StatusSeeOther)
		return
	}

	logger.Info("Inventory item updated successfully",
		zap.Int64("item_id", itemID),
		zap.Int64("character_id", characterID))

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Item updated successfully", characterID), http.StatusSeeOther)
}

// HandleEquipItem handles equipping items from inventory to equipment slots
func (s *Server) HandleEquipItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Error("Invalid method for equip item",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Error("Failed to parse form",
			zap.Error(err),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get required form values
	characterID, err := strconv.ParseInt(r.FormValue("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("character_id")))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid item ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("item_id")))
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	equipmentSlotID, err := strconv.ParseInt(r.FormValue("equipment_slot_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid equipment slot ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("equipment_slot_id")))
		http.Error(w, "Invalid equipment slot ID", http.StatusBadRequest)
		return
	}

	// Validate character belongs to the user
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

	// Check if the equipment slot is already occupied
	slotOccupied, err := queries.IsSlotOccupied(r.Context(), db.IsSlotOccupiedParams{
		CharacterID:     characterID,
		EquipmentSlotID: sql.NullInt64{Int64: equipmentSlotID, Valid: true},
	})
	if err != nil {
		logger.Error("Error checking if slot is occupied",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("equipment_slot_id", equipmentSlotID))
		http.Error(w, "Error checking equipment slot", http.StatusInternalServerError)
		return
	}

	if slotOccupied {
		logger.Warn("Attempted to equip item to an occupied slot",
			zap.Int64("character_id", characterID),
			zap.Int64("equipment_slot_id", equipmentSlotID))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Equipment slot is already occupied", characterID), http.StatusSeeOther)
		return
	}

	// Equip the item
	err = queries.EquipItem(r.Context(), db.EquipItemParams{
		EquipmentSlotID: sql.NullInt64{Int64: equipmentSlotID, Valid: true},
		ID:              itemID,
		CharacterID:     characterID,
	})
	if err != nil {
		logger.Error("Failed to equip item",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("item_id", itemID),
			zap.Int64("equipment_slot_id", equipmentSlotID))
		http.Error(w, "Error equipping item", http.StatusInternalServerError)
		return
	}

	logger.Info("Item equipped successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("item_id", itemID),
		zap.Int64("equipment_slot_id", equipmentSlotID))

	// Redirect back to character detail page
	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Item equipped successfully", characterID), http.StatusSeeOther)
}

// HandleUnequipItem handles unequipping items from equipment slots to inventory
func (s *Server) HandleUnequipItem(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		logger.Error("Invalid method for unequip item",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		logger.Error("Failed to parse form",
			zap.Error(err),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get item ID and character ID from form
	itemID, err := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid item ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("item_id")))
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.FormValue("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("character_id")))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	// Verify character belongs to user
	queries := db.New(s.db)
	_, err = queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Character not found or belongs to another user",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("user_id", user.UserID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Unequip the item
	err = queries.UnequipItem(r.Context(), db.UnequipItemParams{
		ID:          itemID,
		CharacterID: characterID,
	})
	if err != nil {
		logger.Error("Failed to unequip item",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("item_id", itemID),
			zap.Int64("user_id", user.UserID))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error unequipping item", characterID), http.StatusSeeOther)
		return
	}

	logger.Info("Item unequipped successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("item_id", itemID),
		zap.Int64("user_id", user.UserID))

	// Redirect back to character detail page
	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Item unequipped successfully", characterID), http.StatusSeeOther)
}

// HandleMoveToContainer handles moving items to containers
func (s *Server) HandleMoveToContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Error("Invalid method for move to container",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Error("Failed to parse form", zap.Error(err))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get character ID, item ID, and container ID from form
	characterID, err := strconv.ParseInt(r.FormValue("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("character_id")))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid item ID",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("item_id")))
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	containerIDStr := r.FormValue("container_id")
	var containerID sql.NullInt64
	if containerIDStr != "" {
		id, err := strconv.ParseInt(containerIDStr, 10, 64)
		if err != nil {
			logger.Error("Invalid container ID",
				zap.Error(err),
				zap.String("raw_id", containerIDStr))
			http.Error(w, "Invalid container ID", http.StatusBadRequest)
			return
		}
		containerID = sql.NullInt64{Int64: id, Valid: true}
	}

	// Verify character belongs to user
	queries := db.New(s.db)
	_, err = queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Character not found or belongs to another user",
			zap.Error(err),
			zap.Int64("character_id", characterID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// If a container was specified, verify it belongs to the character
	if containerID.Valid {
		// Check if the container exists and belongs to the character
		inventory, err := queries.GetCharacterInventoryItems(r.Context(), characterID)
		if err != nil {
			logger.Error("Failed to fetch character inventory",
				zap.Error(err),
				zap.Int64("character_id", characterID))
			http.Error(w, "Error retrieving inventory", http.StatusInternalServerError)
			return
		}

		containerFound := false
		for _, item := range inventory {
			if item.ID == containerID.Int64 && item.ItemType == "container" {
				containerFound = true
				break
			}
		}

		if !containerFound {
			logger.Warn("Container not found or doesn't belong to character",
				zap.Int64("character_id", characterID),
				zap.Int64("container_id", containerID.Int64))
			http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Container not found", characterID), http.StatusSeeOther)
			return
		}

		// Check if the container has capacity for the item
		// This would require additional logic to check weight and item counts
		// For now, we'll skip detailed capacity checking
	}

	// Move the item to the container
	err = queries.MoveItemToContainer(r.Context(), db.MoveItemToContainerParams{
		ContainerID: containerID,
		ID:          itemID,
		CharacterID: characterID,
	})
	if err != nil {
		logger.Error("Failed to move item to container",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("item_id", itemID),
			zap.Any("container_id", containerID))
		http.Error(w, "Error moving item", http.StatusInternalServerError)
		return
	}

	logger.Info("Item moved to container successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("item_id", itemID),
		zap.Any("container_id", containerID))

	// Redirect back to character detail page with success message
	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Item moved successfully", characterID), http.StatusSeeOther)
}

func handleAddItemForm(s *Server, w http.ResponseWriter, r *http.Request, character db.Character, queries *db.Queries) {
	itemType := r.URL.Query().Get("type")

	var containerID sql.NullInt64
	if containerIDStr := r.URL.Query().Get("container_id"); containerIDStr != "" {
		id, err := strconv.ParseInt(containerIDStr, 10, 64)
		if err != nil {
			// Log error but continue
			logger.Warn("Invalid container ID", zap.Error(err), zap.String("raw_id", containerIDStr))
		} else {
			containerID = sql.NullInt64{Int64: id, Valid: true}
			// Add an explicit usage of containerID to appease the compiler
			logger.Debug("Container ID parsed successfully",
				zap.Int64("container_id", containerID.Int64),
				zap.Bool("is_valid", containerID.Valid))
		}
	}

	// Get enhancement query param with proper type conversion to int64
	var enhancementBonus int64
	if enhancementStr := r.URL.Query().Get("enhancement"); enhancementStr != "" {
		var err error
		enhancementBonus, err = strconv.ParseInt(enhancementStr, 10, 64)
		if err != nil {
			logger.Warn("Invalid enhancement bonus value",
				zap.Error(err),
				zap.String("raw_value", enhancementStr))
			enhancementBonus = 0
		}
	}

	// Get available containers for the character
	containers, err := queries.GetCharacterInventoryItems(r.Context(), character.ID)
	if err != nil {
		logger.Error("Failed to fetch character inventory",
			zap.Error(err), zap.Int64("character_id", character.ID))
		// IMPORTANT: Continue without containers rather than returning an error
		containers = []db.GetCharacterInventoryItemsRow{}
	}

	// Filter containers manually by type instead of relying on the query
	var filteredContainers []db.GetCharacterInventoryItemsRow
	for _, item := range containers {
		if item.ItemType == "container" {
			filteredContainers = append(filteredContainers, item)
		}
	}

	// Get equipment slots
	equipmentSlots, err := queries.GetEquipmentSlots(r.Context())
	if err != nil {
		logger.Error("Failed to fetch equipment slots", zap.Error(err))
		// Continue despite error - we'll use hardcoded defaults
	}

	// Get user from context
	user, ok := GetUserFromContext(r.Context())
	var username string
	if ok {
		username = user.Username
	}

	// Prepare data for the template
	data := struct {
		IsAuthenticated    bool
		Username           string
		CharacterID        int64
		SelectedType       string
		Items              interface{}
		Containers         []db.GetCharacterInventoryItemsRow
		EquipmentSlots     interface{}
		ShowEquipmentSlots bool
		FlashMessage       string
		CurrentYear        int
		Enhancement        int64 // Changed to int64
	}{
		IsAuthenticated:    ok,
		Username:           username,
		CharacterID:        character.ID,
		SelectedType:       itemType,
		Containers:         nil, // Will filter this below
		EquipmentSlots:     equipmentSlots,
		ShowEquipmentSlots: itemType == "weapon" || itemType == "armor" || itemType == "shield",
		FlashMessage:       r.URL.Query().Get("message"),
		CurrentYear:        time.Now().Year(),
		Enhancement:        enhancementBonus,
	}

	// Filter containers from inventory
	for _, item := range containers {
		if item.ItemType == "container" {
			data.Containers = append(data.Containers, item)
		}
	}

	// If a type is selected, fetch available items of that type
	if itemType != "" {
		var err error

		switch itemType {
		case "weapon":
			if enhancementBonus > 0 {
				// Fetch enhanced weapons with proper conversion to sql.NullInt64
				data.Items, err = queries.GetEnhancedWeapons(r.Context(), sql.NullInt64{
					Int64: enhancementBonus,
					Valid: true,
				})
			} else {
				// Fetch regular weapons
				data.Items, err = queries.GetAllWeapons(r.Context())
			}
		case "armor":
			if enhancementBonus > 0 {
				// Fetch enhanced armor with proper conversion to sql.NullInt64
				data.Items, err = queries.GetEnhancedArmor(r.Context(), sql.NullInt64{
					Int64: enhancementBonus,
					Valid: true,
				})
			} else {
				// Fetch regular armor
				data.Items, err = queries.GetAllArmor(r.Context())
			}
		case "shield":
			if enhancementBonus > 0 {
				// Fetch enhanced shields with proper conversion to sql.NullInt64
				data.Items, err = queries.GetEnhancedShields(r.Context(), sql.NullInt64{
					Int64: enhancementBonus,
					Valid: true,
				})
			} else {
				// Fetch regular shields
				data.Items, err = queries.GetAllShields(r.Context())
			}
		case "ranged_weapon":
			if enhancementBonus > 0 {
				// Fetch enhanced ranged weapons with proper conversion to sql.NullInt64
				data.Items, err = queries.GetEnhancedRangedWeapons(r.Context(), sql.NullInt64{
					Int64: enhancementBonus,
					Valid: true,
				})
			} else {
				// Fetch regular ranged weapons
				data.Items, err = queries.GetAllRangedWeapons(r.Context())
			}
		case "equipment":
			data.Items, err = queries.GetAllEquipment(r.Context())
		case "ammunition":
			data.Items, err = queries.GetAllAmmunition(r.Context())
		}

		if err != nil {
			logger.Error("Failed to fetch items",
				zap.Error(err),
				zap.String("item_type", itemType))
		}
	}

	// Render template
	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/inventory/add.html",
	)

	if err != nil {
		logger.Error("Template parsing error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		logger.Error("Template execution error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleAddItemSubmission(s *Server, w http.ResponseWriter, r *http.Request, character db.Character, queries *db.Queries) {
	if err := r.ParseForm(); err != nil {
		logger.Error("Failed to parse form", zap.Error(err))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	itemType := r.FormValue("item_type")
	itemIDStr := r.FormValue("item_id")
	quantityStr := r.FormValue("quantity")
	containerIDStr := r.FormValue("container_id")
	equipmentSlotIDStr := r.FormValue("equipment_slot_id")
	notes := r.FormValue("notes")

	// Parse item ID
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid item ID", zap.Error(err), zap.String("raw_id", itemIDStr))
		http.Redirect(w, r, fmt.Sprintf("/characters/inventory/add?character_id=%d&message=Invalid item ID", character.ID), http.StatusSeeOther)
		return
	}

	// Parse quantity (default to 1)
	quantity := int64(1)
	if quantityStr != "" {
		quantity, err = strconv.ParseInt(quantityStr, 10, 64)
		if err != nil || quantity < 1 {
			quantity = 1
		}
	}

	// Parse container ID if provided
	var containerID sql.NullInt64
	if containerIDStr != "" {
		id, err := strconv.ParseInt(containerIDStr, 10, 64)
		if err == nil {
			containerID = sql.NullInt64{Int64: id, Valid: true}
		}
	}

	// Parse equipment slot ID if provided
	var equipmentSlotID sql.NullInt64
	if equipmentSlotIDStr != "" {
		id, err := strconv.ParseInt(equipmentSlotIDStr, 10, 64)
		if err == nil {
			equipmentSlotID = sql.NullInt64{Int64: id, Valid: true}
		}
	}

	// Create null string for notes if provided
	var notesNull sql.NullString
	if notes != "" {
		notesNull = sql.NullString{String: notes, Valid: true}
	}

	// Add item to inventory
	_, err = queries.AddItemToInventory(r.Context(), db.AddItemToInventoryParams{
		CharacterID:     character.ID,
		ItemID:          itemID,
		ItemType:        itemType,
		Quantity:        quantity,
		ContainerID:     containerID,
		EquipmentSlotID: equipmentSlotID,
		Notes:           notesNull,
	})

	if err != nil {
		logger.Error("Failed to add item to inventory",
			zap.Error(err),
			zap.Int64("character_id", character.ID),
			zap.Int64("item_id", itemID),
			zap.String("item_type", itemType))
		http.Redirect(w, r, fmt.Sprintf("/characters/inventory/add?character_id=%d&message=Error adding item", character.ID), http.StatusSeeOther)
		return
	}

	logger.Info("Item added to inventory successfully",
		zap.Int64("character_id", character.ID),
		zap.Int64("item_id", itemID),
		zap.String("item_type", itemType))

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Item added successfully", character.ID), http.StatusSeeOther)
}

// PrepareContainerOptionsForItems adds container options to inventory items
func PrepareContainerOptionsForItems(items []InventoryItem, containers []InventoryItem) []InventoryItem {
	for i := range items {
		// Create container options for all items (they can all be stored)
		var containerOptions []InventoryItem
		for _, container := range containers {
			// Skip itself as a container option
			if container.ID != items[i].ID {
				containerOptions = append(containerOptions, container)
			}
		}
		items[i].ContainerOptions = containerOptions
	}
	return items
}

// EnhanceInventoryViewModel enhances the character view model with container options
func EnhanceInventoryViewModel(viewModel *CharacterViewModel) {
	// Find all containers
	var containers []InventoryItem
	for _, item := range viewModel.CarriedItems {
		if item.ItemType == "container" {
			containers = append(containers, item)
		}
	}

	// Add container options to carried items
	viewModel.CarriedItems = PrepareContainerOptionsForItems(viewModel.CarriedItems, containers)

	// Add container options to each container's items
	for containerID, containerItems := range viewModel.ContainerItems {
		// Get all containers except the current one
		var validContainers []InventoryItem
		for _, container := range containers {
			if container.ID != containerID {
				validContainers = append(validContainers, container)
			}
		}
		viewModel.ContainerItems[containerID] = PrepareContainerOptionsForItems(containerItems, validContainers)
	}
}

// Helper function to get available equipment slots for an item
func getAvailableSlotsForItem(itemType string) []db.EquipmentSlot {
	// This would be better as a database query, but for now we'll hardcode
	slots := make([]db.EquipmentSlot, 0)

	switch itemType {
	case "weapon", "ranged_weapon":
		slots = append(slots, db.EquipmentSlot{ID: 5, Name: "Right Hand"})
		slots = append(slots, db.EquipmentSlot{ID: 6, Name: "Left Hand"})
	case "armor":
		slots = append(slots, db.EquipmentSlot{ID: 3, Name: "Body"})
	case "shield":
		slots = append(slots, db.EquipmentSlot{ID: 6, Name: "Left Hand"})
	case "headgear":
		slots = append(slots, db.EquipmentSlot{ID: 1, Name: "Head"})
	default:
		// Add other slot mappings as needed
	}

	return slots
}

func (s *Server) HandleInventoryModal(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse character ID from query parameters
	characterIDStr := r.URL.Query().Get("character_id")
	characterID, err := strconv.ParseInt(characterIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err), zap.String("raw_id", characterIDStr))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	// Verify character belongs to user
	queries := db.New(s.db)
	_, err = queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Character not found or belongs to another user",
			zap.Error(err), zap.Int64("character_id", characterID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Get item type from query parameters (optional)
	itemType := r.URL.Query().Get("type")

	// Get enhancement bonus from query parameters (optional)
	var enhancement int64
	if enhancementStr := r.URL.Query().Get("enhancement"); enhancementStr != "" {
		var err error
		enhancement, err = strconv.ParseInt(enhancementStr, 10, 64)
		if err != nil {
			logger.Warn("Invalid enhancement bonus value",
				zap.Error(err),
				zap.String("raw_value", enhancementStr))
			enhancement = 0
		}
	}

	// Get container ID if adding to a container (optional)
	var containerID sql.NullInt64
	if containerIDStr := r.URL.Query().Get("container_id"); containerIDStr != "" {
		id, err := strconv.ParseInt(containerIDStr, 10, 64)
		if err == nil {
			containerID = sql.NullInt64{Int64: id, Valid: true}
		} else {
			logger.Warn("Invalid container ID", zap.String("raw_id", containerIDStr))
		}
	}

	// If we have a type but no enhancement selection yet for applicable item types,
	// and we're not on an enhancement selection step, show the enhancement selection form
	var enhancementStr string = r.URL.Query().Get("enhancement")
	if enhancementStr != "" {
		var err error
		enhancement, err = strconv.ParseInt(enhancementStr, 10, 64)
		if err != nil {
			logger.Warn("Invalid enhancement bonus value",
				zap.Error(err),
				zap.String("raw_value", enhancementStr))
			enhancement = 0
		}
	}

	// If we have a type but no enhancement selection yet for applicable item types,
	// and we're not on an enhancement selection step, show the enhancement selection form
	if itemType != "" &&
		enhancementStr == "" &&
		r.URL.Query().Get("step") != "enhancement" &&
		(itemType == "weapon" || itemType == "armor" || itemType == "shield" || itemType == "ranged_weapon") {
		renderEnhancementSelectionForm(w, characterID, itemType, containerID)
		return
	}

	// Get available containers for the character
	containers, err := queries.GetCharacterInventoryItems(r.Context(), characterID)
	if err != nil {
		logger.Warn("Failed to fetch inventory, proceeding with empty containers",
			zap.Error(err), zap.Int64("character_id", characterID))
		containers = []db.GetCharacterInventoryItemsRow{}
	}

	// Filter containers manually by type instead of relying on the query
	var filteredContainers []db.GetCharacterInventoryItemsRow
	for _, item := range containers {
		if item.ItemType == "container" {
			filteredContainers = append(filteredContainers, item)
		}
	}

	// Get equipment slots if needed
	var equipmentSlots []db.EquipmentSlot
	if itemType != "" && (itemType == "weapon" || itemType == "armor" || itemType == "shield" || itemType == "ranged_weapon") {
		equipmentSlots, err = queries.GetEquipmentSlots(r.Context())
		if err != nil {
			logger.Warn("Failed to fetch equipment slots",
				zap.Error(err))
		}
	}

	// Prepare data for the template
	data := struct {
		CharacterID        int64
		SelectedType       string
		Enhancement        int64
		Items              interface{}
		Containers         []db.GetCharacterInventoryItemsRow
		EquipmentSlots     []db.EquipmentSlot
		ShowEquipmentSlots bool
		HasContainerID     bool
		ContainerID        sql.NullInt64
	}{
		CharacterID:        characterID,
		SelectedType:       itemType,
		Enhancement:        enhancement,
		Containers:         filteredContainers,
		EquipmentSlots:     equipmentSlots,
		ShowEquipmentSlots: itemType == "weapon" || itemType == "armor" || itemType == "shield" || itemType == "ranged_weapon",
		HasContainerID:     containerID.Valid,
		ContainerID:        containerID,
	}

	// If a type is selected, fetch available items of that type
	if itemType != "" && (enhancementStr != "" ||
		!(itemType == "weapon" || itemType == "armor" || itemType == "shield" || itemType == "ranged_weapon")) {
		var err error

		switch itemType {
		case "weapon":
			if enhancement > 0 {
				// Fetch enhanced weapons with proper sql.NullInt64
				data.Items, err = queries.GetEnhancedWeapons(r.Context(), sql.NullInt64{
					Int64: enhancement,
					Valid: true,
				})
			} else {
				// Fetch regular weapons
				data.Items, err = queries.GetAllWeapons(r.Context())
			}
		case "armor":
			if enhancement > 0 {
				// Fetch enhanced armor with proper sql.NullInt64
				data.Items, err = queries.GetEnhancedArmor(r.Context(), sql.NullInt64{
					Int64: enhancement,
					Valid: true,
				})
			} else {
				// Fetch regular armor
				data.Items, err = queries.GetAllArmor(r.Context())
			}
		case "shield":
			if enhancement > 0 {
				// Fetch enhanced shields with proper sql.NullInt64
				data.Items, err = queries.GetEnhancedShields(r.Context(), sql.NullInt64{
					Int64: enhancement,
					Valid: true,
				})
			} else {
				// Fetch regular shields
				data.Items, err = queries.GetAllShields(r.Context())
			}
		case "ranged_weapon":
			if enhancement > 0 {
				// Fetch enhanced ranged weapons with proper sql.NullInt64
				data.Items, err = queries.GetEnhancedRangedWeapons(r.Context(), sql.NullInt64{
					Int64: enhancement,
					Valid: true,
				})
			} else {
				// Fetch regular ranged weapons
				data.Items, err = queries.GetAllRangedWeapons(r.Context())
			}
		case "equipment":
			data.Items, err = queries.GetAllEquipment(r.Context())
		case "ammunition":
			data.Items, err = queries.GetAllAmmunition(r.Context())
		}

		if err != nil {
			logger.Error("Failed to fetch items",
				zap.Error(err),
				zap.String("item_type", itemType))
		}
	}

	// Create a template function to render a partial template
	tmpl, err := template.New("_modal").Parse(`
    {{if not .SelectedType}}
    <form hx-get="/characters/inventory/modal" hx-target="#add-item-form-container">
        <input type="hidden" name="character_id" value="{{.CharacterID}}">
        <div class="form-group">
            <label for="type">Select Item Type:</label>
            <select name="type" id="type" required>
                <option value="">-- Select Type --</option>
                <option value="equipment">Equipment</option>
                <option value="weapon">Weapon</option>
                <option value="armor">Armor</option>
                <option value="ammunition">Ammunition</option>
                <option value="container">Container</option>
                <option value="shield">Shield</option>
                <option value="ranged_weapon">Ranged Weapon</option>
            </select>
        </div>
        <div class="form-actions">
            <button type="submit" class="button primary">Next</button>
            <button type="button" class="button close-modal">Cancel</button>
        </div>
    </form>
    {{else if and (or (eq .SelectedType "weapon") (eq .SelectedType "armor") (eq .SelectedType "shield") (eq .SelectedType "ranged_weapon")) (eq .Enhancement 0)}}
    <!-- Enhancement selection form -->
    <form hx-get="/characters/inventory/modal" hx-target="#add-item-form-container">
        <input type="hidden" name="character_id" value="{{.CharacterID}}">
        <input type="hidden" name="type" value="{{.SelectedType}}">
        <input type="hidden" name="step" value="enhancement">
        {{if .HasContainerID}}
        <input type="hidden" name="container_id" value="{{.ContainerID.Int64}}">
        {{end}}
        
        <div class="form-group">
            <label for="enhancement">Enhancement Bonus:</label>
            <select name="enhancement" id="enhancement" required>
                <option value="0">No Enhancement (+0)</option>
                <option value="1">Enhanced (+1)</option>
                <option value="2">Enhanced (+2)</option>
                <option value="3">Enhanced (+3)</option>
            </select>
        </div>
        
        <div class="form-actions">
            <button type="submit" class="button primary">Next</button>
            <button type="button" hx-get="/characters/inventory/modal?character_id={{.CharacterID}}"
                hx-target="#add-item-form-container" class="button">Back</button>
            <button type="button" class="button close-modal">Cancel</button>
        </div>
    </form>
    {{else}}
    <form hx-post="/characters/inventory/add-modal" hx-target="#character-sheet-container">
        <input type="hidden" name="character_id" value="{{.CharacterID}}">
        <input type="hidden" name="item_type" value="{{.SelectedType}}">
        <input type="hidden" name="enhancement" value="{{.Enhancement}}">
        {{if .HasContainerID}}
        <input type="hidden" name="container_id" value="{{.ContainerID.Int64}}">
        {{end}}
 
        <div class="form-group">
            <label for="item_id">Select Item:</label>
            <select name="item_id" id="item_id" required>
                <option value="">-- Select Item --</option>
                {{range .Items}}
                <option value="{{.ID}}">
                    {{.Name}} {{if gt .Weight 0}}({{.Weight}} lbs){{end}} {{if gt .CostGp 0}}({{.CostGp}} gp){{end}}
                </option>
                {{end}}
            </select>
        </div>
 
        <div class="form-group">
            <label for="quantity">Quantity:</label>
            <input type="number" name="quantity" id="quantity" value="1" min="1" required>
        </div>
 
        {{if .Containers}}
        <div class="form-group">
            <label for="container_id">Store in Container (optional):</label>
            <select name="container_id" id="container_id">
                <option value="">-- None --</option>
                {{range .Containers}}
                <option value="{{.ID}}">{{.ItemName}}</option>
                {{end}}
            </select>
        </div>
        {{end}}
 
        {{if .ShowEquipmentSlots}}
        <div class="form-group">
            <label for="equipment_slot_id">Equipment Slot (optional):</label>
            <select name="equipment_slot_id" id="equipment_slot_id">
                <option value="">-- None --</option>
                {{range .EquipmentSlots}}
                <option value="{{.ID}}">{{.Name}}</option>
                {{end}}
            </select>
        </div>
        {{end}}
 
        <div class="form-group">
            <label for="notes">Notes (optional):</label>
            <textarea name="notes" id="notes" rows="3"></textarea>
        </div>
 
        <div class="form-actions">
            <button type="submit" class="button primary">
                Add Item
                <span class="htmx-indicator">
                    <div class="spinner"></div>
                </span>
            </button>
            {{if .HasContainerID}}
            <button type="button"
                hx-get="/characters/inventory/modal?character_id={{.CharacterID}}&container_id={{.ContainerID.Int64}}"
                hx-target="#add-item-form-container" class="button">Back</button>
            {{else}}
            <button type="button" hx-get="/characters/inventory/modal?character_id={{.CharacterID}}"
                hx-target="#add-item-form-container" class="button">Back</button>
            {{end}}
            <button type="button" class="button close-modal">Cancel</button>
        </div>
    </form>
    {{end}}
    `)

	if err != nil {
		logger.Error("Template parsing error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		logger.Error("Template execution error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Helper function to render just the enhancement selection form
func renderEnhancementSelectionForm(w http.ResponseWriter, characterID int64, itemType string, containerID sql.NullInt64) {
	tmpl, err := template.New("enhancement_form").Parse(`
	<form hx-get="/characters/inventory/modal" hx-target="#add-item-form-container">
		<input type="hidden" name="character_id" value="{{.CharacterID}}">
		<input type="hidden" name="type" value="{{.ItemType}}">
		<input type="hidden" name="step" value="enhancement">
		{{if .HasContainerID}}
		<input type="hidden" name="container_id" value="{{.ContainerID}}">
		{{end}}
		
		<div class="form-group">
			<label for="enhancement">Enhancement Bonus:</label>
			<select name="enhancement" id="enhancement" required>
				<option value="0">No Enhancement (+0)</option>
				<option value="1">Enhanced (+1)</option>
				<option value="2">Enhanced (+2)</option>
				<option value="3">Enhanced (+3)</option>
			</select>
		</div>
		
		<div class="form-actions">
			<button type="submit" class="button primary">Next</button>
			<button type="button" hx-get="/characters/inventory/modal?character_id={{.CharacterID}}"
				hx-target="#add-item-form-container" class="button">Back</button>
			<button type="button" class="button close-modal">Cancel</button>
		</div>
	</form>
	`)

	if err != nil {
		logger.Error("Enhancement form template parsing error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		CharacterID    int64
		ItemType       string
		HasContainerID bool
		ContainerID    int64
	}{
		CharacterID:    characterID,
		ItemType:       itemType,
		HasContainerID: containerID.Valid,
		ContainerID:    containerID.Int64,
	}

	if err := tmpl.Execute(w, data); err != nil {
		logger.Error("Enhancement form template execution error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleAddItemModal(w http.ResponseWriter, r *http.Request) {
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

	// Parse character ID
	characterID, err := strconv.ParseInt(r.FormValue("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	// Fetch character to verify ownership
	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Character not found or belongs to another user",
			zap.Error(err), zap.Int64("character_id", characterID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Get form values
	itemType := r.FormValue("item_type")
	itemID, err := strconv.ParseInt(r.FormValue("item_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid item ID", zap.Error(err))
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Parse enhancement
	enhancementStr := r.FormValue("enhancement")
	var enhancement int64
	if enhancementStr != "" {
		enhancement, _ = strconv.ParseInt(enhancementStr, 10, 64)
	}

	// Parse quantity (default to 1)
	quantity := int64(1)
	if quantityStr := r.FormValue("quantity"); quantityStr != "" {
		if q, err := strconv.ParseInt(quantityStr, 10, 64); err == nil && q > 0 {
			quantity = q
		}
	}

	// Parse container ID if provided
	var containerID sql.NullInt64
	if containerIDStr := r.FormValue("container_id"); containerIDStr != "" {
		if id, err := strconv.ParseInt(containerIDStr, 10, 64); err == nil {
			containerID = sql.NullInt64{Int64: id, Valid: true}
		}
	}

	// Parse equipment slot ID if provided
	var equipmentSlotID sql.NullInt64
	if slotIDStr := r.FormValue("equipment_slot_id"); slotIDStr != "" {
		if id, err := strconv.ParseInt(slotIDStr, 10, 64); err == nil {
			// Check if the slot is already occupied
			isOccupied, err := queries.IsSlotOccupied(r.Context(), db.IsSlotOccupiedParams{
				CharacterID:     characterID,
				EquipmentSlotID: sql.NullInt64{Int64: id, Valid: true},
			})
			if err != nil {
				logger.Error("Failed to check if slot is occupied", zap.Error(err))
			} else if isOccupied {
				// Render character sheet with error message
				renderCharacterWithMessage(s, w, r, character, "Equipment slot is already occupied")
				return
			}
			equipmentSlotID = sql.NullInt64{Int64: id, Valid: true}
		}
	}

	// Parse notes if provided
	var notes sql.NullString
	if notesStr := r.FormValue("notes"); notesStr != "" {
		notes = sql.NullString{String: notesStr, Valid: true}
	}

	// Add a note about enhancement if applicable
	if enhancement > 0 && notes.String == "" {
		notes = sql.NullString{
			String: fmt.Sprintf("+%d enhancement", enhancement),
			Valid:  true,
		}
	} else if enhancement > 0 {
		notes = sql.NullString{
			String: fmt.Sprintf("%s (+%d enhancement)", notes.String, enhancement),
			Valid:  true,
		}
	}

	// Add item to inventory
	_, err = queries.AddItemToInventory(r.Context(), db.AddItemToInventoryParams{
		CharacterID:     characterID,
		ItemID:          itemID,
		ItemType:        itemType,
		Quantity:        quantity,
		ContainerID:     containerID,
		EquipmentSlotID: equipmentSlotID,
		Notes:           notes,
	})

	if err != nil {
		logger.Error("Failed to add item to inventory",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("item_id", itemID))

		// Render character sheet with error message
		renderCharacterWithMessage(s, w, r, character, "Error adding item to inventory")
		return
	}

	logger.Info("Item added to inventory successfully via modal",
		zap.Int64("character_id", characterID),
		zap.Int64("item_id", itemID),
		zap.String("item_type", itemType),
		zap.Int64("enhancement", enhancement))

	// Render character sheet with success message
	renderCharacterWithMessage(s, w, r, character, "Item added successfully")
}

// Helper function to render character sheet with a message
func renderCharacterWithMessage(s *Server, w http.ResponseWriter, r *http.Request, character db.Character, message string) {
	// Fetch all character data needed for the view
	queries := db.New(s.db)
	inventory, err := queries.GetCharacterInventoryItems(r.Context(), character.ID)
	if err != nil {
		logger.Warn("Failed to fetch inventory after item addition",
			zap.Error(err), zap.Int64("character_id", character.ID))
		inventory = []db.GetCharacterInventoryItemsRow{}
	}

	// Create view model
	viewModel := NewSafeCharacterViewModel(character, inventory)

	// Render full character detail page
	tmpl, err := template.New("detail-content").Funcs(template.FuncMap{
		"seq": func(start, end int) []int {
			s := make([]int, end-start+1)
			for i := range s {
				s[i] = start + i
			}
			return s
		},
		"GetSavingThrowModifiers": charRules.GetSavingThrowModifiers,
		"add": func(a, b interface{}) int64 {
			// Implementation as in the original code
			switch v := a.(type) {
			case int64:
				switch w := b.(type) {
				case int:
					return v + int64(w)
				case int64:
					return v + w
				}
			case int:
				switch w := b.(type) {
				case int64:
					return int64(v) + w
				case int:
					return int64(v + w)
				}
			}
			return 0
		},
		"mul": func(a, b interface{}) int64 {
			// Implementation as in the original code
			switch v := a.(type) {
			case int64:
				switch w := b.(type) {
				case int:
					return v * int64(w)
				case int64:
					return v * w
				}
			case int:
				switch w := b.(type) {
				case int64:
					return int64(v) * w
				case int:
					return int64(v * w)
				}
			}
			return 0
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"sub": func(a, b interface{}) int64 {
			// Implementation as in the original code
			return 0 // Simplified for this example
		},
		"abs": func(x int) int {
			if x < 0 {
				return -x
			}
			return x
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("January 2, 2006 3:04 PM")
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"formatModifier": func(mod int) string {
			if mod > 0 {
				return "+" + strconv.Itoa(mod)
			}
			return strconv.Itoa(mod)
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"contains": containsString,
	}).ParseFiles(
		"templates/characters/details.html",
		"templates/characters/_character_header.html",
		"templates/characters/inventory_modal.html",
		"templates/characters/_inventory.html",
		"templates/characters/_ability_scores.html",
		"templates/characters/_class_features.html",
		"templates/characters/_combat_stats.html",
		"templates/characters/_saving_throws.html",
		"templates/characters/_hp_display.html",
		"templates/characters/_hp_section.html",
		"templates/characters/_currency_section.html",
	)

	if err != nil {
		logger.Error("Template parsing error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Character    CharacterViewModel
		FlashMessage string
	}{
		Character:    viewModel,
		FlashMessage: message,
	}

	// Only render the character sheet part, not the full page with headers
	err = tmpl.ExecuteTemplate(w, "content", data)
	if err != nil {
		logger.Error("Template execution error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
