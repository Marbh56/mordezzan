package server

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/logger"
	"github.com/marbh56/mordezzan/internal/rules/ability_scores"
	charRules "github.com/marbh56/mordezzan/internal/rules/character"
	"go.uber.org/zap"
)

func calculateTotalHP(baseHP, level, constitution int64) (int64, error) {
	if baseHP <= 0 {
		return 0, fmt.Errorf("base HP must be positive")
	}

	if level < 1 || level > 20 {
		return 0, fmt.Errorf("level must be between 1 and 20")
	}

	if constitution < 3 || constitution > 18 {
		return 0, fmt.Errorf("constitution must be between 3 and 18")
	}

	// Get constitution modifiers from new package
	conMods := ability_scores.CalculateConstitutionModifiers(constitution)

	// Calculate additional HP from constitution modifier
	constitutionBonus := int64(conMods.HitPointMod) * level

	// Total HP is base HP plus constitution bonus
	totalHP := baseHP + constitutionBonus

	// Ensure minimum HP of 1
	if totalHP < 1 {
		return 1, nil
	}

	return totalHP, nil
}

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
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get parameters
	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseInt(r.Form.Get("amount"), 10, 64)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid amount", characterID), http.StatusSeeOther)
		return
	}

	denomination := r.Form.Get("denomination")
	if !isValidDenomination(denomination) {
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid denomination", characterID), http.StatusSeeOther)
		return
	}

	// Get character
	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Error fetching character: %v", err)
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Update currency based on denomination

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
		PlatinumPieces:   character.PlatinumPieces,
		GoldPieces:       character.GoldPieces,
		ElectrumPieces:   character.ElectrumPieces,
		SilverPieces:     character.SilverPieces,
		CopperPieces:     character.CopperPieces,
	}

	switch denomination {
	case "pp":
		updateParams.PlatinumPieces = character.PlatinumPieces + amount
	case "gp":
		updateParams.GoldPieces = character.GoldPieces + amount
	case "ep":
		updateParams.ElectrumPieces = character.ElectrumPieces + amount
	case "sp":
		updateParams.SilverPieces = character.SilverPieces + amount
	case "cp":
		updateParams.CopperPieces = character.CopperPieces + amount
	}

	// Perform update
	_, err = queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		log.Printf("Error updating character currency: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error updating currency", characterID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Currency updated successfully", characterID), http.StatusSeeOther)
}

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

func (s *Server) HandleRest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	characterID, err := strconv.ParseInt(r.FormValue("character_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Error fetching character: %v", err)
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	progression := charRules.GetClassProgression(character.Class)
	hitDice := progression.GetHitDice(character.Level)

	parts := strings.Split(hitDice, "d")
	if len(parts) != 2 {
		http.Error(w, "Invalid hit dice format", http.StatusInternalServerError)
		return
	}

	diceSize, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Error parsing dice size: %v", err)
		http.Error(w, "Invalid hit dice format", http.StatusInternalServerError)
		return
	}
	// Roll just one die
	total := rand.IntN(diceSize) + 1

	// Add Constitution bonus using new package
	conMods := ability_scores.CalculateConstitutionModifiers(character.Constitution)
	total += conMods.HitPointMod

	// Calculate new HP, not exceeding max
	newHP := character.CurrentHp + int64(total)
	if newHP > character.MaxHp {
		newHP = character.MaxHp
	}

	// Update character
	updateParams := db.UpdateCharacterParams{
		ID:               characterID,
		UserID:           user.UserID,
		Name:             character.Name,
		Class:            character.Class,
		Level:            character.Level,
		MaxHp:            character.MaxHp,
		CurrentHp:        newHP,
		Strength:         character.Strength,
		Dexterity:        character.Dexterity,
		Constitution:     character.Constitution,
		Intelligence:     character.Intelligence,
		Wisdom:           character.Wisdom,
		Charisma:         character.Charisma,
		ExperiencePoints: character.ExperiencePoints,
		PlatinumPieces:   character.PlatinumPieces,
		GoldPieces:       character.GoldPieces,
		ElectrumPieces:   character.ElectrumPieces,
		SilverPieces:     character.SilverPieces,
		CopperPieces:     character.CopperPieces,
	}

	_, err = queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		log.Printf("Error updating character HP: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error during rest", characterID), http.StatusSeeOther)
		return
	}

	message := fmt.Sprintf("Rest complete! Healed for %d HP", total)
	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=%s", characterID, message), http.StatusSeeOther)
}

