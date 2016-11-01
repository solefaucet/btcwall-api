package middlewares

import "testing"

func TestParseIDCombination(t *testing.T) {
	f := func(c string) error {
		_, _, _, err := parseIDCombination(c)
		return err
	}
	data := map[string]bool{
		"1_2":   false,
		"s_1_2": false,
		"1_s_2": false,
		"1_2_s": false,
		"1_2_3": true,
	}
	for c, e := range data {
		if (f(c) == nil) != e {
			t.Errorf("parse %v nil error should be %v but get %v", c, e, !e)
		}
	}
}
