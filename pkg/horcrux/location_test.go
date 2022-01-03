package horcrux

import (
	"testing"
)

func TestParseSecretKeyLocation(t *testing.T) {
	tests := []struct {
		name     string
		stdout   string
		keyID    string
		expected KeyLocation
	}{
		{
			"Unknown Stub",
			`
sec:u:4096:1:1667D82C2BF9F3E1:1549414784:::u:::cESCA:::#:::23::0:
fpr:::::::::E68A304BC1806237B05CD2A21667D82C2BF9F3E1:
grp:::::::::50B2730B2CA442F34623550AA95F8F2F1CC904F7:
uid:u::::1631217630::7F28D1DB1D5FBD697DC5DC6E3ADACE825E370340::Rob Witoff <rob@unit410.com>::::::::::0:
uid:u::::1549414784::45792D1AC904013D8C4B22BC0001112DFDEFEF3C::Rob Witoff <rob@polychain.capital>::::::::::0:
ssb:u:4096:1:FA98AED0B2614592:1549414987::::::s:::D2760001240100000006086221730000:::23:
fpr:::::::::4C840A7FF94779A156CE1600FA98AED0B2614592:
grp:::::::::8BEE4E9964A46F4246E524AD0C856709AAF3F194:
ssb:u:4096:1:52877655EB3E351B:1549415052::::::e:::D2760001240100000006086221730000:::23:
fpr:::::::::21F2004DFD8F1C9D62C3AEC852877655EB3E351B:
grp:::::::::BFFDCF8505F9201EE845EED69AE93B3BA5E3A134:
ssb:u:4096:1:E57C6E0A4A999F55:1549415083::::::a:::D2760001240100000006086221730000:::23:
fpr:::::::::10257E5E0924F47F311E9BD2E57C6E0A4A999F55:
grp:::::::::E0E2152402C0454D238895AB6A1DCF1D55EC16DF:`,
			"1667D82C2BF9F3E1",
			KeyLocationStub,
		},
		{
			"Stub on Card",
			`
sec:u:4096:1:1667D82C2BF9F3E1:1549414784:::u:::cESCA:::#:::23::0:
fpr:::::::::E68A304BC1806237B05CD2A21667D82C2BF9F3E1:
grp:::::::::50B2730B2CA442F34623550AA95F8F2F1CC904F7:
uid:u::::1631217630::7F28D1DB1D5FBD697DC5DC6E3ADACE825E370340::Rob Witoff <rob@unit410.com>::::::::::0:
uid:u::::1549414784::45792D1AC904013D8C4B22BC0001112DFDEFEF3C::Rob Witoff <rob@polychain.capital>::::::::::0:
ssb:u:4096:1:FA98AED0B2614592:1549414987::::::s:::D2760001240100000006086221730000:::23:
fpr:::::::::4C840A7FF94779A156CE1600FA98AED0B2614592:
grp:::::::::8BEE4E9964A46F4246E524AD0C856709AAF3F194:
ssb:u:4096:1:52877655EB3E351B:1549415052::::::e:::D2760001240100000006086221730000:::23:
fpr:::::::::21F2004DFD8F1C9D62C3AEC852877655EB3E351B:
grp:::::::::BFFDCF8505F9201EE845EED69AE93B3BA5E3A134:
ssb:u:4096:1:E57C6E0A4A999F55:1549415083::::::a:::D2760001240100000006086221730000:::23:
fpr:::::::::10257E5E0924F47F311E9BD2E57C6E0A4A999F55:
grp:::::::::E0E2152402C0454D238895AB6A1DCF1D55EC16DF:`,
			"52877655EB3E351B",
			KeyLocationStub,
		},
		{"Local",
			`
sec:-:4096:1:236C461456BD7016:1640994825:::-:::escESC:::+:::23::0:
fpr:::::::::9024D1F036D814AA9FD60C44236C461456BD7016:
grp:::::::::A0F408196BB0DFD89D38124D7567AD1D2F32FECC:
uid:-::::1640994825::082EA91A3AC99BEEB1F641EBDCC3C6B58462421A::test-3@example.com::::::::::0:
sec:-:4096:1:9250FE8E21C53CD0:1640994815:::-:::escESC:::+:::23::0:
fpr:::::::::171EC2DC549DC008B30E30CA9250FE8E21C53CD0:
grp:::::::::4A10CA3B1E573BCD637E0198F82CF0AA5A656251:
uid:-::::1640994815::BC84C34B3EBF5E0FE2AF5550633839CD4395F919::test-1@example.com::::::::::0:
	`,
			"236C461456BD7016",
			KeyLocationLocal,
		},
		{"Local",
			`
sec:-:4096:1:236C461456BD7016:1640994825:::-:::escESC:::+:::23::0:
fpr:::::::::9024D1F036D814AA9FD60C44236C461456BD7016:
grp:::::::::A0F408196BB0DFD89D38124D7567AD1D2F32FECC:
uid:-::::1640994825::082EA91A3AC99BEEB1F641EBDCC3C6B58462421A::test-3@example.com::::::::::0:
sec:-:4096:1:9250FE8E21C53CD0:1640994815:::-:::escESC:::+:::23::0:
fpr:::::::::171EC2DC549DC008B30E30CA9250FE8E21C53CD0:
grp:::::::::4A10CA3B1E573BCD637E0198F82CF0AA5A656251:
uid:-::::1640994815::BC84C34B3EBF5E0FE2AF5550633839CD4395F919::test-1@example.com::::::::::0:
	`,
			"9250FE8E21C53CD0",
			KeyLocationLocal,
		},
		{"Local",
			`sec:-:4096:1:236C461456BD7016:1640994825:::-:::escESC:::+:::23::0:`,
			"ABCDEFABCEDFABCD",
			KeyLocationUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parseSecretKeyLocation([]byte(tt.stdout), tt.keyID)
			if actual != tt.expected {
				t.Errorf("parseSecretKeyLocation() = %v, expected %v", actual, tt.expected)
			}
		})
	}
}
