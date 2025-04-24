package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

// go test -bench=BenchmarkGetDomainStatLargeData -benchmem -count=1 -v -tags bench -run=^$ > result.txt .
func BenchmarkGetDomainStatLargeData(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()

	require.Equal(b, 1, len(r.File))

	data, err := r.File[0].Open()
	require.NoError(b, err)
	defer data.Close()

	domain := "biz"

	for i := 0; i < b.N; i++ {
		_, err := GetDomainStat(data, domain)
		require.NoError(b, err)
	}
}
