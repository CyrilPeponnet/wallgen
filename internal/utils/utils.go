package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseHexColor(v string) (out [3]int, err error) {

	if len(v) != 7 {
		return out, errors.New("hex color must be 7 characters")
	}
	if v[0] != '#' {
		return out, errors.New("hex color must start with '#'")
	}
	var red, redError = strconv.ParseUint(v[1:3], 16, 8)
	if redError != nil {
		return out, errors.New("red component invalid")
	}
	out[0] = int(red)
	var green, greenError = strconv.ParseUint(v[3:5], 16, 8)
	if greenError != nil {
		return out, errors.New("green component invalid")
	}
	out[1] = int(green)
	var blue, blueError = strconv.ParseUint(v[5:7], 16, 8)
	if blueError != nil {
		return out, errors.New("blue component invalid")
	}
	out[2] = int(blue)
	return
}

func ParsePalette(s string) map[string]int {

	p := make(map[string]int)
	for _, chunk := range strings.Split(s, ",") {
		parts := strings.Split(chunk, ":")
		if len(parts) != 2 {
			panic(fmt.Errorf("unable to split palette chunk: %s", chunk))
		}
		w, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		p[parts[0]] = w
	}

	return p
}
