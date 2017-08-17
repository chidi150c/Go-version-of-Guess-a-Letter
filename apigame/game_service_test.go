package apigame_test

import (
	"fmt"
	"testing"
	"user-apiv2/apigame"
)

type GameService struct {
	*apigame.GameService
}

func TestGameService(t *testing.T) {
	var m = []string{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz", "aBaBaBaBaaBaBaaBaBBaBaaBaaBaBaaa", "DDDDDDDDDDDDDDD"}
	var h GameService
	for i := 0; i < 4; i++ {
		gm, err := h.Start()
		if err != nil {
			t.Fatalf("It did not work and here is the error: %+v", err)
		}
		for _, v := range m[i] {
			gm.GuessedLetter = byte(v)
			gm, err = h.Guess(gm)
			if err == nil {
				//fmt.Println(gm)

			} else {
				fmt.Println(err)
				break
			}
		}

		fmt.Println(string(gm.WordSoFar), gm.WrongGuesses)
		if h.WinOrLoss(gm) {
			fmt.Println("Win!!!")
		} else {
			fmt.Println("Loss!!!")
		}
	}
}
