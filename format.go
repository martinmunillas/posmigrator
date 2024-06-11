package posmigrator

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
)

var red = color.New(color.FgRed, color.Bold)
var yellow = color.New(color.FgYellow, color.Bold)
var green = color.New(color.FgGreen, color.Bold)
var blue = color.New(color.FgBlue, color.Bold)
var migrationRegex *regexp.Regexp

const invalidMigrationFileName = "invalid migration file name %s, must be in format m000_description_text.sql"

func init() {
	var err error
	migrationRegex, err = regexp.Compile(`m(\d{3})_(.+)\.sql`)
	if err != nil {
		red.Println("Error compiling regex")
		fmt.Println(err)
		os.Exit(1)
	}
}
