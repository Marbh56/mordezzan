package server

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marbh56/mordezzan/internal/db"
)

type containerDetails struct {
	CapacityWeight int64
	CapacityItems  sql.NullInt64
}

func (s *Server) getContainerDetails(ctx context.Context, containerID int64) (*containerDetails, error) {
	var container containerDetails
	err := s.db.QueryRowContext(ctx, `
        SELECT capacity_weight, capacity_items
        FROM containers
        WHERE id = ?`, containerID).Scan(&container.CapacityWeight, &container.CapacityItems)
	if err != nil {
		return nil, fmt.Errorf("error getting container details: %v", err)
	}
	return &container, nil
}

func (s *Server) validateContainerCapacity(ctx context.Context, containerInvID int64, characterID int64, itemWeight int64, quantity int64) error {
	queries := db.New(s.db)

	// Get the container's inventory item
	containerItem, err := queries.GetItemFromInventory(ctx, db.GetItemFromInventoryParams{
		ID:          containerInvID,
		CharacterID: characterID,
	})
	if err != nil {
		return fmt.Errorf("error getting container: %v", err)
	}

	// Get container details
	container, err := s.getContainerDetails(ctx, containerItem.ItemID)
	if err != nil {
		return err
	}

	// Get current contents
	contents, err := queries.GetContainerContents(ctx, db.GetContainerContentsParams{
		ContainerInventoryID: sql.NullInt64{Int64: containerInvID, Valid: true},
		CharacterID:          characterID,
	})
	if err != nil {
		return fmt.Errorf("error getting container contents: %v", err)
	}

	// Calculate current weight
	var currentWeight int64
	for _, item := range contents {
		currentWeight += int64(item.ItemWeight) * item.Quantity
	}

	// Check weight capacity
	newTotalWeight := currentWeight + (itemWeight * quantity)
	if newTotalWeight > container.CapacityWeight {
		return fmt.Errorf("container cannot hold more than %d lbs (trying to add %d lbs, current %d lbs)",
			container.CapacityWeight, itemWeight*quantity, currentWeight)
	}

	// Check item count if limited
	if container.CapacityItems.Valid {
		if int64(len(contents)+1) > container.CapacityItems.Int64 {
			return fmt.Errorf("container cannot hold more than %d items", container.CapacityItems.Int64)
		}
	}

	return nil
}
