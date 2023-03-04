package choose_age

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
					Age: 0,
				},
			},
			want: true,
		},
		{
			name: "false_age_already_filled",
			args: handler.BotRequest{
				User: user.User{
					Age: 1,
				},
			},
			want: false,
		},
		{
			name: "true_filled_age_chosen_age_18",
			args: handler.BotRequest{
				User: user.User{
					Age: 1,
				},
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: ChooseAgePrefix + "_" + AgeLessThen18,
					},
				},
			},
			want: true,
		},
		{
			name: "true_filled_age_chosen_age_18_25",
			args: handler.BotRequest{
				User: user.User{
					Age: 1,
				},
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: ChooseAgePrefix + "_" + AgeBetween1825,
					},
				},
			},
			want: true,
		},
		{
			name: "true_filled_age_chosen_age_25",
			args: handler.BotRequest{
				User: user.User{
					Age: 1,
				},
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: ChooseAgePrefix + "_" + AgeGreatThen25,
					},
				},
			},
			want: true,
		},
		{
			name: "true_filled_age_chosen_age_25",
			args: handler.BotRequest{
				User: user.User{
					Age: 1,
				},
			},
			want: false,
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
			name: "success_need_choose_age",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					msg := tgbotapi.NewMessage(1, "Выберите ваш возраст")
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("До 18", ChooseAgePrefix+AgeLessThen18),
							tgbotapi.NewInlineKeyboardButtonData("От 18 до 25", ChooseAgePrefix+AgeBetween1825),
							tgbotapi.NewInlineKeyboardButtonData("Больше 25", ChooseAgePrefix+AgeGreatThen25),
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
			name: "success_choose_age_1",
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
						Age: 1,
					})
					return m
				},
			},
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						ID:   "1",
						Data: ChooseAgePrefix + "_" + AgeLessThen18,
					},
				},
				User: user.User{
					ID:  "bbb0a5c2-c4c9-4548-9303-acdd1234af96",
					Age: 0,
				},
			},
		},
		{
			name: "success_choose_age_2",
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
						Age: 2,
					})
					return m
				},
			},
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						ID:   "1",
						Data: ChooseAgePrefix + "_" + AgeBetween1825,
					},
				},
				User: user.User{
					ID:  "bbb0a5c2-c4c9-4548-9303-acdd1234af96",
					Age: 0,
				},
			},
		},
		{
			name: "success_choose_age_1",
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
						Age: 3,
					})
					return m
				},
			},
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						ID:   "1",
						Data: ChooseAgePrefix + "_" + AgeGreatThen25,
					},
				},
				User: user.User{
					ID:  "bbb0a5c2-c4c9-4548-9303-acdd1234af96",
					Age: 0,
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
