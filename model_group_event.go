package zca

import "encoding/json"

type GroupEventType string

const (
	GroupEventTypeJOIN_REQUEST  GroupEventType = "join_request"
	GroupEventTypeJOIN          GroupEventType = "join"
	GroupEventTypeLEAVE         GroupEventType = "leave"
	GroupEventTypeREMOVE_MEMBER GroupEventType = "remove_member"
	GroupEventTypeBLOCK_MEMBER  GroupEventType = "block_member"

	GroupEventTypeUPDATE_SETTING GroupEventType = "update_setting"
	GroupEventTypeUPDATE         GroupEventType = "update"
	GroupEventTypeNEW_LINK       GroupEventType = "new_link"

	GroupEventTypeADD_ADMIN    GroupEventType = "add_admin"
	GroupEventTypeREMOVE_ADMIN GroupEventType = "remove_admin"

	GroupEventTypeNEW_PIN_TOPIC     GroupEventType = "new_pin_topic"
	GroupEventTypeUPDATE_PIN_TOPIC  GroupEventType = "update_pin_topic"
	GroupEventTypeREORDER_PIN_TOPIC GroupEventType = "reorder_pin_topic"

	GroupEventTypeUPDATE_BOARD GroupEventType = "update_board"
	GroupEventTypeREMOVE_BOARD GroupEventType = "remove_board"

	GroupEventTypeUPDATE_TOPIC GroupEventType = "update_topic"
	GroupEventTypeUNPIN_TOPIC  GroupEventType = "unpin_topic"
	GroupEventTypeREMOVE_TOPIC GroupEventType = "remove_topic"

	GroupEventTypeACCEPT_REMIND GroupEventType = "accept_remind"
	GroupEventTypeREJECT_REMIND GroupEventType = "reject_remind"
	GroupEventTypeREMIND_TOPIC  GroupEventType = "remind_topic"

	GroupEventTypeUPDATE_AVATAR GroupEventType = "update_avatar"

	GroupEventTypeUNKNOWN GroupEventType = "unknown"
)

type GroupEventUpdateMember struct {
	ID       string `json:"id"`
	DName    string `json:"dName"`
	Avatar   string `json:"avatar"`
	Type     int64  `json:"type"`
	Avatar25 string `json:"avatar_25"`
}

// TGroupEventBase. Info (GroupEventGroupInfo: group_link, link_expired_time, ...) and
// ExtraData (GroupEventExtraData: featureId, field, ...) are open objects.
type TGroupEventBase struct {
	SubType       int64                    `json:"subType"`
	GroupID       string                   `json:"groupId"`
	CreatorID     string                   `json:"creatorId"`
	GroupName     string                   `json:"groupName"`
	SourceID      string                   `json:"sourceId"`
	UpdateMembers []GroupEventUpdateMember `json:"updateMembers"`
	GroupSetting  *GroupSetting            `json:"groupSetting"`
	GroupTopic    *GroupTopic              `json:"groupTopic"`
	Info          map[string]any           `json:"info"`
	ExtraData     map[string]any           `json:"extraData"`
	Time          string                   `json:"time"`
	Avt           *string                  `json:"avt"`
	FullAvt       *string                  `json:"fullAvt"`
	IsAdd         int64                    `json:"isAdd"`
	HideGroupInfo int64                    `json:"hideGroupInfo"`
	Version       string                   `json:"version"`
	GroupType     int64                    `json:"groupType"`
	ClientID      *int64                   `json:"clientId,omitempty"`
	ErrorMap      map[string]any           `json:"errorMap,omitempty"`
	E2ee          *int64                   `json:"e2ee,omitempty"`
}

type TGroupEventJoinRequest struct {
	UIDs         []string `json:"uids"`
	TotalPending int64    `json:"totalPending"`
	GroupID      string   `json:"groupId"`
	Time         string   `json:"time"`
}

