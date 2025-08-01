package main

import (
  "archive/zip"
  "bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func printUsageAndExit() {
  fmt.Println("Usage: file-transfer-over-powershell FILE")
  os.Exit(1)
}

func must[T any](val T, err error) T {
  if err != nil {
    panic(err)
  }
  return val
}

func main() {
  flag.Parse()

  if flag.NArg() != 1 {
    printUsageAndExit()
  }

  filePath := flag.Arg(0)
  fileName := filepath.Base(filePath)
  fileBytes := must(os.ReadFile(filePath))

  buf := new(bytes.Buffer)
  zipWriter := zip.NewWriter(buf)
  zipFile := must(zipWriter.Create(fileName))
  zipFile.Write(fileBytes)
  must(true, zipWriter.Close())
  zipBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
  zipFileName := fmt.Sprintf("%s.zip", fileName)

  fmt.Printf(`$b64 = '%s'
$filename = "$env:TEMP\%s"
$bytes = [Convert]::FromBase64String($b64)
[IO.File]::WriteAllBytes($filename, $bytes)
explorer.exe "$env:TEMP"
`, zipBase64, zipFileName)
}
