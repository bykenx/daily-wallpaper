package api

import (
	"fmt"

	"github.com/andybrewer/mack"
)

func SetWallpaper(path string) error {
	res, err := mack.Tell(
		"Finder",
		"activate",
		fmt.Sprintf(`set desktop picture to POSIX file "%s"`, path),
	)
	if err != nil {
		return err
	}
	println(res)
	return nil
}
