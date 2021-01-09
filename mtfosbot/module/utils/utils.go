package utils

import (
	"math"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// PageObject -
type PageObject struct {
	Page   int `json:"page" cc:"page"`
	Total  int `json:"total" cc:"total"`
	Offset int `json:"offset" cc:"offset"`
	Limit  int `json:"limit" cc:"limit"`
}

// CalcPage -
func CalcPage(count, page, max int) (po PageObject) {
	if count < 0 {
		count = 0
	}
	if page < 1 {
		page = 1
	}
	if max < 1 {
		max = 1
	}

	total := int(math.Ceil(float64(count) / float64(max)))
	if total < 1 {
		total = 1
	}
	if page > total {
		page = total
	}
	offset := (page - 1) * max
	if offset > count {
		offset = count
	}
	limit := max

	po = PageObject{}
	po.Limit = limit
	po.Page = page
	po.Offset = offset
	po.Total = total

	return
}

// ToMap struct to map[string]interface{}
func ToMap(ss interface{}) map[string]interface{} {
	t := reflect.ValueOf(ss)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	smap := make(map[string]interface{})
	mtag := regexp.MustCompile(`cc:\"(.+)\"`)

	re := regexp.MustCompile(`^[A-Z]`)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := string(t.Type().Field(i).Tag)
		str := mtag.FindStringSubmatch(tag)
		name := t.Type().Field(i).Name
		if !re.Match([]byte(name)) {
			continue
		}
		leveling := false
		if len(str) > 1 {
			strArr := strings.Split(str[1], ",")
			name = strArr[0]
			if len(strArr) > 1 && strArr[1] == "<<" {
				leveling = true
			}
		}

		if f.Kind() == reflect.Slice {
			if name == "-" {
				continue
			}
			if f.Len() > 0 && f.Index(0).Kind() != reflect.Struct {
				smap[name] = f.Interface()
				continue
			}

			tmp := make([]map[string]interface{}, f.Len())
			for i := 0; i < f.Len(); i++ {
				tmp[i] = ToMap(f.Index(i).Interface())
			}
			smap[name] = tmp
			continue
		}

		if f.Kind() == reflect.Struct {
			if name == "-" && leveling == false {
				continue
			}

			if leveling == true {
				tmp := ToMap(f.Interface())
				for k, v := range tmp {
					smap[k] = v
				}
			} else if name != "-" && leveling == false {
				smap[name] = f.Interface()
			}
			continue
		}

		if name != "-" {
			smap[name] = f.Interface()
		}
	}

	return smap
}

// ParsePath - parse file path to absPath
func ParsePath(dst string) string {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}

	if []rune(dst)[0] == '~' {
		home := UserHomeDir()
		if len(home) > 0 {
			dst = strings.Replace(dst, "~", home, -1)
		}
	}

	if path.IsAbs(dst) {
		dst = path.Clean(dst)
		return dst
	}

	str := path.Join(wd, dst)
	str = path.Clean(str)
	return str
}

// UserHomeDir - get user home directory
func UserHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

// CheckExists - check file exists
func CheckExists(filePath string, allowDir bool) bool {
	filePath = ParsePath(filePath)
	stat, err := os.Stat(filePath)
	if err != nil && !os.IsExist(err) {
		return false
	}
	if !allowDir && stat.IsDir() {
		return false
	}
	return true
}

// MemoryUsage -
type MemoryUsage struct {
	Alloc        uint64 `json:"alloc" cc:"alloc"`
	TotalAlloc   uint64 `json:"total_alloc" cc:"total_alloc"`
	Sys          uint64 `json:"sys" cc:"sys"`
	Lookups      uint64 `json:"lookups" cc:"lookups"`
	Mallocs      uint64 `json:"mallocs" cc:"mallocs"`
	Frees        uint64 `json:"frees" cc:"frees"`
	HeapAlloc    uint64 `json:"heap_alloc" cc:"heap_alloc"`
	HeapSys      uint64 `json:"heap_sys" cc:"heap_sys"`
	HeapIdle     uint64 `json:"heap_idle" cc:"heap_idle"`
	HeapInuse    uint64 `json:"heap_inuse" cc:"heap_inuse"`
	HeapReleased uint64 `json:"heap_released" cc:"heap_released"`
	HeapObjects  uint64 `json:"heap_objects" cc:"heap_objects"`
	StackInuse   uint64 `json:"stack_inuse" cc:"stack_inuse"`
	StackSys     uint64 `json:"stack_sys" cc:"stack_sys"`
	GCSys        uint64 `json:"gc_sys" cc:"gc_sys"`
	OtherSys     uint64 `json:"other_sys" cc:"other_sys"`
	NextGC       uint64 `json:"next_gc" cc:"next_gc"`
	LastGC       uint64 `json:"last_gc" cc:"last_gc"`
	NumGC        uint32 `json:"num_gc" cc:"num_gc"`
	NumForcedGC  uint32 `json:"num_forced_gc" cc:"num_forced_gc"`
}

// GetMemoryUsage -
func GetMemoryUsage() MemoryUsage {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	u := MemoryUsage{
		Alloc:        m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		Sys:          m.Sys,
		Lookups:      m.Lookups,
		Mallocs:      m.Mallocs,
		Frees:        m.Frees,
		HeapAlloc:    m.HeapAlloc,
		HeapSys:      m.HeapSys,
		HeapIdle:     m.HeapIdle,
		HeapInuse:    m.HeapInuse,
		HeapReleased: m.HeapReleased,
		HeapObjects:  m.HeapObjects,
		StackSys:     m.StackSys,
		StackInuse:   m.StackInuse,
		GCSys:        m.GCSys,
		OtherSys:     m.OtherSys,
		NextGC:       m.NextGC,
		LastGC:       m.LastGC,
		NumGC:        m.NumGC,
		NumForcedGC:  m.NumForcedGC,
	}
	return u
}
