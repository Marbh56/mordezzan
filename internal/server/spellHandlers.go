package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/marbh56/mordezzan/internal/logger"
	"github.com/marbh56/mordezzan/internal/rules/spells"
	"go.uber.org/zap"
)

// HandleSpellList shows all spells by class and level
func (s *Server) HandleSpellList(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	class := r.URL.Query().Get("class")
	levelStr := r.URL.Query().Get("level")

	// Default to mag if no class specified
	if class == "" {
		class = "mag"
	}

	// Create spell repository
	spellRepo := spells.NewSpellRepository(s.db)

	var spellList []spells.SpellSummary
	var err error

	// If level is specified, get spells for that class and level
	if levelStr != "" {
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			logger.Error("Invalid spell level format",
				zap.Error(err),
				zap.String("level", levelStr))
			http.Error(w, "Invalid level", http.StatusBadRequest)
			return
		}

		spellList, err = spellRepo.ListSpellsByClassAndLevel(r.Context(), class, level)
	} else {
		// Otherwise, get all spells for the class
		spellList, err = spellRepo.ListSpellsByClass(r.Context(), class)
	}

	if err != nil {
		logger.Error("Failed to fetch spell list",
			zap.Error(err),
			zap.String("class", class))
		http.Error(w, "Failed to fetch spells", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := struct {
		IsAuthenticated bool
		Username        string
		Class           string
		Level           string
		Spells          []spells.SpellSummary
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		Class:           class,
		Level:           levelStr,
		Spells:          spellList,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	RenderTemplate(w, "templates/spells/list.html", "base.html", data)
}

// HandleSpellDetail shows details of a specific spell
func (s *Server) HandleSpellDetail(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse spell ID from query
	spellIDStr := r.URL.Query().Get("id")
	if spellIDStr == "" {
		logger.Warn("No spell ID provided")
		http.Redirect(w, r, "/spells", http.StatusSeeOther)
		return
	}

	spellID, err := strconv.ParseInt(spellIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid spell ID format",
			zap.Error(err),
			zap.String("id", spellIDStr))
		http.Error(w, "Invalid spell ID", http.StatusBadRequest)
		return
	}

	// Get spell details
	spellRepo := spells.NewSpellRepository(s.db)
	spell, err := spellRepo.GetSpellByID(r.Context(), spellID)
	if err != nil {
		logger.Error("Failed to fetch spell details",
			zap.Error(err),
			zap.Int64("spell_id", spellID))
		http.Error(w, "Failed to fetch spell details", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := struct {
		IsAuthenticated bool
		Username        string
		Spell           *spells.Spell
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		Spell:           spell,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	RenderTemplate(w, "templates/spells/detail.html", "base.html", data)
}

// HandleSpellCreate handles the creation of a new spell
func (s *Server) HandleSpellCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		// Handle form submission
		if err := r.ParseForm(); err != nil {
			logger.Error("Failed to parse form", zap.Error(err))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Create spell object from form data
		spell := &spells.Spell{
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
			Range:       r.FormValue("range"),
			Duration:    r.FormValue("duration"),
			ClassLevels: make(map[string]int),
		}

		// Process class levels
		// We'll check for common class codes from the form
		classCodes := []string{"mag", "wch", "clr", "drd", "brd", "pal", "rng"}
		for _, class := range classCodes {
			levelStr := r.FormValue("level_" + class)
			if levelStr != "" {
				level, err := strconv.Atoi(levelStr)
				if err != nil {
					logger.Error("Invalid level format",
						zap.Error(err),
						zap.String("class", class),
						zap.String("level", levelStr))
					continue
				}
				if level > 0 {
					spell.ClassLevels[class] = level
				}
			}
		}

		// Save the spell
		spellRepo := spells.NewSpellRepository(s.db)
		spellID, err := spellRepo.AddSpell(r.Context(), spell)
		if err != nil {
			logger.Error("Failed to create spell", zap.Error(err))
			http.Error(w, "Failed to create spell", http.StatusInternalServerError)
			return
		}

		// Redirect to the new spell's detail page
		http.Redirect(w, r, "/spells/detail?id="+strconv.FormatInt(spellID, 10)+"&message=Spell created successfully", http.StatusSeeOther)
		return
	}

	// Display the form
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

	RenderTemplate(w, "templates/spells/create.html", "base.html", data)
}

// HandleSpellEdit handles editing an existing spell
func (s *Server) HandleSpellEdit(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse spell ID from query
	spellIDStr := r.URL.Query().Get("id")
	if spellIDStr == "" {
		logger.Warn("No spell ID provided")
		http.Redirect(w, r, "/spells", http.StatusSeeOther)
		return
	}

	spellID, err := strconv.ParseInt(spellIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid spell ID format",
			zap.Error(err),
			zap.String("id", spellIDStr))
		http.Error(w, "Invalid spell ID", http.StatusBadRequest)
		return
	}

	spellRepo := spells.NewSpellRepository(s.db)

	if r.Method == http.MethodPost {
		// Handle form submission
		if err := r.ParseForm(); err != nil {
			logger.Error("Failed to parse form", zap.Error(err))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Create spell object from form data
		spell := &spells.Spell{
			ID:          spellID,
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
			Range:       r.FormValue("range"),
			Duration:    r.FormValue("duration"),
			ClassLevels: make(map[string]int),
		}

		// Process class levels
		classCodes := []string{"mag", "wch", "clr", "drd", "brd", "pal", "rng"}
		for _, class := range classCodes {
			levelStr := r.FormValue("level_" + class)
			if levelStr != "" {
				level, err := strconv.Atoi(levelStr)
				if err != nil {
					logger.Error("Invalid level format",
						zap.Error(err),
						zap.String("class", class),
						zap.String("level", levelStr))
					continue
				}
				if level > 0 {
					spell.ClassLevels[class] = level
				}
			}
		}

		// Update the spell
		err = spellRepo.UpdateSpell(r.Context(), spell)
		if err != nil {
			logger.Error("Failed to update spell",
				zap.Error(err),
				zap.Int64("spell_id", spellID))
			http.Error(w, "Failed to update spell", http.StatusInternalServerError)
			return
		}

		// Redirect to the spell's detail page
		http.Redirect(w, r, "/spells/detail?id="+spellIDStr+"&message=Spell updated successfully", http.StatusSeeOther)
		return
	}

	// Get existing spell data for the form
	spell, err := spellRepo.GetSpellByID(r.Context(), spellID)
	if err != nil {
		logger.Error("Failed to fetch spell details",
			zap.Error(err),
			zap.Int64("spell_id", spellID))
		http.Error(w, "Failed to fetch spell details", http.StatusInternalServerError)
		return
	}

	// Display the edit form
	data := struct {
		IsAuthenticated bool
		Username        string
		Spell           *spells.Spell
		FlashMessage    string
		CurrentYear     int
	}{
		IsAuthenticated: true,
		Username:        user.Username,
		Spell:           spell,
		FlashMessage:    r.URL.Query().Get("message"),
		CurrentYear:     time.Now().Year(),
	}

	RenderTemplate(w, "templates/spells/edit.html", "base.html", data)
}

// HandleSpellDelete handles deleting a spell
func (s *Server) HandleSpellDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Unauthorized access attempt",
			zap.String("path", r.URL.Path))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		logger.Error("Failed to parse form", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get spell ID
	spellIDStr := r.FormValue("spell_id")
	if spellIDStr == "" {
		logger.Warn("No spell ID provided for deletion")
		http.Redirect(w, r, "/spells?message=No spell specified for deletion", http.StatusSeeOther)
		return
	}

	spellID, err := strconv.ParseInt(spellIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid spell ID format",
			zap.Error(err),
			zap.String("id", spellIDStr))
		http.Error(w, "Invalid spell ID", http.StatusBadRequest)
		return
	}

	// Delete the spell
	spellRepo := spells.NewSpellRepository(s.db)
	err = spellRepo.DeleteSpell(r.Context(), spellID)
	if err != nil {
		logger.Error("Failed to delete spell",
			zap.Error(err),
			zap.Int64("spell_id", spellID))
		http.Error(w, "Failed to delete spell", http.StatusInternalServerError)
		return
	}

	// Redirect to spell list
	http.Redirect(w, r, "/spells?message=Spell deleted successfully", http.StatusSeeOther)
}
