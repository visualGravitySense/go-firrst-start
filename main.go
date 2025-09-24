package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	playerSize   = 32
	platformSize = 100
)

// Player представляет игрока
type Player struct {
	x, y     float64
	vx, vy   float64 // скорость по x и y
	onGround bool
}

// Platform представляет платформу
type Platform struct {
	x, y, width, height float64
}

// Game представляет основную игру
type Game struct {
	player    *Player
	platforms []Platform
	coins     []Coin
	score     int
}

// Coin представляет монету для сбора
type Coin struct {
	x, y float64
	collected bool
}

// NewGame создает новую игру
func NewGame() *Game {
	player := &Player{
		x: 50,
		y: 400,
	}
	
	platforms := []Platform{
		{0, 550, 200, 50},      // земля слева
		{200, 450, 150, 20},    // платформа 1
		{400, 350, 150, 20},    // платформа 2
		{600, 250, 150, 20},    // платформа 3
		{600, 550, 200, 50},    // земля справа
	}
	
	coins := []Coin{
		{250, 400, false},  // монета на платформе 1
		{450, 300, false},  // монета на платформе 2
		{650, 200, false},  // монета на платформе 3
	}
	
	return &Game{
		player:    player,
		platforms: platforms,
		coins:     coins,
		score:     0,
	}
}

// Update обновляет состояние игры
func (g *Game) Update() error {
	// Обработка ввода
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.vx = -3
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.vx = 3
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.player.onGround {
		g.player.vy = -12 // прыжок
		g.player.onGround = false
	}
	
	// Применяем гравитацию
	g.player.vy += 0.5
	
	// Ограничиваем скорость падения
	if g.player.vy > 10 {
		g.player.vy = 10
	}
	
	// Обновляем позицию игрока
	g.player.x += g.player.vx
	g.player.y += g.player.vy
	
	// Замедляем горизонтальное движение
	g.player.vx *= 0.8
	
	// Проверяем коллизии с платформами
	g.player.onGround = false
	for _, platform := range g.platforms {
		if g.checkCollision(g.player.x, g.player.y, playerSize, playerSize,
			platform.x, platform.y, platform.width, platform.height) {
			
			// Если игрок падает на платформу сверху
			if g.player.vy > 0 && g.player.y < platform.y {
				g.player.y = platform.y - playerSize
				g.player.vy = 0
				g.player.onGround = true
			}
		}
	}
	
	// Проверяем сбор монет
	for i, coin := range g.coins {
		if !coin.collected {
			if g.checkCollision(g.player.x, g.player.y, playerSize, playerSize,
				coin.x, coin.y, 20, 20) {
				g.coins[i].collected = true
				g.score++
			}
		}
	}
	
	// Проверяем границы экрана
	if g.player.x < 0 {
		g.player.x = 0
	}
	if g.player.x > screenWidth-playerSize {
		g.player.x = screenWidth - playerSize
	}
	
	return nil
}

// Draw отрисовывает игру
func (g *Game) Draw(screen *ebiten.Image) {
	// Очищаем экран
	screen.Fill(color.RGBA{135, 206, 235, 255}) // небесно-голубой
	
	// Рисуем платформы
	for _, platform := range g.platforms {
		ebitenutil.DrawRect(screen, platform.x, platform.y, platform.width, platform.height, color.RGBA{139, 69, 19, 255})
	}
	
	// Рисуем монеты
	for _, coin := range g.coins {
		if !coin.collected {
			ebitenutil.DrawRect(screen, coin.x, coin.y, 20, 20, color.RGBA{255, 215, 0, 255})
		}
	}
	
	// Рисуем игрока
	ebitenutil.DrawRect(screen, g.player.x, g.player.y, playerSize, playerSize, color.RGBA{255, 0, 0, 255})
	
	// Рисуем счет
	ebitenutil.DebugPrint(screen, "Score: "+string(rune(g.score+'0')))
}

// Layout возвращает размеры экрана
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// checkCollision проверяет коллизию между двумя прямоугольниками
func (g *Game) checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func main() {
	game := NewGame()
	
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Simple Platformer - Go Game")
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
