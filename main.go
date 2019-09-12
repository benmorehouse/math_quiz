package main

import(
	"flag"
	"strings"
	"fmt"
	"log"
	"os"
	"encoding/csv"
	"time"
	"math/rand"
)

func main(){
	csv_filename := flag.String("csv","problems.csv", "a csv file with our data") // this is a pointer to problems.csv
	time_limit := flag.Int("limit",5, "A time limit for the quiz")
	flag.Parse()
	// this also makes the first argument available for use to use in the rest of the code
	file , err := os.Open(*csv_filename)
	if err != nil{
		log.Println("Failed to open ", *csv_filename)
	}
	// we are gonna make a new reader using this csv file and then loop through it and keep track of how many times the user is correct
	r := csv.NewReader(file) // csv takes in an ioreader and returns a new reader for csv files
	lines , err := r.ReadAll() // this will then take the csv file reader and read through and return a string of all the objects that we are trying to touch
	if err != nil{
		log.Println(err)
	}
	problems := parse_lines(lines) // problems is of type problem set
	timer := time.NewTimer(time.Duration(*time_limit) * time.Second) // before it was int * duration... time.Duration returns durations type 
	//<-timer.C now we are waiting on something to fill in here ... timer will then push something back through the channel when it hits the limit


	correct := 0
	rand.Shuffle(len(problems), func(i, j int) {
		 problems[i], problems[j] = problems[j], problems[i]
	})
	for _ , val := range problems{
		fmt.Println(val.question)
		var answerChan chan string
		go func(){  // you are essentially just declaring a function inside
			var answer string
			fmt.Scan(&answer)
			answer = strings.TrimSpace(answer)
			answerChan <- answer
		}()

		select{ // now you have separate function going within for loop to get input from user. That starts and while it is running it will do the select and case portion
		case <-timer.C:
			fmt.Println("Time is up")
			time.Sleep(5 * time.Second)
			return
		case answer := <-answerChan:
			if answer == val.answer{
				correct++
			}
		}
	}
	fmt.Println("correct is",correct)
}

func parse_lines(lines [][]string) []problem_set{
	returnVal := make([]problem_set,len(lines))
	for i, val := range lines{
		returnVal[i] = problem_set{
			question: val[0],
			answer: strings.TrimSpace(val[1]),
		}
	}
	return returnVal
}

type problem_set struct{
	question string
	answer string
	marked bool
}

