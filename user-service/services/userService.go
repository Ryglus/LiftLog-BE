package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
	"user-service/database"
	"user-service/models"
	"user-service/repositories"
)

var ctx = context.Background()

func SearchUsers(query string) ([]models.SearchedUser, error) {
	var users []models.SearchedUser

	// Prepare the search range
	query = strings.ToLower(query)
	start := "[" + query        // Redis range start for ZRANGEBYLEX (inclusive)
	end := "[" + query + "\xff" // End range ensures all matches for the prefix

	// Fetch matching user keys from the sorted set using ZRANGEBYLEX
	autocompleteKeys, err := database.RedisClient.ZRangeByLex(ctx, "users_autocomplete", &redis.ZRangeBy{
		Min: start,
		Max: end,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			log.Printf("No cached results found for query: %s", query)
			return nil, errors.New("no cached results found in Redis")
		}
		log.Printf("Error fetching data from Redis: %v", err)
		return nil, err
	}

	// Log the cached autocomplete keys for debugging
	log.Printf("Autocomplete keys: %v", autocompleteKeys)

	// Fetch the user profiles from the 'user_profiles' hash using the userID from user_id_map
	for _, redisKey := range autocompleteKeys {
		// Get the userID from the Redis hash map (user_id_map) using the redisKey (username/fullname)
		userID, err := database.RedisClient.HGet(ctx, "user_id_map", redisKey).Result()
		if err == redis.Nil {
			log.Printf("No user ID found for key: %s", redisKey)
			continue
		} else if err != nil {
			log.Printf("Error fetching user ID from Redis: %v", err)
			continue
		}

		// Now use the userID to fetch the full user profile from 'user_profiles' hash
		cachedUser, err := database.RedisClient.HGet(ctx, "user_profiles", userID).Result()
		if err == redis.Nil {
			log.Printf("No cached user profile found for userID: %s", userID)
			continue
		} else if err != nil {
			log.Printf("Error fetching user profile from Redis: %v", err)
			continue
		}

		var user models.SearchedUser
		err = json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			log.Printf("Error unmarshaling cached user profile: %v", err)
			continue
		}

		// Append the user to the result slice
		users = append(users, user)
	}

	// Return the matching users
	return users, nil
}

// SearchProfiles searches for profiles by query (username or bio)
func SearchProfiles(query string) ([]models.User, error) {
	return repositories.SearchProfiles(query)
}

// UpdateUserProfile updates the user's profile based on the provided fields
func UpdateUserProfile(userID uint, username, bio, location, profileImage string) error {
	// Fetch the user from the database
	user, err := repositories.FindUserByID(userID)
	if err != nil {
		return err
	}

	// Update fields only if provided
	if username != "" {
		user.Username = strings.ToLower(username)
	}
	if bio != "" {
		user.Bio = bio
	}
	if location != "" {
		user.Location = location
	}
	if profileImage != "" {
		user.ProfileImage = profileImage
	}

	// Save the updated profile in the database
	err = repositories.UpdateUserProfile(user)
	if err != nil {
		return err
	}

	// Add username and fullname to Redis sorted set for autocomplete
	err = database.RedisClient.ZAdd(ctx, "users_autocomplete", &redis.Z{
		Score:  0,                              // Score doesn't matter for autocomplete
		Member: strings.ToLower(user.Username), // Store the username in autocomplete
	}).Err()
	if err != nil {
		return err
	}

	if user.FullName != "" {
		// Add fullname separately to Redis sorted set for autocomplete
		err = database.RedisClient.ZAdd(ctx, "users_autocomplete", &redis.Z{
			Score:  0,                              // Score doesn't matter for autocomplete
			Member: strings.ToLower(user.FullName), // Store the fullname in autocomplete
		}).Err()
		if err != nil {
			return err
		}
	}

	// Map username and fullname to the userID in Redis hash
	err = database.RedisClient.HSet(ctx, "user_id_map", strings.ToLower(user.Username), user.ID).Err()
	if err != nil {
		return err
	}

	if user.FullName != "" {
		err = database.RedisClient.HSet(ctx, "user_id_map", strings.ToLower(user.FullName), user.ID).Err()
		if err != nil {
			return err
		}
	}

	// Cache the full user profile in the 'user_profiles' hash using userID as the field
	searchedUser := models.SearchedUser{
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		ProfileImage: user.ProfileImage,
	}

	userJSON, err := json.Marshal(searchedUser)
	if err != nil {
		return err
	}

	// Store the full profile data using the userID as the field in a single hash
	err = database.RedisClient.HSet(ctx, "user_profiles", fmt.Sprintf("%d", user.ID), userJSON).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserProfile returns a user by ID
func GetUserProfile(id float64) (*models.User, error) {
	return repositories.FindUserByIDStr(id)
}
