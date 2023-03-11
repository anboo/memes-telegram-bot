package config_source

import (
	"context"
	"testing"

	"memes-bot/handler"
	"memes-bot/storage/user"
	"memes-bot/storage/user_source"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
)

func TestHandler_Support(t *testing.T) {
	tests := []struct {
		name string
		args handler.BotRequest
		want bool
	}{
		{
			name: "show_buttons",
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					Message: &tgbotapi.Message{
						Text: "/settings",
						Entities: []tgbotapi.MessageEntity{
							{Type: "bot_command", Length: 9},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "success_callback",
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: callbackKeyPrefix + callbackDisableSourcePrefix,
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := NewHandler(NewMockTelegramAPI(ctrl), NewMockUserSourceRepository(ctrl), "")
			res := h.Support(&tt.args)

			if res != tt.want {
				t.Fatalf("expected want %v got %v", tt.want, res)
			}
		})
	}
}

func TestHandler_Handle(t *testing.T) {
	type fields struct {
		bot                  func(ctrl *gomock.Controller) TelegramAPI
		userSourceRepository func(ctrl *gomock.Controller) UserSourceRepository
	}

	tests := []struct {
		name    string
		fields  fields
		args    handler.BotRequest
		wantErr bool
	}{
		{
			name: "success_user_sources_list",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					msg := tgbotapi.NewMessage(1, "Выберите источники откуда получать мемы:")

					m := NewMockTelegramAPI(ctrl)
					m.EXPECT().Send(msg).Return(tgbotapi.Message{
						MessageID: 23,
					}, nil)
					m.EXPECT().Send(tgbotapi.NewEditMessageReplyMarkup(1, 23, tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("✅ borsh", callbackKeyPrefix+"23_"+callbackEnableSourcePrefix+"borsh"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("❌ memes", callbackKeyPrefix+"23_"+callbackDisableSourcePrefix+"memes"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("❌ golang", callbackKeyPrefix+"23_"+callbackDisableSourcePrefix+"golang"),
						),
					))).Return(tgbotapi.Message{}, nil)
					return m
				},
				userSourceRepository: func(ctrl *gomock.Controller) UserSourceRepository {
					m := NewMockUserSourceRepository(ctrl)
					m.EXPECT().ByUserId(gomock.Any(), "9671f363-d29b-4c3b-8d39-5454cdc09c9b").Return([]user_source.UserSource{
						{Source: "borsh", Enabled: false},
					}, nil)
					return m
				},
			},
			args: handler.BotRequest{
				FromID: 1,
				User: user.User{
					ID: "9671f363-d29b-4c3b-8d39-5454cdc09c9b",
				},
				Update: tgbotapi.Update{
					Message: &tgbotapi.Message{
						Text: "/settings",
						Entities: []tgbotapi.MessageEntity{
							{Type: "bot_command", Length: 9},
						},
					},
				},
			},
		},
		{
			name: "success_disable_button",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					m := NewMockTelegramAPI(ctrl)
					m.EXPECT().Request(tgbotapi.NewCallback("1", "Сохранено"))
					m.EXPECT().Send(tgbotapi.NewEditMessageReplyMarkup(1, 23, tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("✅ borsh", callbackKeyPrefix+"23_"+callbackEnableSourcePrefix+"borsh"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("❌ memes", callbackKeyPrefix+"23_"+callbackDisableSourcePrefix+"memes"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("❌ golang", callbackKeyPrefix+"23_"+callbackDisableSourcePrefix+"golang"),
						),
					))).Return(tgbotapi.Message{}, nil)
					return m
				},
				userSourceRepository: func(ctrl *gomock.Controller) UserSourceRepository {
					m := NewMockUserSourceRepository(ctrl)
					m.EXPECT().Create(gomock.Any(), &user_source.UserSource{
						UserID:  "9671f363-d29b-4c3b-8d39-5454cdc09c9b",
						Source:  "borsh",
						Enabled: false,
					}).Return(nil)
					m.EXPECT().ByUserId(gomock.Any(), "9671f363-d29b-4c3b-8d39-5454cdc09c9b").Return([]user_source.UserSource{
						{Source: "borsh", Enabled: false},
					}, nil)
					return m
				},
			},
			args: handler.BotRequest{
				FromID: 1,
				User: user.User{
					ID: "9671f363-d29b-4c3b-8d39-5454cdc09c9b",
				},
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						ID:   "1",
						Data: callbackKeyPrefix + "23_" + callbackDisableSourcePrefix + "borsh",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := NewHandler(tt.fields.bot(ctrl), tt.fields.userSourceRepository(ctrl), "borsh,memes,golang")
			err := h.Handle(context.Background(), &tt.args)

			if err != nil != tt.wantErr {
				t.Fatalf("expected err %v got %v", tt.wantErr, err)
			}
		})
	}
}
