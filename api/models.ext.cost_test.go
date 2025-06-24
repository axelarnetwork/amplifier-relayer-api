package api_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axelarnetwork/amplifier-relayer-api/api"
	"github.com/axelarnetwork/amplifier-relayer-api/internal/funcs"
)

func TestCostFromToken(t *testing.T) {
	token := api.Token{
		Amount: "123",
	}

	cost := api.CostFromToken(token)

	t.Run("should not unmarshal as fees", func(t *testing.T) {
		fees, err := cost.AsFees()
		assert.Error(t, err)
		assert.Nil(t, fees)
	})

	t.Run("should unmarshal as token", func(t *testing.T) {
		result, err := cost.AsToken()
		assert.NoError(t, err)
		assert.Equal(t, token, result)
	})
}

func TestCost_Validate(t *testing.T) {
	t.Run("when valid token", func(t *testing.T) {
		costJSON := `
{
	"token": {
		"amount": "123"
	}
}
`
		var cost api.Cost
		funcs.MustNoErr(
			json.Unmarshal([]byte(costJSON), &cost),
		)

		result := cost.Validate()

		require.NoError(t, result)
	})

	t.Run("when valid fees", func(t *testing.T) {
		costJSON := `
[
	{
		"id": "fee1",
		"token": {
			"amount": "123"
		}
	},
	{
		"id": "fee2",
		"token": {
			"amount": "456"
		}
	}
]
`
		var cost api.Cost
		funcs.MustNoErr(
			json.Unmarshal([]byte(costJSON), &cost),
		)

		result := cost.Validate()

		require.NoError(t, result)
	})

	t.Run("when empty fees", func(t *testing.T) {
		costJSON := `[]`

		var cost api.Cost
		funcs.MustNoErr(
			json.Unmarshal([]byte(costJSON), &cost),
		)

		result := cost.Validate()

		require.NoError(t, result)
	})

	t.Run("when invalid fees with duplicate IDs", func(t *testing.T) {
		costJSON := `
[
	{
		"id": "same_fee_id",
		"token": {
			"amount": "123"
		}
	},
	{
		"id": "same_fee_id",
		"token": {
			"amount": "456"
		}
	}
]
`
		var cost api.Cost
		funcs.MustNoErr(
			json.Unmarshal([]byte(costJSON), &cost),
		)

		result := cost.Validate()

		require.Error(t, result)
	})

	t.Run("when neither token, nor fees", func(t *testing.T) {
		costJSON := `"value"`

		var cost api.Cost
		funcs.MustNoErr(
			json.Unmarshal([]byte(costJSON), &cost),
		)

		result := cost.Validate()

		require.Error(t, result)
	})
}
