package main

import (
	"net/http"
	"strconv"

	"github.com/GeorgianBadita/backgammon-move-generator/pkg/board"
	"github.com/GeorgianBadita/backgammon-move-generator/pkg/game"
	"github.com/gin-gonic/gin"
)

type Move struct {
	From     int    `json:"from"`
	To       int    `json:"to"`
	MoveType string `json:"move_type"`
}

type MakeMoveData struct {
	BoardStr  string `json:"board_str" binding:"required"`
	Move      Move   `json:"move" binding:"required"`
	EndOfTurn bool   `json:"end_of_turn"`
}

func getMoveRolls(c *gin.Context) {
	boardStr := c.Param("board")
	die1, _ := strconv.Atoi(c.Param("die1"))
	die2, _ := strconv.Atoi(c.Param("die2"))
	moveRolls := game.GetMoveRollsForSerializedBoard(boardStr, board.DieRoll{Die1: die1, Die2: die2})
	movesMap := map[Move]bool{}

	for idx := 0; idx < len(moveRolls); idx++ {
		for jdx := 0; jdx < len(moveRolls[idx]); jdx++ {
			mvType := "NORMAL_MOVE"
			if moveRolls[idx][jdx].Type == board.CHECKER_ON_BAR_MOVE {
				mvType = "CHECKER_ON_BAR_MOVE"
			} else if moveRolls[idx][jdx].Type == board.BEARING_OFF_MOVE {
				mvType = "BEARING_OFF_MOVE"
			}
			movesMap[Move{From: int(moveRolls[idx][jdx].From), To: int(moveRolls[idx][jdx].To), MoveType: mvType}] = true
		}
	}
	movesToRet := []Move{}
	for mv := range movesMap {
		movesToRet = append(movesToRet, mv)
	}

	c.JSON(http.StatusOK, gin.H{"moves": movesToRet})
}

func getMovesForOneDie(c *gin.Context) {
	boardStr := c.Param("board")
	die, _ := strconv.Atoi(c.Param("die"))
	moves := game.GetMovesForSerializedBoard(boardStr, die)
	movesToRet := []Move{}
	for idx := 0; idx < len(moves); idx++ {
		mvType := "NORMAL_MOVE"
		if moves[idx].Type == board.CHECKER_ON_BAR_MOVE {
			mvType = "CHECKER_ON_BAR_MOVE"
		} else if moves[idx].Type == board.BEARING_OFF_MOVE {
			mvType = "BEARING_OFF_MOVE"
		}
		movesToRet = append(movesToRet, Move{From: int(moves[idx].From), To: int(moves[idx].To), MoveType: mvType})
	}
	c.JSON(http.StatusOK, gin.H{"moves": movesToRet})
}

func makeMove(c *gin.Context) {
	var input MakeMoveData

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonMove := input.Move
	mvType := board.NORMAL_MOVE
	if jsonMove.MoveType == "CHECKER_ON_BAR_MOVE" {
		mvType = board.CHECKER_ON_BAR_MOVE
	} else if jsonMove.MoveType == "BEARING_OFF_MOVE" {
		mvType = board.BEARING_OFF_MOVE
	}

	move := board.Move{From: board.PointIndex(jsonMove.From), To: board.PointIndex(jsonMove.To), Type: mvType}
	c.JSON(http.StatusOK, gin.H{"board": game.MakeMoveOnSerializedBoard(input.BoardStr, move, input.EndOfTurn)})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.UseRawPath = true
	// r.UnescapePathValues = false
	r.GET("/move-rolls/:board/:die1/:die2", getMoveRolls)
	r.GET("/moves/:board/:die", getMovesForOneDie)
	r.POST("/move", makeMove)
	r.Run() // listen and serve on 0.0.0.0:8080
}
