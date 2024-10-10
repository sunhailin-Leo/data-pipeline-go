package utils

import "testing"

func TestGenerateSQL(t *testing.T) {
	tests := []struct {
		tableName         string
		placeHolderSymbol string
		columnLen         int
		expectedSQL       string
	}{
		{
			tableName:         "users",
			placeHolderSymbol: "&",
			columnLen:         3,
			expectedSQL:       "INSERT INTO users VALUES (&1, &2, &3)",
		},
		{
			tableName:         "orders",
			placeHolderSymbol: ":",
			columnLen:         2,
			expectedSQL:       "INSERT INTO orders VALUES (:1, :2)",
		},
		{
			tableName:         "products",
			placeHolderSymbol: "",
			columnLen:         4,
			expectedSQL:       "INSERT INTO products VALUES (?, ?, ?, ?)",
		},
		{
			tableName:         "employees",
			placeHolderSymbol: "?",
			columnLen:         0,
			expectedSQL:       "INSERT INTO employees VALUES ()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.tableName, func(t *testing.T) {
			gotSQL := GenerateInsertSQL(tt.tableName, tt.placeHolderSymbol, tt.columnLen)
			if gotSQL != tt.expectedSQL {
				t.Errorf("GenerateSQL(%q, %q, %d) = %q; want %q",
					tt.tableName, tt.placeHolderSymbol, tt.columnLen, gotSQL, tt.expectedSQL)
			}
		})
	}
}
