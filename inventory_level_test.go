package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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
	}

	level, err := client.InventoryLevel.Set(context.Background(), options)
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

	level, err := client.InventoryLevel.Set(context.Background(), options)
	if err != nil {
		t.Errorf("InventoryLevels.Set returned error: %v", err)
	}

	inventoryLevelTests(t, level)
}

func TestInventoryLevelLocationFlags(t *testing.T) {
	setup()
	defer teardown()
	for _, enabled := range []bool{false, true} {
		for _, action := range []string{"connect", "set"} {
			t.Run(fmt.Sprintf("%s/%t", action, enabled), func(t *testing.T) {
				want := map[string]interface{}{"inventory_item_id": float64(1), "location_id": float64(2), "available": float64(0)}
				level := InventoryLevel{InventoryItemId: 1, LocationId: 2, Available: 0}
				if action == "connect" {
					level.RelocateIfNecessary = enabled
					if enabled {
						want["relocate_if_necessary"] = true
					}
				} else {
					level.DisconnectIfNecessary = enabled
					if enabled {
						want["disconnect_if_necessary"] = true
					}
				}
				httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_levels/%s.json", client.pathPrefix, action), func(req *http.Request) (*http.Response, error) {
					var got map[string]interface{}
					if err := json.NewDecoder(req.Body).Decode(&got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, want) {
						t.Errorf("body = %#v, want %#v", got, want)
					}
					return httpmock.NewBytesResponse(200, loadFixture("inventory_level.json")), nil
				})
				call := client.InventoryLevel.Connect
				if action == "set" {
					call = client.InventoryLevel.Set
				}
				got, err := call(context.Background(), level)
				if err != nil {
					t.Fatal(err)
				}
				inventoryLevelTests(t, got)
			})
		}
	}
}
