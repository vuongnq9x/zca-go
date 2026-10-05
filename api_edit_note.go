package zca

import "context"

type EditNoteOptions struct {
	// New note title.
	Title string
	// Topic ID of the note to edit.
	TopicID string
	// Should the note be pinned?
	PinAct bool
}

type EditNoteResponse = NoteDetail

// EditNote edits an existing note in a group.
func (a *API) EditNote(ctx context.Context, options EditNoteOptions, groupID string) (*EditNoteResponse, error) {
	pin := 2
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
		"topicId":   options.TopicID,
		"repeat":    0,
		"imei":      a.IMEI,
		"pinAct":    pin,
	}
	return createNoteCall(ctx, a, a.svc("group_board")+"/api/board/topic/updatev2", params)
}
