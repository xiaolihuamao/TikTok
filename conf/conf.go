/*
author:fuxingyuan
此文件修改需联系作者
*/
package conf

// Secret 密钥
var Secret = "1234abcd"

const ValidComment = 0   //评论状态：有效
const InvalidComment = 1 //评论状态：取消
const DateTime = "2006-01-02 15:04:05"

//const ChanCapacity = 10 //chan管道容量，暂时没定

const IsLike = 0     //点赞的状态
const Unlike = 1     //取消赞的状态
const LikeAction = 1 //点赞的行为
const Attempts = 3   //操作数据库的最大尝试次数

const DefaultRedisValue = -1 //redis中key对应的预设值，防脏读

const IPAndPort = "http://192.168.137.1:8081"
