package main

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UIComponent struct {
	font     rl.Font
	portrait string
}

type Actor struct {
	name string
	ui   UIComponent
}

type DialogueLines struct {
	lines []string
	actor Actor
}

type StoryBeat struct {
	dialogueSequence []DialogueLines
	id               string
}

type UP_OR_DOWN string

const (
	UP   UP_OR_DOWN = "UP"
	DOWN UP_OR_DOWN = "DOWN"
)

type DialogueBoxSettings struct {
	boxPosition             UP_OR_DOWN
	portraitWidthProportion float32
	boxHeightProportion     float32
	spacing                 float32
}

func calculateDialogueBoxParams(settings DialogueBoxSettings) (rl.Rectangle, rl.Rectangle) {
	screenHeight := float32(rl.GetScreenHeight())
	screenWidth := float32(rl.GetScreenWidth())
	spacing := settings.spacing
	portraitRect := rl.Rectangle{}
	boxRect := rl.Rectangle{}
	heightProportion := settings.boxHeightProportion * screenHeight
	widthProportion := settings.portraitWidthProportion * screenWidth
	if settings.boxPosition == UP {
		portraitRect.Y = spacing
		boxRect.Y = spacing
	} else if settings.boxPosition == DOWN {
		portraitRect.Y = screenHeight - heightProportion
		boxRect.Y = screenHeight - heightProportion
	}
	portraitRect.X = settings.spacing
	boxRect.X = 1*spacing + widthProportion

	portraitRect.Width = widthProportion - 2*spacing
	boxRect.Width = screenWidth - widthProportion - 2*spacing

	portraitRect.Height = heightProportion - spacing
	boxRect.Height = heightProportion - spacing

	return portraitRect, boxRect
}

func drawStoryBeatSequence(storyBeat StoryBeat, dialogueBoxSettings DialogueBoxSettings) {
	space := dialogueBoxSettings.spacing
	portraitRect, boxRect := calculateDialogueBoxParams(dialogueBoxSettings)

	rl.DrawRectangleRec(portraitRect, rl.Black)
	rl.DrawRectangleRec(boxRect, rl.Gray)

	portraitRectLine := rl.Rectangle{X: portraitRect.X - space, Y: portraitRect.Y - space, Width: portraitRect.Width + 2*space, Height: portraitRect.Height + 2*space}
	boxRectLine := rl.Rectangle{X: boxRect.X - space, Y: boxRect.Y - space, Width: boxRect.Width + 2*space, Height: boxRect.Height + 2*space}

	rl.DrawRectangleLinesEx(boxRectLine, space, rl.Blue)
	rl.DrawRectangleLinesEx(portraitRectLine, space, rl.Red)
}

// SNIPPET taken from https://gist.github.com/CyrusJavan/48254b3fef4ef3bc27784f3fe798623f
// Justifying Text: Given an input string and an maxWidth integer, return a string that is the input justified to the given maxWidth.
//
// text="The quick brown fox jumps over the lazy dog”
// maxWidth=16
// Output:
// The  quick brown
// fox  jumps  over
func justifyText(text string, maxWidth int) string {
	// Split the text into words
	words := strings.Split(text, " ")
	// Keep track of current line and previous lines
	var lines, current string
	// Iterate through each word
	for _, word := range words {
		// If the word fits, add to current
		if len(current)+len(word) <= maxWidth {
			// All words need a space after, expect for the last word
			// we will trim the space for the last word
			current += word + " "
			continue
		}
		// Otherwise, word did not fit so lets finish this line
		// Trim extra space from the last word
		current = current[:len(current)-1]
		// Distribute any extra spaces and add current line to previous lines
		lines += distributeSpaces(current, maxWidth-len(current))
		// The next iteration current line will be this word that did not fit
		current = word + " "
	}
	// return the previous lines and any current line we were working on
	return lines + current
}

// distributeSpaces will spread extra amount of whitespaces into the gaps
// between the words in s
func distributeSpaces(s string, extra int) string {
	// Nothing to do, return early
	if extra < 1 {
		return s + "\n"
	}
	// Split into tokens as it's easier to work with
	words := strings.Split(s, " ")
	// Edge case where only 1 word is on a line
	// we can just put all extra space at the end then return
	if len(words) == 1 {
		for extra > 0 {
			words[0] += " "
			extra--
		}
		return words[0] + "\n"
	}
	// calculate how many spaces between each word and any extra that dont
	// divide evenly into the gaps
	gaps := len(words) - 1
	spaces := extra + gaps
	spacePerGap, extraSpace := spaces/gaps, spaces%gaps
	// Distribute the spaces that dont fit evenly

	for i := 0; extraSpace > 0 && i < len(words)-1; i++ {
		words[i] += " "
		extraSpace--
	}
	// Make a tmp spacer for the whitespace that does divide evenly
	tmp := ""
	for spacePerGap > 0 {
		tmp += " "
		spacePerGap--
	}
	// Join the words with the tmp spacer
	return strings.Join(words, tmp) + "\n"
}

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")

	rl.SetTargetFPS(60)
	storyBeat := StoryBeat{}
	dialogueBoxSettings := DialogueBoxSettings{boxPosition: DOWN, portraitWidthProportion: 0.3, boxHeightProportion: 0.3, spacing: 20}
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		drawStoryBeatSequence(storyBeat, dialogueBoxSettings)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
