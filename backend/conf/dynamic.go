package conf

import (
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 应用方式参考：horus/golib/config/server

type ReloadableConfig interface {
	UnmarshalFrom(*viper.Viper) error // forwarding to v.Unmarshal() is usually enough.
	OnReloadError(error)              // handles what is returned by UnmarshalFrom() when config changes.
}

type SafeConfig[T any] struct {
	data atomic.Pointer[T]
}

func (c *SafeConfig[T]) Get() *T {
	return c.data.Load()
}

// `v`应当已经被正确配置为可以找到文件。
// 注意！viper使用`mapstructure`（既非yaml也非json）填充结构体，并且（按照设计）忽略key的大小写。
func Watch[T any, PT interface {
	ReloadableConfig
	*T
}](v *viper.Viper) (*SafeConfig[T], error) {
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	d := new(T)
	if err := PT(d).UnmarshalFrom(v); err != nil {
		return nil, err
	}
	c := new(SafeConfig[T])
	c.data.Store(d)
	v.OnConfigChange(func(in fsnotify.Event) {
		d := new(T)
		p := PT(d)
		err := p.UnmarshalFrom(v)
		if err != nil {
			p.OnReloadError(err)
			return
		}
		c.data.Store(d) // previous pointer will be gc'ed after all external reference, if any, is dropped.
	})
	v.WatchConfig()
	return c, nil
}
