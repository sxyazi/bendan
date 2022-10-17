package yes

import "testing"

func TestIsTokenize(t *testing.T) {
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
		{"33是不是", "sub=33"},
		{"是不是傻", "sub=, obj=傻"},
		{"X是不是傻", "sub=X, obj=傻"},
		{"X是不是在自慰", "sub=X, obj=在自慰"},
		{"X是幼女嘛，还是  我编不下去了", "sub=X, obj=幼女, ind=我编不下去了"},

		{input: "你是那个？", want: "sub=你, obj=那个"},
		{input: "你是用XX还是啥嘛", want: "sub=你, obj=用XX"},
		{"你是谁是谁，还是你是你", "sub=你, obj=谁是谁, ind=你是你"},
		{input: "但是我是不是傻", want: "sub=我, obj=傻"},

		{input: "日本是哪啊？", want: ""},
		{input: "我的评价是：别尬黑好吧", want: ""},
		{input: "但是AA还是个BB", want: ""},
		{input: "想XX但是不敢！", want: ""},
	}

	for _, d := range data {
		if got := IsTokenize(d.input); got.String() != d.want {
			t.Errorf("IsTokenize(%q) = %q, want %q", d.input, got, d.want)
		}
	}
}
