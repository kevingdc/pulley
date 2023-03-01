package idconv

import "strconv"

func ToRepoID(toConvert int64) string {
	return strconv.FormatInt(toConvert, 10)
}
