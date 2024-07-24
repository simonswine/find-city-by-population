package main

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testRandom = rand.New(rand.NewSource(0))
	result     []*City
)

// test the happy path of the app
func TestFindCityBy(t *testing.T) {
	dir := t.TempDir()

	f, err := os.Create(filepath.Join(dir, "XX.txt"))
	require.NoError(t, err)

	_, err = f.Write([]byte(`2823044	Thalmässing	Thalmassing	Tal'messing,Talmesing,Thalmassing,Thalmässing,ta er mei xing,tarumesshingu,Талмесинг,Тальмессинг,Тальмессінг,Թալմեսինգ,タールメッシング,塔尔梅兴	49.08834	11.2215	P	PPL	DE		02	095	09576	09576148	5442		415	Europe/Berlin	2013-02-19
2867714	Munich	Munich	Lungsod ng Muenchen,Lungsod ng München,MUC,Minca,Minche,Minga,Minhen,Minhene,Minkhen,Miunchenas,Mjunkhen,Mnichov,Mnichow,Mníchov,Monachium,Monacho,Monaco de Baviera,Monaco di Baviera,Monaco e Baviera,Monacu,Monacu di Baviera,Monacum,Muenchen,Muenegh,Muenhen,Muenih,Munchen,Munhen,Munic,Munich,Munich ed Baviera,Munih,Munike,Munique,Munix,Munkeno,Munkhen,Munîh,Mynihu,Myunxen,Myunxén,Mònacu,Mùnich ëd Baviera,Múnic,Múnich,München,Münegh,Münhen,Münih,mi wnik,mi'unikha,miunkheni,miyunik,mu ni hei,mwinhen,mwnykh,mynkn,myunhen,myunik,myunikha,myunsena,mywnkh,mywnykh,Μόναχο,Минхен,Мюнхен,Мүнхен,Мүнхэн,Мӱнхен,Մյունխեն,מינכן,مونیخ,ميونخ,ميونيخ,میونخ,म्युन्शेन,म्यूनिख,মিউনিখ,மியூனிக்,ಮ್ಯೂನಿಕ್,มิวนิก,မြူးနစ်ချ်မြို့,მიუნხენი,ミュンヘン,慕尼黑,뮌헨	48.13743	11.57549	P	PPLA	DE		02	091	09162	09162000	1260391		524	Europe/Berlin	2023-10-12`))
	require.NoError(t, err)
	require.NoError(t, f.Close())

	mux := http.NewServeMux()
	require.NoError(t, registerHandler(mux, dir))

	// test the happy path
	{
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/by-population/5000", nil)
		mux.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "Thalmässing")
	}

	{
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?population=1000000", nil)
		mux.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "Munich")
	}
}

func BenchmarkFindCityBy(b *testing.B) {
	var err error
	dataFiles, err = findFileByExt(".", ".txt")
	if err != nil {
		b.Errorf("error while finding files: %v", err)
	}

	if len(dataFiles) == 0 {
		b.Skip("No data files found")
	}
	var (
		population int
		r          []*City
	)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		population = testRandom.Intn(100000)
		r, _ = findCityWithPopulation(population)
	}
	result = r
}
