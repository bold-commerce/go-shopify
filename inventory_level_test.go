package goshopify

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func inventoryLevelTests(t *testing.T, item *InventoryLevel) {
	if item == nil {
		t.Errorf("InventoryItem is nil")
		return
	}

	expectedInt := int64(808950810)
	if item.InventoryItemId != expectedInt {
		t.Errorf("InventoryLevel.InventoryItemId returned %+v, expected %+v",
			item.InventoryItemId, expectedInt)
	}

	expectedInt = int64(905684977)
	if item.LocationId != expectedInt {
		t.Errorf("InventoryLevel.LocationId is %+v, expected %+v",
			item.LocationId, expectedInt)
	}

	expectedAvailable := 6
	if item.Available != expectedAvailable {
		t.Errorf("InventoryLevel.Available is %+v, expected %+v",
			item.Available, expectedInt)
	}
}

func inventoryLevelsTests(t *testing.T, levels []InventoryLevel) {
	expectedLen := 4
	if len(levels) != expectedLen {
		t.Errorf("InventoryLevels list lenth is %+v, expected %+v", len(levels), expectedLen)
	}
}

func TestInventoryLevelsList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("inventory_levels.json")))

	levels, err := client.InventoryLevel.List(nil)
	if err != nil {
		t.Errorf("InventoryLevels.List returned error: %v", err)
	}

	inventoryLevelsTests(t, levels)
}

func TestInventoryLevelListWithItemId(t *testing.T) {
	setup()
	defer teardown()

	params := map[string]string{
		"inventory_item_ids": "1,2",
	}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels.json", client.pathPrefix),
		params,
		httpmock.NewBytesResponder(200, loadFixture("inventory_levels.json")),
	)

	options := InventoryLevelListOptions{
		InventoryItemIds: []int64{1, 2},
	}

	levels, err := client.InventoryLevel.List(options)
	if err != nil {
		t.Errorf("InventoryLevels.List returned error: %v", err)
	}

	inventoryLevelsTests(t, levels)
}

func TestInventoryLevelListWithLocationId(t *testing.T) {
	setup()
	defer teardown()

	params := map[string]string{
		"location_ids": "1,2",
	}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels.json", client.pathPrefix),
		params,
		httpmock.NewBytesResponder(200, loadFixture("inventory_levels.json")),
	)

	options := InventoryLevelListOptions{
		LocationIds: []int64{1, 2},
	}

	levels, err := client.InventoryLevel.List(options)
	if err != nil {
		t.Errorf("InventoryLevels.List returned error: %v", err)
	}

	inventoryLevelsTests(t, levels)
}

func TestInventoryLevelAdjust(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels/adjust.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("inventory_level.json")))

	option := InventoryLevelAdjustOptions{
		InventoryItemId: 808950810,
		LocationId:      905684977,
		Adjust:          6,
	}

	adjItem, err := client.InventoryLevel.Adjust(option)
	if err != nil {
		t.Errorf("InventoryLevel.Adjust returned error: %v", err)
	}

	inventoryLevelTests(t, adjItem)
}

func TestInventoryLevelDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels.json", client.pathPrefix),
		httpmock.NewStringResponder(200, "{}"))

	err := client.InventoryLevel.Delete(1, 1)
	if err != nil {
		t.Errorf("InventoryLevel.Delete returned error: %v", err)
	}
}

func TestInventoryLevelConnect(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels/connect.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("inventory_level.json")),
	)

	options := InventoryLevel{
		InventoryItemId: 1,
		LocationId:      1,
	}

	level, err := client.InventoryLevel.Connect(options)
	if err != nil {
		t.Errorf("InventoryLevels.Connect returned error: %v", err)
	}

	inventoryLevelTests(t, level)
}

func TestInventoryLevelSet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels/set.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("inventory_level.json")),
	)

	options := InventoryLevel{
		InventoryItemId: 1,
		LocationId:      1,
	}

	level, err := client.InventoryLevel.Set(options)
	if err != nil {
		t.Errorf("InventoryLevels.Set returned error: %v", err)
	}

	inventoryLevelTests(t, level)
}

func TestInventoryLevelSetZero(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels/set.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("inventory_level.json")),
	)

	options := InventoryLevel{
		InventoryItemId: 1,
		LocationId:      1,
		Available:       0,
	}

	level, err := client.InventoryLevel.Set(options)
	if err != nil {
		t.Errorf("InventoryLevels.Set returned error: %v", err)
	}

	inventoryLevelTests(t, level)
}
