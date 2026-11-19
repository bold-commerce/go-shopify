package goshopify

import (
	"context"
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func inventoryLevelTests(t *testing.T, item *InventoryLevel) {
	if item == nil {
		t.Errorf("InventoryItem is nil")
		return
	}

	expectedInt := uint64(808950810)
	if item.InventoryItemId != expectedInt {
		t.Errorf("InventoryLevel.InventoryItemId returned %+v, expected %+v",
			item.InventoryItemId, expectedInt)
	}

	expectedInt = uint64(905684977)
	if item.LocationId != expectedInt {
		t.Errorf("InventoryLevel.LocationId is %+v, expected %+v",
			item.LocationId, expectedInt)
	}

	expectedAvailable := 6
	if item.Available == nil || *item.Available != expectedAvailable {
		t.Errorf("InventoryLevel.Available is %+v, expected %+v",
			item.Available, expectedAvailable)
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

	levels, err := client.InventoryLevel.List(context.Background(), nil)
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
		InventoryItemIds: []uint64{1, 2},
	}

	levels, err := client.InventoryLevel.List(context.Background(), options)
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
		LocationIds: []uint64{1, 2},
	}

	levels, err := client.InventoryLevel.List(context.Background(), options)
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

	adjItem, err := client.InventoryLevel.Adjust(context.Background(), option)
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

	err := client.InventoryLevel.Delete(context.Background(), 1, 1)
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

	level, err := client.InventoryLevel.Connect(context.Background(), options)
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
		Available:       IntPtr(3),
	}

	level, err := client.InventoryLevel.Set(context.Background(), options)
	if err != nil {
		t.Errorf("InventoryLevels.Set returned error: %v", err)
	}

	inventoryLevelTests(t, level)
}

func TestInventoryLevelLocationFlags(t *testing.T) {
	for _, tc := range []struct {
		name, action, body string
		enabled            bool
		call               func(InventoryLevelService, context.Context, InventoryLevel) (*InventoryLevel, error)
	}{
		{"connect", "connect", `{"inventory_item_id":1,"location_id":2}`, false, InventoryLevelService.Connect},
		{"relocate", "connect", `{"inventory_item_id":1,"location_id":2,"relocate_if_necessary":true}`, true, InventoryLevelService.Connect},
		{"set zero", "set", `{"inventory_item_id":1,"location_id":2,"available":0}`, false, InventoryLevelService.Set},
		{"disconnect", "set", `{"inventory_item_id":1,"location_id":2,"available":0,"disconnect_if_necessary":true}`, true, InventoryLevelService.Set},
	} {
		t.Run(tc.name, func(t *testing.T) {
			level := InventoryLevel{InventoryItemId: 1, LocationId: 2, Available: IntPtr(0),
				RelocateIfNecessary: tc.enabled, DisconnectIfNecessary: tc.enabled, AdminGraphqlApiId: "read-only"}
			testFeatureEndpoint(t, featureEndpoint{
				method: "POST", path: "inventory_levels/" + tc.action + ".json", request: tc.body,
				response: `{"inventory_level":{"location_id":2}}`, want: &InventoryLevel{LocationId: 2},
				call: func(ctx context.Context) (interface{}, error) { return tc.call(client.InventoryLevel, ctx, level) },
			})
		})
	}
}

func TestInventoryLevelSetRequiresQuantity(t *testing.T) {
	if _, err := (&InventoryLevelServiceOp{}).Set(context.Background(), InventoryLevel{}); err == nil {
		t.Fatal("set accepted an unknown quantity")
	}
}
