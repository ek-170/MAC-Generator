package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/exp/slices"
)

// MACAddressGenerator はMACアドレスを生成する構造体
type MACAddressGenerator struct {
	Number   int
	Format   string
	Output   string
	Hyphen   bool
	Surround bool
}

// GenerateMACAddresses は指定された数のMACアドレスを生成する
func (m *MACAddressGenerator) GenerateMACAddresses() ([]string, error) {
	var macAddresses []string
	for i := 0; i < m.Number; i++ {
		mac := generateRandomMAC(m.Hyphen)
    // if m.Surround && (m.Format == "csv") {
    // 	macAddresses = append(macAddresses, mac)
    // } else {
        macAddresses = append(macAddresses, mac)
    // }
	}
	return macAddresses, nil
}

// WriteToFile は生成されたMACアドレスを指定されたファイルに書き込む
func (m *MACAddressGenerator) WriteToFile(macAddresses []string) error {

	formats := strings.Split(m.Format, ",")
	var format string
	if len(formats) == 1 {
		format = formats[0]
	}else if (len(formats)==2 && slices.Contains(formats, "json") && slices.Contains(formats, "csv")){
		format = "json,csv"
	}
	switch format {
	case "json":
		m.writeToJson(macAddresses)
	case "csv":
		m.writeToCsv(macAddresses)
	case "json,csv":
		m.writeToJson(macAddresses)
		m.writeToCsv(macAddresses)
	default:
		return fmt.Errorf("invalid format specified: %s", m.Format)
	}

	return nil
}

func (m *MACAddressGenerator) writeToJson(macAddresses []string) error {
	file, err := os.Create(m.Output + ".json")
	if err != nil {
		return err
	}
	defer file.Close()
	data := map[string][]string{"macAddress": macAddresses}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (m *MACAddressGenerator) writeToCsv(macAddresses []string) error {
	file, err := os.Create(m.Output + ".csv")
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := writer.Write(macAddresses); err != nil {
			return err
	}
	return nil
}

func generateRandomMAC(hyphen bool) string {
    
    rand.Seed(time.Now().UnixNano())
    format := "%x"
    formats := []string{}
    for i := 0; i < 6; i++{
        n := rand.Intn(256)
        v := fmt.Sprintf(format, n)
        if utf8.RuneCountInString(v) == 1{
            v = "0" + v
        }
        formats = append(formats, v)
    }
    separator := ":"
	if hyphen {
        separator = "-"
	}
    return strings.Join(formats, separator)
}

func main() {
	// フラグの定義
	number := flag.Int("number", 10, "Number of MAC addresses to generate. Default is 10")
	format := flag.String("format", "csv", "File format: csv, json. Defalut is csv. Specify comma separated value like \"json,csv\", can multi file output")
	output := flag.String("out", "", "Output file path and name without extension (required)")
	hyphen := flag.Bool("hyphen", false, "Use hyphen (-) as delimiter in MAC addresses")
	// surround := flag.Bool("surround", false, "Surround each MAC address with double quotes")
	help := flag.Bool("help", false, "Show help information")

	// ショートフラグの定義
	flag.IntVar(number, "n", 0, "Number of MAC addresses to generate. Default is 10")
	flag.StringVar(format, "f", "csv", "File format: csv, json. Defalut is csv. Specify comma separated value like \"json,csv\", can multi file output")
	flag.StringVar(output, "o", "", "Output file path and name without extension (required)")
	// flag.BoolVar(surround, "s", false, "Surround each MAC address with double quotes")
	flag.BoolVar(help, "h", false, "Show help information")

	// フラグのパース
	flag.Parse()

	// ヘルプが指定された場合はヘルプを表示
	if *help || flag.NFlag() == 0 {
		showHelp()
		os.Exit(0)
	}

	// 必須の引数が指定されているか確認
	if *number <= 0 {
		fmt.Fprintln(os.Stderr, "Error: Number of MAC addresses is required and must be greater than 0.")
		os.Exit(1)
	}
	if *output == "" {
		fmt.Fprintln(os.Stderr, "Error: Output file path and name are required.")
		os.Exit(1)
	}

	// ファイルのディレクトリが存在するか確認
	outputDir := filepath.Dir(*output)
	if outputDir != "." && outputDir != "" {
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Output directory %s does not exist.\n", outputDir)
			os.Exit(1)
		}
	}

	// MACアドレス生成器の初期化
	macGenerator := &MACAddressGenerator{
		Number:   *number,
		Format:   strings.ToLower(*format),
		Output:   *output,
		Hyphen:   *hyphen,
		// Surround: *surround,
	}

	// MACアドレスの生成
	macAddresses, err := macGenerator.GenerateMACAddresses()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating MAC addresses: %s\n", err)
		os.Exit(1)
	}

	// ファイルへの書き込み
	err = macGenerator.WriteToFile(macAddresses)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %s\n", err)
		os.Exit(1)
	}

	// 正常終了時のメッセージ
	fmt.Println("MAC addresses successfully generated and written to", *output)
}

// showHelp はヘルプメッセージを表示する
func showHelp() {
	fmt.Println("Usage: mac-generator [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
