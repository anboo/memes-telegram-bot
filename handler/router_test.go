package handler

import (
	"context"
	"testing"

	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

func TestRouter_Handle(t *testing.T) {
	type fields struct {
		handlers func(ctrl *gomock.Controller) []Handler
	}

	botReqA := BotRequest{
		FromID:          1,
		Update:          tgbotapi.Update{},
		User:            user.User{},
		IsNewUser:       false,
		StopPropagation: false,
	}

	tests := []struct {
		name    string
		fields  fields
		args    BotRequest
		wantErr bool
	}{
		{
			name: "found_handler_handle_success",
			fields: fields{
				handlers: func(ctrl *gomock.Controller) []Handler {
					m := NewMockHandler(ctrl)
					m.EXPECT().Support(gomock.Any()).Return(true)
					m.EXPECT().Handle(gomock.Any(), &botReqA)
					return []Handler{m}
				},
			},
			args: botReqA,
		},
		{
			name: "not_found_handler_with_non_empty_list",
			fields: fields{
				handlers: func(ctrl *gomock.Controller) []Handler {
					m := NewMockHandler(ctrl)
					m.EXPECT().Support(gomock.Any()).Return(false)
					return []Handler{m}
				},
			},
			args:    botReqA,
			wantErr: true,
		},
		{
			name: "not_found_handler_with_empty_list",
			fields: fields{
				handlers: func(ctrl *gomock.Controller) []Handler {
					return []Handler{}
				},
			},
			args:    botReqA,
			wantErr: true,
		},
		{
			name: "found_handler_handle_fail",
			fields: fields{
				handlers: func(ctrl *gomock.Controller) []Handler {
					m := NewMockHandler(ctrl)
					m.EXPECT().Support(gomock.Any()).Return(true)
					m.EXPECT().Handle(gomock.Any(), &botReqA).Return(errors.New("some_error"))
					return []Handler{m}
				},
			},
			args:    botReqA,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := NewRouter(tt.fields.handlers(ctrl)...)
			err := h.Handle(context.Background(), tt.args)
			if err != nil != tt.wantErr {
				t.Fatalf("expected err %v got %v", tt.wantErr, err)
			}
		})
	}
}
