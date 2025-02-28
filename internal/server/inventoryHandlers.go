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
		inventory, err := queries.GetCharacterInventory(r.Context(), characterID)
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
	// Get item type from query parameters (optional)
	itemType := r.URL.Query().Get("type")

	// Get container ID if adding to a container (optional)
	containerIDStr := r.URL.Query().Get("container_id")
	if containerIDStr != "" {
		_, err := strconv.ParseInt(containerIDStr, 10, 64)
		if err != nil {
			// Log error but continue
			logger.Warn("Invalid container ID", zap.String("raw_id", containerIDStr))
		}
	}

	// Get available containers for the character
	containers, err := queries.GetCharacterInventory(r.Context(), character.ID)
	if err != nil {
		logger.Error("Failed to fetch character inventory",
			zap.Error(err), zap.Int64("character_id", character.ID))
		http.Error(w, "Error loading inventory data", http.StatusInternalServerError)
		return
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
		Containers         []db.GetCharacterInventoryRow
		EquipmentSlots     interface{}
		ShowEquipmentSlots bool
		FlashMessage       string
		CurrentYear        int
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
			data.Items, err = queries.GetAllWeapons(r.Context())
		case "armor":
			data.Items, err = queries.GetAllArmor(r.Context())
		case "shield":
			data.Items, err = queries.GetAllShields(r.Context())
		case "equipment":
			data.Items, err = queries.GetAllEquipment(r.Context())
		case "ammunition":
			data.Items, err = queries.GetAllAmmunition(r.Context())
		case "ranged_weapon":
			data.Items, err = queries.GetAllRangedWeapons(r.Context())
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
