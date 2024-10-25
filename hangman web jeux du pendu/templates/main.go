package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type GameData struct {
	WordToGuess    string   
	GuessedLetters []string 
	LivesLeft      int      
	Message        string   
	GameOver       bool     
	Pseudo         string   
	TimeRemaining  int      
}

var templates = template.Must(template.ParseGlob("templates/*.html"))
var wordToGuess = "EXEMPLE"
var guessedLetters = []string{}
var livesLeft = 6
var pseudo = ""
var startTime time.Time
const gameDuration = 4 * time.Minute 

func init() {
	rand.Seed(time.Now().UnixNano()) 
}


func indexPage(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w, "index.html", nil) 
}

var correctAnswer = "macos"

func verifyAnswer(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        userAnswer := r.FormValue("answer")
        correctAnswer := "macos"

        if strings.ToLower(userAnswer) == correctAnswer {
          
            http.Redirect(w, r, "/menu", http.StatusSeeOther)
            return
        }

        data := struct{ Error string }{Error: "Mauvaise réponse, veuillez réessayer."}
        templates.ExecuteTemplate(w, "index.html", data)
    }
}

func menuPage(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w, "menu.html", nil) 
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "home.html", nil)
	} else if r.Method == "POST" {

		wordToGuess = "EXEMPLE" 
		guessedLetters = []string{}
		livesLeft = 6
		pseudo = r.FormValue("pseudo")
		startTime = time.Now() 
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	}
}


func gamePage(w http.ResponseWriter, r *http.Request) {
	
	elapsed := time.Since(startTime)
	remaining := int(gameDuration.Seconds() - elapsed.Seconds())

	if remaining <= 0 {
		saveScore(pseudo, wordToGuess, "Défaite (temps écoulé)")
		http.Redirect(w, r, "/gameover", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		guess := r.FormValue("guess")

		for _, letter := range guessedLetters {
			if letter == guess {
				templates.ExecuteTemplate(w, "game.html", GameData{
					WordToGuess:    maskWord(wordToGuess, guessedLetters),
					GuessedLetters: guessedLetters,
					LivesLeft:      livesLeft,
					Message:        "Vous avez déjà essayé cette lettre.",
					Pseudo:         pseudo,
					TimeRemaining:  remaining,
				})
				return
			}
		}

		guessedLetters = append(guessedLetters, guess)

		if !isLetterInWord(guess, wordToGuess) {
			livesLeft--
		}

		if livesLeft <= 0 {
			saveScore(pseudo, wordToGuess, "Défaite")
			http.Redirect(w, r, "/gameover", http.StatusSeeOther)
			return
		}

		if !strings.Contains(maskWord(wordToGuess, guessedLetters), "_") {
			saveScore(pseudo, wordToGuess, "Victoire")
			http.Redirect(w, r, "/gameover", http.StatusSeeOther)
			return
		}
	}

	templates.ExecuteTemplate(w, "game.html", GameData{
		WordToGuess:    maskWord(wordToGuess, guessedLetters),
		GuessedLetters: guessedLetters,
		LivesLeft:      livesLeft,
		Pseudo:         pseudo,
		TimeRemaining:  remaining,
	})
}

func gameOverPage(w http.ResponseWriter, r *http.Request) {
	messages := []string{
		"Vous avez perdu !",
		"Vous avez échoué ! Essayez encore !",
		"Le pendu est complet. Retentez votre chance.",
	}
	randomMessage := messages[rand.Intn(len(messages))]
	templates.ExecuteTemplate(w, "gameover.html", struct {
		GameOverMessage string
	}{
		GameOverMessage: randomMessage,
	})
}

func scoresPage(w http.ResponseWriter, r *http.Request) {
	fileContent, err := os.ReadFile("scores.txt")
	if err != nil {
		http.Error(w, "Impossible de lire les scores", http.StatusInternalServerError)
		return
	}
	templates.ExecuteTemplate(w, "scores.html", string(fileContent))
}

func maskWord(word string, guessedLetters []string) string {
	maskedWord := ""
	for _, letter := range word {
		if contains(guessedLetters, string(letter)) {
			maskedWord += string(letter)
		} else {
			maskedWord += "_ "
		}
	}
	return maskedWord
}

func isLetterInWord(letter, word string) bool {
	for _, l := range word {
		if string(l) == letter {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func saveScore(pseudo, word, result string) {
	f, err := os.OpenFile("scores.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Erreur lors de l'ouverture du fichier des scores :", err)
		return
	}
	defer f.Close()

	scoreEntry := fmt.Sprintf("%s a joué avec le mot '%s' et a eu une %s\n", pseudo, word, result)
	if _, err := f.WriteString(scoreEntry); err != nil {
		log.Println("Erreur lors de l'écriture du score :", err)
	}
}

func main() {
    http.HandleFunc("/", indexPage)          
    http.HandleFunc("/verify", verifyAnswer) 
    http.HandleFunc("/menu", menuPage)       

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    port := ":8080"
    log.Println("Serveur démarré sur le port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}