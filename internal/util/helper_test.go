package util

import (
	"testing"
)

func TestSanitizeIdentifier_RejectsInjection(t *testing.T) {
	malicious := []string{
		"users; DROP TABLE users--",
		"users\" WHERE 1=1--",
		"users) UNION SELECT password FROM users--",
		"table name with spaces",
		"1users",
		"users--",
		"u\\sers",
		"`users`",
	}

	for _, name := range malicious {
		if err := sanitizeIdentifier(name); err == nil {
			t.Errorf("expected sanitizeIdentifier to reject %q", name)
		}
	}
}

func TestSanitizeIdentifier_AcceptsValidNames(t *testing.T) {
	valid := []string{
		"users",
		"user_profiles",
		"_internal",
		"users2",
		"a",
	}

	for _, name := range valid {
		if err := sanitizeIdentifier(name); err != nil {
			t.Errorf("expected sanitizeIdentifier to accept %q, got err: %v", name, err)
		}
	}
}

func TestSeedTable_RejectsSQLInjectionBeforeDB(t *testing.T) {
	// sanitization runs before any DB interaction, so a nil DB is sufficient
	err := SeedTable(nil, "users; DROP TABLE users--", []string{"email"}, nil)
	if err == nil {
		t.Fatal("expected SeedTable to reject malicious table name")
	}
}

func TestSeedTable_RejectsMaliciousColumn(t *testing.T) {
	err := SeedTable(nil, "users", []string{"email); DROP TABLE users; --"}, nil)
	if err == nil {
		t.Fatal("expected SeedTable to reject malicious column name")
	}
}

func TestUpdateSeederMap_RejectsInvalidFunctionName(t *testing.T) {
	// should fail validation before touching any file on disk
	if err := UpdateSeederMap("Foo; process.Getpid() // ; DROP TABLE x"); err == nil {
		t.Fatal("expected UpdateSeederMap to reject invalid function name")
	}
}

func TestToPascalCase(t *testing.T) {
	cases := map[string]string{
		"users_seeder": "UsersSeeder",
		"user":         "User",
		"a_b_c":        "ABC",
	}
	for input, want := range cases {
		if got := ToPascalCase(input); got != want {
			t.Errorf("ToPascalCase(%q) = %q, want %q", input, got, want)
		}
	}
}