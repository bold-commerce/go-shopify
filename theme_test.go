package goshopify

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func getTheme() Theme {
	createdAt := time.Date(2017, time.September, 23, 18, 15, 47, 0, time.UTC)
	updatedAt := time.Date(2017, time.September, 23, 18, 15, 47, 0, time.UTC)
	return Theme{
		ID:                1,
		Name:              "launchpad",
		Previewable:       true,
		Processing:        false,
		Role:              "unpublished",
		ThemeStoreID:      1234,
		CreatedAt:         &createdAt,
		UpdatedAt:         &updatedAt,
		AdminGraphQLApiID: "gid://shopify/Theme/1234",
	}
}

func TestThemeList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/themes.json", globalApiPathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"themes": [{"id":1},{"id":2}]}`,
		),
	)

	params := map[string]string{"role": "main"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/themes.json", globalApiPathPrefix),
		params,
		httpmock.NewStringResponder(
			200,
			`{"themes": [{"id":1}]}`,
		),
	)

	themes, err := client.Theme.List(nil)
	if err != nil {
		t.Errorf("Theme.List returned error: %v", err)
	}

	expected := []Theme{{ID: 1}, {ID: 2}}
	if !reflect.DeepEqual(themes, expected) {
		t.Errorf("Theme.List returned %+v, expected %+v", themes, expected)
	}

	themes, err = client.Theme.List(ThemeListOptions{Role: "main"})
	if err != nil {
		t.Errorf("Theme.List returned error: %v", err)
	}

	expected = []Theme{{ID: 1}}
	if !reflect.DeepEqual(themes, expected) {
		t.Errorf("Theme.List returned %+v, expected %+v", themes, expected)
	}
}

func TestThemeGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/%s/1.json", globalApiPathPrefix, themesBasePath),
		httpmock.NewBytesResponder(200, loadFixture("theme.json")))

	theme, err := client.Theme.Get(1, nil)
	if err != nil {
		t.Errorf("Theme.Get returned error: %v", err)
	}

	expectation := getTheme()
	if theme.ID != expectation.ID {
		t.Errorf("Theme.ID returned %+v, expected %+v", theme.ID, expectation.ID)
	}
	if theme.Name != expectation.Name {
		t.Errorf("Theme.Name returned %+v, expected %+v", theme.Name, expectation.Name)
	}
	if theme.Previewable != expectation.Previewable {
		t.Errorf("Theme.Previewable returned %+v, expected %+v", theme.Previewable, expectation.Previewable)
	}
	if theme.Processing != expectation.Processing {
		t.Errorf("Theme.Processing returned %+v, expected %+v", theme.Processing, expectation.Processing)
	}
	if theme.Role != expectation.Role {
		t.Errorf("Theme.Role returned %+v, expected %+v", theme.Role, expectation.Role)
	}
	if theme.ThemeStoreID != expectation.ThemeStoreID {
		t.Errorf("Theme.ThemeStoreID returned %+v, expected %+v", theme.ThemeStoreID, expectation.ThemeStoreID)
	}
	if !theme.CreatedAt.Equal(*expectation.CreatedAt) {
		t.Errorf("Theme.CreatedAt returned %+v, expected %+v", theme.CreatedAt, expectation.CreatedAt)
	}
	if !theme.UpdatedAt.Equal(*expectation.UpdatedAt) {
		t.Errorf("Theme.UpdatedAt returned %+v, expected %+v", theme.UpdatedAt, expectation.UpdatedAt)
	}
	if theme.AdminGraphQLApiID != expectation.AdminGraphQLApiID {
		t.Errorf("Theme.AdminGraphQLApiID returned %+v, expected %+v", theme.AdminGraphQLApiID, expectation.AdminGraphQLApiID)
	}
}

func TestThemeUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/%s/1.json", globalApiPathPrefix, themesBasePath),
		httpmock.NewBytesResponder(200, loadFixture("theme.json")))

	theme := getTheme()
	expectation, err := client.Theme.Update(theme)
	if err != nil {
		t.Errorf("Theme.Update returned error: %v", err)
	}

	expectedThemeID := int64(1)
	if expectation.ID != expectedThemeID {
		t.Errorf("Theme.ID returned %+v expected %+v", expectation.ID, expectedThemeID)
	}
}

func TestThemeCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/%s.json", globalApiPathPrefix, themesBasePath),
		httpmock.NewBytesResponder(200, loadFixture("theme.json")))

	theme := getTheme()
	expectation, err := client.Theme.Create(theme)
	if err != nil {
		t.Errorf("Theme.Create returned error: %v", err)
	}

	expectedThemeID := int64(1)
	if expectation.ID != expectedThemeID {
		t.Errorf("Theme.ID returned %+v expected %+v", expectation.ID, expectedThemeID)
	}
}

func TestThemeDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/%s/1.json", globalApiPathPrefix, themesBasePath),
		httpmock.NewStringResponder(200, ""))

	err := client.Theme.Delete(1)
	if err != nil {
		t.Errorf("Theme.Delete returned error: %v", err)
	}
}
