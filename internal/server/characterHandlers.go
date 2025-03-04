package server

import (
	"database/sql"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/marbh56/mordezzan/internal/currency"
	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/logger"
	"github.com/marbh56/mordezzan/internal/rules"
	"github.com/marbh56/mordezzan/internal/rules/ability_scores"
	charRules "github.com/marbh56/mordezzan/internal/rules/character"
	"github.com/marbh56/mordezzan/internal/rules/combat"
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

// Helper to render updated currency section
func renderCurrencySection(w http.ResponseWriter, character CharacterViewModel, message string) {
	tmpl, err := template.ParseFiles("templates/characters/_currency_section.html")
	if err != nil {
		logger.Error("Template parsing error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Character CharacterViewModel
		Message   string
	}{
		Character: character,
		Message:   message,
	}

	err = tmpl.ExecuteTemplate(w, "_currency_section", data)
	if err != nil {
		logger.Error("Template execution error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleRest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Error("Invalid method for rest handler",
			zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt to rest handler")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	characterID, err := strconv.ParseInt(r.FormValue("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID in rest handler",
			zap.Error(err),
			zap.String("raw_id", r.FormValue("character_id")))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Failed to fetch character for rest",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("user_id", user.UserID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	progression := charRules.GetClassProgression(character.Class)
	hitDice := progression.GetHitDice(character.Level)

	parts := strings.Split(hitDice, "d")
	if len(parts) != 2 {
		logger.Error("Invalid hit dice format",
			zap.String("hit_dice", hitDice),
			zap.Int64("character_id", characterID))
		http.Error(w, "Invalid hit dice format", http.StatusInternalServerError)
		return
	}

	diceSize, err := strconv.Atoi(parts[1])
	if err != nil {
		logger.Error("Error parsing dice size",
			zap.Error(err),
			zap.String("dice_part", parts[1]),
			zap.Int64("character_id", characterID))
		http.Error(w, "Invalid hit dice format", http.StatusInternalServerError)
		return
	}

	total := rand.IntN(diceSize) + 1
	conMods := ability_scores.CalculateConstitutionModifiers(character.Constitution)
	total += conMods.HitPointMod

	newHP := character.CurrentHp + int64(total)
	if newHP > character.MaxHp {
		newHP = character.MaxHp
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
		logger.Error("Failed to update character HP after rest",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("new_hp", newHP),
			zap.Int64("healing", int64(total)))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error during rest", characterID), http.StatusSeeOther)
		return
	}

	message := fmt.Sprintf("Rest complete! Healed for %d HP", total)
	logger.Info("Character rest successful",
		zap.Int64("character_id", characterID),
		zap.Int64("healing", int64(total)),
		zap.Int64("new_hp", newHP))
	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=%s", characterID, message), http.StatusSeeOther)
}

func calculateMinimumXPForLevel(class string, level int64) int64 {
	progression := charRules.GetClassProgression(class)
	for _, levelInfo := range progression.Levels {
		if levelInfo.Level == level {
			return levelInfo.XPRequired
		}
	}
	logger.Debug("No XP requirement found for level",
		zap.String("class", class),
		zap.Int64("level", level))
	return 0
}

func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

func (s *Server) HandleCharacterList(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	queries := db.New(s.db)
	characters, err := queries.ListCharactersByUser(r.Context(), user.UserID)
	if err != nil {
		logger.Error("Failed to fetch characters",
			zap.Error(err),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Debug("Retrieved characters for user",
		zap.String("user_id", strconv.FormatInt(user.UserID, 10)),
		zap.Int("character_count", len(characters)))

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/characters/list.html",
	)
	if err != nil {
		logger.Error("Template parsing failed",
			zap.Error(err),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
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

	if err = tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		logger.Error("Template execution failed",
			zap.Error(err),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Debug("Character list page rendered successfully",
		zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
}

func (s *Server) HandleCharacterCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleCharacterCreateForm(w, r)
	case http.MethodPost:
		s.handleCharacterCreateSubmission(w, r)
	default:
		logger.Error("Invalid method for character creation",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleCharacterCreateForm(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/layout/base.html",
		"templates/characters/create.html",
	)
	if err != nil {
		logger.Error("Template parsing failed",
			zap.Error(err),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
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

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		logger.Error("Template execution failed",
			zap.Error(err),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Debug("Character creation form rendered",
		zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
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
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	characterIDStr := r.URL.Query().Get("id")
	characterID, err := strconv.ParseInt(characterIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", characterIDStr))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)

	switch r.Method {
	case http.MethodGet:
		character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
			ID:     characterID,
			UserID: user.UserID,
		})
		if err != nil {
			logger.Error("Failed to fetch character",
				zap.Error(err),
				zap.Int64("character_id", characterID),
				zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
			http.Error(w, "Character not found", http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles(
			"templates/layout/base.html",
			"templates/characters/edit.html",
		)
		if err != nil {
			logger.Error("Template parsing failed",
				zap.Error(err),
				zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
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

		if err = tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
			logger.Error("Template execution failed",
				zap.Error(err),
				zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Debug("Character edit form rendered",
			zap.Int64("character_id", characterID),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			logger.Error("Failed to parse form",
				zap.Error(err),
				zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Parse and validate ability scores
		abilities := make(map[string]int64)
		for _, field := range []string{"strength", "dexterity", "constitution", "intelligence", "wisdom", "charisma"} {
			score, err := strconv.ParseInt(r.Form.Get(field), 10, 64)
			if err != nil {
				logger.Error("Invalid ability score",
					zap.Error(err),
					zap.String("field", field),
					zap.String("raw_value", r.Form.Get(field)))
				http.Redirect(w, r, fmt.Sprintf("/characters/edit?id=%d&message=Invalid %s value", characterID, field), http.StatusSeeOther)
				return
			}
			if score < 3 || score > 18 {
				logger.Warn("Ability score out of range",
					zap.String("field", field),
					zap.Int64("value", score))
				http.Redirect(w, r, fmt.Sprintf("/characters/edit?id=%d&message=Ability scores must be between 3 and 18", characterID), http.StatusSeeOther)
				return
			}
			abilities[field] = score
		}

		maxHp, _ := strconv.ParseInt(r.Form.Get("max_hp"), 10, 64)
		currentHp, _ := strconv.ParseInt(r.Form.Get("current_hp"), 10, 64)
		level, _ := strconv.ParseInt(r.Form.Get("level"), 10, 64)

		updateParams := db.UpdateCharacterParams{
			ID:           characterID,
			UserID:       user.UserID,
			Name:         r.Form.Get("name"),
			Class:        r.Form.Get("class"),
			Level:        level,
			MaxHp:        maxHp,
			CurrentHp:    currentHp,
			Strength:     abilities["strength"],
			Dexterity:    abilities["dexterity"],
			Constitution: abilities["constitution"],
			Intelligence: abilities["intelligence"],
			Wisdom:       abilities["wisdom"],
			Charisma:     abilities["charisma"],
		}

		_, err = queries.UpdateCharacter(r.Context(), updateParams)
		if err != nil {
			logger.Error("Failed to update character",
				zap.Error(err),
				zap.Int64("character_id", characterID),
				zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
			http.Redirect(w, r, fmt.Sprintf("/characters/edit?id=%d&message=Error updating character", characterID), http.StatusSeeOther)
			return
		}

		logger.Info("Character updated successfully",
			zap.Int64("character_id", characterID),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))

		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Character updated successfully", characterID), http.StatusSeeOther)

	default:
		logger.Error("Invalid method for character edit",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) HandleDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Error("Invalid method for character deletion",
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
			zap.Error(err))
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", r.Form.Get("character_id")))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	// Verify character exists and belongs to user before deletion
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

	// Delete the character
	err = queries.DeleteCharacter(r.Context(), db.DeleteCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})

	if err != nil {
		logger.Error("Failed to delete character",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("user_id", user.UserID))
		http.Error(w, "Error deleting character", http.StatusInternalServerError)
		return
	}

	logger.Info("Character deleted successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("user_id", user.UserID))

	http.Redirect(w, r, "/characters?message=Character deleted successfully", http.StatusSeeOther)
}

func NewSafeCharacterViewModel(c db.Character, inventory []db.GetCharacterInventoryItemsRow) CharacterViewModel {
	vm := CharacterViewModel{
		ID:               c.ID,
		UserID:           c.UserID,
		Name:             c.Name,
		Class:            c.Class,
		Level:            c.Level,
		MaxHp:            c.MaxHp,
		CurrentHp:        c.CurrentHp,
		Strength:         c.Strength,
		Dexterity:        c.Dexterity,
		Constitution:     c.Constitution,
		Intelligence:     c.Intelligence,
		Wisdom:           c.Wisdom,
		Charisma:         c.Charisma,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		ExperiencePoints: c.ExperiencePoints,
		PlatinumPieces:   c.PlatinumPieces,
		GoldPieces:       c.GoldPieces,
		ElectrumPieces:   c.ElectrumPieces,
		SilverPieces:     c.SilverPieces,
		CopperPieces:     c.CopperPieces,

		// Initialize modifiers
		StrengthModifiers:     ability_scores.CalculateStrengthModifiers(c.Strength),
		DexterityModifiers:    ability_scores.CalculateDexterityModifiers(c.Dexterity),
		ConstitutionModifiers: ability_scores.CalculateConstitutionModifiers(c.Constitution),
		IntelligenceModifiers: ability_scores.CalculateIntelligenceModifiers(c.Intelligence),
		WisdomModifiers:       ability_scores.CalculateWisdomModifiers(c.Wisdom),
		CharismaModifiers:     ability_scores.CalculateCharismaModifiers(c.Charisma),

		// Initialize inventory containers
		ContainerItems: make(map[int64][]InventoryItem),
	}

	if classGetsFighterBonus(c.Class) {
		vm.StrengthModifiers.ExtraordinaryFeat += 8
	}

	// Get class progression
	progression := charRules.GetClassProgression(vm.Class)

	// Calculate XP needed for next level
	vm.XPNeeded = progression.GetXPForNextLevel(c.ExperiencePoints)

	// Get XP required for current level
	for _, levelInfo := range progression.Levels {
		if levelInfo.XPRequired > c.ExperiencePoints {
			vm.NextLevelXP = levelInfo.XPRequired
			break
		}
	}

	// Calculate base AC
	baseAC := 9
	var armorAC int64
	var shieldBonus int64

	// Check equipped items for armor and shield with safety checks
	for _, item := range inventory {
		if item.EquipmentSlotID.Valid {
			if item.ItemType == "armor" {
				// For armor, find the armor class value safely
				if acVal, ok := safeGetArmorClass(item); ok {
					armorAC = acVal
				}
			} else if item.ItemType == "shield" {
				// For shields, get defense bonus safely
				if bonus, ok := safeGetDefenseBonus(item.DefenseBonus); ok {
					shieldBonus = bonus
				}
			}
		}
	}

	// If armor is equipped, use its AC instead of base AC
	if armorAC > 0 {
		baseAC = int(armorAC)
	}

	// Apply shield bonus if any
	totalAC := baseAC - int(shieldBonus)

	// Apply Dexterity modifier
	totalAC -= vm.DexterityModifiers.DefenseAdj

	vm.ArmorClass = totalAC

	// Initialize inventory stats with encumbrance thresholds
	encumbranceThresholds := rules.CalculateEncumbranceThresholds(c.Strength, c.Constitution)
	vm.InventoryStats = InventoryStats{
		BaseEncumbered:      encumbranceThresholds.BaseEncumbered,
		BaseHeavyEncumbered: encumbranceThresholds.BaseHeavyEncumbered,
		MaximumCapacity:     encumbranceThresholds.MaximumCapacity,
	}

	// Add coin weight to total weight - using our currency package
	purse := currency.Purse{
		PlatinumPieces: c.PlatinumPieces,
		GoldPieces:     c.GoldPieces,
		ElectrumPieces: c.ElectrumPieces,
		SilverPieces:   c.SilverPieces,
		CopperPieces:   c.CopperPieces,
	}
	coinageWeight := int(currency.GetTotalWeight(&purse) + 0.5) // Round to nearest pound
	vm.InventoryStats.CoinWeight = coinageWeight

	// Process each inventory item with safety checks
	for _, item := range inventory {
		// Build inventory item with safe conversions
		invItem := InventoryItem{
			ID:              item.ID,
			CharacterID:     item.CharacterID,
			ItemType:        item.ItemType,
			ItemID:          item.ItemID,
			ItemName:        safeGetItemName(item.ItemName),
			ItemWeight:      safeGetItemWeight(item.ItemWeight),
			Quantity:        item.Quantity,
			ContainerID:     item.ContainerID,
			EquipmentSlotID: item.EquipmentSlotID,
			SlotName:        item.SlotName,
			Notes:           item.Notes,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
		}

		// Safely handle type-specific fields
		if item.ItemType == "weapon" || item.ItemType == "ranged_weapon" {
			invItem.Damage = safeGetNullString(item.Damage)
			invItem.AttacksPerRound = safeGetNullString(item.AttacksPerRound)
		}

		if item.ItemType == "armor" {
			if val, ok := safeGetInt64(item.MovementRate); ok {
				invItem.MovementRate = sql.NullInt64{Int64: val, Valid: true}
			}
		}

		if item.ItemType == "shield" {
			invItem.DefenseBonus = item.DefenseBonus // This is already an interface{}
		}

		// Calculate total weight for this item
		itemTotalWeight := invItem.ItemWeight * int(invItem.Quantity)

		// Distribute the item to the appropriate collection
		if invItem.EquipmentSlotID.Valid {
			vm.EquippedItems = append(vm.EquippedItems, invItem)
			vm.InventoryStats.EquippedWeight += itemTotalWeight
		} else if invItem.ContainerID.Valid {
			containerID := invItem.ContainerID.Int64
			vm.ContainerItems[containerID] = append(vm.ContainerItems[containerID], invItem)
			vm.InventoryStats.ContainersWeight += itemTotalWeight
		} else {
			vm.CarriedItems = append(vm.CarriedItems, invItem)
			vm.InventoryStats.CarriedWeight += itemTotalWeight
		}
	}

	// Calculate total weight and encumbrance level
	vm.InventoryStats.TotalWeight = vm.InventoryStats.EquippedWeight +
		vm.InventoryStats.CarriedWeight +
		vm.InventoryStats.ContainersWeight +
		vm.InventoryStats.CoinWeight

	// Determine encumbrance level based on TOTAL weight (including coins)
	switch {
	case vm.InventoryStats.TotalWeight > vm.InventoryStats.MaximumCapacity:
		vm.InventoryStats.EncumbranceLevel = "Over"
	case vm.InventoryStats.TotalWeight > vm.InventoryStats.BaseHeavyEncumbered:
		vm.InventoryStats.EncumbranceLevel = "Heavy"
	case vm.InventoryStats.TotalWeight > vm.InventoryStats.BaseEncumbered:
		vm.InventoryStats.EncumbranceLevel = "Encumbered"
	default:
		vm.InventoryStats.EncumbranceLevel = "None"
	}

	// Calculate FA and generate combat matrix row
	fa := combat.CalculateFightingAbility(c.Class, c.Level)
	vm.CombatMatrix = make([]int64, 19) // -9 to 9 AC
	for ac := -9; ac <= 9; ac++ {
		vm.CombatMatrix[ac+9] = combat.GetTargetNumber(fa, int64(ac))
	}

	// Get saving throw value
	progression = charRules.GetClassProgression(vm.Class)
	vm.SavingThrow = progression.GetSavingThrow(vm.Level)

	return vm
}

// safeGetItemName safely extracts an item name from an interface
func safeGetItemName(v interface{}) string {
	if v == nil {
		return "Unknown Item"
	}

	switch val := v.(type) {
	case string:
		return val
	case sql.NullString:
		if val.Valid {
			return val.String
		}
	}

	// Fallback: convert to string representation
	return fmt.Sprintf("%v", v)
}

// safeGetItemWeight safely extracts an item weight from an interface
func safeGetItemWeight(v interface{}) int {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case sql.NullInt64:
		if val.Valid {
			return int(val.Int64)
		}
	}

	return 0
}

// safeGetNullString safely converts an interface to sql.NullString
func safeGetNullString(v interface{}) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}

	switch val := v.(type) {
	case string:
		return sql.NullString{String: val, Valid: val != ""}
	case sql.NullString:
		return val
	}

	// Try to convert to string
	str := fmt.Sprintf("%v", v)
	return sql.NullString{String: str, Valid: str != ""}
}

// safeGetInt64 safely extracts an int64 from an interface
func safeGetInt64(v interface{}) (int64, bool) {
	if v == nil {
		return 0, false
	}

	switch val := v.(type) {
	case int:
		return int64(val), true
	case int64:
		return val, true
	case float64:
		return int64(val), true
	case sql.NullInt64:
		return val.Int64, val.Valid
	}

	return 0, false
}

// safeGetDefenseBonus safely extracts a defense bonus from an interface
func safeGetDefenseBonus(v interface{}) (int64, bool) {
	if v == nil {
		return 0, false
	}

	switch val := v.(type) {
	case int:
		return int64(val), true
	case int64:
		return val, true
	case float64:
		return int64(val), true
	}

	return 0, false
}

// safeGetArmorClass safely extracts armor class from an item
func safeGetArmorClass(item db.GetCharacterInventoryItemsRow) (int64, bool) {
	// This is a placeholder - in your actual code you would need to
	// determine how armor class is stored in your inventory items
	// This might be a property directly on the item or derivable from other properties

	// Example implementation:
	if item.ItemType != "armor" {
		return 0, false
	}

	// Try to find armor class in item properties
	// This would depend on your data structure
	return 7, true // Default value
}

func (s *Server) HandleCharacterDetail(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	characterIDStr := r.URL.Query().Get("id")
	characterID, err := strconv.ParseInt(characterIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.String("raw_id", characterIDStr),
			zap.Error(err))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Error fetching character",
			zap.Int64("character_id", characterID),
			zap.Int64("user_id", user.UserID),
			zap.Error(err))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Check if character has any inventory items first
	var count int
	countErr := s.db.QueryRowContext(r.Context(),
		"SELECT COUNT(*) FROM character_inventory WHERE character_id = ?",
		characterID).Scan(&count)

	var inventory []db.GetCharacterInventoryItemsRow

	if countErr != nil {
		logger.Warn("Failed to check inventory count, proceeding with empty inventory",
			zap.Error(countErr),
			zap.Int64("character_id", characterID))
		inventory = []db.GetCharacterInventoryItemsRow{}
	} else if count == 0 {
		// No inventory items, use empty slice
		logger.Info("Character has no inventory items",
			zap.Int64("character_id", characterID))
		inventory = []db.GetCharacterInventoryItemsRow{}
	} else {
		// Fetch inventory with error handling
		inventory, err = queries.GetCharacterInventoryItems(r.Context(), characterID)
		if err != nil {
			logger.Error("Error fetching inventory, proceeding with empty inventory",
				zap.Int64("character_id", characterID),
				zap.Error(err))
			inventory = []db.GetCharacterInventoryItemsRow{}
		}
	}

	// Create a robust view model using the safe function with panic recovery
	var viewModel CharacterViewModel
	func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic recovered while creating character view model",
					zap.Any("panic", r),
					zap.Int64("character_id", characterID))
				// We'll create a minimal view model when we return from this function
			}
		}()

		viewModel = NewSafeCharacterViewModel(character, inventory)
	}()

	// If viewModel is empty (due to panic), create a minimal one
	if viewModel.ID == 0 {
		logger.Warn("Creating minimal view model after error",
			zap.Int64("character_id", characterID))
		viewModel = CharacterViewModel{
			ID:               character.ID,
			UserID:           character.UserID,
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
			CreatedAt:        character.CreatedAt,
			UpdatedAt:        character.UpdatedAt,
			ArmorClass:       9, // Default AC
			ContainerItems:   make(map[int64][]InventoryItem),
			InventoryStats: InventoryStats{
				EncumbranceLevel: "None",
			},
		}
	}

	funcMap := template.FuncMap{
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
		"seq": func(start, end int) []int {
			s := make([]int, end-start+1)
			for i := range s {
				s[i] = start + i
			}
			return s
		},
		"GetSavingThrowModifiers": charRules.GetSavingThrowModifiers,
		"add": func(a, b interface{}) int64 {
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
		"contains": containsString,
	}

	// Try to parse templates with error handling
	tmpl, err := template.New("base.html").Funcs(funcMap).ParseFiles(
		"templates/layout/base.html",
		"templates/characters/details.html",
		"templates/characters/_inventory.html",
		"templates/characters/_ability_scores.html",
		"templates/characters/_class_features.html",
		"templates/characters/_combat_stats.html",
		"templates/characters/_saving_throws.html",
		"templates/characters/_character_header.html",
		"templates/characters/_currency_management.html",
		"templates/characters/_hp_display.html",
		"templates/characters/_hp_section.html",
		"templates/characters/_currency_section.html",
		"templates/characters/inventory_modal.html",
	)

	if err != nil {
		logger.Error("Template parsing error",
			zap.Error(err))
		http.Error(w, "Internal Server Error: Template parsing failed", http.StatusInternalServerError)
		return
	}

	data := struct {
		IsAuthenticated bool
		Username        string
		Character       CharacterViewModel
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		Character:       viewModel,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	// Add a special message if we had inventory issues but are still rendering
	if len(inventory) == 0 && count > 0 {
		data.FlashMessage = "There was an issue loading some inventory items. Please contact support if this persists."
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		logger.Error("Template execution error",
			zap.Error(err))
		http.Error(w, "Internal Server Error: Template execution failed", http.StatusInternalServerError)
		return
	}

	logger.Info("Character detail rendered successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("user_id", user.UserID),
		zap.Int("inventory_count", len(inventory)))
}

func (s *Server) HandleHPForm(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	characterID, err := strconv.ParseInt(r.URL.Query().Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	// Get character from database to verify ownership
	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Error fetching character", zap.Error(err))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Render the HP update form template
	tmpl, err := template.ParseFiles("templates/characters/_hp_update_form.html")
	if err != nil {
		logger.Error("Template error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "hp_update_form", character)
	if err != nil {
		logger.Error("Template execution error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleMaxHPForm renders the form to update maximum HP
func (s *Server) HandleMaxHPForm(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	characterID, err := strconv.ParseInt(r.URL.Query().Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	// Get character from database to verify ownership
	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Error fetching character", zap.Error(err))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	// Render the Max HP update form template
	tmpl, err := template.ParseFiles("templates/characters/_maxhp_update_form.html")
	if err != nil {
		logger.Error("Template error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "maxhp_update_form", character)
	if err != nil {
		logger.Error("Template execution error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleHPCancel clears the form container
func (s *Server) HandleHPCancel(w http.ResponseWriter, r *http.Request) {
	// Just return an empty response to clear the form container
	w.WriteHeader(http.StatusOK)
}

// UpdatedHandleUpdateHP handles updating the character's current HP
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
		logger.Error("Failed to parse form", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Error fetching character", zap.Error(err))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	hpChange, err := strconv.ParseInt(r.Form.Get("hp_change"), 10, 64)
	if err != nil {
		logger.Warn("Invalid HP change value", zap.Error(err))
		renderHPSection(w, character, "Invalid HP value")
		return
	}

	newHP := character.CurrentHp + hpChange
	if newHP > character.MaxHp {
		newHP = character.MaxHp
		logger.Info("HP change capped at max HP", zap.Int64("character_id", characterID))
	}
	if newHP < 0 {
		newHP = 0
		logger.Info("HP change floored at 0", zap.Int64("character_id", characterID))
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

	updatedCharacter, err := queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		logger.Error("Failed to update character HP", zap.Error(err))
		renderHPSection(w, character, "Error updating HP")
		return
	}

	logger.Info("Character HP updated successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("old_hp", character.CurrentHp),
		zap.Int64("new_hp", newHP))

	message := fmt.Sprintf("HP updated by %+d", hpChange)
	renderHPSection(w, updatedCharacter, message)
}

func renderHPSection(w http.ResponseWriter, character db.Character, message string) {
	tmpl, err := template.ParseFiles(
		"templates/characters/_hp_display.html",
		"templates/characters/_hp_section.html",
	)
	if err != nil {
		logger.Error("Template parsing error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Character    db.Character
		Message      string
		FlashMessage string // Add this field to match what the template expects
	}{
		Character:    character,
		Message:      message,
		FlashMessage: message, // Set FlashMessage to the message value
	}

	err = tmpl.ExecuteTemplate(w, "hp_display_section", data)
	if err != nil {
		logger.Error("Template execution error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Modified function to handle updating maximum HP
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
		logger.Error("Failed to parse form", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID", zap.Error(err))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)
	character, err := queries.GetCharacter(r.Context(), db.GetCharacterParams{
		ID:     characterID,
		UserID: user.UserID,
	})
	if err != nil {
		logger.Error("Error fetching character", zap.Error(err))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	maxHPChange, err := strconv.ParseInt(r.Form.Get("max_hp_change"), 10, 64)
	if err != nil {
		logger.Warn("Invalid max HP change value", zap.Error(err))
		renderHPSection(w, character, "Invalid HP value")
		return
	}

	newMaxHP := character.MaxHp + maxHPChange
	if newMaxHP < 1 {
		logger.Warn("Attempted to set max HP below 1",
			zap.Int64("character_id", characterID),
			zap.Int64("attempted_max_hp", newMaxHP))
		renderHPSection(w, character, "Maximum HP cannot be less than 1")
		return
	}

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

	updatedCharacter, err := queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		logger.Error("Failed to update character max HP", zap.Error(err))
		renderHPSection(w, character, "Error updating maximum HP")
		return
	}

	logger.Info("Character max HP updated successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("old_max_hp", character.MaxHp),
		zap.Int64("new_max_hp", newMaxHP))

	message := fmt.Sprintf("Maximum HP changed by %+d", maxHPChange)
	renderHPSection(w, updatedCharacter, message)
}

func (s *Server) renderCharacterDetail(w http.ResponseWriter, r *http.Request, user *db.GetSessionRow, character db.Character, viewModel CharacterViewModel) error {
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
	}

	// Parse templates - make sure to include the inventory_modal.html template
	tmpl, err := template.New("base.html").Funcs(funcMap).ParseFiles(
		"templates/layout/base.html",
		"templates/characters/details.html",
		"templates/characters/_inventory.html",
		"templates/characters/_ability_scores.html",
		"templates/characters/_class_features.html",
		"templates/characters/_combat_stats.html",
		"templates/characters/_saving_throws.html",
		"templates/characters/_character_header.html",
		"templates/characters/_currency_management.html",
		"templates/characters/_hp_display.html",
		"templates/characters/_hp_section.html",
		"templates/characters/_currency_section.html",
		"templates/characters/inventory_modal.html", // Add this line to include the inventory modal template
	)

	if err != nil {
		logger.Error("Template parsing error", zap.Error(err))
		return err
	}

	data := struct {
		IsAuthenticated bool
		Username        string
		Character       CharacterViewModel
		FlashMessage    string
		Message         string // Add this field for backward compatibility
		CurrentYear     int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		Character:       viewModel,
		FlashMessage:    r.URL.Query().Get("message"),
		Message:         r.URL.Query().Get("message"), // Set to the same value as FlashMessage
		CurrentYear:     time.Now().Year(),
	}

	// Add a special message if we had inventory issues but are still rendering
	if count, err := getCharacterInventoryCount(s.db, character.ID); err == nil && count > 0 && len(viewModel.CarriedItems) == 0 {
		data.FlashMessage = "There was an issue loading some inventory items. Please contact support if this persists."
	}

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		logger.Error("Template execution error", zap.Error(err))
		return err
	}

	logger.Info("Character detail rendered successfully",
		zap.Int64("character_id", character.ID),
		zap.Int64("user_id", user.UserID))

	return nil
}

// Helper function to get inventory count for a character
func getCharacterInventoryCount(db *sql.DB, characterID int64) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM character_inventory WHERE character_id = ?", characterID).Scan(&count)
	return count, err
}
