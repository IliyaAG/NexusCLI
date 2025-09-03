package output

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "text/tabwriter"

    "gopkg.in/yaml.v3"
)

func Render(
    items []map[string]interface{},
    format string,
    headers []string,
    colorFn func(map[string]interface{}),
) {
    switch strings.ToLower(format) {
    case "json":
        data, _ := json.MarshalIndent(items, "", "  ")
        fmt.Println(string(data))
        return

    case "yaml", "yml":
        data, _ := yaml.Marshal(items)
        fmt.Println(string(data))
        return

    case "color":
        if colorFn != nil {
            for _, it := range items {
                colorFn(it)
            }
            return
        }
        printTable(items, headers)
        return

    default: // "table"
        printTable(items, headers)
    }
}

func printTable(items []map[string]interface{}, headers []string) {
    if len(items) == 0 {
        return
    }

    if len(headers) == 0 {
        for k := range items[0] {
            headers = append(headers, k)
        }
    }

    w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

    fmt.Fprintln(w, strings.Join(headers, "\t"))

    for _, it := range items {
        row := make([]string, len(headers))
        for i, h := range headers {
            row[i] = toString(it[h])
        }
        fmt.Fprintln(w, strings.Join(row, "\t"))
    }
    w.Flush()
}

func toString(v interface{}) string {
    switch t := v.(type) {
    case nil:
        return ""
    case string:
        return t
    default:
        return fmt.Sprintf("%v", t)
    }
}
