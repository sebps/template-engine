package filtering

import (
	"fmt"
	"testing"
)

func TestIsJsonPathCompliant(t *testing.T) {
	type testIsJsonPathCompliant struct {
		input string
	}

	tests := []struct {
		args testIsJsonPathCompliant
		want bool
	}{
		{
			args: testIsJsonPathCompliant{
				input: "$[?(@.sku in (record2,record3)]",
			},
			want: true,
		},
	}

	for i, tc := range tests {
		out := IsJsonPathCompliant(tc.args.input)
		if out != tc.want {
			t.Error(fmt.Sprintf("test #%d failed expected result \n want : %#v \n have : %#v", i+1, tc.want, out))
		}
	}

}
