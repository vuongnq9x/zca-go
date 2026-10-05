package zca

import "encoding/json"

type Reactions string

const (
	ReactionsHEART        Reactions = "/-heart"
	ReactionsLIKE         Reactions = "/-strong"
	ReactionsHAHA         Reactions = ":>"
	ReactionsWOW          Reactions = ":o"
	ReactionsCRY          Reactions = ":-(("
	ReactionsANGRY        Reactions = ":-h"
	ReactionsKISS         Reactions = ":-*"
	ReactionsTEARS_OF_JOY Reactions = ":')"
	ReactionsSHIT         Reactions = "/-shit"
	ReactionsROSE         Reactions = "/-rose"
	ReactionsBROKEN_HEART Reactions = "/-break"
	ReactionsDISLIKE      Reactions = "/-weak"
	ReactionsLOVE         Reactions = ";xx"
	ReactionsCONFUSED     Reactions = ";-/"
	ReactionsWINK         Reactions = ";-)"
	ReactionsFADE         Reactions = "/-fade"
	ReactionsSUN          Reactions = "/-li"
	ReactionsBIRTHDAY     Reactions = "/-bd"
	ReactionsBOMB         Reactions = "/-bome"
	ReactionsOK           Reactions = "/-ok"
	ReactionsPEACE        Reactions = "/-v"
	ReactionsTHANKS       Reactions = "/-thanks"
	ReactionsPUNCH        Reactions = "/-punch"
	ReactionsSHARE        Reactions = "/-share"
	ReactionsPRAY         Reactions = "_()_"
	ReactionsNO           Reactions = "/-no"
	ReactionsBAD          Reactions = "/-bad"
	ReactionsLOVE_YOU     Reactions = "/-loveu"
	ReactionsSAD          Reactions = "--b"
	ReactionsVERY_SAD     Reactions = ":(("
	ReactionsCOOL         Reactions = "x-)"
	ReactionsNERD         Reactions = "8-)"
	ReactionsBIG_SMILE    Reactions = ";-d"
	ReactionsSUNGLASSES   Reactions = "b-)"
	ReactionsNEUTRAL      Reactions = ":--|"
	ReactionsSAD_FACE     Reactions = "p-("
	ReactionsBYE          Reactions = ":-bye"
	ReactionsSLEEPY       Reactions = "|-)"
	ReactionsWIPE         Reactions = ":wipe"
	ReactionsDIG          Reactions = ":-dig"
	ReactionsANGUISH      Reactions = "&-("
	ReactionsHANDCLAP     Reactions = ":handclap"
	ReactionsANGRY_FACE   Reactions = ">-|"
	ReactionsF_CHAIR      Reactions = ":-f"
	ReactionsL_CHAIR      Reactions = ":-l"
	ReactionsR_CHAIR      Reactions = ":-r"
	ReactionsSILENT       Reactions = ";-x"
	ReactionsSURPRISE     Reactions = ":-o"
	ReactionsEMBARRASSED  Reactions = ";-s"
	ReactionsAFRAID       Reactions = ";-a"
	ReactionsSAD2         Reactions = ":-<"
	ReactionsBIG_LAUGH    Reactions = ":))"
	ReactionsRICH         Reactions = "$-)"
	ReactionsBEER         Reactions = "/-beer"
	ReactionsNONE         Reactions = ""
)

type TReactionRMsg struct {
	GMsgID  StringOrNumber `json:"gMsgID"`
	CMsgID  StringOrNumber `json:"cMsgID"`
	MsgType int64          `json:"msgType"`
}

type TReactionContent struct {
	RMsg   []TReactionRMsg `json:"rMsg"`
	RIcon  Reactions       `json:"rIcon"`
	RType  int64           `json:"rType"`
	Source int64           `json:"source"`
}

// UnmarshalJSON also accepts content sent as a JSON-encoded string (listen.ts parses it).
func (c *TReactionContent) UnmarshalJSON(b []byte) error {
	type plain TReactionContent
	return json.Unmarshal(unwrapJSONString(b), (*plain)(c))
}

type TReaction struct {
	ActionID string           `json:"actionId"`
	MsgID    string           `json:"msgId"`
	CliMsgID string           `json:"cliMsgId"`
	MsgType  string           `json:"msgType"`
	UIDFrom  string           `json:"uidFrom"`
	IDTo     string           `json:"idTo"`
	DName    string           `json:"dName,omitempty"`
	Content  TReactionContent `json:"content"`
	TS       string           `json:"ts"`
	TTL      int64            `json:"ttl"`
}

// Reaction is TS class Reaction; isGroup is expressed as Type == ThreadTypeGroup.
type Reaction struct {
	Type     ThreadType `json:"type"`
	Data     TReaction  `json:"data"`
	ThreadID string     `json:"threadId"`
	IsSelf   bool       `json:"isSelf"`
}

func NewReaction(uid string, data TReaction, isGroup bool) *Reaction {
	r := &Reaction{Type: ThreadTypeUser, ThreadID: data.UIDFrom, IsSelf: data.UIDFrom == "0"}
	if isGroup {
		r.Type = ThreadTypeGroup
	}
	if isGroup || data.UIDFrom == "0" {
		r.ThreadID = data.IDTo
	}
	if data.IDTo == "0" {
		data.IDTo = uid
	}
	if data.UIDFrom == "0" {
		data.UIDFrom = uid
	}
	r.Data = data
	return r
}
