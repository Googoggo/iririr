package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

//声明新切片类型
type units []uint32

//返回切片长度
func (x units) Len() int{
	return len(x)
}

//
func (x units) Less(i,j int) bool{
	return x[i] <x[j]
}

func (x units) Swap(i, j int){
	x[i],x[j] = x[j],x[i]
}

//当哈希环上没有数据出错
var errEmpty = errors.New("Hash 环没有数据")
// 创建结构体保存一致性hash信息
type Consistent struct {
	//hash环，key为哈希值，值存放节点的信息
	circle map[uint32]string
	//已经排序的节点hash切片
	sortedHashes units
	//虚拟节点, 增加hash一致平衡
	VirtualNode int
	//map读写锁
	sync.RWMutex
}

//更新排序，方便查找
func (c *Consistent) updateSortedHashes(){
	hashes := c.sortedHashes[:0]
	//判断切片容量是否过大
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle){
		hashes = nil
	}
	for k := range c.circle{
		hashes = append(hashes,k)
	}
	//对所有节点hash值进行排序，方便之后二分查找
	sort.Sort(hashes)
	//重新辅助
	c.sortedHashes = hashes
}
func NewConsistent() *Consistent{
	return &Consistent{
		//初始化变量
		circle: make(map[uint32]string),
		VirtualNode: 20,

	}
}
//自动生成Key值
func (c *Consistent) generateKey(element string, index int)string{
	return element+strconv.Itoa(index)
}

func (c *Consistent) hashKey(key string) uint32{
	if len(key) < 64{
		//声明一个数组长度为64
		var srcatch [64]byte
		//拷贝数据到数组中
		copy(srcatch[:],key)
		//使用IEEE多项式返回数据的CRC-32
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

//向哈希环添加节点
func (c *Consistent) Add(element string){
	//枷锁
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) add(element string){
	//循环虚拟节点，设置副本
	for i:=0;i<c.VirtualNode;i++{
		c.circle[c.hashKey(c.generateKey(element,i))] = element
	}
	//更新排序
	c.updateSortedHashes()
}

func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
		c.updateSortedHashes()
	}
}

//删除节点
func (c *Consistent) Remove (element string){
		c.Lock()
		defer c.Unlock()
		c.remove(element)
}

//顺时针查找最近的节点
func (c *Consistent) search (key uint32) int{
	//查找算法
	f := func(x int) bool{
		return c.sortedHashes[x]>key
	}
	//使用二分查找,算出来搜索指定切片满足条件最小值
	i := sort.Search(len(c.sortedHashes),f)
	//如果超出范围
	if i>=len(c.sortedHashes){
		i = 0
	}
	return i
}

//获取最近节点
func (c *Consistent) Get(name string)(string ,error){
	c.RLock()
	defer c.Unlock()
	//如果为零
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	//计算hash值
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}

