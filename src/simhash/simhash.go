// Copyright 2013 Matthew Fonda. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// simhash package implements Charikar's simhash algorithm to generate a 64-bit
// fingerprint of a given document.
//
// simhash指纹具有类似文档具有类似指纹的属性。 因此，如果文档相似，则两个指纹之间的汉明距离将很小
package simhash

import (
	"bytes"
	"golang.org/x/text/unicode/norm"
	"hash/fnv"
	"regexp"
)

type Vector [64]int

//Feature包括64位哈希和权重
type Feature interface {
	// Sum 返回此功能(feature)的64位总和
	Sum() uint64

	// Weight returns the weight of this feature
	Weight() int
}

// FeatureSet表示给定文档中的一组要素
type FeatureSet interface {
	GetFeatures() []Feature
}

// Vectorize在给定一组特征的情况下生成64维向量。
//向量初始化为零。 然后，如果设置了特征的第i个比特，则向量的第i个元素按第i个特征的权重递增，否则递减第i个特征的权重。
func Vectorize(features []Feature) Vector {
	var v Vector
	for _, feature := range features {
		sum := feature.Sum()
		weight := feature.Weight()
		for i := uint8(0); i < 64; i++ {
			bit := ((sum >> i) & 1)
			if bit == 1 {
				v[i] += weight
			} else {
				v[i] -= weight
			}
		}
	}
	return v
}

//VectorizeBytes在给定一组[] []字节的情况下生成64维向量，其中每个[]字节是具有均匀权重的要素。

func VectorizeBytes(features [][]byte) Vector {
	var v Vector
	h := fnv.New64()
	for _, feature := range features {
		h.Reset()
		h.Write(feature)
		sum := h.Sum64()
		for i := uint8(0); i < 64; i++ {
			bit := ((sum >> i) & 1)
			if bit == 1 {
				v[i]++
			} else {
				v[i]--
			}
		}
	}
	return v
}

//指纹返回给定向量的64位指纹。
//给定64维向量v的指纹f定义如下：
//   f[i] = 1 if v[i] >= 0
//   f[i] = 0 if v[i] < 0
func Fingerprint(v Vector) uint64 {
	var f uint64
	for i := uint8(0); i < 64; i++ {
		if v[i] >= 0 {
			f |= (1 << i)
		}
	}
	return f
}

type feature struct {
	sum    uint64
	weight int
}

// Sum returns the 64-bit hash of this feature
// Sum返回此功能的64位哈希值
func (f feature) Sum() uint64 {
	return f.sum
}

// Weight returns the weight of this feature
func (f feature) Weight() int {
	return f.weight
}

// Returns a new feature representing the given byte slice, using a weight of 1
//使用权重1返回表示给定字节切片的新要素
func NewFeature(f []byte) feature {
	h := fnv.New64()
	h.Write(f)
	return feature{h.Sum64(), 1}
}

// Returns a new feature representing the given byte slice with the given weight
//返回一个新特征，表示具有给定权重的给定字节切片
func NewFeatureWithWeight(f []byte, weight int) feature {
	fw := NewFeature(f)
	fw.weight = weight
	return fw
}

// Compare计算两个64位整数之间的汉明距离
//目前，这是使用Kernighan方法[1]计算的。 存在其他方法可能更有效并且在某些方面值得探索
// [1] http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetKernighan
func Compare(a uint64, b uint64) uint8 {
	v := a ^ b
	println("V",v)
	var c uint8
	for c = 0; v != 0; c++ {
		v &= v - 1
	}
	return c
}

//返回给定功能集的64位simhash
func Simhash(fs FeatureSet) uint64 {
	return Fingerprint(Vectorize(fs.GetFeatures()))
}//

// 返回给定字节的64位simhash
func SimhashBytes(b [][]byte) uint64 {
	return Fingerprint(VectorizeBytes(b))
}

// WordFeatureSet是一个功能集，其中每个单词都是一个功能，所有功能都相等。
type WordFeatureSet struct {
	b []byte
}

func NewWordFeatureSet(b []byte) *WordFeatureSet {
	fs := &WordFeatureSet{b}
	fs.normalize()
	return fs
}

func (w *WordFeatureSet) normalize() {
	w.b = bytes.ToLower(w.b)//所有字母大写转到小写
}

var boundaries = regexp.MustCompile(`[\w']+(?:\://[\w\./]+){0,1}`)
var unicodeBoundaries = regexp.MustCompile(`[\pL-_']+`)

// Returns a []Feature representing each word in the byte slice
//返回表示字节切片中每个字的[]特征
func (w *WordFeatureSet) GetFeatures() []Feature {
	return getFeatures(w.b, boundaries)
}

// UnicodeWordFeatureSet是一个功能集，其中每个单词都是一个功能，所有功能都相等。
//
// See: http://blog.golang.org/normalization
// See: https://groups.google.com/forum/#!topic/golang-nuts/YyH1f_qCZVc
type UnicodeWordFeatureSet struct {
	b []byte
	f norm.Form
}

func NewUnicodeWordFeatureSet(b []byte, f norm.Form) *UnicodeWordFeatureSet {
	fs := &UnicodeWordFeatureSet{b, f}
	fs.normalize()
	return fs
}

func (w *UnicodeWordFeatureSet) normalize() {
	b := bytes.ToLower(w.f.Append(nil, w.b...))
	w.b = b
}

// Returns a []Feature representing each word in the byte slice
//返回表示字节切片中每个字的[]特征
func (w *UnicodeWordFeatureSet) GetFeatures() []Feature {
	return getFeatures(w.b, unicodeBoundaries)
}

// 使用给定的regexp拆分给定的[]字节，然后返回包含由regexp匹配的每个部分构造的Feature的切片
func getFeatures(b []byte, r *regexp.Regexp) []Feature {
	words := r.FindAll(b, -1)
	features := make([]Feature, len(words))
	for i, w := range words {
		features[i] = NewFeature(w)
	}
	return features
}

// Shingle返回给定字节集的w-shingling。 For example, if the given
// input was {"this", "is", "a", "test"}, this returns {"this is", "is a", "a test"}
func Shingle(w int, b [][]byte) [][]byte {
	if w < 1 {
		// TODO: use error here instead of panic?
		panic("simhash.Shingle(): k must be a positive integer")
	}

	if w == 1 {
		return b
	}

	if w > len(b) {
		w = len(b)
	}

	count := len(b) - w + 1
	shingles := make([][]byte, count)
	for i := 0; i < count; i++ {
		shingles[i] = bytes.Join(b[i:i+w], []byte(" "))
	}
	return shingles
}
