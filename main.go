package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = int32(800)
	screenHeight = int32(450)
)

type GameScreen int

const (
	LOGO = iota
	TITLE
	GAMEPLAY
	ENDING
)

type GameObject struct {
	posX    int32
	posY    int32
	width   int32
	height  int32
	color   rl.Color
	texture rl.Texture2D
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Birds")
	rl.SetTargetFPS(60)

	var currentScreen GameScreen
	currentScreen = LOGO
	framesCountes := 0

	rl.InitAudioDevice()        // <--- Inicialização do dispositivo de áudio
	defer rl.CloseAudioDevice() // <--- Adicionado defer para fechar o dispositivo de áudio

	soundEat := rl.LoadSound("sound/eat.wav") // <--- Carregar o som
	if soundEat.Stream.Buffer == nil {
		fmt.Println("Failed to load sound")
		return
	}
	defer rl.UnloadSound(soundEat) // <--- Adicionado defer para descarregar o som

	rand.New(rand.NewSource(time.Now().UnixNano()))

	bird_up := rl.LoadImage("assets/bird-up.png")
	bird_down := rl.LoadImage("assets/bird-down.png")
	angry_bird_mid := rl.LoadImage("assets/redbird-midflap.png")
	angry_bird_up := rl.LoadImage("assets/redbird-upflap.png")
	angry_bird_down := rl.LoadImage("assets/redbird-downflap.png")
	fruit := rl.LoadImage("assets/fruit.png")

	rl.ImageFlipHorizontal(angry_bird_mid)
	rl.ImageFlipHorizontal(angry_bird_up)
	rl.ImageFlipHorizontal(angry_bird_down)

	textureBird := rl.LoadTextureFromImage(bird_up)
	textureBackGround := rl.LoadTexture("assets/background.jpg")
	textureFruit := rl.LoadTextureFromImage(fruit)
	textureAngryBird := rl.LoadTextureFromImage(angry_bird_mid)

	AngryBirds := []GameObject{}
	Fruits := []GameObject{}
	addFruit(&Fruits, textureFruit)
	addAngryBird(&AngryBirds, textureAngryBird)

	x_coords, y_coords, score, gameOver := startGame(textureBird)
	highestScore := 0

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		switch currentScreen {
		case LOGO:
			framesCountes++
			if framesCountes > 120 {
				currentScreen = TITLE
			}
		case TITLE:
			if rl.IsKeyPressed(rl.KeyEnter) {
				currentScreen = GAMEPLAY
			}
		case GAMEPLAY:
			if gameOver {
				currentScreen = ENDING
			}
		case ENDING:
			if rl.IsKeyPressed(rl.KeyEnter) {
				textureBird = rl.LoadTextureFromImage(bird_up)
				textureBackGround = rl.LoadTexture("assets/background.jpg")
				textureFruit = rl.LoadTextureFromImage(fruit)
				textureAngryBird = rl.LoadTextureFromImage(angry_bird_mid)

				if score > highestScore {
					highestScore = score
				}
				gameOver = false
				AngryBirds, Fruits = []GameObject{}, []GameObject{}
				x_coords, y_coords, score, gameOver = startGame(textureBird)
				addFruit(&Fruits, textureFruit)
				addAngryBird(&AngryBirds, textureAngryBird)
				currentScreen = GAMEPLAY
			}
		}

		rl.BeginDrawing()
		switch currentScreen {
		case LOGO:
			rl.DrawText("LOADING GAME", 20, 20, 40, rl.LightGray)
			rl.DrawText("Wait for 2 seconds...", 290, 220, 20, rl.LightGray)
		case TITLE:
			rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Green)
			rl.DrawText("B.I.R.D.S", 20, 20, 40, rl.DarkGray)
			rl.DrawText("PRESS ENTER TO START", 120, 220, 20, rl.DarkGray)
		case ENDING:
			rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.DarkBrown)
			rl.DrawText("Your final score is: "+strconv.Itoa(score), 20, 20, 40, rl.Red)
			rl.DrawText("Your highest score is: "+strconv.Itoa(highestScore), 20, 60, 40, rl.Red)
			rl.DrawText("PRESS ENTER TO TRY AGAIN", 120, 220, 20, rl.DarkGray)
		case GAMEPLAY:
			{
				rl.DrawTexture(textureBackGround, 0, 0, rl.White)
				rl.DrawTexture(textureBird, x_coords, y_coords, rl.White)
				rl.DrawText("Current Score: "+strconv.Itoa(score), 0, 0, 30, rl.LightGray)
				rl.DrawText("Highest Score: "+strconv.Itoa(highestScore), 350, 0, 30, rl.LightGray)

				//Bird
				if rl.IsKeyDown(rl.KeySpace) {
					textureBird = rl.LoadTextureFromImage(bird_up)
					y_coords -= 5
				} else {
					textureBird = rl.LoadTextureFromImage(bird_down)
					y_coords += 5
				}
				if y_coords > screenHeight || y_coords < 0 {
					gameOver = true
					Fruits, AngryBirds = nil, nil
					rl.UnloadTexture(textureBird)
					rl.UnloadTexture(textureAngryBird)
					rl.UnloadTexture(textureBackGround)
				}

				//Angry Birds
				for io, current_angryBird := range AngryBirds {
					rl.DrawTexture(textureAngryBird, current_angryBird.posX, current_angryBird.posY, rl.White)
					AngryBirds[io].posX = AngryBirds[io].posX - 5
					switch rand.Intn(3) {
					case 1:
						if AngryBirds[io].posY > 30 && AngryBirds[io].posY < screenHeight-30 {
							AngryBirds[io].posY = AngryBirds[io].posY + 5
							textureAngryBird = rl.LoadTextureFromImage(angry_bird_up)
						}
					case 2:
						if AngryBirds[io].posY > 30 && AngryBirds[io].posY < screenHeight-30 {
							AngryBirds[io].posY = AngryBirds[io].posY - 5
							textureAngryBird = rl.LoadTextureFromImage(angry_bird_down)
						}
					default:
						textureAngryBird = rl.LoadTextureFromImage(angry_bird_mid)
					}
					if current_angryBird.posX < 0 {
						AngryBirds[io].posX = 800
						AngryBirds[io].posY = int32(rand.Intn(450-2+1) - 2)
					}
					if rl.CheckCollisionRecs(rl.NewRectangle(float32(x_coords), float32(y_coords), float32(34), float32(24)),
						rl.NewRectangle(float32(current_angryBird.posX), float32(current_angryBird.posY), float32(current_angryBird.width), float32(current_angryBird.height))) {
						gameOver = true
					}
				}

				//Fruits
				for io, current_apple := range Fruits {
					rl.DrawTexture(textureFruit, current_apple.posX, current_apple.posY, rl.White)
					Fruits[io].posX = Fruits[io].posX - 5
					if current_apple.posX < 0 {
						Fruits[io].posX = 800
						Fruits[io].posY = int32(rand.Intn(450-2+1) - 2)
						if score > 0 {
							score--
						}
					}
					if rl.CheckCollisionRecs(rl.NewRectangle(float32(x_coords), float32(y_coords), float32(34), float32(24)),
						rl.NewRectangle(float32(current_apple.posX), float32(current_apple.posY), float32(current_apple.width), float32(current_apple.height))) {
						Fruits[io].posX = 800
						Fruits[io].posY = int32(rand.Intn(450-2+1) - 2)
						rl.PlaySound(soundEat)
						score++
					}
				}

				if len(AngryBirds) < score {
					addAngryBird(&AngryBirds, textureAngryBird)
				}

			}

		}
		rl.EndDrawing()
		time.Sleep(10000000)

	}
	rl.UnloadTexture(textureBird)
	rl.UnloadTexture(textureAngryBird)
	rl.UnloadTexture(textureBackGround)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func addAngryBird(angryBirds *[]GameObject, texture rl.Texture2D) {
	angryBird_loc := rand.Intn(450-2+1) - 2
	newBird := GameObject{screenWidth, int32(angryBird_loc), 32, 24, rl.Red, texture}
	*angryBirds = append(*angryBirds, newBird)
}

func addFruit(fruits *[]GameObject, texture rl.Texture2D) {
	fruit_loc := rand.Intn(450-2+1) - 2
	newFruit := GameObject{screenWidth, int32(fruit_loc), 25, 24, rl.Red, texture}
	*fruits = append(*fruits, newFruit)
}

func startGame(textureBird rl.Texture2D) (x_coords, y_coords int32, score int, gameOver bool) {
	x_coords = screenWidth/2 - textureBird.Width/2
	y_coords = screenHeight/2 - textureBird.Height/2 - 40
	score = 0
	gameOver = false
	return x_coords, y_coords, score, gameOver
}
