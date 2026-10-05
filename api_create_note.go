package zca

import (
	"context"
	"encoding/json"
	"net/http"
)

type CreateNoteOptions struct {
	Title  string
	PinAct bool
}

type CreateNoteResponse = NoteDetail

// CreateNote creates a note in a group.
func (a *API) CreateNote(ctx context.Context, options CreateNoteOptions, groupID string) (*CreateNoteResponse, error) {
	pin := 0
	if options.PinAct {
		pin = 1
	}
	params := map[string]any{
		"grid":      groupID,
		"type":      0,
		"color":     -16777216,
		"emoji":     "",
		"startTime": -1,
		"duration":  -1,
		"params":    mustJSON(map[string]string{"title": options.Title}),
		"repeat":    0,
		"src":       1,
		"imei":      a.IMEI,
		"pinAct":    pin,
	}
	return createNoteCall(ctx, a, a.svc("group_board")+"/api/board/topic/createv2", params)
}

// createNoteCall posts params and decodes a NoteDetail whose params may be a JSON string.
func createNoteCall(ctx context.Context, a *API, base string, params map[string]any) (*NoteDetail, error) {
	raw, err := call[json.RawMessage](ctx, a.Session, http.MethodPost, a.MakeURL(base, nil, true), params)
	if err != nil {
		return nil, err
	}
	var v struct {
		NoteDetail
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	v.NoteDetail.Params, err = decodeData[NoteDetailParams](addReactionUnquote(v.Params))
	if err != nil {
		return nil, err
	}
	return &v.NoteDetail, nil
}
