package main

import (
	"strings"
	"strconv"
)

/**
 * Get Id of an object from a Key
 *
 * @params key string - key of an object i.e. "UserAsset12345"
 * @params prefix string - prefix for an object key i.e. "UserAsset"
 *
 * @return id int64 - i.e. 12345
 */
func GetId(key string, prefix string) int64 {
	id, _ := strconv.Atoi(strings.Split(key, prefix)[1])

	return int64(id)
}

/**
 * Get Index of a element in a array
 *
 * @params key string - key of an object i.e. "UserAsset1"
 * @params data string - array of keys i.e. ["UserAsset1", "UserAsset2"]
 *
 * @return id int - i.e. 0 or -1 if not found
 */
func indexOf(key string, data []string) (int) {
	for k, v := range data {
		if key == v {
			return k
		}
	}
	return -1
}

/**
 * Remove element from an array using an index
 *
 * @params s string - array of keys i.e. ["UserAsset1", "UserAsset2"]
 * @params i int - index of the object you want to delete i.e. index=1 for UserAsset2
 *
 * @return keys array - new array with the index element removed
 */
func removeKey(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}