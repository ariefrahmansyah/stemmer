package stemmer

import "testing"

func TestInit(t *testing.T) {

}

func TestStemm(t *testing.T) {
	testCases := []struct {
		word     string
		baseWord string
	}{
		{"mei", "mei"},
		{"bui", "bui"},
		{"nilai", "nilai"},
		{"hancurlah", "hancur"},
		{"benarkah", "benar"},
		{"apatah", "apa"},
		{"siapapun", "siapa"},
		{"jubahku", "jubah"},
		{"bajumu", "baju"},
		{"celananya", "celana"},
		{"hantui", "hantu"},
		{"belikan", "beli"},
		{"jualan", "jual"},
		{"bukumukah", "buku"},
		{"miliknyalah", "milik"},
		{"kulitkupun", "kulit"},
		{"berikanku", "beri"},
		{"sakitimu", "sakit"},
		{"beriannya", "beri"},
		{"kasihilah", "kasih"},
		{"dibuang", "buang"},
		{"kesakitan", "sakit"},
		{"sesuap", "suap"},
		{"teriakanmu", "teriak"},
		{"beradu", "adu"},
		{"berambut", "rambut"},
		{"bersuara", "suara"},
		{"berdaerah", "daerah"},
		{"belajar", "ajar"},
		{"bekerja", "kerja"},
		{"beternak", "ternak"},
		{"terasing", "asing"},
		{"teraup", "raup"},
		{"tergerak", "gerak"},
		{"terpuruk", "puruk"},
		{"teterbang", "terbang"},
		{"melipat", "lipat"},
		{"meringkas", "ringkas"},
		{"mewarnai", "warna"},
		{"meyakinkan", "yakin"},
		{"membangun", "bangun"},
		{"memfitnah", "fitnah"},
		{"memvonis", "vonis"},
		{"memperbarui", "baru"},
		{"mempelajari", "ajar"},
		{"meminum", "minum"},
		{"memukul", "pukul"},
		{"mencinta", "cinta"},
		{"mendua", "dua"},
		{"menjauh", "jauh"},
		{"menziarah", "ziarah"},
		{"menuklir", "nuklir"},
		{"menangkap", "tangkap"},
		{"menggila", "gila"},
		{"menghajar", "hajar"},
		{"mengqasar", "qasar"},
		{"mengudara", "udara"},
		{"mengupas", "kupas"},
		{"menyuarakan", "suara"},
		{"mempopulerkan", "populer"},
		{"pewarna", "warna"},
		{"peyoga", "yoga"},
		{"peradilan", "adil"},
		{"perumahan", "rumah"},
		{"permuka", "muka"},
		{"perdaerah", "daerah"},
		{"pembangun", "bangun"},
		{"pemfitnah", "fitnah"},
		{"pemvonis", "vonis"},
		{"peminum", "minum"},
		{"pemukul", "pukul"},
		{"pencinta", "cinta"},
		{"pendahulu", "dahulu"},
		{"penjarah", "jarah"},
		{"penziarah", "ziarah"},
		{"penasihat", "nasihat"},
		{"penangkap", "tangkap"},
		{"penggila", "gila"},
		{"penghajar", "hajar"},
		{"pengqasar", "qasar"},
		{"pengudara", "udara"},
		{"pengupas", "kupas"},
		{"penyuara", "suara"},
		{"pelajar", "ajar"},
		{"pelabuhan", "labuh"},
		{"petarung", "tarung"},
		{"terpercaya", "percaya"},
		{"pekerja", "kerja"},
		{"peserta", "serta"},
		{"mempengaruhi", "pengaruh"},
		{"mengkritik", "kritik"},
		{"bersekolah", "sekolah"},
		{"bertahan", "tahan"},
		{"mencapai", "capai"},
		{"dimulai", "mulai"},
		{"petani", "tani"},
		{"terabai", "abai"},
		{"mensyaratkan", "syarat"},
		{"mensyukuri", "syukur"},
		{"mengebom", "bom"},
		{"mempromosikan", "promosi"},
		{"memproteksi", "proteksi"},
		{"memprediksi", "prediksi"},
		{"pengkajian", "kaji"},
		{"pengebom", "bom"},
		{"bersembunyi", "sembunyi"},
		{"bersembunyilah", "sembunyi"},
		{"pelanggan", "langgan"},
		{"pelaku", "laku"},
		{"pelangganmukah", "langgan"},
		{"pelakunyalah", "laku"},
		{"perbaikan", "baik"},
		{"kebaikannya", "baik"},
		{"bisikan", "bisik"},
		{"menerangi", "terang"},
		{"berimanlah", "iman"},
		{"memuaskan", "puas"},
		{"berpelanggan", "langgan"},
		{"bermakanan", "makan"},
		{"menyala", "nyala"},
		{"menyanyikan", "nyanyi"},
		{"menyatakannya", "nyata"},
		{"penyanyi", "nyanyi"},
		{"penyawaan", "nyawa"},
		{"bertebaran", "tebar"},
		{"terasingkan", "asing"},
		{"membangunkan", "bangun"},
		{"mencintai", "cinta"},
		{"menduakan", "dua"},
		{"menjauhi", "jauh"},
		{"menggilai", "gila"},
		{"pembangunan", "bangun"},
		{"memberdayakan", "daya"},
		{"persemakmuran", "makmur"},
		{"keberuntunganmu", "untung"},
		{"Perekonomian", "ekonomi"},
		{"menahan", "tahan"},
		{"peranan", "peran"},
		{"memberikan", "beri"},
		{"medannya", "medan"},
		{"sebagai", "bagai"},
		{"bagian", "bagi"},
		{"berbadan", "badan"},
		{"budayawan", "budaya"},
		{"karyawati", "karya"},
		{"pendaratan", "darat"},
		{"penstabilan", "stabil"},
		{"pentranskripsi", "transkripsi"},
		{"mentaati", "taat"},
		{"melewati", "lewat"},
		{"menganga", "nganga"},
	}

	s := New()
	for k, tc := range testCases {
		out := s.Stemm(tc.word)
		if out[0] != tc.baseWord {
			t.Error(k+1, tc.word, tc.baseWord, out[0])
			return
		}
	}
}

