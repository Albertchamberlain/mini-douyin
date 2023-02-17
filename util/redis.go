package util

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
)

type redisUtil struct {
	server       *redis.Client
	ctx          context.Context
	requestGroup singleflight.Group
	// singleflight类的使用方法就新建一个singleflight.Group，使用其方法Do或者DoChan来包装方法，
	// 被包装的方法在对于同一个key，只会有一个协程执行，其他协程等待那个协程执行结束后，拿到同样的结果。
}

var Redis *redisUtil
var RedisClient1 *redis.Client

func InitRedis(options *redis.Options) {
	Redis = new(redisUtil)
	Redis.server = redis.NewClient(options)
	Redis.ctx = context.TODO()
	RedisClient1 = Redis.server
	if _, err := Redis.server.Ping(Redis.ctx).Result(); err != nil {
		panic(err)
	}
}

func (r *redisUtil) Like(userId uint, videoId uint) error {
	if err := r.server.ZAdd(r.ctx, videoLikedKey(videoId), newZWithNowTime(userId)).Err(); err != nil {
		return err
	}
	if err := r.server.ZAdd(r.ctx, userLikeKey(userId), newZWithNowTime(videoId)).Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisUtil) Unlike(userId uint, videoId uint) error {
	if err := r.server.ZRem(r.ctx, videoLikedKey(videoId), userId).Err(); err != nil {
		return err
	}
	if err := r.server.ZRem(r.ctx, userLikeKey(userId), videoId).Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisUtil) CountLiked(videoId uint) (int64, error) {
	return r.server.ZCard(r.ctx, videoLikedKey(videoId)).Result()
}

func (r *redisUtil) CountLike(userId uint) (int64, error) {
	return r.server.ZCard(r.ctx, userLikeKey(userId)).Result()
}

func (r *redisUtil) RangeLiked(videoId uint, start, stop int64) ([]*uint64, error) {
	return r.range0(videoLikedKey(videoId), start, stop)
}

func (r *redisUtil) RangeLike(userId uint, start, stop int64) ([]*uint64, error) {
	return r.range0(userLikeKey(userId), start, stop)
}

func (r *redisUtil) range0(key string, start, stop int64) ([]*uint64, error) {
	data, err := r.server.ZRange(r.ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	res := make([]*uint64, len(data))
	for i, v := range data {
		u64, _ := strconv.ParseUint(v, 10, 32)
		res[i] = &u64
	}
	return res, nil
}

func (r *redisUtil) IsLike(videoId, userId uint) (bool, error) {
	_, err := r.server.ZRank(r.ctx, videoLikedKey(videoId), strconv.FormatUint(uint64(userId), 10)).Result()
	return err != redis.Nil, err
}

func (r *redisUtil) Follow(who uint, target uint) error {
	if err := r.server.ZAdd(r.ctx, userFollowerKey(target), newZWithNowTime(who)).Err(); err != nil {
		return err
	}
	if err := r.server.ZAdd(r.ctx, userFolloweeKey(who), newZWithNowTime(target)).Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisUtil) UnFollow(who uint, target uint) error {
	if err := r.server.ZRem(r.ctx, userFollowerKey(target), who).Err(); err != nil {
		return err
	}
	if err := r.server.ZRem(r.ctx, userFolloweeKey(who), target).Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisUtil) CountFollower(id uint) (int64, error) {
	return r.server.ZCard(r.ctx, userFollowerKey(id)).Result()
}

func (r *redisUtil) CountFollowee(id uint) (int64, error) {
	return r.server.ZCard(r.ctx, userFolloweeKey(id)).Result()
}

func (r *redisUtil) RangeFollower(id uint, start, stop int64) ([]*uint64, error) {
	return r.range0(userFollowerKey(id), start, stop)
}

func (r *redisUtil) RangeFollowee(id uint, start, stop int64) ([]*uint64, error) {
	return r.range0(userFolloweeKey(id), start, stop)
}

func (r *redisUtil) IsFollow(who uint, target uint) (bool, error) {
	_, err := r.server.ZRank(r.ctx, userFolloweeKey(who), strconv.FormatUint(uint64(target), 10)).Result()
	return err != redis.Nil, err
}

func (r *redisUtil) IncrComment(id uint) error {
	return r.server.Incr(r.ctx, commentKey(id)).Err()
}

func (r *redisUtil) DecrComment(id uint) error {
	return r.server.Decr(r.ctx, commentKey(id)).Err()
}

func (r *redisUtil) CountComment(id uint) (interface{}, error) {
	v, err, _ := r.requestGroup.Do(string(rune(id)), func() (interface{}, error) {
		return r.server.Get(r.ctx, commentKey(id)).Int64()
	})
	if err != nil {
		return v, err
	}
	return v, err
}

func newZWithNowTime(member interface{}) *redis.Z {
	return &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: member,
	}
}

const (
	cacheName        = "douyin"
	userLikeDomain   = "user_like"
	videoLikedDomain = "video_like"
	followerDomain   = "follower"
	followeeDomain   = "followee"
	commentDomain    = "comment"
)

func videoLikedKey(videoId uint) string {
	return cacheName + ":" + videoLikedDomain + ":" + strconv.FormatUint(uint64(videoId), 10)
}

func userLikeKey(userId uint) string {
	return cacheName + ":" + userLikeDomain + ":" + strconv.FormatUint(uint64(userId), 10)
}

func userFollowerKey(id uint) string {
	return cacheName + ":" + followerDomain + ":" + strconv.FormatUint(uint64(id), 10)
}

func userFolloweeKey(id uint) string {
	return cacheName + ":" + followeeDomain + ":" + strconv.FormatUint(uint64(id), 10)
}

func commentKey(videoId uint) string {
	return cacheName + ":" + commentDomain + ":" + strconv.FormatUint(uint64(videoId), 10)
}
