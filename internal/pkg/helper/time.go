package helper

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func TimeZoneOffsetToSeconds(offset string) (int, error) {
	validFormat := regexp.MustCompile(`^[+-]\d{2}:\d{2}$`)
	if !validFormat.MatchString(offset) {
		return 0, errors.New("invalid time zone offset format")
	}

	sign := 1
	if offset[0] == '-' {
		sign = -1
	}
	parts := strings.Split(offset[1:], ":")
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	return sign * (hours*3600 + minutes*60), nil
}
