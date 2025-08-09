package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tobiashort/clap-go"
	. "github.com/tobiashort/utils-go/must"
)

type Args struct {
	File string `clap:"positional,mandatory,description='The file.'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	filePath := args.File
	fileName := filepath.Base(filePath)
	fileBytes := Must2(os.ReadFile(filePath))

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	zipFile := Must2(zipWriter.Create(fileName))
	zipFile.Write(fileBytes)
	Must2(true, zipWriter.Close())
	zipBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	zipFileName := fmt.Sprintf("%s.zip", fileName)

	fmt.Printf(`$b64 = '%s'
$filename = "$env:TEMP\%s"
$bytes = [Convert]::FromBase64String($b64)
[IO.File]::WriteAllBytes($filename, $bytes)
explorer.exe "$env:TEMP"
`, zipBase64, zipFileName)
}
