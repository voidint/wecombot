package wecombot

import "testing"

func TestExtractKey(t *testing.T) {
	type args struct {
		webhookURL string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
	}{
		{
			name:    "URL中包含key参数",
			args:    args{webhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=aaa-bbb-ccc-ddd-eee-fff"},
			wantKey: "aaa-bbb-ccc-ddd-eee-fff",
		},
		{
			name:    "URL中包含大写的KEY参数",
			args:    args{webhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?KEY=aaa-bbb-ccc-ddd-eee-fff"},
			wantKey: "",
		},
		{
			name:    "URL中key参数为空",
			args:    args{webhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="},
			wantKey: "",
		},
		{
			name:    "URL中不包含key参数",
			args:    args{webhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?hello=world"},
			wantKey: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotKey := ExtractKey(tt.args.webhookURL); gotKey != tt.wantKey {
				t.Errorf("ExtractKey() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
