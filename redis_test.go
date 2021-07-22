package main

import (
	"testing"

	"github.com/alicebob/miniredis/v2"

	"github.com/stretchr/testify/assert"
	redis "gopkg.in/redis.v5"
)

func newTestRedis() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}
func TestInitTarget(t *testing.T) {
	mclient := &Client{newTestRedis()}
	coordinates := mclient.InitTarget()

	assert.True(t, coordinates.X < 21)
	assert.True(t, coordinates.Y < 21)
	assert.True(t, coordinates.X >= 0)
	assert.True(t, coordinates.Y >= 0)
}

func TestInitUser(t *testing.T) {
	mclient := Client{newTestRedis()}
	coordinates := mclient.InitUser()

	assert.Equal(t, 10, coordinates.X)
	assert.Equal(t, 10, coordinates.Y)
}

func TestGetTarget(t *testing.T) {
	mclient := Client{newTestRedis()}
	setTarget(Coordinates{X: 10, Y: 10}, mclient)
	coordinates := mclient.GetTarget()
	assert.Equal(t, 10, coordinates.X)
	assert.Equal(t, 10, coordinates.Y)
}

func TestGetUser(t *testing.T) {
	mclient := Client{newTestRedis()}
	setUser(Coordinates{X: 9, Y: 8}, mclient)
	coordinates := mclient.GetUser()
	assert.Equal(t, 9, coordinates.X)
	assert.Equal(t, 8, coordinates.Y)
}

func TestShot(t *testing.T) {
	mclient := Client{newTestRedis()}
	coordinates := Coordinates{X: 10, Y: 10}
	setTarget(coordinates, mclient)
	assert.Equal(t, "miss", mclient.Shot(Coordinates{X: 0, Y: 0}))
	assert.Equal(t, "touch", mclient.Shot(coordinates))
	assert.Equal(t, "touch", mclient.Shot(coordinates))
	assert.Equal(t, "kill", mclient.Shot(coordinates))
}

func TestMove(t *testing.T) {
	mclient := Client{newTestRedis()}

	_ = mclient.InitUser()
	coordinates := mclient.Move("up")
	assert.Equal(t, 9, coordinates.Y)

	_ = mclient.InitUser()
	coordinates = mclient.Move("down")
	assert.Equal(t, 11, coordinates.Y)

	_ = mclient.InitUser()
	coordinates = mclient.Move("left")
	assert.Equal(t, 9, coordinates.X)

	_ = mclient.InitUser()
	coordinates = mclient.Move("right")
	assert.Equal(t, 11, coordinates.X)

	_ = setUser(Coordinates{X: 20, Y: 20}, mclient)
	coordinates = mclient.Move("down")
	assert.Equal(t, 20, coordinates.Y)
	coordinates = mclient.Move("right")
	assert.Equal(t, 20, coordinates.X)

	_ = setUser(Coordinates{X: 0, Y: 0}, mclient)
	coordinates = mclient.Move("up")
	assert.Equal(t, 0, coordinates.Y)
	coordinates = mclient.Move("left")
	assert.Equal(t, 0, coordinates.X)
}