func (s *Server) HandleUpdateXP(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	// Get the character to verify ownership and current XP
	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Error fetching character: %v", err)
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Parse XP change
	xpChange, err := strconv.ParseInt(r.Form.Get("xp_change"), 10, 64)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid XP value", characterID), http.StatusSeeOther)
		return
	}

	newXP := character.ExperiencePoints + xpChange
	if newXP < 0 {
		newXP = 0
	}

	// Get class progression
	progression := charRules.GetClassProgression(character.Class)

	// Calculate appropriate level for new XP
	newLevel := progression.GetLevelForXP(newXP)

	// Update character with new XP and level
	updateParams := db.UpdateCharacterParams{
		ID:               characterID,
		UserID:           user.UserID,
		Name:             character.Name,
		Class:            character.Class,
		Level:            newLevel, // Update level based on XP
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

	_, err = queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		log.Printf("Error updating character XP: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error updating XP", characterID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=XP and level updated successfully", characterID), http.StatusSeeOther)
}

func calculateMinimumXPForLevel(class string, level int64) int64 {
	progression := charRules.GetClassProgression(class)
	for _, levelInfo := range progression.Levels {
		if levelInfo.Level == level {
			return levelInfo.XPRequired
		}
	}
	return 0
}

func (s *Server) HandleCharacterDetail(w http.ResponseWriter, r *http.Request) {
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
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Error fetching character: %v", err)
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	inventory, err := queries.GetCharacterInventory(r.Context(), characterID)
	if err != nil {
		log.Printf("Error fetching inventory: %v", err)
		http.Error(w, "Error loading character inventory", http.StatusInternalServerError)
		return
	}

	viewModel := NewCharacterViewModel(character, inventory)

	// Create a function map and add our functions
	funcMap := template.FuncMap{
		"seq": func(start, end int) []int {
			s := make([]int, end-start+1)
			for i := range s {
				s[i] = start + i
			}
			return s
		},
		"GetSavingThrowModifiers": charRules.GetSavingThrowModifiers,
		"add": func(a, b interface{}) int64 {
			// Handle int64 + int
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
			// Handle multiplication between different numeric types
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
			// Handle int64 - int
			switch v := a.(type) {
			case int64:
				switch w := b.(type) {
				case int:
					return v - int64(w)
				case int64:
					return v - w
				}
			case int:
				switch w := b.(type) {
				case int64:
					return int64(v) - w
				case int:
					return int64(v - w)
				}
			}
			return 0
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
		// Add the dict function here:
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
	}

	tmpl, err := template.New("base.html").Funcs(funcMap).ParseFiles(
		"templates/layout/base.html",
		"templates/layout/navbar.html",
		"templates/characters/details.html",
		"templates/characters/_inventory.html",
		"templates/characters/_ability_scores.html",
		"templates/characters/_class_features.html",
		"templates/characters/_combat_stats.html",
		"templates/characters/_saving_throws.html",
		"templates/characters/_character_header.html",
		"templates/characters/_currency_management.html",
	)

	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get weapon masteries if character is a fighter
	var weaponMasteries []db.GetCharacterWeaponMasteriesRow
	if character.Class == "Fighter" {
		weaponMasteries, err = queries.GetCharacterWeaponMasteries(r.Context(), character.ID)
		if err != nil {
			log.Printf("Error fetching weapon masteries: %v", err)
		}
	}

	data := struct {
		IsAuthenticated bool
		Username        string
		Character       CharacterViewModel
		WeaponMasteries []db.GetCharacterWeaponMasteriesRow
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		Character:       viewModel,
		WeaponMasteries: weaponMasteries,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleUpdateHP(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Error fetching character: %v", err)
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	hpChange, err := strconv.ParseInt(r.Form.Get("hp_change"), 10, 64)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid HP value", characterID), http.StatusSeeOther)
		return
	}

	newHP := character.CurrentHp + hpChange
	if newHP > character.MaxHp {
		newHP = character.MaxHp
	}
	if newHP < 0 {
		newHP = 0
	}

	updateParams := db.UpdateCharacterParams{
		ID:               characterID,
		UserID:           user.UserID,
		Name:             character.Name,
		Class:            character.Class,
		Level:            character.Level,
		MaxHp:            character.MaxHp,
		CurrentHp:        newHP,
		Strength:         character.Strength,
		Dexterity:        character.Dexterity,
		Constitution:     character.Constitution,
		Intelligence:     character.Intelligence,
		Wisdom:           character.Wisdom,
		Charisma:         character.Charisma,
		ExperiencePoints: character.ExperiencePoints,
		PlatinumPieces:   character.PlatinumPieces,
		GoldPieces:       character.GoldPieces,
		ElectrumPieces:   character.ElectrumPieces,
		SilverPieces:     character.SilverPieces,
		CopperPieces:     character.CopperPieces,
	}

	_, err = queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		log.Printf("Error updating character HP: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error updating HP", characterID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d", characterID), http.StatusSeeOther)
}

func (s *Server) HandleCharacterList(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by AuthMiddleware)
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all characters for the user
	queries := db.New(s.db)
	characters, err := queries.ListCharactersByUser(r.Context(), user.UserID)
	if err != nil {
		log.Printf("Error fetching characters: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/characters/list.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsAuthenticated bool
		Username        string
		Characters      []db.Character
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		Characters:      characters,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleCharacterCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleCharacterCreateForm(w, r)
	case http.MethodPost:
		s.handleCharacterCreateSubmission(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleCharacterCreateForm(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/characters/create.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsAuthenticated bool
		Username        string
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleCharacterCreateSubmission(w http.ResponseWriter, r *http.Request) {
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

	// Parse and validate form fields
	params := db.CreateCharacterParams{
		UserID: user.UserID,
		Name:   r.Form.Get("name"),
		Class:  r.Form.Get("class"),
	}

	logger.Debug("Processing character creation",
		zap.String("character_name", params.Name),
		zap.String("class", params.Class),
		zap.String("user_id", strconv.FormatInt(user.UserID, 10)))

	// Validate class
	validClasses := map[string]bool{
		"Fighter":     true,
		"Cleric":      true,
		"Magician":    true,
		"Thief":       true,
		"Barbarian":   true,
		"Ranger":      true,
		"Paladin":     true,
		"Druid":       true,
		"Assassin":    true,
		"Necromancer": true,
	}

	if !validClasses[params.Class] {
		logger.Warn("Invalid character class attempted",
			zap.String("attempted_class", params.Class),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Redirect(w, r, "/characters/create?message=Invalid character class", http.StatusSeeOther)
		return
	}

	// Parse level
	level, err := strconv.ParseInt(r.Form.Get("level"), 10, 64)
	if err != nil {
		logger.Error("Invalid level value",
			zap.Error(err),
			zap.String("raw_level", r.Form.Get("level")),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Redirect(w, r, "/characters/create?message=Invalid level value", http.StatusSeeOther)
		return
	}
	params.Level = level

	minimumXP := calculateMinimumXPForLevel(params.Class, level)
	params.ExperiencePoints = minimumXP

	// Parse base HP and constitution
	baseHP, err := strconv.ParseInt(r.Form.Get("max_hp"), 10, 64)
	if err != nil {
		logger.Error("Invalid HP value",
			zap.Error(err),
			zap.String("raw_hp", r.Form.Get("max_hp")),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Redirect(w, r, "/characters/create?message=Invalid HP value", http.StatusSeeOther)
		return
	}

	constitution, err := strconv.ParseInt(r.Form.Get("constitution"), 10, 64)
	if err != nil {
		logger.Error("Invalid constitution value",
			zap.Error(err),
			zap.String("raw_constitution", r.Form.Get("constitution")),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Redirect(w, r, "/characters/create?message=Invalid constitution value", http.StatusSeeOther)
		return
	}
	params.Constitution = constitution

	// Calculate total HP using the rules package
	totalHP, err := calculateTotalHP(baseHP, level, constitution)
	if err != nil {
		logger.Error("Failed to calculate total HP",
			zap.Error(err),
			zap.Int64("base_hp", baseHP),
			zap.Int64("level", level),
			zap.Int64("constitution", constitution))
		http.Redirect(w, r, "/characters/create?message="+err.Error(), http.StatusSeeOther)
		return
	}
	params.MaxHp = totalHP
	params.CurrentHp = totalHP

	// Parse other ability scores
	str, err := strconv.ParseInt(r.Form.Get("strength"), 10, 64)
	if err != nil {
		logger.Error("Invalid strength value",
			zap.Error(err),
			zap.String("raw_strength", r.Form.Get("strength")))
		http.Redirect(w, r, "/characters/create?message=Invalid strength value", http.StatusSeeOther)
		return
	}
	params.Strength = str

	dex, err := strconv.ParseInt(r.Form.Get("dexterity"), 10, 64)
	if err != nil {
		logger.Error("Invalid dexterity value",
			zap.Error(err),
			zap.String("raw_dexterity", r.Form.Get("dexterity")))
		http.Redirect(w, r, "/characters/create?message=Invalid dexterity value", http.StatusSeeOther)
		return
	}
	params.Dexterity = dex

	intel, err := strconv.ParseInt(r.Form.Get("intelligence"), 10, 64)
	if err != nil {
		logger.Error("Invalid intelligence value",
			zap.Error(err),
			zap.String("raw_intelligence", r.Form.Get("intelligence")))
		http.Redirect(w, r, "/characters/create?message=Invalid intelligence value", http.StatusSeeOther)
		return
	}
	params.Intelligence = intel

	wis, err := strconv.ParseInt(r.Form.Get("wisdom"), 10, 64)
	if err != nil {
		logger.Error("Invalid wisdom value",
			zap.Error(err),
			zap.String("raw_wisdom", r.Form.Get("wisdom")))
		http.Redirect(w, r, "/characters/create?message=Invalid wisdom value", http.StatusSeeOther)
		return
	}
	params.Wisdom = wis

	cha, err := strconv.ParseInt(r.Form.Get("charisma"), 10, 64)
	if err != nil {
		logger.Error("Invalid charisma value",
			zap.Error(err),
			zap.String("raw_charisma", r.Form.Get("charisma")))
		http.Redirect(w, r, "/characters/create?message=Invalid charisma value", http.StatusSeeOther)
		return
	}
	params.Charisma = cha

	// Validate ability scores
	abilities := []int64{params.Strength, params.Dexterity, params.Constitution,
		params.Intelligence, params.Wisdom, params.Charisma}
	for _, score := range abilities {
		if score < 3 || score > 18 {
			logger.Warn("Invalid ability score attempted",
				zap.Int64("score", score),
				zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
			http.Redirect(w, r, "/characters/create?message=Ability scores must be between 3 and 18",
				http.StatusSeeOther)
			return
		}
	}

	// Create character in database
	queries := db.New(s.db)
	character, err := queries.CreateCharacter(r.Context(), params)
	if err != nil {
		logger.Error("Failed to create character in database",
			zap.Error(err),
			zap.Any("params", params))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("Character created successfully",
		zap.Int64("character_id", character.ID),
		zap.String("character_name", character.Name),
		zap.String("user_id", strconv.FormatInt(user.UserID, 10)))

	http.Redirect(w, r, "/characters?message=Character created successfully", http.StatusSeeOther)
}

func (s *Server) HandleCharacterEdit(w http.ResponseWriter, r *http.Request) {
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

	switch r.Method {
	case http.MethodGet:
		// Get character from database
		character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
			ID:     characterID,
			UserID: user.UserID,
		})
		if err != nil {
			log.Printf("Error fetching character: %v", err)
			http.Error(w, "Character not found", http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles(
			"templates/layout/base.html",
			"templates/characters/edit.html",
		)
		if err != nil {
			log.Printf("Template parsing error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := struct {
			IsAuthenticated bool
			Username        string
			Character       db.Character
			FlashMessage    string
			CurrentYear     int
		}{
			IsAuthenticated: true,
			Username:        user.Username,
			Character:       character,
			FlashMessage:    r.URL.Query().Get("message"),
			CurrentYear:     time.Now().Year(),
		}

		err = tmpl.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Printf("Template execution error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Parse form values
		strength, _ := strconv.ParseInt(r.Form.Get("strength"), 10, 64)
		dexterity, _ := strconv.ParseInt(r.Form.Get("dexterity"), 10, 64)
		constitution, _ := strconv.ParseInt(r.Form.Get("constitution"), 10, 64)
		intelligence, _ := strconv.ParseInt(r.Form.Get("intelligence"), 10, 64)
		wisdom, _ := strconv.ParseInt(r.Form.Get("wisdom"), 10, 64)
		charisma, _ := strconv.ParseInt(r.Form.Get("charisma"), 10, 64)
		maxHp, _ := strconv.ParseInt(r.Form.Get("max_hp"), 10, 64)
		currentHp, _ := strconv.ParseInt(r.Form.Get("current_hp"), 10, 64)
		level, _ := strconv.ParseInt(r.Form.Get("level"), 10, 64)

		// Validate ability scores
		abilities := []int64{strength, dexterity, constitution, intelligence, wisdom, charisma}
		for _, score := range abilities {
			if score < 3 || score > 18 {
				http.Redirect(w, r, fmt.Sprintf("/characters/edit?id=%d&message=Ability scores must be between 3 and 18", characterID), http.StatusSeeOther)
				return
			}
		}

		// Update character
		_, err = queries.UpdateCharacter(r.Context(), db.UpdateCharacterParams{
			ID:           characterID,
			UserID:       user.UserID,
			Name:         r.Form.Get("name"),
			Class:        r.Form.Get("class"),
			Level:        level,
			MaxHp:        maxHp,
			CurrentHp:    currentHp,
			Strength:     strength,
			Dexterity:    dexterity,
			Constitution: constitution,
			Intelligence: intelligence,
			Wisdom:       wisdom,
			Charisma:     charisma,
		})

		if err != nil {
			log.Printf("Error updating character: %v", err)
			http.Redirect(w, r, fmt.Sprintf("/characters/edit?id=%d&message=Error updating character", characterID), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Character updated successfully", characterID), http.StatusSeeOther)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) HandleUpdateMaxHP(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		log.Printf("Error fetching character: %v", err)
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	maxHPChange, err := strconv.ParseInt(r.Form.Get("max_hp_change"), 10, 64)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid HP value", characterID), http.StatusSeeOther)
		return
	}

	newMaxHP := character.MaxHp + maxHPChange
	if newMaxHP < 1 {
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Maximum HP cannot be less than 1", characterID), http.StatusSeeOther)
		return
	}

	// Ensure current HP doesn't exceed new max HP
	newCurrentHP := character.CurrentHp
	if newCurrentHP > newMaxHP {
		newCurrentHP = newMaxHP
	}

	updateParams := db.UpdateCharacterParams{
		ID:               characterID,
		UserID:           user.UserID,
		Name:             character.Name,
		Class:            character.Class,
		Level:            character.Level,
		MaxHp:            newMaxHP,
		CurrentHp:        newCurrentHP,
		Strength:         character.Strength,
		Dexterity:        character.Dexterity,
		Constitution:     character.Constitution,
		Intelligence:     character.Intelligence,
		Wisdom:           character.Wisdom,
		Charisma:         character.Charisma,
		ExperiencePoints: character.ExperiencePoints,
		PlatinumPieces:   character.PlatinumPieces,
		GoldPieces:       character.GoldPieces,
		ElectrumPieces:   character.ElectrumPieces,
		SilverPieces:     character.SilverPieces,
		CopperPieces:     character.CopperPieces,
	}

	_, err = queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		log.Printf("Error updating character max HP: %v", err)
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error updating maximum HP", characterID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d", characterID), http.StatusSeeOther)
}
