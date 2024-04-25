package constants

const (
	LOGIN_CODE_KEY        = "login:code:"
	LOGIN_CODE_TTL        = 2
	LOGIN_USER_KEY        = "login:token:"
	LOGIN_USER_TTL        = 36000
	CACHE_NULL_TTL        = 2
	CACHE_SHOP_TTL        = 30
	CACHE_SHOP_KEY        = "cache:shop:"
	CACHE_SHOP_TYPE_KEY   = "cache:shoptype"
	LOCK_SHOP_KEY         = "lock:shop:"
	LOCK_SHOP_TTL         = 10
	SECKILL_STOCK_KEY     = "seckill:stock:"
	SECKILL               = "seckill:"
	BLOG_LIKED_KEY        = "blog:liked:"
	FEED_KEY              = "feed:"
	SHOP_GEO_KEY          = "shop:geo:"
	USER_SIGN_KEY         = "sign:"
	ACCESS_LIMIT_IP_KEY   = "access:limit:ip:"
	ACCESS_LIMIT_USER_KEY = "access:limit:user:"
	// SUBMIT_ORDER_TOKEN_KEY 提交订单令牌的缓存 key
	SUBMIT_ORDER_TOKEN_KEY = "order:submit:%s:%s"
	START_TIMESTAMP        = 1640995200
	COUNT_BITS             = 32
)
