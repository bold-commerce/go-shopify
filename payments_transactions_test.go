package goshopify

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"reflect"
	"testing"
	"time"
)

func TestPaymentsTransactionsList(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", fmt.Sprintf("https://fooshop.myshopify.com/%s/shopify_payments/balance/transactions.json", client.pathPrefix),
		httpmock.NewBytesResponder(200, loadFixture("payments_transactions.json")))
	date1 := OnlyDate{time.Date(2013, 11, 01, 0, 0, 0, 0, time.UTC)}
	paymentsTransactions, err := client.PaymentsTransactions.List(PaymentsTransactionsListOptions{PayoutId: 623721858})
	if err != nil {
		t.Errorf("PaymentsTransactions.List returned error: %v", err)
	}

	expected := []PaymentsTransactions{
		{
			ID:                       699519475,
			Type:                     PaymentsTransactionsDebit,
			Test:                     false,
			PayoutID:                 623721858,
			PayoutStatus:             PayoutStatusPaid,
			Currency:                 "USD",
			Amount:                   "-50.00",
			Fee:                      "0.00",
			Net:                      "-50.00",
			SourceID:                 460709370,
			SourceType:               "adjustment",
			SourceOrderID:            0,
			SourceOrderTransactionID: 0,
			ProcessedAt:              date1,
		},
		{
			ID:                       77412310,
			Type:                     PaymentsTransactionsCredit,
			Test:                     false,
			PayoutID:                 623721858,
			PayoutStatus:             PayoutStatusPaid,
			Currency:                 "USD",
			Amount:                   "50.00",
			Fee:                      "0.00",
			Net:                      "50.00",
			SourceID:                 374511569,
			SourceType:               "Payments::Balance::AdjustmentReversal",
			SourceOrderID:            0,
			SourceOrderTransactionID: 0,
			ProcessedAt:              date1,
		},
		{
			ID:                       1006917261,
			Type:                     PaymentsTransactionsRefund,
			Test:                     false,
			PayoutID:                 623721858,
			PayoutStatus:             PayoutStatusPaid,
			Currency:                 "USD",
			Amount:                   "-3.45",
			Fee:                      "0.00",
			Net:                      "-3.45",
			SourceID:                 1006917261,
			SourceType:               "Payments::Refund",
			SourceOrderID:            217130470,
			SourceOrderTransactionID: 1006917261,
			ProcessedAt:              date1,
		},
	}
	if !reflect.DeepEqual(paymentsTransactions, expected) {
		t.Errorf("PaymentsTransactions.List returned %+v, expected %+v", paymentsTransactions, expected)
	}
}
