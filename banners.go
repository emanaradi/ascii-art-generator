package asciiartweb

import "fmt"

func GetBanner(style string) (Banner, error) {
	switch style {
	case "shadow":
		return Loader("banners/shadow.txt")

	case "standard":
		return Loader("banners/standard.txt")

	case "thinkertoy":
		return Loader("banners/thinkertoy.txt")

	default:
		return nil, fmt.Errorf("invalid style: %s", style)
	}
}
