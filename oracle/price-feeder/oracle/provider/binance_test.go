package provider

import (
	"context"
	"testing"

	"github.com/bluelink-lab/blk-chain/oracle/price-feeder/config"
	"github.com/bluelink-lab/blk-chain/oracle/price-feeder/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestBinanceProvider_GetTickerPrices(t *testing.T) {
	server := NewMockProviderServer()
	server.Start()
	defer server.Close()

	p, err := NewBinanceProvider(
		context.TODO(),
		zerolog.Nop(),
		config.ProviderEndpoint{
			Name:      config.ProviderBinance,
			Rest:      "",
			Websocket: server.GetBaseURL(),
		},
		types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
	)
	require.NoError(t, err)

	t.Run("valid_request_single_ticker", func(t *testing.T) {
		lastPrice := "34.69000000"
		volume := "2396974.02000000"

		tickerMap := map[string]BinanceTicker{}
		tickerMap["ATOMUSDT"] = BinanceTicker{
			Symbol:    "ATOMUSDT",
			LastPrice: lastPrice,
			Volume:    volume,
		}

		p.tickers = tickerMap

		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "ATOM", Quote: "USDT"})
		require.NoError(t, err)
		require.Len(t, prices, 1)
		require.Equal(t, sdk.MustNewDecFromStr(lastPrice), prices["ATOMUSDT"].Price)
		require.Equal(t, sdk.MustNewDecFromStr(volume), prices["ATOMUSDT"].Volume)
	})

	t.Run("valid_request_multi_ticker", func(t *testing.T) {
		lastPriceAtom := "34.69000000"
		lastPriceShe := "41.35000000"
		volume := "2396974.02000000"

		tickerMap := map[string]BinanceTicker{}
		tickerMap["ATOMUSDT"] = BinanceTicker{
			Symbol:    "ATOMUSDT",
			LastPrice: lastPriceAtom,
			Volume:    volume,
		}

		tickerMap["SHEUSDT"] = BinanceTicker{
			Symbol:    "SHEUSDT",
			LastPrice: lastPriceShe,
			Volume:    volume,
		}

		p.tickers = tickerMap
		prices, err := p.GetTickerPrices(
			types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
			types.CurrencyPair{Base: "BLT", Quote: "USDT"},
		)
		require.NoError(t, err)
		require.Len(t, prices, 2)
		require.Equal(t, sdk.MustNewDecFromStr(lastPriceAtom), prices["ATOMUSDT"].Price)
		require.Equal(t, sdk.MustNewDecFromStr(volume), prices["ATOMUSDT"].Volume)
		require.Equal(t, sdk.MustNewDecFromStr(lastPriceShe), prices["SHEUSDT"].Price)
		require.Equal(t, sdk.MustNewDecFromStr(volume), prices["SHEUSDT"].Volume)
	})

	t.Run("invalid_request_invalid_ticker", func(t *testing.T) {
		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "FOO", Quote: "BAR"})
		require.NoError(t, err)
		require.Zero(t, len(prices))
	})
}

func TestBinanceProvider_SubscribeCurrencyPairs(t *testing.T) {
	server := NewMockProviderServer()
	server.Start()
	defer server.Close()

	p, err := NewBinanceProvider(
		context.TODO(),
		zerolog.Nop(),
		config.ProviderEndpoint{
			Name:      config.ProviderBinance,
			Rest:      "",
			Websocket: server.GetBaseURL(),
		},
		types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
	)
	require.NoError(t, err)

	t.Run("invalid_subscribe_channels_empty", func(t *testing.T) {
		err = p.SubscribeCurrencyPairs([]types.CurrencyPair{}...)
		require.ErrorContains(t, err, "currency pairs is empty")
	})
}

func TestBinanceCurrencyPairToBinancePair(t *testing.T) {
	cp := types.CurrencyPair{Base: "ATOM", Quote: "USDT"}
	binanceSymbol := currencyPairToBinanceTickerPair(cp)
	require.Equal(t, binanceSymbol, "atomusdt@ticker")
}
