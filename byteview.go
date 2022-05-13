package atticuscache

// ByteView 存储真实的缓存值, byte 类型可以支持任意的数据类型存储，例如字符串、图片等
type ByteView struct {
	b []byte
}

// Len 实现了 Value 的 len 方法
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 转换为 切片，byteview 是只读的，返回一个 拷贝，防止值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// 把 v 的数据转换为字符串
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