type TGroupEventPinTopic struct {
	OldBoardVersion int64      `json:"oldBoardVersion"`
	BoardVersion    int64      `json:"boardVersion"`
	Topic           GroupTopic `json:"topic"`
	ActorID         string     `json:"actorId"`
	GroupID         string     `json:"groupId"`
}

type TGroupEventReorderPinTopicItem struct {
	TopicID   string `json:"topicId"`
	TopicType int64  `json:"topicType"`
}

type TGroupEventReorderPinTopic struct {
	OldBoardVersion int64                            `json:"oldBoardVersion"`
	ActorID         string                           `json:"actorId"`
	Topics          []TGroupEventReorderPinTopicItem `json:"topics"`
	GroupID         string                           `json:"groupId"`
	BoardVersion    int64                            `json:"boardVersion"`
	Topic           any                              `json:"topic"` // always null
}

type TGroupEventBoard struct {
	SourceID  string `json:"sourceId"`
	GroupName string `json:"groupName"`
	// GroupTopic is (GroupTopic | ReminderGroup) with params as a string.
	GroupTopic map[string]any `json:"groupTopic"`
	GroupID    string         `json:"groupId"`
	CreatorID  string         `json:"creatorId"`

	SubType       int64                    `json:"subType,omitempty"`
	UpdateMembers []GroupEventUpdateMember `json:"updateMembers,omitempty"`
	GroupSetting  *GroupSetting            `json:"groupSetting,omitempty"`
	Info          map[string]any           `json:"info,omitempty"`
	ExtraData     map[string]any           `json:"extraData,omitempty"`
	Time          string                   `json:"time,omitempty"`
	Avt           any                      `json:"avt,omitempty"`
	FullAvt       any                      `json:"fullAvt,omitempty"`
	IsAdd         int64                    `json:"isAdd,omitempty"`
	HideGroupInfo int64                    `json:"hideGroupInfo,omitempty"`
	Version       string                   `json:"version,omitempty"`
	GroupType     int64                    `json:"groupType,omitempty"`
}

type TGroupEventRemindRespond struct {
	TopicID       string   `json:"topicId"`
	UpdateMembers []string `json:"updateMembers"`
	GroupID       string   `json:"groupId"`
	Time          string   `json:"time"`
}

type TGroupEventRemindTopic struct {
	Msg        string `json:"msg"`
	EditorID   string `json:"editorId"`
	Color      string `json:"color"`
	Emoji      string `json:"emoji"`
	CreatorID  string `json:"creatorId"`
	EditTime   int64  `json:"editTime"`
	Type       int64  `json:"type"`
	Duration   int64  `json:"duration"`
	GroupID    string `json:"group_id"`
	CreateTime int64  `json:"createTime"`
	Repeat     int64  `json:"repeat"`
	StartTime  int64  `json:"startTime"`
	Time       int64  `json:"time"`
	RemindType int64  `json:"remindType"`
}

// GroupEvent is the TS GroupEvent union. Data holds, by Type:
//   - JOIN_REQUEST: *TGroupEventJoinRequest
//   - NEW_PIN_TOPIC, UNPIN_TOPIC, UPDATE_PIN_TOPIC: *TGroupEventPinTopic
//   - REORDER_PIN_TOPIC: *TGroupEventReorderPinTopic
//   - UPDATE_BOARD, REMOVE_BOARD: *TGroupEventBoard
//   - ACCEPT_REMIND, REJECT_REMIND: *TGroupEventRemindRespond
//   - REMIND_TOPIC: *TGroupEventRemindTopic
//   - everything else: *TGroupEventBase
type GroupEvent struct {
	Type     GroupEventType `json:"type"`
	Data     any            `json:"data"`
	Act      string         `json:"act"`
	ThreadID string         `json:"threadId"`
	IsSelf   bool           `json:"isSelf"`
}

