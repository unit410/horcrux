package internal

import (
	"bytes"
	"log"

	"gitlab.com/polychainlabs/vault-shamir/shamir"
)

func combinationRecurse(original []byte, parts [][]byte, tmp [][]byte, start int, end int, idx int, r int) bool {
	if idx == r {
		result, err := shamir.Combine(tmp)
		if err != nil {
			log.Fatal(err)
		}

		return bytes.Equal(result, original)
	}

	for i := start; i <= end && end-i+1 >= r-idx; i++ {
		tmp[idx] = parts[i]
		if !combinationRecurse(original, parts, tmp, i+1, end, idx+1, r) {
			return false
		}
	}

	return true
}

// CheckAllCombinations validates all combinations of parts against original using threshold
// size groups from parts
// return true if all combinations are sucessfully usable to create the original
func CheckAllCombinations(original []byte, parts [][]byte, threshold int) bool {
	tmp := make([][]byte, threshold)
	return combinationRecurse(original, parts, tmp, 0, len(parts)-1, 0, threshold)
}
