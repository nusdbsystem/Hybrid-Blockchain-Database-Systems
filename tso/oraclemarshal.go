package tso

import (
	"io"
	"sync"
)

func (t *GetTS) BinarySize() (nBytes int, sizeKnown bool) {
	return 4, true
}

type GetTSCache struct {
	mu    sync.Mutex
	cache []*GetTS
}

func NewGetTSCache() *GetTSCache {
	c := &GetTSCache{}
	c.cache = make([]*GetTS, 0)
	return c
}

func (p *GetTSCache) Get() *GetTS {
	var t *GetTS
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &GetTS{}
	}
	return t
}
func (p *GetTSCache) Put(t *GetTS) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *GetTS) Marshal(wire io.Writer) {
	var b [4]byte
	var bs []byte
	bs = b[:4]
	tmp32 := t.Num
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	wire.Write(bs)
}

func (t *GetTS) Unmarshal(wire io.Reader) error {
	var b [4]byte
	var bs []byte
	bs = b[:4]
	if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
		return err
	}
	t.Num = int32(uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24))
	return nil
}

func (t *ReplyTS) BinarySize() (nBytes int, sizeKnown bool) {
	return 8, true
}

type ReplyTSCache struct {
	mu    sync.Mutex
	cache []*ReplyTS
}

func NewReplyTSCache() *ReplyTSCache {
	c := &ReplyTSCache{}
	c.cache = make([]*ReplyTS, 0)
	return c
}

func (p *ReplyTSCache) Get() *ReplyTS {
	var t *ReplyTS
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &ReplyTS{}
	}
	return t
}
func (p *ReplyTSCache) Put(t *ReplyTS) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *ReplyTS) Marshal(wire io.Writer) {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	tmp64 := t.Timestamp
	bs[0] = byte(tmp64)
	bs[1] = byte(tmp64 >> 8)
	bs[2] = byte(tmp64 >> 16)
	bs[3] = byte(tmp64 >> 24)
	bs[4] = byte(tmp64 >> 32)
	bs[5] = byte(tmp64 >> 40)
	bs[6] = byte(tmp64 >> 48)
	bs[7] = byte(tmp64 >> 56)
	wire.Write(bs)
}

func (t *ReplyTS) Unmarshal(wire io.Reader) error {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.Timestamp = int64(uint64(bs[0]) | (uint64(bs[1]) << 8) | (uint64(bs[2]) << 16) | (uint64(bs[3]) << 24) | (uint64(bs[4]) << 32) | (uint64(bs[5]) << 40) | (uint64(bs[6]) << 48) | (uint64(bs[7]) << 56))
	return nil
}

func (t *LogTS) BinarySize() (nBytes int, sizeKnown bool) {
	return 12, true
}

type LogTSCache struct {
	mu    sync.Mutex
	cache []*LogTS
}

func NewLogTSCache() *LogTSCache {
	c := &LogTSCache{}
	c.cache = make([]*LogTS, 0)
	return c
}

func (p *LogTSCache) Get() *LogTS {
	var t *LogTS
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &LogTS{}
	}
	return t
}
func (p *LogTSCache) Put(t *LogTS) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *LogTS) Marshal(wire io.Writer) {
	var b [12]byte
	var bs []byte
	bs = b[:12]
	tmp32 := t.crc
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp64 := t.ts
	bs[4] = byte(tmp64)
	bs[5] = byte(tmp64 >> 8)
	bs[6] = byte(tmp64 >> 16)
	bs[7] = byte(tmp64 >> 24)
	bs[8] = byte(tmp64 >> 32)
	bs[9] = byte(tmp64 >> 40)
	bs[10] = byte(tmp64 >> 48)
	bs[11] = byte(tmp64 >> 56)
	wire.Write(bs)
}

func (t *LogTS) Unmarshal(wire io.Reader) error {
	var b [12]byte
	var bs []byte
	bs = b[:12]
	if _, err := io.ReadAtLeast(wire, bs, 12); err != nil {
		return err
	}
	t.crc = uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)
	t.ts = int64(uint64(bs[4]) | (uint64(bs[5]) << 8) | (uint64(bs[6]) << 16) | (uint64(bs[7]) << 24) | (uint64(bs[8]) << 32) | (uint64(bs[9]) << 40) | (uint64(bs[10]) << 48) | (uint64(bs[11]) << 56))
	return nil
}
