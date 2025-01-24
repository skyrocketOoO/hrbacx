package usecase

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisUsecase struct {
	SessionPool *redis.Client
}

func NewRedisUsecase(sessionPool *redis.Client) *RedisUsecase {
	return &RedisUsecase{
		SessionPool: sessionPool,
	}
}

func (u *RedisUsecase) AddLeader(leaderID string, roleID string) error {
	ctx := context.Background()
	return u.SessionPool.SAdd(ctx, "role_"+leaderID, "role_"+roleID).Err()
}

// AssignPermission assigns a permission to a role for an object.
func (u *RedisUsecase) AssignPermission(objectID, permissionType, roleID string) error {
	ctx := context.Background()
	return u.SessionPool.SAdd(ctx, "role_"+roleID, "obj_"+objectID).Err()
}

// AssignRole assigns a role to a user.
func (u *RedisUsecase) AssignRole(userID, roleID string) error {
	ctx := context.Background()
	// Using SADD to associate userID with roleID
	return u.SessionPool.SAdd(ctx, "user_"+userID, "role_"+roleID).Err()
}

// CheckPermission checks if a user has the specified permission on an object.
func (u *RedisUsecase) CheckPermission(userID, permissionType, objectID string) (
	ok bool, err error,
) {
	userID = "user_" + userID
	objectID = "obj_" + objectID
	// Initialize queue and visited map
	queue := []string{}
	visited := make(map[string]bool)
	var current string

	members, err := u.SessionPool.SMembers(context.Background(), userID).Result()
	for _, member := range members {
		queue = append(queue, member)
	}

	for len(queue) > 0 {
		current = queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}

		visited[current] = true

		ok, err := u.SessionPool.SIsMember(context.Background(), current, objectID).Result()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}

		members, err := u.SessionPool.SMembers(context.Background(), current).Result()
		if err != nil {
			return false, err
		}

		for _, member := range members {
			if len(member) >= 5 && member[:5] == "role_" && visited[member] == false {
				queue = append(queue, member)
			}
		}
	}

	// If no permission found, return false
	return false, nil
}

// ClearAll clears all Redis data related to roles, permissions, etc.
func (u *RedisUsecase) ClearAll() error {
	ctx := context.Background()
	return u.SessionPool.FlushAll(ctx).Err()
}
