package uregexp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchPunctuation(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "chinese",
			args: args{
				text: "要创建一个能够过滤所有语言标点符号的正则表达式，可以使用Unicode属性来匹配标点符号。不同语言的标点符号涵盖了广泛的Unicode代码点范围，因此需要使用Unicode属性来准确匹配它们。以下是一个高级写法的正则表达式，可以过滤掉所有语言的标点符号28.96：",
			},
			want: []string{"，", "。", "，", "。", "，", "："},
		},
		{
			name: "english",
			args: args{
				text: "She had been shopping with her Mom in Wal-Mart. She must have been 6 years old, this beautiful brown haired, freckle-faced image of innocence. It was pouring outside. The kind of rain that gushes over the top of rain gutters, so much in a hurry to hit the Earth, it has no time to flow down the spout!",
			},
			want: []string{"-", ".", ",", ",", "-", ".", ".", ",", ",", "!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Punctuations(tt.args.text)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSentenceEndPunctuation(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "en have",
			args: args{
				text: "This is a sample text with some incomplete sentences, 分句是一个多语言文本处理的复杂任务， Esto es un ejemplo de texto. Contiene múltiples oraciones. これはサンプルテキストです。これにはいくつかの文が含まれています.",
			},
			want: true,
		},
		{
			name: "en lack",
			args: args{
				text: "This is a sample text with some incomplete sentences, 分句是一个多语言文本处理的复杂任务， Esto es un ejemplo de texto. Contiene múltiples oraciones",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := SentenceEndPunctuation(tt.args.text)
			assert.Equalf(t, tt.want, ok, "SentenceEndPunctuation(%v)", tt.args.text)
		})
	}
}
