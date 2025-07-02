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
 @Time    : 2024/11/4 -- 17:54
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: election.go
*/

package xetcd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/xneogo/Z1ON0101/xlog"
	etcd "go.etcd.io/etcd/client/v2"
)

const (
	LockPrefixPath = "/nerv/election"

	defaultTTL    = 3 * time.Minute
	etcdOpTimeout = 3 * time.Second
)

var (
	ERRServiceNameNull = errors.New("service name not be empty")
	ERRNameNull        = errors.New("election name not be empty")
	ERRPreLane         = errors.New("pre environment cannot acquire the lock")
	ERREtcdClientNil   = errors.New("etcd client is nil")
)

type Election struct {
	Client etcd.KeysAPI

	serviceName string

	lock        sync.Mutex
	stopManager map[string]chan struct{}
}

type ElectionOptions struct {
	// 必填项：标识符
	Name string
	// 设置etcd path value的标识符，用于区分不同节点。默认值：hostName
	Value string
	// 占用当前锁的过期时间, 心跳时间为TTL/3
	TTL time.Duration
	// 泳道信息，直接传入对应roc的Lane即可，为了使pre泳道不参与选主，并一直阻塞
	Lane string
}

/*
	NewElection 需要手动传入etcd endpoints地址来实现选主
*/
func NewElection(ctx context.Context, serviceName string, endpoints []string) (*Election, error) {
	etcdClient, err := etcd.New(etcd.Config{
		Endpoints:               endpoints,
		HeaderTimeoutPerRequest: etcdOpTimeout,
	})
	if err != nil {
		xlog.Errorf(ctx, "init default etcd client error, err: %v", err)
		return nil, err
	}
	return NewElectionWithClient(ctx, serviceName, etcdClient)
}

func NewElectionWithClient(ctx context.Context, serviceName string, etcdClient etcd.Client) (*Election, error) {
	if serviceName == "" {
		return nil, ERRServiceNameNull
	}
	if etcdClient == nil {
		return nil, ERREtcdClientNil
	}
	return &Election{
		Client:      etcd.NewKeysAPI(etcdClient),
		stopManager: make(map[string]chan struct{}),
		serviceName: serviceName,
	}, nil
}

// Campaign
//
//  根据提供的option，发起一个关于自定义name的选举.
//
//  name: 本次选举的名称，必填，不可为空
//  value: 当前选举所用的值，默认值：hostname，下次使用相同值的节点可优先获取
//  tll: 节点存活时间，即：leader节点不在保持心跳，其最多占用的时间，默认30s
//
//	会将符合条件的value赋值到etcd的path使之成为leader，
// 	多个相同ServiceName，Name参数的Election会参与同个选举，
// 	并且只有一个会成为leader，其余将阻塞，直到成为leader
//  当context被cancel或者timeout将不再阻塞
func (p *Election) Campaign(ctx context.Context, opts *ElectionOptions) error {
	err := validateOptions(opts)
	if errors.Is(err, ERRPreLane) {
		<-(chan int)(nil)
	}
	if err != nil {
		return err
	}
	// 获取锁
	err = p.setValue(ctx, p.path(opts.Name), opts.Value, opts.TTL)
	if err != nil {
		return err
	}

	p.startHeart(ctx, opts)
	return nil
}

// TryCampaign
//
// 尝试参与一次选举，非阻塞的，成功则返回true，失败则返回false
func (p *Election) TryCampaign(ctx context.Context, opts *ElectionOptions) (bool, error) {
	err := validateOptions(opts)
	if err != nil {
		return false, err
	}

	path := p.path(opts.Name)
	if err := p.setPrevValue(ctx, path, opts.Value, opts.TTL); err == nil {
		p.startHeart(ctx, opts)
		return true, nil
	}

	// key不存在则设置值
	err = p.setPrevNoExist(ctx, path, opts.Value, opts.TTL)
	if err == nil {
		p.startHeart(ctx, opts)
		return true, nil
	}
	return false, nil
}

// Resign
// 	 提供对应发起选举的option
//	 让出leader开始一次新的选举
//
//   name：需要辞职选举的名称
//   value：给出对应leader的value，否则辞职失败
func (p *Election) Resign(ctx context.Context, opts *ElectionOptions) error {
	err := validateOptions(opts)
	if err != nil {
		return err
	}
	// 停止心跳
	p.lock.Lock()

	stop := p.stopManager[opts.Name]
	stop <- struct{}{}
	delete(p.stopManager, opts.Name)

	p.lock.Unlock()
	// 删除value
	return nil
}

