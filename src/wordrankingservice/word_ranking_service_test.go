package wordrankingservice

import "testing"



func Test_word_ranking_generation(t *testing.T) {

	user1 := "lagigliaivan"

	rankingService := WordRankingService{}

	rankingService.calculateWordRankingFor(user1)

	t.FailNow()
}