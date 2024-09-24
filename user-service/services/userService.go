package services

import (
	"log"
	"strings"
	"user-service/models"
	"user-service/repositories"
)

func SearchUsers(query string) ([]models.SearchedUser, error) {
	var users []models.SearchedUser

	// Fetch autocomplete keys (prioritized usernames/fullnames)
	autocompleteKeys, err := repositories.GetAutocompleteKeys(query)
	if err != nil {
		return nil, err
	}

	// Log the cached autocomplete keys for debugging
	log.Printf("Autocomplete keys: %v", autocompleteKeys)

	// Fetch user profiles for each key
	for _, redisKey := range autocompleteKeys {
		// Fetch the userIDs (array) from the hash (since there can be multiple IDs for non-unique names)
		userIDs, err := repositories.GetUserIDsFromKey(redisKey)
		if err != nil {
			log.Printf("Error fetching userIDs for key: %s", redisKey)
			continue
		}

		// Iterate over each userID and fetch the corresponding user profile
		for _, userID := range userIDs {
			user, err := repositories.GetUserProfileByID(userID)
			if err != nil {
				log.Printf("Error fetching user profile for userID: %s", userID)
				continue
			}

			// Append the user to the result slice
			users = append(users, *user)
		}
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

	// Store the old username and fullname before updating
	oldUsername := user.Username
	oldFullname := user.FullName

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

	// Remove old username and fullname from Redis if they have changed
	if oldUsername != "" && oldUsername != user.Username {
		err = repositories.RemoveFromAutocomplete(strings.ToLower(oldUsername))
		if err != nil {
			return err
		}
		err = repositories.RemoveUserMapping(strings.ToLower(oldUsername))
		if err != nil {
			return err
		}
	}

	if oldFullname != "" && oldFullname != user.FullName {
		err = repositories.RemoveFromAutocomplete(strings.ToLower(oldFullname))
		if err != nil {
			return err
		}
		err = repositories.RemoveUserMapping(strings.ToLower(oldFullname))
		if err != nil {
			return err
		}
	}

	// Add new username to Redis sorted set for autocomplete with higher priority (score = 0)
	err = repositories.AddToAutocompleteWithPriority(strings.ToLower(user.Username), 1)
	if err != nil {
		return err
	}

	// Add fullname to Redis sorted set with lower priority (score = -1)
	if user.FullName != "" {
		err = repositories.AddToAutocompleteWithPriority(strings.ToLower(user.FullName), 0)
		if err != nil {
			return err
		}
	}

	// Map username and fullname to userID in Redis (allow multiple IDs)
	err = repositories.AddUserToMap(strings.ToLower(user.Username), user.ID)
	if err != nil {
		return err
	}

	if user.FullName != "" {
		err = repositories.AddUserToMap(strings.ToLower(user.FullName), user.ID)
		if err != nil {
			return err
		}
	}

	// Cache the full user profile in Redis
	searchedUser := models.SearchedUser{
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		ProfileImage: user.ProfileImage,
	}

	err = repositories.CacheUserProfile(&searchedUser)
	if err != nil {
		return err
	}

	return nil
}

// GetUserProfile returns a user by ID
func GetUserProfile(id float64) (*models.User, error) {
	return repositories.FindUserByIDStr(id)
}
