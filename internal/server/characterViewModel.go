package server

import (
	"database/sql"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/rules"
	"github.com/marbh56/mordezzan/internal/rules/ability_scores"
	charRules "github.com/marbh56/mordezzan/internal/rules/character"
	"github.com/marbh56/mordezzan/internal/rules/combat"
)

func interfaceToNullString(v interface{}) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	if s, ok := v.(string); ok {
		return sql.NullString{
			String: s,
			Valid:  true,
		}
	}
	return sql.NullString{}
}

func interfaceToNullInt64(v interface{}) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	switch i := v.(type) {
	case int64:
		return sql.NullInt64{Int64: i, Valid: true}
	case int:
		return sql.NullInt64{Int64: int64(i), Valid: true}
	case sql.NullInt64:
		return i
	}
	return sql.NullInt64{}
}

// Represents a single item in a character's inventory
type InventoryItem struct {
	ID               int64           `json:"id"`
	CharacterID      int64           `json:"character_id"`
	ItemType         string          `json:"item_type"`
	ItemID           int64           `json:"item_id"`
	ItemName         string          `json:"item_name"`
	ItemWeight       int             `json:"item_weight"`
	Quantity         int64           `json:"quantity"`
	ContainerID      sql.NullInt64   `json:"container_id"`
	EquipmentSlotID  sql.NullInt64   `json:"equipment_slot_id"`
	SlotName         sql.NullString  `json:"slot_name"`
	CustomName       sql.NullString  `json:"custom_name"`
	CustomNotes      sql.NullString  `json:"custom_notes"`
	IsIdentified     bool            `json:"is_identified"`
	Charges          sql.NullInt64   `json:"charges"`
	Condition        string          `json:"condition"`
	Damage           sql.NullString  `json:"damage"`
	AttacksPerRound  sql.NullString  `json:"attacks_per_round"`
	MovementRate     sql.NullInt64   `json:"movement_rate"`
	DefenseBonus     interface{}     `json:"defense_bonus"`
	EnhancementBonus sql.NullInt64   `json:"enhancement_bonus,omitempty"`
	Notes            sql.NullString  `json:"notes"`
	ContainerOptions []InventoryItem `json:"container_options,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// Contains inventory statistics and calculated values
type InventoryStats struct {
	TotalWeight         int    `json:"total_weight"`
	EquippedWeight      int    `json:"equipped_weight"`
	CarriedWeight       int    `json:"carried_weight"`
	ContainersWeight    int    `json:"containers_weight"`
	CoinWeight          int    `json:"coin_weight"`
	BaseEncumbered      int    `json:"base_encumbered"`
	BaseHeavyEncumbered int    `json:"base_heavy_encumbered"`
	MaximumCapacity     int    `json:"maximum_capacity"`
	EncumbranceLevel    string `json:"encumbrance_level"`
}

func classGetsFighterBonus(class string) bool {
	fighterClasses := map[string]bool{
		"Fighter":   true,
		"Barbarian": true,
		"Beserker":  true,
		"Huntsman":  true,
		"Paladin":   true,
		"Ranger":    true,
		"Warlock":   true,
	}
	return fighterClasses[class]
}

// Complete character view model including inventory
type CharacterViewModel struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	Name       string `json:"name"`
	Class      string `json:"class"`
	Level      int64  `json:"level"`
	MaxHp      int64  `json:"max_hp"`
	CurrentHp  int64  `json:"current_hp"`
	ArmorClass int    `json:"armor_class"`

	// Ability scores and modifiers
	Strength          int64                            `json:"strength"`
	StrengthModifiers ability_scores.StrengthModifiers `json:"strength_modifiers"`

	Dexterity          int64                             `json:"dexterity"`
	DexterityModifiers ability_scores.DexterityModifiers `json:"dexterity_modifiers"`

	Constitution          int64                                `json:"constitution"`
	ConstitutionModifiers ability_scores.ConstitutionModifiers `json:"constitution_modifiers"`

	Intelligence          int64                                `json:"intelligence"`
	IntelligenceModifiers ability_scores.IntelligenceModifiers `json:"intelligence_modifiers"`

	Wisdom          int64                          `json:"wisdom"`
	WisdomModifiers ability_scores.WisdomModifiers `json:"wisdom_modifiers"`

	Charisma          int64                            `json:"charisma"`
	CharismaModifiers ability_scores.CharismaModifiers `json:"charisma_modifiers"`

	// Combat information
	CombatMatrix []int64 `json:"combat_matrix"`
	SavingThrow  int64   `json:"saving_throw"`

	// Inventory organization
	EquippedItems  []InventoryItem           `json:"equipped_items"`
	CarriedItems   []InventoryItem           `json:"carried_items"`
	ContainerItems map[int64][]InventoryItem `json:"container_items"`

	// Calculated inventory statistics
	InventoryStats InventoryStats `json:"inventory_stats"`
	PlatinumPieces int64
	GoldPieces     int64
	ElectrumPieces int64
	SilverPieces   int64
	CopperPieces   int64

	// Experience points and level progression
	ExperiencePoints int64 `json:"experience_points"`
	NextLevelXP      int64 `json:"next_level_xp"`
	XPNeeded         int64 `json:"xp_needed"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Helper function to safely convert interface{} to string
func interfaceToString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// Helper function to safely convert interface{} to int
func interfaceToInt(v interface{}) int {
	if v == nil {
		return 0
	}
	if i, ok := v.(int64); ok {
		return int(i)
	}
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}

// Creates a new character view model with inventory data
func NewCharacterViewModel(c db.Character, inventory []db.GetCharacterInventoryItemsRow) CharacterViewModel {
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

	// Add coin weight to total weight
	coinageWeight := calculateCoinWeight(c.PlatinumPieces, c.GoldPieces, c.ElectrumPieces, c.SilverPieces, c.CopperPieces)
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

	// Determine encumbrance level
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

// Helper function to safely convert interface{} to int64
func interfaceToInt64(v interface{}) (int64, bool) {
	if v == nil {
		return 0, false
	}

	switch value := v.(type) {
	case int64:
		return value, true
	case int:
		return int64(value), true
	case float64:
		return int64(value), true
	case sql.NullInt64:
		return value.Int64, value.Valid
	}

	return 0, false
}
