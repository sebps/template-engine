package utils

import "testing"

func TestGenerateWrapperRegexp(t *testing.T) {
	type testGenerateWrapperRegexp struct {
		wrapperStart     string
		wrapperEnd       string
		wrapperGroup     string
		wrapperLineBreak bool
	}

	tests := []struct {
		args testGenerateWrapperRegexp
		want string
	}{
		{
			args: testGenerateWrapperRegexp{
				wrapperStart:     "(",
				wrapperEnd:       ")",
				wrapperGroup:     "variable",
				wrapperLineBreak: false,
			},
			want: "\\((?P<variable>([^\\)])*)\\)",
		},
		{
			args: testGenerateWrapperRegexp{
				wrapperStart:     "(",
				wrapperEnd:       ")",
				wrapperGroup:     "variable",
				wrapperLineBreak: true,
			},
			want: "\\s*\\(\n\\s*(?P<variable>([^\\)])*)\n\\s*\\)",
		},
		{
			args: testGenerateWrapperRegexp{
				wrapperStart:     "[",
				wrapperEnd:       "]",
				wrapperGroup:     "content",
				wrapperLineBreak: false,
			},
			want: "\\[(?P<content>([^\\]])*)\\]",
		},
		{
			args: testGenerateWrapperRegexp{
				wrapperStart:     "[",
				wrapperEnd:       "]",
				wrapperGroup:     "content",
				wrapperLineBreak: true,
			},
			want: "\\s*\\[\n\\s*(?P<content>([^\\]])*)\n\\s*\\]",
		},
		{
			args: testGenerateWrapperRegexp{
				wrapperStart:     "[[",
				wrapperEnd:       "]]",
				wrapperGroup:     "content",
				wrapperLineBreak: false,
			},
			want: "\\[\\[(?P<content>([^\\]]|\\][^\\]])*)\\]\\]",
		},
		{
			args: testGenerateWrapperRegexp{
				wrapperStart:     "[[",
				wrapperEnd:       "]]",
				wrapperGroup:     "content",
				wrapperLineBreak: true,
			},
			want: "\\s*\\[\\[\n\\s*(?P<content>([^\\]]|\\][^\\]])*)\n\\s*\\]\\]",
		},
		{
			args: testGenerateWrapperRegexp{
				wrapperStart:     "#{",
				wrapperEnd:       "}#",
				wrapperGroup:     "content",
				wrapperLineBreak: false,
			},
			want: "#\\{(?P<content>([^\\}]|\\}[^#])*)\\}#",
		},
		{
			args: testGenerateWrapperRegexp{
				wrapperStart:     "#{",
				wrapperEnd:       "}#",
				wrapperGroup:     "content",
				wrapperLineBreak: true,
			},
			want: "\\s*#\\{\n\\s*(?P<content>([^\\}]|\\}[^#])*)\n\\s*\\}#",
		},
	}
	for i, tc := range tests {
		out := GenerateWrapperRegexp(tc.args.wrapperStart, tc.args.wrapperEnd, tc.args.wrapperGroup, tc.args.wrapperLineBreak)
		if out != tc.want {
			t.Errorf("test #%d failed expected result \n want : %s \n have : %s", i+1, tc.want, out)
		}
	}
}
