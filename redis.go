package main

import (
	"math/rand"
	"strconv"

	redis "gopkg.in/redis.v5"
)

// Client is the Redis client structure
type Client struct {
	*redis.Client
}

// Cache interface
type Cache interface {
	InitUser() Coordinate
	InitTarget() Coordinate
	GetTarget() Coordinate
	GetUser() Coordinate
	Shot(coordinate Coordinate) string
	Move(direction string) Coordinate
}

func getCoordinates(pattern string, client Client) Coordinate {
	data, _ := client.Get(pattern + "x").Result()
	x, _ := strconv.Atoi(data)
	data, _ = client.Get(pattern + "y").Result()
	y, _ := strconv.Atoi(data)

	return Coordinate{X: x, Y: y}
}

// GetUser returns the user coordinates
func (client Client) GetUser() Coordinate {
	return getCoordinates("pccg:user:", client)
}

// GetTarget returns the target coordinates
// Move moves the user in the specified direction
func (client Client) GetTarget() Coordinate {
	return getCoordinates("pccg:target:", client)
}

// Move moves the user in the specified direction
func (client Client) Move(direction string) Coordinate {
	coordinates := client.GetUser()
	x := coordinates.X
	y := coordinates.Y

	switch direction {
	case "up":
		if y > 0 {
			y = y - 1
		}
	case "down":
		if y < 20 {
			y = y + 1
		}
	case "left":
		if x > 0 {
			x = x - 1
		}
	case "right":
		if x < 20 {
			x = x + 1
		}
	}

	_ = client.Set("pccg:user:x", x, 0)
	_ = client.Set("pccg:user:y", y, 0)

	return Coordinate{X: x, Y: y}
}

// InitUser set the starting coordinates for the user & reset the score
func (client Client) InitUser() Coordinate {
	coordinate := Coordinate{X: 10, Y: 10}
	_ = client.Del("pccg:user:score").Err()
	_ = client.Set("pccg:user:x", coordinate.X, 0)
	_ = client.Set("pccg:user:y", coordinate.Y, 0)
	return coordinate
}

// InitTarget set the starting coordinates for the target
func (client Client) InitTarget() Coordinate {
	coordinate := Coordinate{X: rand.Intn(21), Y: rand.Intn(21)}
	_ = client.Set("pccg:target:x", coordinate.X, 0)
	_ = client.Set("pccg:target:y", coordinate.Y, 0)
	return coordinate
}

// Shot shots at the target
func (client Client) Shot(coordinate Coordinate) string {
	position := getCoordinates("pccg:target:", client)
	if position.X != coordinate.X || position.Y != coordinate.Y {
		return "miss"
	}

	data, _ := client.Get("pccg:user:score").Result()
	score, _ := strconv.Atoi(data)
	if score == 2 {
		return "kill"
	}

	_ = client.Set("pccg:user:score", score+1, 0)

	return "touch"
}

// New creates a new Redis Client
func NewRedis(URL string) (*Client, error) {
	options, err := redis.ParseURL(URL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}