func (p *Election) setValue(ctx context.Context, path string, value string, ttl time.Duration) error {
	fun := "Election.setValue"

	// 设置成功返回nil，不存在，已有值返回err
	if err := p.setPrevValue(ctx, path, value, ttl); err == nil {
		return nil
	}

	for {
		// key不存在则设置值
		err := p.setPrevNoExist(ctx, path, value, ttl)
		if err == nil {
			return nil
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		r, err := p.Client.Get(ctx, path, &etcd.GetOptions{})
		if err != nil {
			xlog.Warnf(ctx, "%s little rate get check path:%s resp:%v err:%v", fun, path, r, err)
			continue
		}

		wop := &etcd.WatcherOptions{
			AfterIndex: r.Index,
		}
		watcher := p.Client.Watcher(path, wop)
		if watcher == nil {
			xlog.Errorf(ctx, "%s get watcher get check path:%s err:%v", fun, path, err)
			return fmt.Errorf("get wather err")
		}

		r, err = watcher.Next(ctx)
		xlog.Infof(ctx, "%s watchnext check path:%s resp:%v err:%v", fun, path, r, err)
	}
}

func (p *Election) delValue(ctx context.Context, path string, value string) error {
	return p.delPrevValue(ctx, path, value)
}

func (p *Election) heart(ctx context.Context, path, value string, ttl time.Duration, stop chan struct{}) {
	fun := "Election.heart"
	tick := time.NewTicker(ttl / 3.0)

	for {
		select {
		case <-tick.C:
			xlog.Infof(ctx, "%s heart check path:%s", fun, path)
			p.refresh(ctx, path, value, ttl)
		case <-stop:
			xlog.Infof(ctx, "%s stop path:%s", fun, path)
			err := p.delPrevValue(ctx, path, value)
			if err != nil {
				xlog.Errorf(ctx, "%s del value, path: %s, value: %s", fun, path, value)
			}
			return
		}
	}
}

func (p *Election) refresh(ctx context.Context, path, value string, ttl time.Duration) {
	fun := "Election.refresh"
	// 刷新节点ttl，
	gr, err := p.Client.Get(ctx, path, &etcd.GetOptions{
		Recursive: true,
	})
	if err != nil {
		xlog.Errorf(ctx, "%s get path value, err: %v", fun, err)
		return
	}
	// 节点值已经变换，表示当前节点已经不是leader，不能再继续刷新，会直接退出程序
	if gr.Node.Value != value {
		xlog.Errorf(ctx, "%s leader change, %s", fun, gr.Node.Value)
		return
	}
	r, err := p.Client.Set(context.Background(), path, "", &etcd.SetOptions{
		PrevExist: etcd.PrevExist,
		TTL:       ttl,
		Refresh:   true,
	})

	if err != nil {
		xlog.Errorf(ctx, "%s noexist heart path: %s resp: %v err: %v", fun, path, r, err)
	} else {
		xlog.Infof(ctx, "%s noexist heartpath: %s resp: %v", fun, path, r)
	}
}

func (p *Election) setPrevValue(ctx context.Context, path, value string, ttl time.Duration) error {
	fun := "Election.setPrevValue"
	// 指定节点值必须是value才可设置成功。设置为空，则忽略当前值
	res, err := p.Client.Set(ctx, path, value, &etcd.SetOptions{
		PrevValue: value,
		TTL:       ttl,
	})
	if err != nil {
		// 已经存在
		xlog.Infof(ctx, "%s exist check path: %s resp: %v err: %v", fun, path, res, err)
	} else {
		xlog.Infof(ctx, "%s set value success", fun)
	}
	return err
}

func (p *Election) setPrevNoExist(ctx context.Context, path, value string, ttl time.Duration) error {
	fun := "Election.setPrevNoExist"
	// 指定节点值必须不存在
	res, err := p.Client.Set(ctx, path, value, &etcd.SetOptions{
		PrevExist: etcd.PrevNoExist,
		TTL:       ttl,
	})
	if err != nil {
		xlog.Infof(ctx, "%s noexist check path: %s resp: %v err: %v", fun, path, res, err)
	} else {
		xlog.Infof(ctx, "%s noexist check path: %s resp: %v", fun, path, res)
	}
	return err
}

func (p *Election) delPrevValue(ctx context.Context, path, value string) error {
	fun := "Election.delPrevValue"
	r, err := p.Client.Delete(ctx, path, &etcd.DeleteOptions{
		PrevValue: value,
	})
	if err != nil {
		xlog.Errorf(ctx, "%s unlock path: %s resp: %v err: %v", fun, path, r, err)
	} else {
		xlog.Infof(ctx, "%s unlock path: %s resp: %v", fun, path, r)
	}

	return err
}

func (p *Election) startHeart(ctx context.Context, opts *ElectionOptions) {
	stop := make(chan struct{})
	// 拿到锁后保持心跳
	go p.heart(context.Background(), p.path(opts.Name), opts.Value, opts.TTL, stop)

	p.lock.Lock()
	defer p.lock.Unlock()
	p.stopManager[opts.Name] = stop
	return
}

/*
	关闭选主组件

	其会关闭使用当前组件发起的所有选举活动
*/
func (p *Election) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, stop := range p.stopManager {
		stop <- struct{}{}
	}
}

func getValue() (string, error) {
	return os.Hostname()
}

func (p *Election) path(name string) string {
	// /seaweed/lock/{ServiceName}/{name}
	return fmt.Sprintf("%s/%s/%s", LockPrefixPath, p.serviceName, name)
}

func validateOptions(opts *ElectionOptions) error {
	if opts.Name == "" {
		return ERRNameNull
	}
	if opts.Value == "" {
		value, err := getValue()
		if err != nil {
			return fmt.Errorf("new election get default value, err: %v", err)
		}
		opts.Value = value
	}
	if opts.TTL == 0 {
		opts.TTL = defaultTTL
	}
	// 这块暂时用字符串，不引入roc依赖，同roc逻辑
	if opts.Lane == "pre" {
		return ERRPreLane
	}
	return nil
}
