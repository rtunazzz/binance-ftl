package ftl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPosDir(t *testing.T) {
	tests := []struct {
		rp  rawPosition
		exp PositionDirection
		msg string
	}{
		{
			rp:  rawPosition{EntryPrice: 1000, MarkPrice: 2000, Pnl: 1000},
			exp: Long,
			msg: "should be long with positive PNL",
		},
		{
			rp:  rawPosition{EntryPrice: 1000, MarkPrice: 500, Pnl: -500},
			exp: Long,
			msg: "should be long with negative PNL",
		},
		{
			rp:  rawPosition{EntryPrice: 1000, MarkPrice: 500, Pnl: 500},
			exp: Short,
			msg: "should be short with positive PNL",
		},
		{
			rp:  rawPosition{EntryPrice: 1000, MarkPrice: 1500, Pnl: -500},
			exp: Short,
			msg: "should be short with negative PNL",
		},
	}

	for _, tc := range tests {
		dir := getPosDir(tc.rp)
		assert.Equal(t, tc.exp, dir, tc.msg)
	}
}

func TestHashEquality(t *testing.T) {
	tests := []struct {
		p1  Position
		p2  Position
		msg string
	}{
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			msg: "two exactly same positions",
		},
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: AddedTo, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20*10 ^ 3, Amount: 2},
			msg: "added to position at the same price",
		},
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: AddedTo, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 25 * 10e3, Amount: 2},
			msg: "added to position at a different price",
		},
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: PartiallyClosed, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 0.5},
			msg: "partially closed position",
		},
	}

	for _, tc := range tests {
		h1, err := tc.p1.hash()
		assert.Nil(t, err, "hashing errored out")

		h2, err := tc.p2.hash()
		assert.Nil(t, err, "hashing errored out")

		assert.Equal(t, h1, h2, tc.msg)
	}
}

func TestHashInequality(t *testing.T) {
	tests := []struct {
		p1  Position
		p2  Position
		msg string
	}{
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: Opened, Direction: Long, Ticker: "ETHSDT", EntryPrice: 20 * 10e3, Amount: 1},
			msg: "two differnt tickers",
		},
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: Opened, Direction: Short, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			msg: "two differnt directions",
		},
	}

	for _, tc := range tests {
		h1, err := tc.p1.hash()
		assert.Nil(t, err, "hashing errored out")

		h2, err := tc.p2.hash()
		assert.Nil(t, err, "hashing errored out")

		assert.NotEqual(t, h1, h2, tc.msg)
	}
}