package send_mem

import (
	"context"
	"testing"

	"memes-bot/handler"
	"memes-bot/storage/mem"
	"memes-bot/storage/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestHandler_Support(t *testing.T) {
	tests := []struct {
		name string
		args handler.BotRequest
		want bool
	}{
		{
			name: "success_send_mem",
			args: handler.BotRequest{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			l := zerolog.Nop()
			h := NewHandler(NewMockTelegramAPI(ctrl), NewMockMemRepository(ctrl), &l)
			res := h.Support(&tt.args)

			if res != tt.want {
				t.Fatalf("expected want %v got %v", tt.want, res)
			}
		})
	}
}

func TestHandler_Handle(t *testing.T) {
	type fields struct {
		bot           func(ctrl *gomock.Controller) TelegramAPI
		memRepository func(ctrl *gomock.Controller) MemRepository
	}

	tests := []struct {
		name    string
		fields  fields
		args    handler.BotRequest
		wantErr bool
	}{
		{
			name: "success_need_send_mem",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					p := tgbotapi.NewPhoto(1, tgbotapi.FileURL("https://vk.com/img_c8145378-5f2e-4729-87f6-e5ee443a4e3c.png"))
					p.Caption = "Caption"
					p.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("👍", "up_c8145378-5f2e-4729-87f6-e5ee443a4e3c"),
							tgbotapi.NewInlineKeyboardButtonData("👎", "down_c8145378-5f2e-4729-87f6-e5ee443a4e3c"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("🆘", "sos_c8145378-5f2e-4729-87f6-e5ee443a4e3c"),
						),
					)

					m := NewMockTelegramAPI(ctrl)
					m.EXPECT().Send(p).Return(tgbotapi.Message{}, nil)
					return m
				},
				memRepository: func(ctrl *gomock.Controller) MemRepository {
					m := NewMockMemRepository(ctrl)
					m.EXPECT().FindRelevantMemForUser(gomock.Any(), user.User{
						ID: "260a18f0-71d4-423e-b66d-fd0c437cf570",
					}).Return(mem.Mem{
						ID:   "c8145378-5f2e-4729-87f6-e5ee443a4e3c",
						Text: "Caption",
						Img:  "https://vk.com/img_c8145378-5f2e-4729-87f6-e5ee443a4e3c.png",
					}, nil)

					m.EXPECT().ReserveNewMem(gomock.Any(), user.User{
						ID: "260a18f0-71d4-423e-b66d-fd0c437cf570",
					}, mem.Mem{
						ID:   "c8145378-5f2e-4729-87f6-e5ee443a4e3c",
						Text: "Caption",
						Img:  "https://vk.com/img_c8145378-5f2e-4729-87f6-e5ee443a4e3c.png",
					})
					return m
				},
			},
			args: handler.BotRequest{
				FromID: 1,
				User: user.User{
					ID: "260a18f0-71d4-423e-b66d-fd0c437cf570",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			l := zerolog.Nop()
			h := NewHandler(tt.fields.bot(ctrl), tt.fields.memRepository(ctrl), &l)
			err := h.Handle(context.Background(), &tt.args)

			if err != nil != tt.wantErr {
				t.Fatalf("expected err %v got %v", tt.wantErr, err)
			}
		})
	}
}
