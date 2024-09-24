package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"strings"
	"user-service/database"
	"user-service/models"
)

var ctx = context.Background()

// CacheUserProfile caches the full user profile in Redis using userID
func CacheUserProfile(user *models.SearchedUser) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = database.RedisClient.HSet(ctx, "user_profiles", fmt.Sprintf("%d", user.ID), userJSON).Err()
	if err != nil {
		log.Printf("Error caching user profile: %v", err)
		return err
	}
	return nil
}

// GetAutocompleteKeys fetches matching user keys from the sorted set using ZRANGEBYSCORE
func GetAutocompleteKeys(query string) ([]string, error) {
	// Define score range (0 for usernames, negative for full names)
	rangeBy := &redis.ZRangeBy{
		Min: "-inf", // Get all elements with any score (fullnames and usernames)
		Max: "+inf",
	}

	// Fetch matching user keys from the sorted set using ZRANGEBYSCORE
	autocompleteKeys, err := database.RedisClient.ZRangeByScore(ctx, "users_autocomplete", rangeBy).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Printf("No cached results found for query: %s", query)
			return nil, errors.New("no cached results found in Redis")
		}
		log.Printf("Error fetching autocomplete keys from Redis: %v", err)
		return nil, err
	}

	log.Printf("Fetched keys: %v", autocompleteKeys)

	// Filter the results to match the query (prefix search)
	var filteredKeys []string
	for _, key := range autocompleteKeys {
		if strings.HasPrefix(key, query) {
			filteredKeys = append(filteredKeys, key)
		}
	}

	log.Printf("Filtered autocomplete keys: %v", filteredKeys)
	return filteredKeys, nil
}

// GetUserIDsFromKey fetches the userIDs (array) from the Redis hash (user_id_map)
func GetUserIDsFromKey(redisKey string) ([]uint, error) {
	userIDsJSON, err := database.RedisClient.HGet(ctx, "user_id_map", redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Printf("No user IDs found for key: %s", redisKey)
			return nil, errors.New("no user IDs found")
		}
		log.Printf("Error fetching user IDs from Redis: %v", err)
		return nil, err
	}

	// Unmarshal the JSON array of user IDs into a slice of integers
	var userIDs []uint // Or []int if you prefer int for user IDs
	err = json.Unmarshal([]byte(userIDsJSON), &userIDs)
	if err != nil {
		log.Printf("Error unmarshaling user IDs: %v", err)
		return nil, err
	}

	return userIDs, nil
}

// GetUserProfileByID fetches the user profile from 'user_profiles' hash by userID
func GetUserProfileByID(userID uint) (*models.SearchedUser, error) {
	cachedUser, err := database.RedisClient.HGet(ctx, "user_profiles", strconv.Itoa(int(userID))).Result()
	if errors.Is(err, redis.Nil) {
		log.Printf("No cached user profile found for userID: %s", userID)
		return nil, errors.New("no cached user profile found")
	} else if err != nil {
		log.Printf("Error fetching user profile from Redis: %v", err)
		return nil, err
	}

	var user models.SearchedUser
	err = json.Unmarshal([]byte(cachedUser), &user)
	if err != nil {
		log.Printf("Error unmarshaling cached user profile: %v", err)
		return nil, err
	}

	return &user, nil
}

// RemoveFromAutocomplete removes a name from the Redis sorted set for autocomplete
func RemoveFromAutocomplete(name string) error {
	err := database.RedisClient.ZRem(ctx, "users_autocomplete", name).Err()
	if err != nil {
		log.Printf("Error removing name from autocomplete: %v", err)
		return err
	}
	return nil
}

// RemoveUserMapping removes a name-to-userID mapping from Redis
func RemoveUserMapping(name string) error {
	err := database.RedisClient.HDel(ctx, "user_id_map", name).Err()
	if err != nil {
		log.Printf("Error removing user mapping from Redis: %v", err)
		return err
	}
	return nil
}

// AddToAutocompleteWithPriority adds a name to the Redis sorted set for autocomplete with priority (score)
func AddToAutocompleteWithPriority(name string, priorityScore float64) error {
	err := database.RedisClient.ZAdd(ctx, "users_autocomplete", &redis.Z{
		Score:  priorityScore,
		Member: name,
	}).Err()
	if err != nil {
		log.Printf("Error adding name to autocomplete with priority: %v", err)
		return err
	}
	return nil
}

// AddUserToMap adds a userID to the user_id_map hash, handling multiple IDs for the same name
func AddUserToMap(name string, userID uint) error {
	// Fetch existing IDs for this name
	existingIDs, err := database.RedisClient.HGet(ctx, "user_id_map", name).Result()
	var userIDs []uint
	if err == redis.Nil {
		userIDs = []uint{}
	} else if err != nil {
		return err
	} else {
		err = json.Unmarshal([]byte(existingIDs), &userIDs)
		if err != nil {
			return err
		}
	}

	// Add the new userID if it's not already in the list
	for _, id := range userIDs {
		if id == userID {
			return nil // userID already exists
		}
	}
	userIDs = append(userIDs, userID)

	// Store the updated list of userIDs back in Redis
	userIDsJSON, err := json.Marshal(userIDs)
	if err != nil {
		return err
	}

	err = database.RedisClient.HSet(ctx, "user_id_map", name, userIDsJSON).Err()
	if err != nil {
		log.Printf("Error adding userID to map: %v", err)
		return err
	}
	return nil
}
