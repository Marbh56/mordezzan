package spells

import (
	"context"
	"database/sql"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
)

// Spell represents a complete spell with all its details
type Spell struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Range       string    `json:"range"`
	Duration    string    `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Additional calculated fields
	ClassLevels map[string]int `json:"class_levels,omitempty"`
}

// SpellSummary represents a simplified version of a spell for listings
type SpellSummary struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level,omitempty"`
}

// SpellRepository provides methods to interact with spell data
type SpellRepository struct {
	queries *db.Queries
	db      *sql.DB // Added a direct reference to the database connection
}

// NewSpellRepository creates a new spell repository
func NewSpellRepository(dbConn *sql.DB) *SpellRepository {
	return &SpellRepository{
		queries: db.New(dbConn),
		db:      dbConn, // Store the database connection
	}
}

// GetSpellByID retrieves a complete spell by its ID
func (r *SpellRepository) GetSpellByID(ctx context.Context, id int64) (*Spell, error) {
	// Get basic spell data
	spellData, err := r.queries.GetSpellByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get spell levels for each class
	spellLevels, err := r.queries.GetSpellLevels(ctx, id)
	if err != nil {
		return nil, err
	}

	classLevels := make(map[string]int)
	for _, sl := range spellLevels {
		classLevels[sl.Class] = int(sl.Level)
	}

	// Construct the complete Spell object
	spell := &Spell{
		ID:          spellData.ID,
		Name:        spellData.Name,
		Description: spellData.Description,
		Range:       spellData.Range,
		Duration:    spellData.Duration,
		CreatedAt:   spellData.CreatedAt,
		UpdatedAt:   spellData.UpdatedAt,
		ClassLevels: classLevels,
	}

	return spell, nil
}

// ListSpellsByClass retrieves all spells available to a specific class
func (r *SpellRepository) ListSpellsByClass(ctx context.Context, class string) ([]SpellSummary, error) {
	spellsData, err := r.queries.ListSpellsByClass(ctx, class)
	if err != nil {
		return nil, err
	}

	var spells []SpellSummary
	for _, s := range spellsData {
		spells = append(spells, SpellSummary{
			ID:    s.ID,
			Name:  s.Name,
			Level: int(s.Level),
		})
	}

	return spells, nil
}

// ListSpellsByClassAndLevel retrieves all spells for a specific class and level
func (r *SpellRepository) ListSpellsByClassAndLevel(ctx context.Context, class string, level int) ([]SpellSummary, error) {
	spellsData, err := r.queries.ListSpellsByClassAndLevel(ctx, db.ListSpellsByClassAndLevelParams{
		Class: class,
		Level: int64(level),
	})
	if err != nil {
		return nil, err
	}

	var spells []SpellSummary
	for _, s := range spellsData {
		spells = append(spells, SpellSummary{
			ID:    s.ID,
			Name:  s.Name,
			Level: level, // We already know the level from the query params
		})
	}

	return spells, nil
}

// AddSpell creates a new spell with its class levels
func (r *SpellRepository) AddSpell(
	ctx context.Context,
	spell *Spell,
) (int64, error) {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil) // Use the stored db connection
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	// Insert the basic spell data
	spellData, err := qtx.AddSpell(ctx, db.AddSpellParams{
		Name:        spell.Name,
		Description: spell.Description,
		Range:       spell.Range,
		Duration:    spell.Duration,
	})
	if err != nil {
		return 0, err
	}

	// Add class levels
	for class, level := range spell.ClassLevels {
		err = qtx.AddSpellLevel(ctx, db.AddSpellLevelParams{
			SpellID: spellData.ID,
			Class:   class,
			Level:   int64(level),
		})
		if err != nil {
			return 0, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return spellData.ID, nil
}

// UpdateSpell updates an existing spell
func (r *SpellRepository) UpdateSpell(
	ctx context.Context,
	spell *Spell,
) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil) // Use the stored db connection instead of trying to access r.queries.db
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	// Update the basic spell data
	_, err = qtx.UpdateSpell(ctx, db.UpdateSpellParams{
		ID:          spell.ID,
		Name:        spell.Name,
		Description: spell.Description,
		Range:       spell.Range,
		Duration:    spell.Duration,
	})
	if err != nil {
		return err
	}

	// Delete existing class levels
	err = qtx.DeleteSpellLevels(ctx, spell.ID)
	if err != nil {
		return err
	}

	// Add class levels
	for class, level := range spell.ClassLevels {
		err = qtx.AddSpellLevel(ctx, db.AddSpellLevelParams{
			SpellID: spell.ID,
			Class:   class,
			Level:   int64(level),
		})
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// DeleteSpell removes a spell and all its related data
func (r *SpellRepository) DeleteSpell(ctx context.Context, id int64) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil) // Use the stored db connection
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)

	// Delete spell levels first (due to foreign key constraints)
	err = qtx.DeleteSpellLevels(ctx, id)
	if err != nil {
		return err
	}

	// Delete the spell itself
	err = qtx.DeleteSpell(ctx, id)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// ListAllSpells retrieves all spells in the database
func (r *SpellRepository) ListAllSpells(ctx context.Context) ([]SpellSummary, error) {
	spellsData, err := r.queries.ListAllSpells(ctx)
	if err != nil {
		return nil, err
	}

	var spells []SpellSummary
	for _, s := range spellsData {
		spells = append(spells, SpellSummary{
			ID:   s.ID,
			Name: s.Name,
		})
	}

	return spells, nil
}
