package safety

import "testing"

func TestValidateReadOnly_Allows(t *testing.T) {
	ok := []string{
		"SELECT TOP 10 MB001 FROM INVMB",
		"  select * from COPTC where TC003 = '20220101'  ",
		"WITH c AS (SELECT 1 AS x) SELECT * FROM c",
		"SELECT * FROM COPTC;",                      // single trailing semicolon allowed
		"SELECT MB001 -- a comment\nFROM INVMB",     // comment stripped
	}
	for _, q := range ok {
		if err := ValidateReadOnly(q); err != nil {
			t.Errorf("expected allow, got error for %q: %v", q, err)
		}
	}
}

func TestValidateReadOnly_Denies(t *testing.T) {
	bad := []string{
		"",
		"DELETE FROM COPTC",
		"UPDATE INVMB SET MB002='x'",
		"INSERT INTO COPTC VALUES (1)",
		"DROP TABLE COPTC",
		"SELECT * INTO newtbl FROM COPTC",
		"SELECT 1; DROP TABLE COPTC",                // multi-statement
		"EXEC sp_who",
		"TRUNCATE TABLE COPTC",
		"SELECT 1 /* still */; DELETE FROM x",        // comment hides nothing
	}
	for _, q := range bad {
		if err := ValidateReadOnly(q); err == nil {
			t.Errorf("expected deny, got nil for %q", q)
		}
	}
}
