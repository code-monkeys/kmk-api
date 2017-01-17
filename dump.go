// http://stackoverflow.com/questions/19991541/dumping-mysql-tables-to-json-with-golang

import (
    "encoding/json"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)


func dumpTable(w io.Writer, table) {
    rows, err := Query(db, fmt.Sprintf("SELECT * FROM %s", table))
    checkError(err)
    columns, err := rows.Columns()
    checkError(err)

    scanArgs := make([]interface{}, len(columns))
    values   := make([]interface{}, len(columns))

    for i := range values {
        scanArgs[i] = &values[i]
    }

    for rows.Next() {
        err = rows.Scan(scanArgs...)
        checkError(err)

        record := make(map[string]interface{})

        for i, col := range values {
            if col != nil {
                fmt.Printf("\n%s: type= %s\n", columns[i], reflect.TypeOf(col))

                switch t := col.(type) {
                default:
                    fmt.Printf("Unexpected type %T\n", t)
                case bool:
                    fmt.Printf("bool\n")
                    record[columns[i]] = col.(bool)
                case int:
                    fmt.Printf("int\n")
                    record[columns[i]] = col.(int)
                case int64:
                    fmt.Printf("int64\n")
                    record[columns[i]] = col.(int64)
                case float64:
                    fmt.Printf("float64\n")
                    record[columns[i]] = col.(float64)
                case string:
                    fmt.Printf("string\n")
                    record[columns[i]] = col.(string)
                case []byte:   // -- all cases go HERE!
                    fmt.Printf("[]byte\n")
                    record[columns[i]] = string(col.([]byte))
                case time.Time:
                    // record[columns[i]] = col.(string)
                }
            }
        }

        s, _ := json.Marshal(record)
        w.Write(s)
        io.WriteString(w, "\n")
    }
}
