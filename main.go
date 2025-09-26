package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1200
	screenHeight = 800
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
	player              *Player
	platforms           []Platform
	coins               []Coin
	enemies             []Enemy
	strongEnemies       []StrongEnemy
	bullets             []Bullet
	portal              Portal
	score               int
	lives               int
	gameWon             bool
	gameTime            float64 // время игры в секундах
	strongEnemiesActive bool    // флаг появления сильных врагов
}

// Coin представляет монету для сбора
type Coin struct {
	x, y      float64
	collected bool
}

// Enemy представляет врага
type Enemy struct {
	x, y                    float64
	vx                      float64 // скорость по x
	width, height           float64
	patrolLeft, patrolRight float64 // границы патрулирования
}

// Portal представляет портал
type Portal struct {
	x, y          float64
	width, height float64
	active        bool
}

// Bullet представляет пулю
type Bullet struct {
	x, y   float64
	vx, vy float64
	active bool
	width  float64
	height float64
}

// StrongEnemy представляет сильного врага
type StrongEnemy struct {
	x, y                    float64
	vx, vy                  float64 // скорость по x и y
	width, height           float64
	patrolLeft, patrolRight float64 // границы патрулирования
	health                  int     // здоровье
	lastShot                float64 // время последнего выстрела
}

// NewGame создает новую игру
func NewGame() *Game {
	player := &Player{
		x: 50,
		y: 700,
	}

	platforms := []Platform{
		// Основные платформы
		{0, 750, 300, 50},   // земля слева
		{300, 650, 200, 20}, // платформа 1
		{550, 550, 150, 20}, // платформа 2
		{750, 450, 150, 20}, // платформа 3
		{950, 350, 150, 20}, // платформа 4
		{1150, 250, 50, 20}, // платформа 5
		{900, 750, 300, 50}, // земля справа

		// Дополнительные платформы для сложности
		{400, 400, 100, 20},  // промежуточная платформа
		{700, 300, 100, 20},  // промежуточная платформа
		{1000, 200, 100, 20}, // промежуточная платформа
		{150, 500, 80, 20},   // маленькая платформа
		{850, 600, 80, 20},   // маленькая платформа
		{500, 200, 80, 20},   // высокая платформа
		{1100, 100, 80, 20},  // самая высокая платформа
	}

	coins := []Coin{
		// Основные монеты
		{350, 600, false},  // монета на платформе 1
		{600, 500, false},  // монета на платформе 2
		{800, 400, false},  // монета на платформе 3
		{1000, 300, false}, // монета на платформе 4
		{1175, 200, false}, // монета на платформе 5

		// Дополнительные монеты
		{450, 350, false},  // монета на промежуточной платформе
		{750, 250, false},  // монета на промежуточной платформе
		{1050, 150, false}, // монета на промежуточной платформе
		{190, 450, false},  // монета на маленькой платформе
		{890, 550, false},  // монета на маленькой платформе
		{540, 150, false},  // монета на высокой платформе
		{1140, 50, false},  // монета на самой высокой платформе

		// Секретные монеты
		{200, 300, false},  // секретная монета
		{1000, 500, false}, // секретная монета
	}

	enemies := []Enemy{
		// Основные враги
		{350, 600, 1, 30, 30, 300, 500},    // враг на платформе 1
		{600, 500, -1, 30, 30, 550, 700},   // враг на платформе 2
		{800, 400, 1, 30, 30, 750, 900},    // враг на платформе 3
		{1000, 300, -1, 30, 30, 950, 1100}, // враг на платформе 4

		// Дополнительные враги
		{450, 350, 1, 30, 30, 400, 500},  // враг на промежуточной платформе
		{750, 250, -1, 30, 30, 700, 800}, // враг на промежуточной платформе
		{190, 450, 1, 30, 30, 150, 230},  // враг на маленькой платформе
		{890, 550, -1, 30, 30, 850, 930}, // враг на маленькой платформе

		// Сложные враги
		{540, 150, 1, 30, 30, 500, 580},    // враг на высокой платформе
		{1140, 50, -1, 30, 30, 1100, 1180}, // враг на самой высокой платформе
	}

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Портал изначально неактивен
	portal := Portal{
		x: 0, y: 0, width: 40, height: 40, active: false,
	}

	return &Game{
		player:              player,
		platforms:           platforms,
		coins:               coins,
		enemies:             enemies,
		strongEnemies:       []StrongEnemy{},
		bullets:             []Bullet{},
		portal:              portal,
		score:               0,
		lives:               3,
		gameWon:             false,
		gameTime:            0,
		strongEnemiesActive: false,
	}
}

