package horcrux

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestGetEntityNames(t *testing.T) {
	bytes, _ := hex.DecodeString("c6c14d0461cf9809011000e7e86de4f7f670c2be8cefc9da2d9b709f022b8a69d1755a69f20f98ccebc1cb6c8cfbd68237f3f5339c95186f98a6ce94dd1b9c80b48b98abc6eb3140b7b0c0209a196574ecbfdbf8bc58b67260ac51a8f2fa8798509adb85e007869ec890d7f65896f772ebc1e611728f2c63e503c1c01896685d9b97858496c32c43ca128c0fc6fa58d12f7ce19181923667ebdbb67536e23e310d3df682acf6c512da462263e9c17d5d578d2a6ce93b587640bb9318bcaff7d471d17766a3d78a3cf4ab457cf32584990305f3bfbbb199428301ac30b714bc5d37152ee288bc95c1ee6df2f817e25e135588275fc8be8884c1cc3c37d8a95f1040af0b51fa641e5727610d35dc298366d1d8a028cbeb957d9723c9c037553b53b31886e8c97a4dc9e950dbdb3c681c9ba6050c794ef7c2c35475d4c668d38fdb2d566a4e4d500e610636d9c94383cc211bfdae3ea383895d257eb1292cf698bfedc81dd38b7b26d775d766327d770c8188aa5e1418814694c8d6af2eeb5229d12dfd7304250aa619d23442a4cf810f5d9e3e3ef0e1320119b82f5f86163a5655395f063170a76434c9c38886edcb7d76daa948f6361ffa1f5acaf2017ed62ffe8a2f64790fcdb9724949df460f45463908512f093c950d115cb0e9bc26aef42a0c25815b941d05db78e13c90b3ffe23cfa070d8300ef5139a6b0a677bebd1053aeeeb46dacb4f376d33ba70011010001cd12746573742d33406578616d706c652e636f6dc2c19204130108003c1621049024d1f036d814aa9fd60c44236c461456bd7016050261cf9809021b0f050b090807020322020106150a09080b020416020301021e07021780000a0910236c461456bd70163aa40fff595d0375f5f2e1d790f3a1162518ad03fb66339f6d3d1d67928ad95e483d2606c287cd0187aa6c524878db896719a8fa41f28148fa7fd21118a4844097690dcbb3e602cd7254bb80a2b7fb0eac57b119318faac91201f4482b26f94703ffa022cf580e18b18b43e9da54cd66de6c0c35892be8f8e0dc643eb8b52afb4b9cc7cc3e036f8c66cfd00ef189f0fca697711d88ddd66ea251a6ae9f59630ad904a88f16d741cee764447c1c1a625f42b6d5e9d450fe272f18d9515417e457d2c3fb0c0ab7606fefa5146f67b38d94b566118f0370c6a9071193775bd083c9ebc5a4df564dbccecfe89e7e931026a1536eed837b6d118be183d512dc5b36be8400ca0e5a12255c4774d4c12e856acaf055ac89176a531cf1f487b2501d1bb711396f54732b58ee50f45043261c68c037dda8e8e1ea510520f9658eca02d60471443450637642da57a98d941a30f5ffd4275f953ce094b8257633404eb8cd0664dc32027aaa141f83bf58b574bd06684a3125de236c59b0ff9c3e61ee9da0a733bda662a3e9be86340dab5c0741e1d1b6ac460b52a8015ec9473938df99eba187470b9bf23e98e53a8679714ba8be90b2349edf3dddd17073e9554d9f15ee0a3ebf5f5fcecd231a4f9166c89b3f85ac01564bbc41a82c3221d846c1766c0a9d6e46c1f2405d6e29c07c5380204de27cc5e9b9009efb76236d64ba9c3b25d85dae3d84b1")
	keyID := "236C461456BD7016"
	expected := []string{"test-3@example.com"}
	actual := getEntityNames(bytes, keyID)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("getEntityNames() = %v, expected %v", actual, expected)
	}
}
