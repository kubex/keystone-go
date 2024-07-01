package keystone

import "time"

type EntityLocker interface {
	SetKeystoneLockResult(*EntityLockInfo)
}

type EntityLock struct {
	lockInfo *EntityLockInfo
}

type EntityLockInfo struct {
	ID           string
	LockedUntil  time.Time
	Message      string
	LockAcquired bool
}

func (e *EntityLock) LockData() *EntityLockInfo                      { return e.lockInfo }
func (e *EntityLock) SetKeystoneLockResult(lockInfo *EntityLockInfo) { e.lockInfo = lockInfo }
func (e *EntityLock) AcquiredLock() bool {
	if e.lockInfo == nil {
		return false
	}
	return e.lockInfo.LockAcquired
}
