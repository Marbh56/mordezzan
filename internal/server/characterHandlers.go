package server

import (
	"fmt"
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
		logger.Warn("Invalid HTTP method",
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

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.String("raw_id", r.Form.Get("character_id")),
			zap.Error(err))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseInt(r.Form.Get("amount"), 10, 64)
	if err != nil {
		logger.Warn("Invalid currency amount",
			zap.String("raw_amount", r.Form.Get("amount")),
			zap.Int64("character_id", characterID))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid amount", characterID), http.StatusSeeOther)
		return
	}

	denomination := r.Form.Get("denomination")
	if !isValidDenomination(denomination) {
		logger.Warn("Invalid currency denomination",
			zap.String("denomination", denomination),
			zap.Int64("character_id", characterID))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid denomination", characterID), http.StatusSeeOther)
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
			zap.Error(err))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

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

	var oldAmount, newAmount int64
	switch denomination {
	case "pp":
		oldAmount = character.PlatinumPieces
		updateParams.PlatinumPieces = character.PlatinumPieces + amount
		newAmount = updateParams.PlatinumPieces
	case "gp":
		oldAmount = character.GoldPieces
		updateParams.GoldPieces = character.GoldPieces + amount
		newAmount = updateParams.GoldPieces
	case "ep":
		oldAmount = character.ElectrumPieces
		updateParams.ElectrumPieces = character.ElectrumPieces + amount
		newAmount = updateParams.ElectrumPieces
	case "sp":
		oldAmount = character.SilverPieces
		updateParams.SilverPieces = character.SilverPieces + amount
		newAmount = updateParams.SilverPieces
	case "cp":
		oldAmount = character.CopperPieces
		updateParams.CopperPieces = character.CopperPieces + amount
		newAmount = updateParams.CopperPieces
	}

	_, err = queries.UpdateCharacter(r.Context(), updateParams)
	if err != nil {
		logger.Error("Error updating character currency",
			zap.Int64("character_id", characterID),
			zap.String("denomination", denomination),
			zap.Int64("amount_change", amount),
			zap.Error(err))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error updating currency", characterID), http.StatusSeeOther)
		return
	}

	logger.Info("Character currency updated",
		zap.Int64("character_id", characterID),
		zap.String("denomination", denomination),
		zap.Int64("old_amount", oldAmount),
		zap.Int64("new_amount", newAmount),
		zap.Int64("change", amount))

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

func (s *Server) HandleUpdateXP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", r.Form.Get("character_id")))
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
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("user_id", user.UserID))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	xpChange, err := strconv.ParseInt(r.Form.Get("xp_change"), 10, 64)
	if err != nil {
		logger.Warn("Invalid XP value provided",
			zap.Error(err),
			zap.String("raw_xp", r.Form.Get("xp_change")),
			zap.Int64("character_id", characterID))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid XP value", characterID), http.StatusSeeOther)
		return
	}

	newXP := character.ExperiencePoints + xpChange
	if newXP < 0 {
		newXP = 0
	}

	progression := charRules.GetClassProgression(character.Class)
	newLevel := progression.GetLevelForXP(newXP)

	updateParams := db.UpdateCharacterParams{
		ID:               characterID,
		UserID:           user.UserID,
		Name:             character.Name,
		Class:            character.Class,
		Level:            newLevel,
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
		logger.Error("Failed to update character XP",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("new_xp", newXP),
			zap.Int64("new_level", newLevel))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error updating XP", characterID), http.StatusSeeOther)
		return
	}

	logger.Info("Character XP updated successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("old_xp", character.ExperiencePoints),
		zap.Int64("new_xp", newXP),
		zap.Int64("old_level", character.Level),
		zap.Int64("new_level", newLevel))

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=XP and level updated successfully", characterID), http.StatusSeeOther)
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

	inventory, err := queries.GetCharacterInventory(r.Context(), characterID)
	if err != nil {
		logger.Error("Error fetching inventory",
			zap.Int64("character_id", characterID),
			zap.Error(err))
		http.Error(w, "Error loading character inventory", http.StatusInternalServerError)
		return
	}

	viewModel := NewCharacterViewModel(character, inventory)

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
		logger.Error("Template parsing error",
			zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var weaponMasteries []db.GetCharacterWeaponMasteriesRow
	if character.Class == "Fighter" {
		weaponMasteries, err = queries.GetCharacterWeaponMasteries(r.Context(), character.ID)
		if err != nil {
			logger.Error("Error fetching weapon masteries",
				zap.Int64("character_id", character.ID),
				zap.Error(err))
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
		logger.Error("Template execution error",
			zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleUpdateHP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Error("Invalid method for HP update",
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

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", r.Form.Get("character_id")),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
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
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	hpChange, err := strconv.ParseInt(r.Form.Get("hp_change"), 10, 64)
	if err != nil {
		logger.Warn("Invalid HP change value",
			zap.Error(err),
			zap.String("raw_value", r.Form.Get("hp_change")),
			zap.Int64("character_id", characterID))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid HP value", characterID), http.StatusSeeOther)
		return
	}

	newHP := character.CurrentHp + hpChange
	if newHP > character.MaxHp {
		newHP = character.MaxHp
		logger.Info("HP change capped at max HP",
			zap.Int64("character_id", characterID),
			zap.Int64("attempted_hp", newHP),
			zap.Int64("max_hp", character.MaxHp))
	}
	if newHP < 0 {
		newHP = 0
		logger.Info("HP change floored at 0",
			zap.Int64("character_id", characterID),
			zap.Int64("attempted_hp", newHP))
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
		logger.Error("Failed to update character HP",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("new_hp", newHP),
			zap.String("user_id", strconv.FormatInt(user.UserID, 10)))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error updating HP", characterID), http.StatusSeeOther)
		return
	}

	logger.Info("Character HP updated successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("old_hp", character.CurrentHp),
		zap.Int64("new_hp", newHP),
		zap.Int64("hp_change", hpChange),
		zap.String("user_id", strconv.FormatInt(user.UserID, 10)))

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d", characterID), http.StatusSeeOther)
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

func (s *Server) HandleUpdateMaxHP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Error("Invalid method for max HP update",
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

	characterID, err := strconv.ParseInt(r.Form.Get("character_id"), 10, 64)
	if err != nil {
		logger.Error("Invalid character ID",
			zap.Error(err),
			zap.String("raw_id", r.Form.Get("character_id")))
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	queries := db.New(s.db)
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

	maxHPChange, err := strconv.ParseInt(r.Form.Get("max_hp_change"), 10, 64)
	if err != nil {
		logger.Warn("Invalid max HP change value",
			zap.Error(err),
			zap.String("raw_value", r.Form.Get("max_hp_change")))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Invalid HP value", characterID), http.StatusSeeOther)
		return
	}

	newMaxHP := character.MaxHp + maxHPChange
	if newMaxHP < 1 {
		logger.Warn("Attempted to set max HP below 1",
			zap.Int64("character_id", characterID),
			zap.Int64("attempted_max_hp", newMaxHP))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Maximum HP cannot be less than 1", characterID), http.StatusSeeOther)
		return
	}

	newCurrentHP := character.CurrentHp
	if newCurrentHP > newMaxHP {
		logger.Info("Adjusting current HP to new max HP",
			zap.Int64("character_id", characterID),
			zap.Int64("old_current_hp", newCurrentHP),
			zap.Int64("new_max_hp", newMaxHP))
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
		logger.Error("Failed to update character max HP",
			zap.Error(err),
			zap.Int64("character_id", characterID),
			zap.Int64("new_max_hp", newMaxHP))
		http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d&message=Error updating maximum HP", characterID), http.StatusSeeOther)
		return
	}

	logger.Info("Character max HP updated successfully",
		zap.Int64("character_id", characterID),
		zap.Int64("old_max_hp", character.MaxHp),
		zap.Int64("new_max_hp", newMaxHP),
		zap.Int64("old_current_hp", character.CurrentHp),
		zap.Int64("new_current_hp", newCurrentHP))

	http.Redirect(w, r, fmt.Sprintf("/characters/detail?id=%d", characterID), http.StatusSeeOther)
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
