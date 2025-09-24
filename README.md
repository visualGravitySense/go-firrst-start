# Simple Platformer - First Go Game

This is a simple platformer game written in Go using the Ebiten library.

## Game Description

- **Objective**: Collect all coins by jumping on platforms
- **Controls**: 
  - `A` or `←` - move left
  - `D` or `→` - move right  
  - `Space` - jump
- **Physics**: Realistic gravity and collisions

## Installation and Running

1. Make sure you have Go installed (version 1.21 or higher)

2. Clone or download the project

3. Install dependencies:
```bash
go mod tidy
```

4. Run the game:
```bash
go run main.go
```

## Features

- Simple physics with gravity
- Collision system for platforms
- Coin collection with score counting
- Beautiful sky-blue background
- Red player and golden coins

## Code Structure

- `Player` - player structure with position and velocity
- `Platform` - platforms for jumping
- `Coin` - coins to collect
- `Game` - main game logic

This is an excellent first game for learning Go and game development!
