package zca

type ThreadType int

const (
	ThreadTypeUser  ThreadType = 0
	ThreadTypeGroup ThreadType = 1
)

type DestType int

const (
	DestTypeGroup DestType = 1
	DestTypeUser  DestType = 3
	DestTypePage  DestType = 5
)

type Gender int

const (
	GenderMale   Gender = 0
	GenderFemale Gender = 1
)

type AvatarSize int

const (
	AvatarSizeSmall AvatarSize = 120
	// Experimental: use only if you know what you're doing.
	AvatarSizeMedium AvatarSize = 160 // zca-js PR #375: 180 -> 160
	AvatarSizeLarge  AvatarSize = 240
	// Experimental: use only if you know what you're doing.
	AvatarSizeExtraLarge AvatarSize = 360
)
