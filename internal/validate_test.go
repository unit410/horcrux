package internal

import (
	"testing"

	"gitlab.com/polychainlabs/vault-shamir/shamir"
)

func TestCheckAllCombinations(test *testing.T) {
	// 2 of 2
	{
		original := []byte{3, 7, 9}
		numParts := 2
		threshold := 2
		parts, err := shamir.Split(original, numParts, threshold)
		if err != nil {
			test.Fatal(err)
		}

		if !CheckAllCombinations(original, parts, 2) {
			test.Fatal("Check combinations failed - expected pass")
		}

		// editing any part should fail check combinations
		{
			old := parts[0][0]
			parts[0][0] = 1
			if CheckAllCombinations(original, parts, 2) {
				test.Fatal("Check combinations passed - expected fail")
			}
			parts[0][0] = old
		}

		if !CheckAllCombinations(original, parts, 2) {
			test.Fatal("Check combinations failed - expected pass")
		}

		{
			old := parts[1][0]
			parts[1][0] = 1
			if CheckAllCombinations(original, parts, 2) {
				test.Fatal("Check combinations passed - expected fail")
			}
			parts[1][0] = old
		}
	}

	// split into 2 of 3
	{
		original := []byte{3, 7, 9}
		numParts := 3
		threshold := 2
		parts, err := shamir.Split(original, numParts, threshold)
		if err != nil {
			test.Fatal(err)
		}

		// all sets of 2 combinations should pass
		if !CheckAllCombinations(original, parts, 2) {
			test.Fatal("Check combinations failed - expected pass")
		}

		// all sets of 3 combinations should also pass
		if !CheckAllCombinations(original, parts, 3) {
			test.Fatal("Check combinations failed - expected pass")
		}

		// editing any part should fail check combinations even tho only two parts are needed
		{
			old := parts[0][0]
			parts[0][0] = 1
			if CheckAllCombinations(original, parts, 2) {
				test.Fatal("Check combinations passed - expected fail")
			}
			parts[0][0] = old
		}

		if !CheckAllCombinations(original, parts, 2) {
			test.Fatal("Check combinations failed - expected pass")
		}

		{
			old := parts[1][0]
			parts[1][0] = 1
			if CheckAllCombinations(original, parts, 2) {
				test.Fatal("Check combinations passed - expected fail")
			}
			parts[1][0] = old
		}

		if !CheckAllCombinations(original, parts, 2) {
			test.Fatal("Check combinations failed - expected pass")
		}

		{
			old := parts[2][0]
			parts[2][0] = 1
			if CheckAllCombinations(original, parts, 2) {
				test.Fatal("Check combinations passed - expected fail")
			}
			parts[2][0] = old
		}
	}
}

func TestCheckAllCombinations_fail(test *testing.T) {
	// 3 of 3
	{
		original := []byte{3, 7, 9}
		numParts := 3
		threshold := 3
		parts, err := shamir.Split(original, numParts, threshold)
		if err != nil {
			test.Fatal(err)
		}

		// should not pass because we specified to check only 2 part combinations
		if CheckAllCombinations(original, parts, 2) {
			test.Fatal("Check combinations passed - expected fail")
		}
	}
}
