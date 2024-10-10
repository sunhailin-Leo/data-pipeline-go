package utils

import (
	"fmt"
	"strings"
)

// GenerateInsertSQL generate insert into sql with different symbol like &, :, or ?
func GenerateInsertSQL(tableName, placeHolderSymbol string, columnLen int) string {
	// Generate placeholders &1, &2, ... , &n or :1, :2, ... , :n or ? , ? , ... , ?
	placeholders := make([]string, columnLen)
	for i := 0; i < columnLen; i++ {
		if placeHolderSymbol == "" {
			placeholders[i] = "?"
		} else {
			placeholders[i] = fmt.Sprintf("%s%d", placeHolderSymbol, i+1)
		}
	}
	// Splicing column names and placeholders
	placeholdersStr := strings.Join(placeholders, ", ")
	// Building the final SQL insert statement
	return fmt.Sprintf("INSERT INTO %s VALUES (%s)", tableName, placeholdersStr)
}
