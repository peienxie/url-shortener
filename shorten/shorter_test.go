package shorten

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItoaLeadingZeros(t *testing.T) {
	testCases := []struct {
		name   string
		value  uint64
		width  int
		expect string
	}{
		{
			"exactly width",
			102938190248,
			12,
			"102938190248",
		},
		{
			"less width, no padding",
			102938190248,
			12,
			"102938190248",
		},
		{
			"greater width, padding leading zeros upto width",
			102938190248,
			15,
			"000102938190248",
		},
		{
			"padding leading zeros upto 20 width for MaxUint64",
			math.MaxUint64,
			20,
			"18446744073709551615",
		},
		{
			"padding leading zeros upto 21 width for MaxUint64",
			math.MaxUint64,
			21,
			"018446744073709551615",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := string(itoaLeadingZeros(tc.value, tc.width))
			require.Equal(t, tc.expect, output)
		})
	}
}

func TestShortenByInt(t *testing.T) {
	testCases := []struct {
		name   string
		value  uint64
		expect string
	}{
		{
			"base62 encode 102938190248",
			102938190248,
			"4QjMwkTM4MTOyATMwADMwADMwAD",
		},
		{
			"base62 encode 1",
			1,
			"xADMwADMwADMwADMwADMwADMwAD",
		},
		{
			"base62 encode 0",
			0,
			"wADMwADMwADMwADMwADMwADMwAD",
		},
		{
			"base62 encode MaxUint64",
			math.MaxUint64,
			"1EjNxUTN5AzNzcDM0QzN2QDN4ED",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := string(ShortenByInt(tc.value))
			require.Equal(t, tc.expect, output)
		})
	}
}