// Update обновляет состояние игры
func (g *Game) Update() error {
	// Обновляем время игры (60 FPS = 1/60 секунды за кадр)
	g.gameTime += 1.0 / 60.0

	// Активируем сильных врагов через минуту
	if g.gameTime >= 60.0 && !g.strongEnemiesActive {
		g.strongEnemiesActive = true
		g.spawnStrongEnemies()
	}

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

	// Стрельба (клавиша X)
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.shoot()
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

	// Обновляем врагов
	for i := range g.enemies {
		// Движение врага
		g.enemies[i].x += g.enemies[i].vx

		// Разворот на границах патрулирования
		if g.enemies[i].x <= g.enemies[i].patrolLeft || g.enemies[i].x >= g.enemies[i].patrolRight {
			g.enemies[i].vx = -g.enemies[i].vx
		}
	}

	// Обновляем сильных врагов
	for i := range g.strongEnemies {
		if g.strongEnemies[i].x >= 0 { // только живые враги
			// Движение врага
			g.strongEnemies[i].x += g.strongEnemies[i].vx

			// Разворот на границах патрулирования
			if g.strongEnemies[i].x <= g.strongEnemies[i].patrolLeft || g.strongEnemies[i].x >= g.strongEnemies[i].patrolRight {
				g.strongEnemies[i].vx = -g.strongEnemies[i].vx
			}
		}
	}

	// Обновляем пули
	for i := range g.bullets {
		if g.bullets[i].active {
			g.bullets[i].x += g.bullets[i].vx
			g.bullets[i].y += g.bullets[i].vy

			// Удаляем пули, которые вышли за границы экрана
			if g.bullets[i].x < 0 || g.bullets[i].x > screenWidth ||
				g.bullets[i].y < 0 || g.bullets[i].y > screenHeight {
				g.bullets[i].active = false
			}
		}
	}

	// Проверяем сбор монет
	allCoinsCollected := true
	for i, coin := range g.coins {
		if !coin.collected {
			allCoinsCollected = false
			if g.checkCollision(g.player.x, g.player.y, playerSize, playerSize,
				coin.x, coin.y, 20, 20) {
				g.coins[i].collected = true
				g.score++
			}
		}
	}

	// Активируем портал, если все монеты собраны
	if allCoinsCollected && !g.portal.active {
		g.portal.active = true
		// Размещаем портал в случайном месте на экране
		g.portal.x = float64(rand.Intn(screenWidth - int(g.portal.width)))
		g.portal.y = float64(rand.Intn(screenHeight - int(g.portal.height)))
	}

	// Проверяем коллизии пуль с врагами
	for i := range g.bullets {
		if g.bullets[i].active {
			// Коллизии с обычными врагами
			for j := range g.enemies {
				if g.checkCollision(g.bullets[i].x, g.bullets[i].y, g.bullets[i].width, g.bullets[i].height,
					g.enemies[j].x, g.enemies[j].y, g.enemies[j].width, g.enemies[j].height) {
					// Убиваем врага и пулю
					g.bullets[i].active = false
					g.enemies[j].x = -1000 // перемещаем врага за экран (убиваем)
					g.score += 10          // даем очки за убийство врага
				}
			}

			// Коллизии с сильными врагами
			for j := range g.strongEnemies {
				if g.strongEnemies[j].x >= 0 && g.checkCollision(g.bullets[i].x, g.bullets[i].y, g.bullets[i].width, g.bullets[i].height,
					g.strongEnemies[j].x, g.strongEnemies[j].y, g.strongEnemies[j].width, g.strongEnemies[j].height) {
					// Уменьшаем здоровье сильного врага
					g.strongEnemies[j].health--
					g.bullets[i].active = false

					// Если здоровье закончилось, убиваем врага
					if g.strongEnemies[j].health <= 0 {
						g.strongEnemies[j].x = -1000
						g.score += 25 // больше очков за сильного врага
					}
				}
			}
		}
	}

	// Проверяем коллизии с врагами (только с живыми)
	for _, enemy := range g.enemies {
		if enemy.x >= 0 && g.checkCollision(g.player.x, g.player.y, playerSize, playerSize,
			enemy.x, enemy.y, enemy.width, enemy.height) {
			// Игрок теряет жизнь и респавнится
			g.lives--
			g.player.x = 50
			g.player.y = 700
			g.player.vx = 0
			g.player.vy = 0
		}
	}

	// Проверяем коллизии с сильными врагами
	for _, enemy := range g.strongEnemies {
		if enemy.x >= 0 && g.checkCollision(g.player.x, g.player.y, playerSize, playerSize,
			enemy.x, enemy.y, enemy.width, enemy.height) {
			// Игрок теряет жизнь и респавнится
			g.lives--
			g.player.x = 50
			g.player.y = 700
			g.player.vx = 0
			g.player.vy = 0
		}
	}

	// Проверяем коллизию с порталом
	if g.portal.active && g.checkCollision(g.player.x, g.player.y, playerSize, playerSize,
		g.portal.x, g.portal.y, g.portal.width, g.portal.height) {
		g.gameWon = true
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

	// Рисуем врагов (только живых)
	for _, enemy := range g.enemies {
		if enemy.x >= 0 {
			ebitenutil.DrawRect(screen, enemy.x, enemy.y, enemy.width, enemy.height, color.RGBA{255, 0, 255, 255})
		}
	}

	// Рисуем сильных врагов (только живых)
	for _, enemy := range g.strongEnemies {
		if enemy.x >= 0 {
			// Сильные враги красного цвета
			ebitenutil.DrawRect(screen, enemy.x, enemy.y, enemy.width, enemy.height, color.RGBA{255, 0, 0, 255})
		}
	}

	// Рисуем пули
	for _, bullet := range g.bullets {
		if bullet.active {
			ebitenutil.DrawRect(screen, bullet.x, bullet.y, bullet.width, bullet.height, color.RGBA{255, 255, 0, 255})
		}
	}

	// Рисуем портал, если он активен
	if g.portal.active {
		ebitenutil.DrawRect(screen, g.portal.x, g.portal.y, g.portal.width, g.portal.height, color.RGBA{0, 255, 255, 255})
	}

	// Рисуем игрока
	ebitenutil.DrawRect(screen, g.player.x, g.player.y, playerSize, playerSize, color.RGBA{255, 0, 0, 255})

	// Рисуем счет, жизни и время
	timeStr := fmt.Sprintf("%.1f", g.gameTime)
	scoreStr := fmt.Sprintf("%d", g.score)
	livesStr := fmt.Sprintf("%d", g.lives)
	ebitenutil.DebugPrint(screen, "Score: "+scoreStr+" Lives: "+livesStr+" Time: "+timeStr+"s")

	// Показываем экран победы
	if g.gameWon {
		ebitenutil.DebugPrintAt(screen, "YOU WIN! CONGRATULATIONS!", screenWidth/2-100, screenHeight/2)
	}
}

// Layout возвращает размеры экрана
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// checkCollision проверяет коллизию между двумя прямоугольниками
func (g *Game) checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

// shoot создает новую пулю
func (g *Game) shoot() {
	// Проверяем, есть ли уже активные пули (ограничиваем количество)
	activeBullets := 0
	for _, bullet := range g.bullets {
		if bullet.active {
			activeBullets++
		}
	}

	// Максимум 3 пули одновременно
	if activeBullets < 3 {
		bullet := Bullet{
			x:      g.player.x + playerSize/2 - 2, // центрируем пулю
			y:      g.player.y + playerSize/2 - 2,
			vx:     8, // скорость пули
			vy:     0,
			active: true,
			width:  4,
			height: 4,
		}
		g.bullets = append(g.bullets, bullet)
	}
}

// spawnStrongEnemies создает сильных врагов
func (g *Game) spawnStrongEnemies() {
	// Создаем 5 сильных врагов в разных местах
	strongEnemies := []StrongEnemy{
		{400, 600, 2, 0, 40, 40, 300, 500, 3, 0},   // сильный враг 1
		{700, 400, -2, 0, 40, 40, 650, 750, 3, 0},  // сильный враг 2
		{1000, 250, 2, 0, 40, 40, 950, 1050, 3, 0}, // сильный враг 3
		{200, 300, -2, 0, 40, 40, 150, 250, 3, 0},  // сильный враг 4
		{1100, 50, 2, 0, 40, 40, 1050, 1150, 3, 0}, // сильный враг 5
	}

	g.strongEnemies = strongEnemies
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Simple Platformer - Go Game")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
