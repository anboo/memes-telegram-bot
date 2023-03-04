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

func TestHandler_Support(t *testing.T) {
	tests := []struct {
		name string
		args handler.BotRequest
		want bool
	}{
		{
			name: "empty_request",
			args: handler.BotRequest{},
			want: false,
		},
		{
			name: "empty_callback_query",
			args: handler.BotRequest{
				Update: tgbotapi.Update{},
			},
			want: false,
		},
		{
			name: "empty_callback_data",
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{},
				},
			},
			want: false,
		},
		{
			name: "up_vote",
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: UpPrefix + "_792356d3-a7e1-4f67-a153-54a11849afc3",
					},
				},
			},
			want: true,
		},
		{
			name: "down_vote",
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: DownPrefix + "_792356d3-a7e1-4f67-a153-54a11849afc3",
					},
				},
			},
			want: true,
		},
		{
			name: "sos_vote",
			args: handler.BotRequest{
				Update: tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: SosPrefix + "_792356d3-a7e1-4f67-a153-54a11849afc3",
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

			l := zerolog.Nop()
			h := NewHandler(NewMockTelegramAPI(ctrl), NewMockMemRepository(ctrl), NewMockVoteRepository(ctrl), &l)
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
