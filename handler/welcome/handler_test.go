package welcome

import (
	"context"
	"testing"

	"memes-bot/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
)

func TestHandler_Handle(t *testing.T) {
	type fields struct {
		bot func(ctrl *gomock.Controller) TelegramAPI
	}

	tests := []struct {
		name    string
		fields  fields
		args    handler.BotRequest
		wantErr bool
	}{
		{
			name: "success_vote",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					m := NewMockTelegramAPI(ctrl)
					m.EXPECT().Send(tgbotapi.NewMessage(
						1,
						"Вам нужно оценить минимум 20 мемов, прежде чем рекомендации станут более менее релевантными.",
					))
					return m
				},
			},
			args: handler.BotRequest{
				FromID: 1,
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						ID:   "id",
						Data: "up_a195e90a-25c5-4b49-8251-0109c4c8904c",
					},
				},
				IsNewUser: true,
			},
		},
		{
			name: "is_not_new_user",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					m := NewMockTelegramAPI(ctrl)
					return m
				},
			},
			args: handler.BotRequest{
				FromID: 1,
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						ID:   "id",
						Data: "up_a195e90a-25c5-4b49-8251-0109c4c8904c",
					},
				},
				IsNewUser: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := NewHandler(tt.fields.bot(ctrl))
			err := h.Handle(context.Background(), &tt.args)

			if err != nil != tt.wantErr {
				t.Fatalf("expected err %v got %v", tt.wantErr, err)
			}
		})
	}
}
