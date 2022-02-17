package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logCommonKeyUID    = "uid"
	logCommonKeyRoomID = "roomID"
	logCommonKeyHostID = "hostID"
	logCommonKeyGoID   = "goID"
)

// UID 通用的uid field
func UID(uid uint64) zapcore.Field {
	return zap.Uint64(logCommonKeyUID, uid)
}

// UIDInt64 通用的uid field int64
func UIDInt64(uid int64) zapcore.Field {
	return zap.Int64(logCommonKeyUID, uid)
}

// RoomID 通用的房间ID
func RoomID(roomID uint64) zapcore.Field {
	return zap.Uint64(logCommonKeyRoomID, roomID)
}

// HostID 通用的房间主播ID
func HostID(hostID uint64) zapcore.Field {
	return zap.Uint64(logCommonKeyHostID, hostID)
}

// GoID 协程ID
func GoID(goID int64) zapcore.Field {
	return zap.Int64(logCommonKeyGoID, goID)
}
