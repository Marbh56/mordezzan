package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
)

type ItemData struct {
	ID     int64
	Name   string
	Weight int64
	CostGP float64
}

func (s *Server) getEquipmentSlots() ([]db.EquipmentSlot, error) {
	queries := db.New(s.db)
	slots, err := queries.ListEquipmentSlots(context.Background())
	if err != nil {
		return nil, err
	}
	return slots, nil
}

func (s *Server) HandleAddInventoryItem(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	characterIDStr := r.URL.Query().Get("character_id")
	if characterIDStr == "" {
		characterIDStr = r.FormValue("character_id") // Also check POST form data
	}

	if characterIDStr == "" {
		http.Error(w, "Character ID is required", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(characterIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)
	_, err = queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Character not found", http.StatusNotFound)
		} else {
			log.Printf("Error verifying character ownership: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleAddItemForm(w, r, characterID)
	case http.MethodPost:
		s.handleAddItemSubmission(w, r, characterID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleAddItemForm(w http.ResponseWriter, r *http.Request, characterID int64) {
	selectedType := r.URL.Query().Get("type")
	user, _ := GetUserFromContext(r.Context())

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/inventory/add.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsAuthenticated    bool
		Username           string
		CharacterID        int64
		SelectedType       string
		Items              []ItemData
		Containers         []db.GetCharacterInventoryRow
		EquipmentSlots     []db.EquipmentSlot
		ShowEquipmentSlots bool
		FlashMessage       string
		CurrentYear        int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		CharacterID:     characterID,
		SelectedType:    selectedType,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	if selectedType != "" {
		// Load available items based on type
		items, err := s.getItemsByType(selectedType)
		if err != nil {
			log.Printf("Error loading items: %v", err)
			http.Error(w, "Error loading items", http.StatusInternalServerError)
			return
		}
		data.Items = items

		// Load containers for this character
		queries := db.New(s.db)
		containers, err := queries.GetCharacterInventory(r.Context(), characterID)
		if err != nil {
			log.Printf("Error loading containers: %v", err)
		} else {
			// Filter to only show containers
			var containerItems []db.GetCharacterInventoryRow
			for _, item := range containers {
				if item.ItemType == "container" {
					containerItems = append(containerItems, item)
				}
			}
			data.Containers = containerItems
		}

		// Load equipment slots if needed
		data.ShowEquipmentSlots = isEquippableType(selectedType)
		if data.ShowEquipmentSlots {
			slots, err := s.getEquipmentSlots()
			if err != nil {
				log.Printf("Error loading equipment slots: %v", err)
			} else {
				data.EquipmentSlots = slots
			}
		}
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) handleAddItemSubmission(w http.ResponseWriter, r *http.Request, characterID int64) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse basic item info
	itemType := r.Form.Get("item_type")
	itemID, err := strconv.ParseInt(r.Form.Get("item_id"), 10, 64)
	quantity, err := strconv.ParseInt(r.Form.Get("quantity"), 10, 64)
	if err != nil || quantity < 1 {
		http.Redirect(w, r, fmt.Sprintf("/characters/inventory/add?character_id=%d&message=Invalid quantity", characterID), http.StatusSeeOther)
		return
	}

	// Handle container
	var containerID sql.NullInt64
	if contID := r.Form.Get("container_inventory_id"); contID != "" {
		id, err := strconv.ParseInt(contID, 10, 64)
		if err == nil {
			containerID.Int64 = id
			containerID.Valid = true

			// Get item weight based on type
			var itemWeight int64
			switch itemType {
			case "equipment":
				err = s.db.QueryRowContext(r.Context(),
					"SELECT COALESCE(weight, 0) as weight FROM equipment WHERE id = ?",
					itemID).Scan(&itemWeight)
			case "weapon":
				err = s.db.QueryRowContext(r.Context(),
					"SELECT weight FROM weapons WHERE id = ?",
					itemID).Scan(&itemWeight)
			case "armor":
				err = s.db.QueryRowContext(r.Context(),
					"SELECT weight FROM armor WHERE id = ?",
					itemID).Scan(&itemWeight)
			}

			if err != nil {
				log.Printf("Error getting item weight: %v", err)
				http.Redirect(w, r, fmt.Sprintf("/characters/inventory/add?character_id=%d&message=Error validating container capacity", characterID), http.StatusSeeOther)
				return
			}

			// Validate container capacity
			if err := s.validateContainerCapacity(r.Context(), containerID.Int64, characterID, itemWeight, quantity); err != nil {
				http.Redirect(w, r, fmt.Sprintf("/characters/inventory/add?character_id=%d&message=%s", characterID, err.Error()), http.StatusSeeOther)
				return
			}
		}
	}

	// Continue with equipment slot logic
	var equipmentSlotID sql.NullInt64
	if slotID := r.Form.Get("equipment_slot_id"); slotID != "" {
		id, err := strconv.ParseInt(slotID, 10, 64)
		if err == nil {
			equipmentSlotID.Int64 = id
			equipmentSlotID.Valid = true
		}
	}

	// Add notes if any
	var notes sql.NullString
	if noteText := r.Form.Get("notes"); noteText != "" {
		notes.String = noteText
		notes.Valid = true
	}

	// Add the item
	queries := db.New(s.db)
	_, err = queries.AddItemToInventory(r.Context(), db.AddItemToInventoryParams{
		CharacterID:          characterID,
		ItemType:             itemType,
		ItemID:               itemID,
		Quantity:             quantity,
		ContainerInventoryID: containerID,
		EquipmentSlotID:      equipmentSlotID,
		Notes:                notes,
	})

	if err != nil {
		log.Printf("Error adding item to inventory: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/inventory/add?character_id=%d&message=Error adding item to inventory", characterID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Item added successfully", characterID), http.StatusSeeOther)
}

// Helper method to get items by type
func (s *Server) getItemsByType(itemType string) ([]ItemData, error) {
	queries := db.New(s.db)
	var items []ItemData

	switch itemType {
	case "equipment":
		equipmentItems, err := queries.GetEquipmentItems(context.Background())
		if err != nil {
			return nil, err
		}
		for _, item := range equipmentItems {
			items = append(items, ItemData{
				ID:     item.ID,
				Name:   item.Name,
				Weight: item.Weight,
				CostGP: item.CostGp,
			})
		}
	case "weapon":
		weaponItems, err := queries.GetWeaponItems(context.Background())
		if err != nil {
			return nil, err
		}
		for _, item := range weaponItems {
			costGP := float64(item.CostGp)
			items = append(items, ItemData{
				ID:     item.ID,
				Name:   item.Name,
				Weight: item.Weight,
				CostGP: costGP,
			})
		}
	case "armor":
		armorItems, err := queries.GetArmorItems(context.Background())
		if err != nil {
			return nil, err
		}
		for _, item := range armorItems {
			items = append(items, ItemData{
				ID:     item.ID,
				Name:   item.Name,
				Weight: item.Weight,
				CostGP: float64(item.CostGp),
			})
		}
	case "ammunition":
		ammoItems, err := queries.GetAmmunitionItems(context.Background())
		if err != nil {
			return nil, err
		}
		for _, item := range ammoItems {
			items = append(items, ItemData{
				ID:     item.ID,
				Name:   item.Name,
				Weight: item.Weight,
				CostGP: item.CostGp,
			})
		}
	case "container":
		containerItems, err := queries.GetContainerItems(context.Background())
		if err != nil {
			return nil, err
		}
		for _, item := range containerItems {
			items = append(items, ItemData{
				ID:     item.ID,
				Name:   item.Name,
				Weight: item.Weight,
				CostGP: item.CostGp,
			})
		}
	case "shield":
		shieldItems, err := queries.GetShieldItems(context.Background())
		if err != nil {
			return nil, err
		}
		for _, item := range shieldItems {
			items = append(items, ItemData{
				ID:     item.ID,
				Name:   item.Name,
				Weight: item.Weight,
				CostGP: float64(item.CostGp),
			})
		}
	case "ranged_weapon":
		rangedItems, err := queries.GetRangedWeaponItems(context.Background())
		if err != nil {
			return nil, err
		}
		for _, item := range rangedItems {
			var costGP float64
			if item.CostGp.Valid {
				costGP = float64(item.CostGp.Int64)
			}
			items = append(items, ItemData{
				ID:     item.ID,
				Name:   item.Name,
				Weight: item.Weight,
				CostGP: costGP,
			})
		}
	default:
		return nil, fmt.Errorf("unknown item type: %s", itemType)
	}

	return items, nil
}

// Helper function to determine if an item type can be equipped
func isEquippableType(itemType string) bool {
	equippableTypes := map[string]bool{
		"weapon":        true,
		"armor":         true,
		"shield":        true,
		"ranged_weapon": true,
	}
	return equippableTypes[itemType]
}

func (s *Server) HandleRemoveInventoryItem(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form values
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get required parameters
	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(r.Form.Get("item_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Verify character ownership
	queries := db.New(s.db)
	_, err = queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Error verifying character ownership: %v", err)
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Verify item belongs to character
	item, err := queries.GetItemFromInventory(r.Context(), db.GetItemFromInventoryParams{
		ID:          itemID,
		CharacterID: characterID,
	})
	if err != nil {
		log.Printf("Error verifying item ownership: %v", err)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Remove item from inventory
	err = queries.RemoveItemFromInventory(r.Context(), db.RemoveItemFromInventoryParams{
		ID:          item.ID,
		CharacterID: characterID,
	})
	if err != nil {
		log.Printf("Error removing item from inventory: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error removing item", characterID), http.StatusSeeOther)
		return
	}

	// Redirect back to character detail page with success message
	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Item removed successfully", characterID), http.StatusSeeOther)
}
