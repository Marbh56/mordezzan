package server

import (
	"database/sql"
	"time"

	"github.com/marbh56/mordezzan/internal/db"
	"github.com/marbh56/mordezzan/internal/rules"
)

// Represents a single item in a character's inventory
type InventoryItem struct {
    ID                   int64          `json:"id"`
    CharacterID          int64          `json:"character_id"`
    ItemType            string          `json:"item_type"`
    ItemID              int64          `json:"item_id"`
    ItemName            string          `json:"item_name"`
    ItemWeight          int            `json:"item_weight"`
    Quantity            int64          `json:"quantity"`
    ContainerInventoryID sql.NullInt64  `json:"container_inventory_id"`
    EquipmentSlotID     sql.NullInt64  `json:"equipment_slot_id"`
    SlotName            sql.NullString  `json:"slot_name"`
    Notes               sql.NullString  `json:"notes"`
    CreatedAt           time.Time       `json:"created_at"`
    UpdatedAt           time.Time       `json:"updated_at"`
}

// Contains inventory statistics and calculated values
type InventoryStats struct {
    TotalWeight           int     `json:"total_weight"`
    EquippedWeight        int     `json:"equipped_weight"`
    CarriedWeight         int     `json:"carried_weight"`
    ContainersWeight      int     `json:"containers_weight"`
    BaseEncumbered        int     `json:"base_encumbered"`
    BaseHeavyEncumbered   int     `json:"base_heavy_encumbered"`
    MaximumCapacity       int     `json:"maximum_capacity"`
    EncumbranceLevel      string  `json:"encumbrance_level"` // "None", "Encumbered", "Heavy", "Over"
}

// Complete character view model including inventory
type CharacterViewModel struct {
    ID        int64  `json:"id"`
    UserID    int64  `json:"user_id"`
    Name      string `json:"name"`
    MaxHp     int64  `json:"max_hp"`
    CurrentHp int64  `json:"current_hp"`

    // Ability scores with modifiers
    Strength          int64                   `json:"strength"`
    StrengthModifiers rules.StrengthModifiers `json:"strength_modifiers"`

    Dexterity          int64                    `json:"dexterity"`
    DexterityModifiers rules.DexterityModifiers `json:"dexterity_modifiers"`

    Constitution          int64                       `json:"constitution"`
    ConstitutionModifiers rules.ConstitutionModifiers `json:"constitution_modifiers"`

    Intelligence int64 `json:"intelligence"`
    Wisdom       int64 `json:"wisdom"`
    Charisma     int64 `json:"charisma"`

    // Inventory organization
    EquippedItems  []InventoryItem `json:"equipped_items"`
    CarriedItems   []InventoryItem `json:"carried_items"`
    ContainerItems map[int64][]InventoryItem `json:"container_items"` // Key is container inventory ID

    // Calculated inventory statistics
    InventoryStats InventoryStats `json:"inventory_stats"`

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
func NewCharacterViewModel(c db.Character, inventory []db.GetCharacterInventoryRow) CharacterViewModel {
    vm := CharacterViewModel{
        ID:          c.ID,
        UserID:      c.UserID,
        Name:        c.Name,
        MaxHp:       c.MaxHp,
        CurrentHp:   c.CurrentHp,
        Strength:    c.Strength,
        Dexterity:   c.Dexterity,
        Constitution: c.Constitution,
        Intelligence: c.Intelligence,
        Wisdom:      c.Wisdom,
        Charisma:    c.Charisma,
        CreatedAt:   c.CreatedAt,
        UpdatedAt:   c.UpdatedAt,

        // Initialize modifiers
        StrengthModifiers:     rules.CalculateStrengthModifiers(c.Strength),
        DexterityModifiers:    rules.CalculateDexterityModifiers(c.Dexterity),
        ConstitutionModifiers: rules.CalculateConstitutionModifiers(c.Constitution),

        // Initialize inventory containers
        ContainerItems: make(map[int64][]InventoryItem),
    }

    // Initialize inventory stats with encumbrance thresholds
    encumbranceThresholds := rules.CalculateEncumbranceThresholds(c.Strength, c.Constitution)
    vm.InventoryStats = InventoryStats{
        BaseEncumbered:      encumbranceThresholds.BaseEncumbered,
        BaseHeavyEncumbered: encumbranceThresholds.BaseHeavyEncumbered,
        MaximumCapacity:     encumbranceThresholds.MaximumCapacity,
    }

    // Sort inventory items and calculate weights
    for _, item := range inventory {
        invItem := InventoryItem{
            ID:                  item.ID,
            CharacterID:         item.CharacterID,
            ItemType:           item.ItemType,
            ItemID:            item.ItemID,
            ItemName:          interfaceToString(item.ItemName),
            ItemWeight:        interfaceToInt(item.ItemWeight),
            Quantity:          item.Quantity,
            ContainerInventoryID: item.ContainerInventoryID,
            EquipmentSlotID:    item.EquipmentSlotID,
            SlotName:          item.SlotName,
            Notes:             item.Notes,
            CreatedAt:         item.CreatedAt,
            UpdatedAt:         item.UpdatedAt,
        }

        // Calculate total weight for this item
        itemTotalWeight := invItem.ItemWeight * int(invItem.Quantity)

        if invItem.EquipmentSlotID.Valid {
            vm.EquippedItems = append(vm.EquippedItems, invItem)
            vm.InventoryStats.EquippedWeight += itemTotalWeight
        } else if invItem.ContainerInventoryID.Valid {
            containerID := invItem.ContainerInventoryID.Int64
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
                                   vm.InventoryStats.ContainersWeight

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

    return vm
}
