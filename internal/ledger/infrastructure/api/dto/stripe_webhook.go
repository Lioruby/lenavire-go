package dto

type StripeWebhookRequest struct {
	Data struct {
		Object struct {
			ID              string `json:"id"`
			ObjectType      string `json:"object"`
			AdaptivePricing struct {
				Enabled bool `json:"enabled"`
			} `json:"adaptive_pricing"`
			AmountSubtotal int `json:"amount_subtotal"`
			AmountTotal    int `json:"amount_total"`
			AutomaticTax   struct {
				Enabled   bool   `json:"enabled"`
				Liability string `json:"liability"`
				Status    string `json:"status"`
			} `json:"automatic_tax"`
			CustomFields []struct {
				Key   string `json:"key"`
				Label struct {
					Custom string `json:"custom"`
					Type   string `json:"type"`
				} `json:"label"`
				Text struct {
					Value string `json:"value"`
				} `json:"text"`
				Dropdown struct {
					Options []struct {
						Label string `json:"label"`
						Value string `json:"value"`
					} `json:"options"`
					Value string `json:"value"`
				} `json:"dropdown"`
			} `json:"custom_fields"`
			CustomerDetails struct {
				Address struct {
					City       *string `json:"city"`
					Country    string  `json:"country"`
					Line1      *string `json:"line1"`
					Line2      *string `json:"line2"`
					PostalCode *string `json:"postal_code"`
					State      *string `json:"state"`
				} `json:"address"`
				Email     string   `json:"email"`
				Name      string   `json:"name"`
				Phone     *string  `json:"phone"`
				TaxExempt string   `json:"tax_exempt"`
				TaxIDs    []string `json:"tax_ids"`
			} `json:"customer_details"`
		} `json:"object"`
	} `json:"data"`
}
