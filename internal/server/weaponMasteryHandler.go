package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/rules"
)

type WeaponMasteryData struct {
	Character        db.Character
	CurrentMasteries []db.GetCharacterWeaponMasteriesRow
	AvailableWeapons []db.GetWeaponItemsRow
	AvailableSlots   int
	HasGrandMastery  bool
	FlashMessage     string
	CurrentYear      int
	IsAuthenticated  bool
	Username         string
}

func (s *Server) HandleWeaponMastery(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	characterIDStr := r.URL.Query().Get("id")
	characterID, err := strconv.ParseInt(characterIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)

	// Get character data and verify ownership
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Error fetching character: %v", err)
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Only fighters and fighter subclasses can access this
	if !isFighterClass(character.Class) {
		http.Redirect(w, r, "/characters/detail?id="+characterIDStr, http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleWeaponMasteryForm(w, r, character, user)
	case http.MethodPost:
		action := r.URL.Query().Get("action")
		switch action {
		case "add":
			s.handleAddMastery(w, r)
		case "remove":
			s.handleRemoveMastery(w, r)
		case "upgrade":
			s.handleUpgradeMastery(w, r)
		default:
			http.Error(w, "Invalid action", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleWeaponMasteryForm(w http.ResponseWriter, r *http.Request, character db.Character, user *db.GetSessionRow) {
	queries := db.New(s.db)

	// Get current masteries
	currentMasteries, err := queries.GetCharacterWeaponMasteries(r.Context(), character.ID)
	if err != nil {
		log.Printf("Error fetching masteries: %v", err)
		http.Error(w, "Error loading masteries", http.StatusInternalServerError)
		return
	}

	// Get available weapons
	availableWeapons, err := queries.GetWeaponItems(r.Context())
	if err != nil {
		log.Printf("Error fetching weapons: %v", err)
		http.Error(w, "Error loading weapons", http.StatusInternalServerError)
		return
	}

	// Calculate available slots based on level
	availableSlots := rules.GetAvailableMasterySlots(character.Level) - len(currentMasteries)
	if availableSlots < 0 {
		availableSlots = 0
	}

	// Check if character already has a grand mastery
	hasGrandMastery := false
	for _, mastery := range currentMasteries {
		if mastery.MasteryLevel == "grand_mastery" {
			hasGrandMastery = true
			break
		}
	}

	data := WeaponMasteryData{
		Character:        character,
		CurrentMasteries: currentMasteries,
		AvailableWeapons: availableWeapons,
		AvailableSlots:   availableSlots,
		HasGrandMastery:  hasGrandMastery,
		FlashMessage:     r.URL.Query().Get("message"),
		CurrentYear:      time.Now().Year(),
		IsAuthenticated:  true,
		Username:         user.Username,
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/weapons/mastery.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Template data: %+v", data)
	log.Printf("Parsing templates: base.html and weapons/mastery.html")

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func isFighterClass(class string) bool {
	fighterClasses := map[string]bool{
		"Fighter": true,
	}
	return fighterClasses[class]
}

func (s *Server) handleAddMastery(w http.ResponseWriter, r *http.Request) {
	log.Printf("Form values: %v", r.Form)
	log.Printf("URL values: %v", r.URL.Query())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	weaponID, err := strconv.ParseInt(r.Form.Get("weapon_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid weapon ID", http.StatusBadRequest)
		return
	}

	masteryLevel := r.Form.Get("mastery_level")
	if masteryLevel != "mastered" && masteryLevel != "grand_mastery" {
		http.Error(w, "Invalid mastery level", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)

	// Verify character has available slots
	currentMasteries, err := queries.GetCharacterWeaponMasteries(r.Context(), characterID)
	if err != nil {
		log.Printf("Error checking current masteries: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Error checking masteries", characterID), http.StatusSeeOther)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Error fetching character: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Error verifying character", characterID), http.StatusSeeOther)
		return
	}

	availableSlots := rules.GetAvailableMasterySlots(character.Level) - len(currentMasteries)
	if availableSlots <= 0 {
		http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=No available mastery slots", characterID), http.StatusSeeOther)
		return
	}

	// Verify not already at grand mastery if trying to add one
	if masteryLevel == "grand_mastery" {
		for _, mastery := range currentMasteries {
			if mastery.MasteryLevel == "grand_mastery" {
				http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Already have a grand mastery", characterID), http.StatusSeeOther)
				return
			}
		}
	}

	// Add the mastery
	_, err = queries.AddWeaponMastery(r.Context(), db.AddWeaponMasteryParams{
		CharacterID:  characterID,
		WeaponID:     weaponID,
		MasteryLevel: masteryLevel,
	})

	if err != nil {
		log.Printf("Error adding weapon mastery: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Error adding mastery", characterID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Mastery added successfully", characterID), http.StatusSeeOther)
}

func (s *Server) handleRemoveMastery(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	weaponID, err := strconv.ParseInt(r.Form.Get("weapon_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid weapon ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)

	err = queries.RemoveWeaponMastery(r.Context(), db.RemoveWeaponMasteryParams{
		CharacterID: characterID,
		WeaponID:    weaponID,
	})

	if err != nil {
		log.Printf("Error removing weapon mastery: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Error removing mastery", characterID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Mastery removed successfully", characterID), http.StatusSeeOther)
}

func (s *Server) handleUpgradeMastery(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	weaponID, err := strconv.ParseInt(r.Form.Get("weapon_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid weapon ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)

	// Verify no existing grand mastery
	currentMasteries, err := queries.GetCharacterWeaponMasteries(r.Context(), characterID)
	if err != nil {
		log.Printf("Error checking current masteries: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Error checking masteries", characterID), http.StatusSeeOther)
		return
	}

	for _, mastery := range currentMasteries {
		if mastery.MasteryLevel == "grand_mastery" {
			http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Already have a grand mastery", characterID), http.StatusSeeOther)
			return
		}
	}

	_, err = queries.UpdateWeaponMastery(r.Context(), db.UpdateWeaponMasteryParams{
		CharacterID:  characterID,
		WeaponID:     weaponID,
		MasteryLevel: "grand_mastery",
	})

	if err != nil {
		log.Printf("Error upgrading weapon mastery: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Error upgrading mastery", characterID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/characters/masteries?id=%d&message=Mastery upgraded successfully", characterID), http.StatusSeeOther)
}
