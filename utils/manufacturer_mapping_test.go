package utils

import "testing"

func TestTranslateManufacturerNameToEnglish(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "translates known hebrew brand",
			input: "טויוטה",
			want:  "Toyota",
		},
		{
			name:  "normalizes alias before translation",
			input: "דאציה",
			want:  "Dacia",
		},
		{
			name:  "keeps latin names with diacritics",
			input: "Peugéot",
			want:  "Peugéot",
		},
		{
			name:  "returns empty for unmapped non latin",
			input: "יצרן-לא-מוכר",
			want:  "",
		},
		{
			name:  "correct leapmotor spelling",
			input: "ליפמוטור",
			want:  "Leapmotor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TranslateManufacturerNameToEnglish(tt.input)
			if got != tt.want {
				t.Fatalf("TranslateManufacturerNameToEnglish(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestConvertManufacturerToEnglish(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "slugifies mapped hebrew manufacturer",
			input: "טויוטה",
			want:  "toyota",
		},
		{
			name:  "uses legacy slug override",
			input: "מרצדס",
			want:  "mercedes-benz",
		},
		{
			name:  "slugifies latin manufacturer with symbols",
			input: "Lynk & Co",
			want:  "lynk-and-co",
		},
		{
			name:  "returns empty for unmapped hebrew manufacturer",
			input: "יצרן-לא-מוכר",
			want:  "",
		},
		{
			name:  "returns empty for symbol only input",
			input: "///",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertManufacturerToEnglish(tt.input)
			if got != tt.want {
				t.Fatalf("ConvertManufacturerToEnglish(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
