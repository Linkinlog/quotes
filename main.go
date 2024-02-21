package main

import (
	"log/slog"

	"github.com/Linkinlog/quotes/db"
	"github.com/Linkinlog/quotes/handlers"
	"github.com/Linkinlog/quotes/models"
	"github.com/Linkinlog/quotes/repository"
)

const addr = ":8080"

func main() {
	store := db.NewInMemoryStore(makeQuotes())
	repo := repository.NewQuoteRepository(store)

	err := handlers.NewHandler(repo).HandleRoutes(addr)
	if err != nil {
		slog.Error(err.Error())
	}
}

func makeQuotes() []*models.Quote {
	return []*models.Quote{
		models.NewQuote("Hello, World.", "Log"),
		models.NewQuote("Heaven or hell, love or hate, No matter where I turn I meet myself. Holding life precious is Just living with all intensity Holding life precious.", "Opening the Hand of Thought p.81- Kōshō Uchiyama"),
		models.NewQuote("Regardless of which component you begin at, it is impossible to follow the dependency relationships and wind up back at that component. This structure has no cycles. It is a directed acyclic graph(DAG)", "Clean Architecture Ch 14 - Robert C. Martin"),
		models.NewQuote("The key point here is our programmers are Googlers, they’re not researchers. They’re typically, fairly young, fresh out of school, probably learned Java, maybe learned C or C++, probably learned Python. They’re not capable of understanding a brilliant language but we want to use them to build good software. So, the language that we give them has to be easy for them to understand and easy to adopt.", "Rob Pike"),
		models.NewQuote("If you give me a program that works perfectly but is impossible to change, then it won’t work when the requirements change, and I won’t be able to make it work. Therefore the program will become useless.", "Clean Architecture - Robert C. Martin"),
		models.NewQuote("If you give me a program that does not work but is easy to change, then I can make it work, and keep it working as requirements change. Therefore the program will remain continually useful.", "Clean Architecture - Robert C. Martin"),
		models.NewQuote("Testing shows the presence, not the absence, of bugs.", "Djikstra"),
		models.NewQuote("A software artifact should be open for extension but closed for modification.", "Bertrand Meyer"),
		models.NewQuote("A module should be responsible to one, and only one, actor", "Clean Architecture - Robert C. Martin"),
		models.NewQuote("The only way to go fast, is to go well.", "Clean Architecture - Robert C. Martin"),
		models.NewQuote("Their overconfidence will drive the redesign into the same mess as the original project.", "Clean Architecture - Robert C. Martin"),
		models.NewQuote("[..] the Buddha does not offer us palliatives that leave the underlying maladies untouched beneath the surface;", "In The Buddhas Words"),
		models.NewQuote("Things don’t just happen, they are made to happen", "John F. Kennedy"),
		models.NewQuote("Things don’t just happen, they are made to happen", "John F. Kennedy"),
	}
}
