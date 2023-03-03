package vote

import (
	"context"
	"testing"

	"memes-bot/handler"
	"memes-bot/storage/mem"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func TestHandler_Handle(t *testing.T) {
	type fields struct {
		bot            func(ctrl *gomock.Controller) TelegramAPI
		memRepository  func(ctrl *gomock.Controller) MemRepository
		voteRepository func(ctrl *gomock.Controller) VoteRepository
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
					m.EXPECT().Request(tgbotapi.NewCallback("id", "Спасибо за вашу оценку!"))
					return m
				},
				memRepository: func(ctrl *gomock.Controller) MemRepository {
					m := NewMockMemRepository(ctrl)
					m.EXPECT().Find(gomock.Any(), "a195e90a-25c5-4b49-8251-0109c4c8904c").Return(mem.Mem{ID: "a195e90a-25c5-4b49-8251-0109c4c8904c"}, nil)
					m.EXPECT().UpdateRating(gomock.Any(), "a195e90a-25c5-4b49-8251-0109c4c8904c", 1)
					return m
				},
				voteRepository: func(ctrl *gomock.Controller) VoteRepository {
					m := NewMockVoteRepository(ctrl)
					m.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
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
			},
		},
		{
			name: "error_mem_not_found",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					m := NewMockTelegramAPI(ctrl)
					m.EXPECT().Request(tgbotapi.NewCallback("id", "Мем не найден"))
					return m
				},
				memRepository: func(ctrl *gomock.Controller) MemRepository {
					m := NewMockMemRepository(ctrl)
					m.EXPECT().Find(gomock.Any(), "a195e90a-25c5-4b49-8251-0109c4c8904c").Return(mem.Mem{}, errors.New("not found"))
					return m
				},
				voteRepository: func(ctrl *gomock.Controller) VoteRepository {
					m := NewMockVoteRepository(ctrl)
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			l := zerolog.Nop()
			h := NewHandler(tt.fields.bot(ctrl), tt.fields.memRepository(ctrl), tt.fields.voteRepository(ctrl), &l)
			err := h.Handle(context.Background(), &tt.args)

			if err != nil != tt.wantErr {
				t.Fatalf("expected err %v got %v", tt.wantErr, err)
			}
		})
	}
}
