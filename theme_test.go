package goshopify

import (
	"context"
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
		Id:                1,
		Name:              "launchpad",
		Previewable:       true,
		Processing:        false,
		Role:              "unpublished",
		ThemeStoreId:      1234,
		CreatedAt:         &createdAt,
		UpdatedAt:         &updatedAt,
		AdminGraphqlApiId: "gid://shopify/Theme/1234",
	}
}

func TestThemeList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/themes.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"themes": [{"id":1},{"id":2}]}`,
		),
	)

	params := map[string]string{"role": "main"}
	httpmock.RegisterResponderWithQuery(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/themes.json", client.pathPrefix),
		params,
		httpmock.NewStringResponder(
			200,
			`{"themes": [{"id":1}]}`,
		),
	)

	themes, err := client.Theme.List(context.Background(), nil)
	if err != nil {
		t.Errorf("Theme.List returned error: %v", err)
	}

	expected := []Theme{{Id: 1}, {Id: 2}}
	if !reflect.DeepEqual(themes, expected) {
		t.Errorf("Theme.List returned %+v, expected %+v", themes, expected)
	}

	themes, err = client.Theme.List(context.Background(), ThemeListOptions{Role: "main"})
	if err != nil {
		t.Errorf("Theme.List returned error: %v", err)
	}

	expected = []Theme{{Id: 1}}
	if !reflect.DeepEqual(themes, expected) {
		t.Errorf("Theme.List returned %+v, expected %+v", themes, expected)
	}
}

func TestThemeGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/%s/1.json", client.pathPrefix, themesBasePath),
		httpmock.NewBytesResponder(200, loadFixture("theme.json")))

	theme, err := client.Theme.Get(context.Background(), 1, nil)
	if err != nil {
		t.Errorf("Theme.Get returned error: %v", err)
	}

	expectation := getTheme()
	if theme.Id != expectation.Id {
		t.Errorf("Theme.Id returned %+v, expected %+v", theme.Id, expectation.Id)
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
	if theme.ThemeStoreId != expectation.ThemeStoreId {
		t.Errorf("Theme.ThemeStoreId returned %+v, expected %+v", theme.ThemeStoreId, expectation.ThemeStoreId)
	}
	if !theme.CreatedAt.Equal(*expectation.CreatedAt) {
		t.Errorf("Theme.CreatedAt returned %+v, expected %+v", theme.CreatedAt, expectation.CreatedAt)
	}
	if !theme.UpdatedAt.Equal(*expectation.UpdatedAt) {
		t.Errorf("Theme.UpdatedAt returned %+v, expected %+v", theme.UpdatedAt, expectation.UpdatedAt)
	}
	if theme.AdminGraphqlApiId != expectation.AdminGraphqlApiId {
		t.Errorf("Theme.AdminGraphqlApiId returned %+v, expected %+v", theme.AdminGraphqlApiId, expectation.AdminGraphqlApiId)
	}
}

func TestThemeUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/%s/1.json", client.pathPrefix, themesBasePath),
		httpmock.NewBytesResponder(200, loadFixture("theme.json")))

	theme := getTheme()
	expectation, err := client.Theme.Update(context.Background(), theme)
	if err != nil {
		t.Errorf("Theme.Update returned error: %v", err)
	}

	expectedThemeId := uint64(1)
	if expectation.Id != expectedThemeId {
		t.Errorf("Theme.Id returned %+v expected %+v", expectation.Id, expectedThemeId)
	}
}

func TestThemeCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/%s.json", client.pathPrefix, themesBasePath),
		httpmock.NewBytesResponder(200, loadFixture("theme.json")))

	theme := getTheme()
	expectation, err := client.Theme.Create(context.Background(), theme)
	if err != nil {
		t.Errorf("Theme.Create returned error: %v", err)
	}

	expectedThemeId := uint64(1)
	if expectation.Id != expectedThemeId {
		t.Errorf("Theme.Id returned %+v expected %+v", expectation.Id, expectedThemeId)
	}
}

func TestThemeDelete(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/%s/1.json", client.pathPrefix, themesBasePath),
		httpmock.NewStringResponder(200, ""))

	err := client.Theme.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("Theme.Delete returned error: %v", err)
	}
}
