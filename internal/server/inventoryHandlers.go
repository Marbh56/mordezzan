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
	// TODO: Implement inventory update functionality
}

// HandleEquipItem handles equipping items from inventory to equipment slots
func (s *Server) HandleEquipItem(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement equip item functionality
}

// HandleUnequipItem handles unequipping items from equipment slots to inventory
func (s *Server) HandleUnequipItem(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement unequip item functionality
}

// HandleMoveToContainer handles moving items to containers
func (s *Server) HandleMoveToContainer(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement move to container functionality
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
