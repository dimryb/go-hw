//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(""), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		invalidData := `{"Id":1,"Name":"Invalid","Username":"invalid","Email":"invalid@invalid.com","Phone":"123-45-67","Password":"invalid","Address":"Invalid Street 1"}`
		invalidData += "\n" + `{"Id":2,"Name":"Invalid","Username":"invalid","Email":"invalid@invalid.com","Phone":"123-45-67","Password":"invalid","Address":"Invalid Street 1"`
		_, err := GetDomainStat(bytes.NewBufferString(invalidData), "com")
		require.Error(t, err)
	})

	t.Run("case-insensitive domain matching", func(t *testing.T) {
		mixedCaseData := `{"Id":1,"Name":"Test User","Username":"testuser","Email":"test@MiXeDcAsE.CoM","Phone":"123-45-67","Password":"test","Address":"Test Street 1"}`
		result, err := GetDomainStat(bytes.NewBufferString(mixedCaseData), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"mixedcase.com": 1}, result)
	})

	t.Run("large input", func(t *testing.T) {
		largeData := strings.Repeat(
			`{"Id":1,"Name":"Test User","Username":"testuser","Email":"test@largeinput.com","Phone":"123-45-67","Password":"test","Address":"Test Street 1"}`+"\n",
			100000,
		)
		result, err := GetDomainStat(bytes.NewBufferString(largeData), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"largeinput.com": 100000}, result)
	})

	t.Run("empty email field", func(t *testing.T) {
		emptyEmailData := `{"Id":1,"Name":"Test User","Username":"testuser","Email":"","Phone":"123-45-67","Password":"test","Address":"Test Street 1"}`
		result, err := GetDomainStat(bytes.NewBufferString(emptyEmailData), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}
