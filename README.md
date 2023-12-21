# MAC-Generator
Generate MAC address

## Usage

```shell
# sage: mac-generator [options]
# Options:
# -format -f string  File format: csv, or json (default "csv") Specify comma separated value like \"json,csv\", can multi file output"
# -help -h           Show help information
# -hyphen            Use hyphen (-) as delimiter in MAC addresses
# -number -n int     Number of MAC addresses to generate. Default is 10
# -out -o string     Output file path and name without extension (required).if specified dir is not exist, shell will fail.
# -surround -s       Surround each MAC address with single quotes  
```

## Sample

```shell
# generate 500 MacAddress in CSV to current directory.output file name is macAddress.csv
  $ mac-generator -n 500 -o ./macAddress

# generate 500 MacAddress in JSON to specified directory.output file name is macAddress.json
  $ mac-generator -n 500 -o ./some/any/macAddress -f json

# generate 500 MacAddress in JSON and CSV to specified directory.output file name is macAddress.json, macAddress.csv
  $ mac-generator -n 500 -o ./some/any/macAddress -f json,csv

# generate 10 MacAddress using hyphen in CSV to current directory. eg. 86-56-d0-7f-af-69
  $ mac-generator -o ./macAddress -hyphen

# generate 10 MacAddress using hyphen in CSV to current directory. eg. '7c:4e:46:e2:9c:ae'
  $ mac-generator -o ./macAddress -s
```