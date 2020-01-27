package goshopify

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func inventoryItemTests(t *testing.T, item *InventoryItem) {
	if item == nil {
		t.Errorf("InventoryItem is nil")
		return
	}

	expectedInt := int64(808950810)
	if item.ID != expectedInt {
		t.Errorf("InventoryItem.ID returned %+v, expected %+v", item.ID, expectedInt)
	}

	expectedSKU := "new sku"
	if item.SKU != expectedSKU {
		t.Errorf("InventoryItem.SKU sku is %+v, expected %+v", item.SKU, expectedSKU)
	}

	if item.Cost == nil {
		t.Errorf("InventoryItem.Cost is nil")
		return
	}

	expectedCost := 25.00
	costFloat, _ := item.Cost.Float64()
	if costFloat != expectedCost {
		t.Errorf("InventoryItem.Cost (float) is %+v, expected %+v", costFloat, expectedCost)
	}
}

func inventoryItemsTests(t *testing.T, items []InventoryItem) {
	expectedLen := 3
	if len(items) != expectedLen {
		t.Errorf("InventoryItems list lenth is %+v, expected %+v", len(items), expectedLen)
	}
}

func TestInventoryItemsList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_items.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("inventory_items.json")))

	items, err := client.InventoryItem.List(nil)
	if err != nil {
		t.Errorf("InventoryItems.List returned error: %v", err)
	}

	inventoryItemsTests(t, items)
}

func TestInventoryItemsListWithIDs(t *testing.T) {
	setup()
	defer teardown()

	params := map[string]string{
		"ids": "1,2",
	}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_items.json", globalApiPathPrefix),
		params,
		httpmock.NewBytesResponder(200, loadFixture("inventory_items.json")),
	)

	options := ListOptions{
		IDs: []int64{1, 2},
	}

	items, err := client.InventoryItem.List(options)
	if err != nil {
		t.Errorf("InventoryItems.List returned error: %v", err)
	}

	inventoryItemsTests(t, items)
}

func TestInventoryItemGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_items/1.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("inventory_item.json")))

	item, err := client.InventoryItem.Get(1, nil)
	if err != nil {
		t.Errorf("InventoryItem.Get returned error: %v", err)
	}

	inventoryItemTests(t, item)
}
func TestInventoryItemUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT", fmt.Sprintf("https://fooshop.myshopify.com/%s/inventory_items/1.json", globalApiPathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("inventory_item.json")))

	item := InventoryItem{
		ID: 1,
	}

	updatedItem, err := client.InventoryItem.Update(item)
	if err != nil {
		t.Errorf("InentoryItem.Update returned error: %v", err)
	}

	inventoryItemTests(t, updatedItem)
}
