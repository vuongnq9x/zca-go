package zca

import (
	"context"
	"net/http"
	"sync"
)

type StickerDetailResponse = []StickerDetail

// GetStickersDetail gets sticker details by IDs; failed lookups are skipped (Promise.allSettled).
func (a *API) GetStickersDetail(ctx context.Context, stickerIDs []int64) (StickerDetailResponse, error) {
	if len(stickerIDs) == 0 {
		return nil, newError("Missing sticker id")
	}
	u := a.MakeURL(a.svc("sticker")+"/api/message/sticker/sticker_detail", nil, true)
	results := make([]*StickerDetail, len(stickerIDs))
	var wg sync.WaitGroup
	for i, id := range stickerIDs {
		wg.Go(func() {
			d, err := call[*StickerDetail](ctx, a.Session, http.MethodGet, u, map[string]any{"sid": id})
			if err == nil {
				results[i] = d
			}
		})
	}
	wg.Wait()
	stickers := StickerDetailResponse{}
	for _, d := range results {
		if d != nil {
			stickers = append(stickers, *d)
		}
	}
	return stickers, nil
}