func TestInitRootWordsSuccess(t *testing.T) {
	InitRootWords()
}

func TestIsRootWord(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{"cinta", true},
		{"mencintai", false},
	}

	s := New()
	for _, tc := range testCases {
		if s.IsRootWord([]byte(tc.in)) != tc.out {
			t.Error(tc.in)
		}
	}
}

func TestIsRulePrecedence(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{"berkawanlah", true},
		{"dicintai", true},
		{"mencintai", true},
		{"dicinta", false},
		{"mencinta", false},
	}

	s := New()
	for _, tc := range testCases {
		if s.isRulePrecedence([]byte(tc.in)) != tc.out {
			t.Error(tc.in)
		}
	}
}

func TestIsDisallowedPrefixSuffixes(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{"berlari", true},
		{"dikeabadian", true},
		{"kembalikan", true},
		{"merasakan", true},
		{"selagi", true},
		{"tekaan", true},

		{"bercinta", false},
		{"dicinta", false},
		{"mencinta", false},
	}

	s := New()
	for _, tc := range testCases {
		if s.isDisallowedPrefixSuffixes([]byte(tc.in)) != tc.out {
			t.Error(tc.in)
		}
	}
}

func TestRemoveInflectionSuffixes(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{"apalah", "apa"},
		{"apakah", "apa"},
		{"apatah", "apa"},
		{"apapun", "apa"},

		{"cintaku", "cinta"},
		{"cintamu", "cinta"},
		{"cintanya", "cinta"},
	}

	s := New()
	for _, tc := range testCases {
		out := s.removeInflectionSuffixes([]byte(tc.in))
		if string(out) != tc.out {
			t.Error(tc.in)
		}
	}
}

func TestRemoveDerivationSuffixes(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{"cintai", "cinta"},
		{"tekaan", "teka"},
	}

	s := New()
	for _, tc := range testCases {
		out := s.removeDerivationSuffixes([]byte(tc.in))
		if string(out) != tc.out {
			t.Error(tc.in)
		}
	}
}

func TestRemoveDerivationPeople(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{"budayawan", "budaya"},
		{"karyawan", "karya"},
		{"karyawati", "karya"},
		{"ilmuwan", "ilmu"},
		{"seniman", "seni"},
	}

	s := New()
	for _, tc := range testCases {
		out := s.removeDerivationPeople([]byte(tc.in))
		if string(out) != tc.out {
			t.Error(tc.in)
		}
	}
}
