package choose_sex

import (
	"context"
	"testing"

	"memes-bot/handler"
	"memes-bot/storage/user"

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
			name: "success",
			args: handler.BotRequest{
				User: user.User{
					Sex: "",
				},
			},
			want: true,
		},
		{
			name: "false_sex_already_filled",
			args: handler.BotRequest{
				User: user.User{
					Sex: SexGirl,
				},
			},
			want: false,
		},
		{
			name: "true_filled_sex_chosen_sex_18",
			args: handler.BotRequest{
				User: user.User{
					Sex: SexGirl,
				},
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: ChooseSexPrefix + "_" + SexFish,
					},
				},
			},
			want: true,
		},
		{
			name: "true_filled_sex_chosen_sex_18_25",
			args: handler.BotRequest{
				User: user.User{
					Sex: SexGirl,
				},
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: ChooseSexPrefix + "_" + SexMen,
					},
				},
			},
			want: true,
		},
		{
			name: "true_filled_sex_chosen_sex_25",
			args: handler.BotRequest{
				User: user.User{
					Sex: SexGirl,
				},
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: ChooseSexPrefix + "_" + SexGirl,
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

			h := NewHandler(NewMockTelegramAPI(ctrl), NewMockUserRepository(ctrl))
			res := h.Support(&tt.args)

			if res != tt.want {
				t.Fatalf("expected want %v got %v", tt.want, res)
			}
		})
	}
}

func TestHandler_Handle(t *testing.T) {
	type fields struct {
		bot            func(ctrl *gomock.Controller) TelegramAPI
		userRepository func(ctrl *gomock.Controller) UserRepository
	}

	tests := []struct {
		name    string
		fields  fields
		args    handler.BotRequest
		wantErr bool
	}{
		{
			name: "success_need_choose_sex",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					msg := tgbotapi.NewMessage(1, "Выберите ваш пол")
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("👱🏼‍♂️Мужчина", ChooseSexPrefix+SexMen),
							tgbotapi.NewInlineKeyboardButtonData("👱🏼‍♀️Девушка", ChooseSexPrefix+SexGirl),
							tgbotapi.NewInlineKeyboardButtonData("🐟 My son?? Where is my son??", ChooseSexPrefix+SexFish),
						),
					)

					m := NewMockTelegramAPI(ctrl)
					m.EXPECT().Send(msg)
					return m
				},
				userRepository: func(ctrl *gomock.Controller) UserRepository {
					return NewMockUserRepository(ctrl)
				},
			},
			args: handler.BotRequest{
				FromID: 1,
			},
		},
		{
			name: "success_choose_sex_1",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					m := NewMockTelegramAPI(ctrl)
					m.EXPECT().Request(tgbotapi.NewCallback("1", "Спасибо"))
					return m
				},
				userRepository: func(ctrl *gomock.Controller) UserRepository {
					m := NewMockUserRepository(ctrl)
					m.EXPECT().Upsert(gomock.Any(), user.User{
						ID:  "bbb0a5c2-c4c9-4548-9303-acdd1234af96",
						Sex: SexMen,
					})
					return m
				},
			},
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						ID:   "1",
						Data: ChooseSexPrefix + SexMen,
					},
				},
				User: user.User{
					ID: "bbb0a5c2-c4c9-4548-9303-acdd1234af96",
				},
			},
		},
		{
			name: "success_choose_sex_2",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					m := NewMockTelegramAPI(ctrl)
					m.EXPECT().Request(tgbotapi.NewCallback("1", "Спасибо"))
					return m
				},
				userRepository: func(ctrl *gomock.Controller) UserRepository {
					m := NewMockUserRepository(ctrl)
					m.EXPECT().Upsert(gomock.Any(), user.User{
						ID:  "bbb0a5c2-c4c9-4548-9303-acdd1234af96",
						Sex: SexGirl,
					})
					return m
				},
			},
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						ID:   "1",
						Data: ChooseSexPrefix + SexGirl,
					},
				},
				User: user.User{
					ID: "bbb0a5c2-c4c9-4548-9303-acdd1234af96",
				},
			},
		},
		{
			name: "success_choose_sex_1",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					m := NewMockTelegramAPI(ctrl)
					m.EXPECT().Request(tgbotapi.NewCallback("1", "Спасибо"))
					return m
				},
				userRepository: func(ctrl *gomock.Controller) UserRepository {
					m := NewMockUserRepository(ctrl)
					m.EXPECT().Upsert(gomock.Any(), user.User{
						ID:  "bbb0a5c2-c4c9-4548-9303-acdd1234af96",
						Sex: SexFish,
					})
					return m
				},
			},
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						ID:   "1",
						Data: ChooseSexPrefix + SexFish,
					},
				},
				User: user.User{
					ID: "bbb0a5c2-c4c9-4548-9303-acdd1234af96",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := NewHandler(tt.fields.bot(ctrl), tt.fields.userRepository(ctrl))
			err := h.Handle(context.Background(), &tt.args)

			if err != nil != tt.wantErr {
				t.Fatalf("expected err %v got %v", tt.wantErr, err)
			}
		})
	}
}
