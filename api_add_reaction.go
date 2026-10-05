package zca

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type AddReactionResponse struct {
	MsgIDs []int64 `json:"msgIds"`
}

type CustomReaction struct {
	RType  int64
	Source int64
	Icon   string
}

type AddReactionDestination struct {
	Data struct {
		MsgID    string
		CliMsgID string
	}
	ThreadID string
	Type     ThreadType
}

// AddReactionIcon is a Reactions value or a CustomReaction.
type AddReactionIcon interface {
	addReactionParams() (rType, source int64, icon string)
}

var addReactionRTypes = map[Reactions]int64{
	ReactionsHAHA: 0, ReactionsLIKE: 3, ReactionsHEART: 5, ReactionsWOW: 32, ReactionsCRY: 2,
	ReactionsANGRY: 20, ReactionsKISS: 8, ReactionsTEARS_OF_JOY: 7, ReactionsSHIT: 66, ReactionsROSE: 120,
	ReactionsBROKEN_HEART: 65, ReactionsDISLIKE: 4, ReactionsLOVE: 29, ReactionsCONFUSED: 51, ReactionsWINK: 45,
	ReactionsFADE: 121, ReactionsSUN: 67, ReactionsBIRTHDAY: 126, ReactionsBOMB: 127, ReactionsOK: 68,
	ReactionsPEACE: 69, ReactionsTHANKS: 70, ReactionsPUNCH: 71, ReactionsSHARE: 72, ReactionsPRAY: 73,
	ReactionsNO: 131, ReactionsBAD: 132, ReactionsLOVE_YOU: 133, ReactionsSAD: 1, ReactionsVERY_SAD: 16,
	ReactionsCOOL: 21, ReactionsNERD: 22, ReactionsBIG_SMILE: 23, ReactionsSUNGLASSES: 26, ReactionsNEUTRAL: 30,
	ReactionsSAD_FACE: 35, ReactionsBYE: 36, ReactionsSLEEPY: 38, ReactionsWIPE: 39, ReactionsDIG: 42,
	ReactionsANGUISH: 44, ReactionsHANDCLAP: 46, ReactionsANGRY_FACE: 47, ReactionsF_CHAIR: 48, ReactionsL_CHAIR: 49,
	ReactionsR_CHAIR: 50, ReactionsSILENT: 52, ReactionsSURPRISE: 53, ReactionsEMBARRASSED: 54, ReactionsAFRAID: 60,
	ReactionsSAD2: 61, ReactionsBIG_LAUGH: 62, ReactionsRICH: 63, ReactionsBEER: 99,
}

func (r Reactions) addReactionParams() (int64, int64, string) {
	t, ok := addReactionRTypes[r]
	if !ok {
		t = -1
	}
	return t, 6, string(r)
}

func (c CustomReaction) addReactionParams() (int64, int64, string) { return c.RType, c.Source, c.Icon }

// addReactionParseInt mimics JS parseInt: invalid -> null.
func addReactionParseInt(s string) any {
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}
	return nil
}

// AddReaction reacts to a message.
func (a *API) AddReaction(ctx context.Context, icon AddReactionIcon, dest AddReactionDestination) (*AddReactionResponse, error) {
	if dest.Type != ThreadTypeUser && dest.Type != ThreadTypeGroup {
		return nil, newError("Thread type is invalid")
	}
	if icon == nil {
		return nil, newError("Invalid reaction")
	}
	rType, source, rIcon := icon.addReactionParams()
	msg := mustJSON(map[string]any{
		"rMsg": []map[string]any{{
			"gMsgID":  addReactionParseInt(dest.Data.MsgID),
			"cMsgID":  addReactionParseInt(dest.Data.CliMsgID),
			"msgType": 1,
		}},
		"rIcon":  rIcon,
		"rType":  rType,
		"source": source,
	})
	params := map[string]any{
		"react_list": []map[string]any{{"message": msg, "clientId": nowMs()}},
	}
	var u string
	switch dest.Type {
	case ThreadTypeUser:
		u = a.svc("reaction") + "/api/message/reaction"
		params["toid"] = dest.ThreadID
	case ThreadTypeGroup:
		u = a.svc("reaction") + "/api/group/reaction"
		params["grid"] = dest.ThreadID
		params["imei"] = a.IMEI
	}
	raw, err := call[struct {
		MsgIDs json.RawMessage `json:"msgIds"`
	}](ctx, a.Session, http.MethodPost, a.MakeURL(u, nil, true), params)
	if err != nil {
		return nil, err
	}
	ids, err := decodeData[[]int64](addReactionUnquote(raw.MsgIDs))
	if err != nil {
		return nil, err
	}
	return &AddReactionResponse{MsgIDs: ids}, nil
}

// addReactionUnquote unwraps a JSON-encoded string (Zalo sometimes double-encodes fields).
func addReactionUnquote(raw json.RawMessage) json.RawMessage {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return json.RawMessage(s)
	}
	return raw
}