// InitializeGroupEvent ports initializeGroupEvent. data is the raw control content.data:
// a JSON string holding JSON is unwrapped first, as listen.ts does.
// ponytail: decode errors are ignored (best-effort, like the TS casts).
func InitializeGroupEvent(uid string, data json.RawMessage, typ GroupEventType, act string) *GroupEvent {
	raw := unwrapJSONString(data)

	var ids struct {
		GroupIDSnake *string `json:"group_id"`
		GroupID      string  `json:"groupId"`
	}
	_ = json.Unmarshal(raw, &ids)
	ev := &GroupEvent{Type: typ, Act: act, ThreadID: ids.GroupID}
	if ids.GroupIDSnake != nil {
		ev.ThreadID = *ids.GroupIDSnake
	}

	switch typ {
	case GroupEventTypeJOIN_REQUEST:
		var d TGroupEventJoinRequest
		_ = json.Unmarshal(raw, &d)
		ev.Data = &d
	case GroupEventTypeNEW_PIN_TOPIC, GroupEventTypeUNPIN_TOPIC, GroupEventTypeUPDATE_PIN_TOPIC:
		var d TGroupEventPinTopic
		_ = json.Unmarshal(raw, &d)
		ev.Data, ev.IsSelf = &d, d.ActorID == uid
	case GroupEventTypeREORDER_PIN_TOPIC:
		var d TGroupEventReorderPinTopic
		_ = json.Unmarshal(raw, &d)
		ev.Data, ev.IsSelf = &d, d.ActorID == uid
	case GroupEventTypeUPDATE_BOARD, GroupEventTypeREMOVE_BOARD:
		var d TGroupEventBoard
		_ = json.Unmarshal(raw, &d)
		ev.Data, ev.IsSelf = &d, d.SourceID == uid
	case GroupEventTypeACCEPT_REMIND, GroupEventTypeREJECT_REMIND:
		var d TGroupEventRemindRespond
		_ = json.Unmarshal(raw, &d)
		ev.Data = &d
		for _, id := range d.UpdateMembers {
			if id == uid {
				ev.IsSelf = true
				break
			}
		}
	case GroupEventTypeREMIND_TOPIC:
		var d TGroupEventRemindTopic
		_ = json.Unmarshal(raw, &d)
		ev.Data, ev.IsSelf = &d, d.CreatorID == uid
	default:
		var d TGroupEventBase
		_ = json.Unmarshal(raw, &d)
		ev.Data, ev.IsSelf = &d, d.SourceID == uid
		for _, m := range d.UpdateMembers {
			if m.ID == uid {
				ev.IsSelf = true
				break
			}
		}
	}
	return ev
}

func getGroupEventType(act string) GroupEventType {
	switch t := GroupEventType(act); t {
	case GroupEventTypeJOIN_REQUEST, GroupEventTypeJOIN, GroupEventTypeLEAVE, GroupEventTypeREMOVE_MEMBER,
		GroupEventTypeBLOCK_MEMBER, GroupEventTypeUPDATE_SETTING, GroupEventTypeUPDATE_AVATAR, GroupEventTypeUPDATE,
		GroupEventTypeNEW_LINK, GroupEventTypeADD_ADMIN, GroupEventTypeREMOVE_ADMIN,
		GroupEventTypeNEW_PIN_TOPIC, GroupEventTypeUPDATE_PIN_TOPIC, GroupEventTypeUPDATE_TOPIC,
		GroupEventTypeUPDATE_BOARD, GroupEventTypeREMOVE_BOARD, GroupEventTypeREORDER_PIN_TOPIC,
		GroupEventTypeUNPIN_TOPIC, GroupEventTypeREMOVE_TOPIC, GroupEventTypeACCEPT_REMIND,
		GroupEventTypeREJECT_REMIND, GroupEventTypeREMIND_TOPIC:
		// Every act string equals its GroupEventType value.
		return t
	}
	return GroupEventTypeUNKNOWN
}
