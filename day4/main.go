package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

// Constants for the game
const (
    width  = 20
    height = 10
    foodChar = '@'
    snakeChar = 'O'
    emptyChar = ' '
    directionUp = iota
    directionDown
    directionLeft
    directionRight
)

// Snake represents the snake's body
type Snake struct {
    body []Coordinate
    direction int
}

// Coordinate represents a position on the grid
type Coordinate struct {
    x int
    y int
}

// GameState holds the game's data
type GameState struct {
    snake Snake
    food Coordinate
    score int
    gameOver bool
}

// InitializeGame sets up the initial game state
func InitializeGame() GameState {
    snake := Snake{
        body: []Coordinate{{x: width / 2, y: height / 2}},
        direction: directionRight,
    }

    food := Coordinate{x: rand.Intn(width), y: rand.Intn(height)}

    return GameState{
        snake: snake,
        food: food,
        score: 0,
        gameOver: false,
    }
}

// GenerateFood generates a new food coordinate that isn't on the snake
func GenerateFood(snake Snake) Coordinate {
    for {
        food := Coordinate{x: rand.Intn(width), y: rand.Intn(height)}
        onSnake := false
        for _, bodyPart := range snake.body {
            if bodyPart.x == food.x && bodyPart.y == food.y {
                onSnake = true
                break
            }
        }
        if !onSnake {
            return food
        }
    }
}

// MoveSnake updates the snake's position based on its direction
func (gs *GameState) MoveSnake() {
    head := gs.snake.body[0]
    newHead := Coordinate{}

    switch gs.snake.direction {
    case directionUp:
        newHead = Coordinate{x: head.x, y: (head.y - 1 + height) % height}
    case directionDown:
        newHead = Coordinate{x: head.x, y: (head.y + 1) % height}
    case directionLeft:
        newHead = Coordinate{x: (head.x - 1 + width) % width, y: head.y}
    case directionRight:
        newHead = Coordinate{x: (head.x + 1) % width, y: head.y}
    }

    gs.snake.body = append([]Coordinate{newHead}, gs.snake.body[:len(gs.snake.body)-1]...)

    // Check for collision with food
    if newHead.x == gs.food.x && newHead.y == gs.food.y {
        gs.score++
        gs.food = GenerateFood(gs.snake)
    }
}

// CheckCollision checks if the snake has collided with itself or the walls
func (gs *GameState) CheckCollision() bool {
    // Check for wall collision
    if gs.snake.body[0].x < 0 || gs.snake.body[0].x >= width || gs.snake.body[0].y < 0 || gs.snake.body[0].y >= height {
        return true
    }

    // Check for self-collision
    for i := 1; i < len(gs.snake.body); i++ {
        if gs.snake.body[0].x == gs.snake.body[i].x && gs.snake.body[0].y == gs.snake.body[i].y {
            return true
        }
    }

    return false
}

// DrawGame renders the game board to the console
func (gs *GameState) DrawGame() {
    // Clear the screen (platform-dependent)
    clearScreen()

    fmt.Printf("Score: %d\n", gs.score)

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            var char string
            if x == gs.food.x && y == gs.food.y {
                char = string(foodChar)
            } else {
                onSnake := false
                for _, bodyPart := range gs.snake.body {
                    if bodyPart.x == x && bodyPart.y == y {
                        onSnake = true
                        break
                    }
                }
                if onSnake {
                    char = string(snakeChar)
                } else {
                    char = string(emptyChar)
                }
            }
            fmt.Print(char)
        }
        fmt.Println()
    }
}

// clearScreen clears the console screen.  Platform-dependent.
func clearScreen() {
    switch os.Getenv("OS") {
    case "windows":
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
    default:
        exec.Command("clear").Run()
    }
}

// HandleInput updates the snake's direction based on user input
func (gs *GameState) HandleInput() {
    var input string
    fmt.Scanln(&input)

    switch input {
    case "w":
        if gs.snake.direction != directionDown {
            gs.snake.direction = directionUp
        }
    case "s":
        if gs.snake.direction != directionUp {
            gs.snake.direction = directionDown
        }
    case "a":
        if gs.snake.direction != directionRight {
            gs.snake.direction = directionLeft
        }
    case "d":
        if gs.snake.direction != directionLeft {
            gs.snake.direction = directionRight
        }
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())

    game := InitializeGame()

    for !game.gameOver {
        game.DrawGame()
        game.HandleInput()
        game.MoveSnake()
        if game.CheckCollision() {
            game.gameOver = true
        }
        time.Sleep(time.Duration(100) * time.Millisecond) // Adjust speed
    }

    fmt.Println("Game Over!")
    fmt.Printf("Final Score: %d\n", game.score)
}
