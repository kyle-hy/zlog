package chanmgr

import (
	"sync/atomic"
)

// NewChanMgr .
func NewChanMgr(sliceSize, chanSize uint64) *ChanMgr {

	sliceSize = nextPow2(sliceSize)
	chans := make([]chan []byte, sliceSize, sliceSize)
	for i := 0; i < int(sliceSize); i++ {
		chans[i] = make(chan []byte, chanSize)
	}
	return &ChanMgr{
		size:  sliceSize,
		chans: chans,
	}
}

// ChanMgr channel管道分片管理器，轮询读写管道，减少锁竞争
type ChanMgr struct {
	chans    []chan []byte // 管道切片
	size     uint64        // 切片大小
	writeIdx uint64        // 写索引
	readIdx  uint64        // 读索引
}

// NextWrite 切到下一个管道写
func (cm *ChanMgr) NextWrite() (chan []byte, uint64) {
	idx := atomic.AddUint64(&cm.writeIdx, 1)
	return cm.chans[modPow2(idx, cm.size)], idx
}

// NextRead 切到下一个管道读
func (cm *ChanMgr) NextRead() (chan []byte, uint64) {
	idx := atomic.AddUint64(&cm.readIdx, 1)
	return cm.chans[modPow2(idx, cm.size)], idx
}

// Len 获取管道的长度
func (cm *ChanMgr) Len(idx uint64) int {
	return len(cm.chans[modPow2(idx, cm.size)])
}

// Close 关闭所有管道
func (cm *ChanMgr) Close() {
	for _, c := range cm.chans {
		close(c)
	}
}

// 临近较大的2的整数次幂的32位整数
func nextPow2(n uint64) (k uint64) {
	if n == 0 {
		k = uint64(1)
	} else if (n & (n - 1)) == 0 {
		return n // is a power of 2
	} else {
		k = n - 1
		k |= k >> 1
		k |= k >> 2
		k |= k >> 4
		k |= k >> 8
		k |= k >> 16
		k |= k >> 32
		k++
	}
	return k
}

// 对2的整数次幂取模
func modPow2(a, b uint64) uint64 {
	return a & (b - 1)
}

// 是否是2的整数次幂
func isPow2(n uint64) bool {
	return (n != 0) && (n&(n-1)) == 0
}
