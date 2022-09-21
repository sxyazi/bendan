package yes

import "testing"

func TestTokenize(t *testing.T) {
	var data = []struct {
		input string
		want  string
	}{
		{"你是谁", ""},
		{"你是谁吗", "sub=你, obj=谁"},
		{"你是谁吗?", "sub=你, obj=谁"},
		{"我是笨蛋？", "sub=我, obj=笨蛋"},
		{"我是笨蛋！", "sub=我, obj=笨蛋"},

		{"是X嘛", "sub=, obj=X"},
		{"你是AI嘛？", "sub=你, obj=AI"},
		{"是要下雨吗", "sub=, obj=要下雨"},
		{"是不是要下雨", "sub=, obj=要下雨"},
		{"是生还是死", "sub=, obj=生, ind=死"},
		{"是猫娘，还是狐娘呀啊？", "sub=, obj=猫娘, ind=狐娘"},
		{"33是不是", "sub=33, obj="},
		{"是不是傻", "sub=, obj=傻"},
		{"X是不是傻", "sub=X, obj=傻"},
		{"X是不是在自慰", "sub=X, obj=在自慰"},
		{"X是幼女嘛，还是  我编不下去了", "sub=X, obj=幼女, ind=我编不下去了"},

		{"你是谁是谁，还是你是你", "sub=你, obj=谁是谁, ind=你是你"},
	}

	for _, d := range data {
		if got := Tokenize(d.input); got.String() != d.want {
			t.Errorf("Tokenize(%q) = %q, want %q", d.input, got, d.want)
		}
	}
}
