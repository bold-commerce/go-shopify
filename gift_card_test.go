package goshopify

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGiftCardGet(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("gift_card/get.json"),
		),
	)

	giftCard, err := client.GiftCard.Get(1)
	if err != nil {
		t.Errorf("GiftCard.Get returned error: %v", err)
	}

	expected := GiftCard{ID: 1}
	if expected.ID != giftCard.ID {
		t.Errorf("GiftCard.Get returned %+v, expected %+v", giftCard, expected)
	}
}

func TestGiftCardList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("gift_card/list.json"),
		),
	)

	giftCard, err := client.GiftCard.List()
	if err != nil {
		t.Errorf("GiftCard.List returned error: %v", err)
	}

	expected := []GiftCard{{ID: 1}}
	if expected[0].ID != giftCard[0].ID {
		t.Errorf("GiftCard.List returned %+v, expected %+v", giftCard, expected)
	}
}

func TestGiftCardCreate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("gift_card/get.json"),
		),
	)

	giftCard, err := client.GiftCard.Create(GiftCard{})
	if err != nil {
		t.Errorf("GiftCard.Create returned error: %v", err)
	}

	expected := GiftCard{ID: 1}
	if expected.ID != giftCard.ID {
		t.Errorf("GiftCard.Create returned %+v, expected %+v", giftCard, expected)
	}
}

func TestGiftCardUpdate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"PUT",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/1.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("gift_card/get.json"),
		),
	)

	giftCard, err := client.GiftCard.Update(GiftCard{ID: 1})
	if err != nil {
		t.Errorf("GiftCard.Update returned error: %v", err)
	}

	expected := GiftCard{ID: 1}
	if expected.ID != giftCard.ID {
		t.Errorf("GiftCard.Update returned %+v, expected %+v", giftCard, expected)
	}
}

func TestGiftCardDisable(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/1/disable.json", client.pathPrefix),
		httpmock.NewBytesResponder(
			200,
			loadFixture("gift_card/get.json"),
		),
	)

	giftCard, err := client.GiftCard.Disable(1)
	if err != nil {
		t.Errorf("GiftCard.Disable returned error: %v", err)
	}

	expected := []GiftCard{{ID: 1}}
	if expected[0].ID != giftCard.ID {
		t.Errorf("GiftCard.Disable returned %+v, expected %+v", giftCard, expected)
	}
}

func TestGiftCardCount(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards/count.json", client.pathPrefix),
		httpmock.NewStringResponder(
			200,
			`{"count": 5}`,
		),
	)

	cnt, err := client.GiftCard.Count(nil)
	if err != nil {
		t.Errorf("GiftCard.Count returned error: %v", err)
	}

	expected := 5
	if cnt != expected {
		t.Errorf("GiftCard.Count returned %d, expected %d", cnt, expected)
	}

}
