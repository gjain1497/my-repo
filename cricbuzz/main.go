package main

import "fmt"

//## 🏏 Ready for Cricbuzz (Observer Pattern)?

// Let's move to **Week 1, System 2: Cricbuzz**!

// ### What You'll Learn:
// 1. **Observer Pattern** (main focus)
// 2. Event-driven architecture
// 3. Real-time updates/notifications
// 4. Multiple subscribers to events
// 5. Decoupling publishers from subscribers

// ### System Overview:
// ```
// Cricbuzz Requirements:
// - Track live cricket match scores
// - Multiple scoreboards can watch same match
// - When score updates, ALL scoreboards get notified
// - Different types of events: wicket, boundary, over complete
// - Multiple matches can run simultaneously

type MatchStatus string

const (
	NotStarted MatchStatus = "NOT_STARTED"
	InProgress MatchStatus = "IN_PROGRESS"
	Completed  MatchStatus = "COMPLETED"
)

type Observer interface {
	Update(m *Match)
	Display()
}

//Create observer types

type Scoreboard struct {
	name  string
	match *Match
}

func (s *Scoreboard) Update(match *Match) {
	s.match = match
	s.Display()
}

func (s *Scoreboard) Display() {
	if s.match == nil || s.match.TeamBatting == nil {
		return
	}
	fmt.Printf("[%s] %s: %d/%d (%.1f overs) 🏏\n",
		s.name,
		s.match.TeamBatting.Name,
		s.match.TeamBatting.Score.Runs,
		s.match.TeamBatting.Score.Wickets,
		s.match.TeamBatting.Score.Overs)
}

//Commentary Box

type CommentaryBox struct {
	name  string
	match *Match
}

func (c *CommentaryBox) Update(match *Match) {
	c.match = match
	c.Display()
}

func (c *CommentaryBox) Display() {
	if c.match == nil || c.match.TeamBatting == nil {
		return
	}
	score := c.match.TeamBatting.Score
	fmt.Printf("[%s] 📢 That's the score! %s currently at %d runs with %d wickets down!\n",
		c.name,
		c.match.TeamBatting.Name,
		score.Runs,
		score.Wickets)
}

// Mobile App
type MobileApp struct {
	name  string
	match *Match
}

func (m *MobileApp) Update(match *Match) {
	m.match = match
	m.Display()
}

func (m *MobileApp) Display() {
	if m.match == nil || m.match.TeamBatting == nil {
		return
	}
	fmt.Printf("[%s] 📱 Score Update: %s %d/%d (%.1f)\n",
		m.name,
		m.match.TeamBatting.Name,
		m.match.TeamBatting.Score.Runs,
		m.match.TeamBatting.Score.Wickets,
		m.match.TeamBatting.Score.Overs)
}

type Match struct {
	MatchID     string
	Teams       [2]Team
	TeamBatting *Team
	Status      MatchStatus
	Observers   []Observer
}

// subscribe - add an observe
func (m *Match) Subscribe(observer Observer) {
	m.Observers = append(m.Observers, observer)
}

// notify all observers about the change
func (m *Match) NotifyAll() {
	for _, observer := range m.Observers {
		observer.Update(m)
	}
}

func (m *Match) startMatch() {
	m.Status = InProgress
	m.TeamBatting = &m.Teams[0]
}

func (m *Match) addRuns(runs int) {
	teamBatting := m.TeamBatting
	teamBatting.Score.Runs += runs
	m.addBall()
	m.NotifyAll()
}

func (m *Match) addBall() {
	//tracks overs
	//incrment overs correctly
	currOvers := m.TeamBatting.Score.Overs

	balls := int((currOvers - float64(int(currOvers))) * 10)

	if balls >= 5 {
		//completed an over now
		m.TeamBatting.Score.Overs = float64(int(currOvers) + 1)
	} else {
		//add 0.1
		m.TeamBatting.Score.Overs += 0.1
	}
}

func (m *Match) addWicket() {
	m.addBall()
	m.TeamBatting.Score.Wickets += 1

	if m.TeamBatting.Score.Wickets >= 10 {
		m.displayBoard()
		m.switchInnings()
	}
	m.NotifyAll()
}

func (m *Match) switchInnings() {
	if m.TeamBatting == &m.Teams[0] {
		m.TeamBatting = &m.Teams[1]
	} else {
		m.endMatch()
	}
}

func (m *Match) displayBoard() {
	teamBatting := m.TeamBatting.Name
	currentScore := m.TeamBatting.Score

	fmt.Printf("%s: %d/%d (%.1f overs)\n",
		teamBatting, currentScore.Runs, currentScore.Wickets, currentScore.Overs)
}

func (m *Match) endMatch() {
	m.Status = Completed

	team1Score := m.Teams[0].Score.Runs
	team2Score := m.Teams[1].Score.Runs

	if team1Score > team2Score {
		fmt.Printf("Winner: %s\n", m.Teams[0].Name)
	} else if team2Score > team1Score {
		fmt.Printf("Winner: %s\n", m.Teams[1].Name)
	} else {
		fmt.Println("Match Tied!")
	}
}

type Team struct {
	Players []Player
	Name    string
	Score   Score
}

type Player struct {
	Name string
	Role string
	// RunsScored int //for batsmen
	// Wickets    int //for bowlers
}

type Score struct {
	Runs    int
	Wickets int
	Overs   float64
}

func main() {
	india := &Team{
		Players: []Player{
			{
				Name: "Virat",
				Role: "Batsman",
			},
			{
				Name: "Rohit",
				Role: "Batsman",
			},
		},
		Name: "India",
	}

	australia := &Team{
		Players: []Player{
			{
				Name: "Simon",
				Role: "Batsman",
			},
			{
				Name: "Shane",
				Role: "Batsman",
			},
		},
		Name: "Australia",
	}

	match := &Match{
		MatchID: "IND-vs-AUS-2024",
		Teams:   [2]Team{*india, *australia},
		Status:  NotStarted,
	}

	//create observers
	stadiumDisplay := &Scoreboard{name: "Stadium Display"}
	commentary := &CommentaryBox{name: "Commentary"}
	mobileApp := &MobileApp{name: "Mobile App"}

	//subscribe
	match.Subscribe(stadiumDisplay)
	match.Subscribe(commentary)
	match.Subscribe(mobileApp)

	//start match
	match.startMatch()
	fmt.Println("Match started")

	//simulate gamePlay()
	match.addRuns(4)
	match.addRuns(6)
	match.addRuns(1)

	// match.displayBoard() no need now due to observer patteren

	match.addWicket()
	match.addRuns(4)

	// match.displayBoard() no need now due to observer patteren

	// Simulate all-out (9 more wickets)
	for i := 0; i < 9; i++ {
		match.addWicket()
	}
	// match.displayBoard() no need now due to observer patteren

	// Second innings
	match.addRuns(6)
	// match.displayBoard() no need now due to observer patteren

	// Complete match (10 wickets)
	for i := 0; i < 10; i++ {
		match.addWicket()
	}

}
