package data

import (
	"github.com/coreos/etcd/client"
	"sync"
)

var (
	SeckillConfCtx = newSeckillConf()
)

const (
	ProductStatusNormal       = 0 //商品正常状态
	ProductStatusSaleout      = 1 //商品售罄
	ProductStatusForceSaleout = 2 //商品强制售罄
)



//etcd config
type EtcdConf struct {
	EtcdConn          *client.Client //连接
	EtcdSecProductKey string         //商品键
}

//visit limit
type AccessLimitConf struct {
	IPSecAccessLimit   int //每秒IP限制数量
	UserSecAccessLimit int // 每秒用户限制数量
	IpMinAccessLimit   int //每分钟IP限制数量
	UserMinAccessLimit int //每分钟用户限制数量
}

//goods info config
type SecProductInfoConf struct {
	ProductId int   //商品ID
	StartTime int64 //开始时间
	EndTime   int64 //结束时间
	Status    int   //商品状态
	Total     int   //总计
	Left      int   //剩余
}

//sec result config
type SecResult struct {
	ProductId int    //商品ID
	UserId    int    //用户ID
	Token     string //token
	TokenTime int64  //token生成时间
	Code      int    //状态码
}

//request config
type SecRequest struct {
	ProductId     int // 商品ID
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	UserId        int
	UserAuthSign  string //用户授权签名
	AccessTime    int64
	ClientAddr    string
	ClientRefence string
	CloseNotify   <-chan bool
	ResultChan    chan *SecResult
}

//sec kill config
type SecKillConf struct {
	RedisConf *RedisConf
	EtcdConf  *EtcdConf

	SecProductInfoMap  map[int]*SecProductInfoConf
	secProductInfoLock sync.RWMutex

	CookieSecretKey string

	ReferWhiteList []string //白名单

	IPBlackMap  map[string]bool
	IDBlackMap  map[int]bool
	RWBlackLock sync.RWMutex

	AccessLimitConf
	WriteProxyToLayerGoroutineNum int
	ReadProxyToLayerGoroutineNum  int

	SecReqChan     chan *SecRequest
	SecReqChanSize int

	UserConnMap  map[string]chan *SecResult
	UserConnLock sync.Mutex
}


func newSeckillConf() (*SecKillConf) {
	return &SecKillConf{
		IPBlackMap: make(map[string]bool,10000),
		IDBlackMap:make(map[int]bool,10000),
	}
}
