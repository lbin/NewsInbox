package paper_test

import (
	"learnerai/golib/app/context"
	"learnerai/logic/paper"
	"testing"
)

func TestComputeSimScore(t *testing.T) {
	// We want to compute the similarity between the query sentence
	query := "A man is eating pasta."
	// With all sentences in the corpus
	corpus := []string{
		"A man is eating food.",
		"A man is eating a piece of bread.",
		"The girl is carrying a baby.",
		"A man is riding a horse.",
		"A woman is playing violin.",
		"Two men pushed carts through the woods.",
		"A man is riding a white horse on an enclosed ground.",
		"A monkey is playing drums.",
		"A cheetah is running behind its prey.",
	}

	scores, err := paper.ComputeSimScore(context.GetGinContextWithRequestId(), query, corpus)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(scores)

	idx, sorted := paper.SortSimScore(scores)
	t.Log(idx, sorted)
	for i := range idx {
		t.Logf("rank=%v, score=%v: %s", i, sorted[i], corpus[idx[i]])
	}
	// expected:
	// query='A man is eating pasta.'
	// rank=0, score=1.9004740715026855: A man is eating food.
	// rank=1, score=1.4803965091705322: A man is eating a piece of bread.
	// rank=2, score=-7.088953971862793: A man is riding a horse.
	// rank=3, score=-8.904168128967285: A man is riding a white horse on an enclosed ground.
	// rank=4, score=-10.762833595275879: A monkey is playing drums.
	// rank=5, score=-10.765618324279785: A woman is playing violin.
	// rank=6, score=-11.05501937866211: Two men pushed carts through the woods.
	// rank=7, score=-11.07595443725586: The girl is carrying a baby.
	// rank=8, score=-11.099100112915039: A cheetah is running behind its prey.
}
