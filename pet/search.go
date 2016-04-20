package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type query_part struct {
	name     string
	id       int
	value    string
	required bool
}

func TableSearch(query string, filename string) {
	fmt.Println("Call to search with:", filename, "and query", query)

    var extended_query []string = strings.Split(query, "=")
    var target_column string
    var target_value string
    if len(extended_query) < 2 {
		fmt.Println("Error: malformed search query; ")
		return
    }

    target_column = strings.Trim(extended_query[0], " \t\n")
    target_value = strings.Trim(strings.Join(extended_query[1:], "="), " \t\n")

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("Error: file `", filename, "` does not exist...")
		return
	}

	var file []string

	f, err := os.Open(filename)
	s := bufio.NewScanner(f)
	for s.Scan() {
		file = append(file, s.Text())
	}

	if file[0][0] != '[' {
		fmt.Println("Error: malformed file. Unknown character `", file[0][0], "` at line 0 position 0.")
		return
	}

	var header []string = strings.Split(file[0], "][")
	var columns int
	var records int

	header[0] = header[0][1:]
	header[len(header)-1] = header[len(header)-1][0 : len(header[len(header)-1])-1]

	columns, err = strconv.Atoi(header[0])
	if err != nil {
		fmt.Println("Fatal Error: malformed file. Cannot parse column count as integer:", err)
		return
	}

	if len(header)-2 != columns {
		fmt.Println("Fatal Error: malformed file. Number of column does not match header column count:", len(header)-2, "!=", columns)
		return
	}

	records, err = strconv.Atoi(header[len(header)-1])
	if err != nil {
		fmt.Println("Fatal Error: malformed file. Cannot parse record count as integer:", err)
		return
	}

	if len(file)-1 != records {
		fmt.Println("Recoverable Error: Number of records do not match header record count. Using number in file:", len(file)-1, "vs", records)

		records = len(file) - 1
	}

	var attribute_names []string
	var attribute_types []int

	for i := range header {
		if i == 0 || i == len(header)-1 {
			continue
		}

		var item []string = strings.Split(header[i], ":")
		if len(item) != 2 {
			fmt.Println("Fatal Error: malformed header. Expected two attributes in column", i, ": got", len(item))
			return
		}

		var attribute_name string = item[0]
		var attribute_type int = 0

		attribute_type, err = strconv.Atoi(item[1])
		if err != nil || attribute_type < 1 || attribute_type > 4 {
			fmt.Println("Fatal Error: malformed header. In column", i, ": cannot parse `", item[1], "` as integer. Error: ", err)
			return
		}

		attribute_names = append(attribute_names, attribute_name)
		attribute_types = append(attribute_types, attribute_type)
	}


    var found_name int = -1;
    for i := range(attribute_names) {
        if attribute_names[i] == target_column {
            found_name = i
			if attribute_types[i] == 1 {
				_, err = strconv.Atoi(target_value)

				if err != nil {
					fmt.Println("Fatal Error: Unable to convert search query to integer", err)
					return
				}
			} else if attribute_types[i] == 2 {
				_, err = strconv.ParseFloat(target_value, 64)

				if err != nil {
					fmt.Println("Fatal Error: Unable to convert search query to double", err)
					return
				}
			} else if attribute_types[i] == 3 {
				target_value = strings.ToUpper(target_value)

				if target_value != "T" && target_value != "F" {
					fmt.Println("Fatal Error: search query unknown boolean value: expected either T or F.")
					return
				}
			} else if attribute_types[i] == 4 {
				if strings.Contains(target_value, "|") || strings.Contains(target_value, "{") || strings.Contains(target_value, "}") {
					fmt.Println("Invalid character in search query string value. Invalid characters are '|', '{'. and '}'.")
					return
				}
			}
        }
    }

    if found_name == -1 {
        fmt.Println("Fatal error: Unknown column name `", target_column, "`!")
        return
    }

    file = file[1:]
    for i := range(file) {
        var line string = file[i][1:len(file[i])-1]
        var values []string = strings.Split(line, "|")


		if attribute_types[found_name] == 1 {
			tv, err := strconv.Atoi(target_value)
			cv, err2 := strconv.Atoi(values[found_name])

			if err != nil {
				fmt.Println("Fatal Error: Unable to convert search query to integer:", err)
				return
			} else if err2 != nil {
				fmt.Println("Invalid integer in table data; ignoring row:", err2)
				continue
			} else {
		        if (cv == tv) {
		            fmt.Println("====rid:", i, "matches!====")

					var row []string = strings.Split(file[i][1:len(file[i])-1], "|")
					if len(row) != columns {
						fmt.Println("Fatal Error: mismatched number of columns: have", len(row), ", expected:", columns)
					}

					for i := range row {
						fmt.Println(attribute_names[i], "("+columTypeToName[attribute_types[i]]+"): "+row[i])
					}

					fmt.Print("\n\n")
		        }
			}
		} else if attribute_types[found_name] == 2 {
			tv, err := strconv.ParseFloat(target_value, 64)
			cv, err2 := strconv.ParseFloat(values[found_name], 64)

			if err != nil {
				fmt.Println("Fatal Error: Unable to convert search query to double", err)
				return
			} else if err2 != nil {
				fmt.Println("Invalid double in table data; ignoring row:", err2)
				continue
			} else {
		        if (cv == tv) {
		            fmt.Println("====rid:", i, "matches!====")

					var row []string = strings.Split(file[i][1:len(file[i])-1], "|")
					if len(row) != columns {
						fmt.Println("Fatal Error: mismatched number of columns: have", len(row), ", expected:", columns)
					}

					for i := range row {
						fmt.Println(attribute_names[i], "("+columTypeToName[attribute_types[i]]+"): "+row[i])
					}

					fmt.Print("\n\n")
		        }
			}
		} else if attribute_types[found_name] == 3 {
			target_value = strings.ToUpper(target_value)
			values[found_name] = strings.ToUpper(values[found_name])

			if values[found_name] != "T" && values[found_name] != "F" {
				fmt.Println("Invalid boolean in table data; ignoring row: expected either T or F.")
				continue
			} else {
		        if (values[found_name] == target_value) {
		            fmt.Println("====rid:", i, "matches!====")

					var row []string = strings.Split(file[i][1:len(file[i])-1], "|")
					if len(row) != columns {
						fmt.Println("Fatal Error: mismatched number of columns: have", len(row), ", expected:", columns)
					}

					for i := range row {
						fmt.Println(attribute_names[i], "("+columTypeToName[attribute_types[i]]+"): "+row[i])
					}

					fmt.Print("\n\n")
		        }
			}
		} else if attribute_types[found_name] == 4 {
			if strings.Contains(values[found_name], "|") || strings.Contains(values[found_name], "{") || strings.Contains(values[found_name], "}") {
				fmt.Println("Invalid string character in table data; ignoring row.")
				continue
			} else {
		        if (values[found_name] == target_value) {
		            fmt.Println("====rid:", i, "matches!====")

					var row []string = strings.Split(file[i][1:len(file[i])-1], "|")
					if len(row) != columns {
						fmt.Println("Fatal Error: mismatched number of columns: have", len(row), ", expected:", columns)
					}

					for i := range row {
						fmt.Println(attribute_names[i], "("+columTypeToName[attribute_types[i]]+"): "+row[i])
					}

					fmt.Print("\n\n")
		        }
			}
		}
    }

	fmt.Println("Successfully searched in table `", filename, "`!")
}
