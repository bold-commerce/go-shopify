package goshopify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGiftCardCustomerId(t *testing.T) {
	setup()
	defer teardown()
	id := uint64(207119551)
	httpmock.RegisterResponder("POST", fmt.Sprintf("https://fooshop.myshopify.com/%s/gift_cards.json", client.pathPrefix), func(req *http.Request) (*http.Response, error) {
		var body struct {
			GiftCard struct {
				CustomerId uint64 `json:"customer_id"`
			} `json:"gift_card"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.GiftCard.CustomerId != id {
			t.Errorf("customer_id = %d, want %d", body.GiftCard.CustomerId, id)
		}
		return httpmock.NewStringResponse(201, `{"gift_card":{"id":1,"customer_id":207119551}}`), nil
	})
	card, err := client.GiftCard.Create(context.Background(), GiftCard{CustomerId: &id})
	if err != nil {
		t.Fatal(err)
	}
	if card.CustomerId == nil || *card.CustomerId != id {
		t.Fatalf("customer_id = %v", card.CustomerId)
	}
	for _, data := range []string{`{}`, `{"customer_id":null}`} {
		var card GiftCard
		if err := json.Unmarshal([]byte(data), &card); err != nil {
			t.Fatal(err)
		}
		if card.CustomerId != nil {
			t.Fatalf("expected no customer for %s", data)
		}
		got, err := json.Marshal(card)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != `{}` {
			t.Fatalf("unset card = %s", got)
		}
	}
}

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

	giftCard, err := client.GiftCard.Get(context.Background(), 1)
	if err != nil {
		t.Errorf("GiftCard.Get returned error: %v", err)
	}

	expected := GiftCard{Id: 1}
	if expected.Id != giftCard.Id {
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

	giftCard, err := client.GiftCard.List(context.Background())
	if err != nil {
		t.Errorf("GiftCard.List returned error: %v", err)
	}

	expected := []GiftCard{{Id: 1}}
	if expected[0].Id != giftCard[0].Id {
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

	giftCard, err := client.GiftCard.Create(context.Background(), GiftCard{})
	if err != nil {
		t.Errorf("GiftCard.Create returned error: %v", err)
	}

	expected := GiftCard{Id: 1}
	if expected.Id != giftCard.Id {
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

	giftCard, err := client.GiftCard.Update(context.Background(), GiftCard{Id: 1})
	if err != nil {
		t.Errorf("GiftCard.Update returned error: %v", err)
	}

	expected := GiftCard{Id: 1}
	if expected.Id != giftCard.Id {
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

	giftCard, err := client.GiftCard.Disable(context.Background(), 1)
	if err != nil {
		t.Errorf("GiftCard.Disable returned error: %v", err)
	}

	expected := []GiftCard{{Id: 1}}
	if expected[0].Id != giftCard.Id {
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

	cnt, err := client.GiftCard.Count(context.Background(), nil)
	if err != nil {
		t.Errorf("GiftCard.Count returned error: %v", err)
	}

	expected := 5
	if cnt != expected {
		t.Errorf("GiftCard.Count returned %d, expected %d", cnt, expected)
	}
}
