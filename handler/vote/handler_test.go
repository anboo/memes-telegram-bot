package vote

import (
	"context"
	"testing"

	"memes-bot/handler"

	"github.com/golang/mock/gomock"
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
			name: "success_need_choose_sex",
			fields: fields{
				bot: func(ctrl *gomock.Controller) TelegramAPI {
					m := NewMockTelegramAPI(ctrl)
					return m
				},
				memRepository: func(ctrl *gomock.Controller) MemRepository {
					return NewMockMemRepository(ctrl)
				},
				voteRepository: func(ctrl *gomock.Controller) VoteRepository {
					return NewMockVoteRepository(ctrl)
				},
			},
			args: handler.BotRequest{
				FromID: 1,
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
