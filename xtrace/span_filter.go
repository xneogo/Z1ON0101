/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/10/18 -- 16:25
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: span_filter.go
*/

package xtrace

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/xneogo/Z1ON0101/xconfig/xapollo"
	"github.com/xneogo/Z1ON0101/xlog"
	"github.com/xneogo/matrix/mconfig/mobserver"
)

const filterUrls = "span_filter_urls"

const listConfigSep = ","

func InitTraceSpanFilter() error {
	fun := "TraceSpanFilter.init --> "
	ctx := context.Background()

	if err := initApolloCenter(ctx); err != nil {
		return err
	}

	urls, ok := apolloCenter.GetStringWithNamespace(ctx, xapollo.DefaultApolloTraceNamespace, filterUrls)
	if !ok {
		return fmt.Errorf("not get %s from apollo namespace %s", filterUrls, xapollo.DefaultApolloTraceNamespace)
	}
	xlog.Infof(ctx, "%s get %s from apollo: %s", fun, filterUrls, urls)

	urlList := strings.Split(urls, listConfigSep)

	apolloSpanFilterConfig = &spanFilterConfig{
		urls: urlList,
	}

	observer := mobserver.NewConfigObserver(apolloSpanFilterConfig.handleChangeEvent)
	apolloCenter.RegisterObserver(ctx, observer)
	return nil
}

type spanFilterConfig struct {
	mu sync.RWMutex

	urls []string
}

func (m *spanFilterConfig) updateUrls(urls []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.urls = urls
}

func (m *spanFilterConfig) filterUrl(url string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.urls {
		if u == url {
			return false
		}
	}

	return true
}

func (m *spanFilterConfig) handleChangeEvent(ctx context.Context, event *mobserver.ChangeEvent) {
	fun := "spanFilterConfig.HandleChangeEvent --> "

	if event.Namespace != xapollo.DefaultApolloTraceNamespace {
		return
	}

	for key, change := range event.Changes {
		if key == filterUrls {
			xlog.Infof(ctx, "%s get key %s from apollo, old value: %s, new value: %s", fun, key, change.OldValue, change.NewValue)
			urlList := strings.Split(change.NewValue, listConfigSep)
			m.updateUrls(urlList)
		}
	}
}

func UrlSpanFilter(r *http.Request) bool {
	if apolloSpanFilterConfig != nil {
		return apolloSpanFilterConfig.filterUrl(r.URL.Path)
	}

	return true
}
