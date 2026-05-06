package retrieval

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"
)

type ContentRole string

const ContentRoleData ContentRole = "retrieved_data"

type Document struct {
	SourceID     string
	Title        string
	Workflow     string
	Scope        string
	RevisionDate time.Time
	Body         string
}

type Query struct {
	Text     string
	Workflow string
	Scope    string
	Limit    int
}

type Citation struct {
	SourceID     string
	Title        string
	Workflow     string
	Scope        string
	RevisionDate time.Time
	CitationRef  string
	Snippet      string
	ContentRole  ContentRole
}

type Result struct {
	Matches []Citation
}

type Retriever struct {
	docs []Document
}

func NewRetriever(docs []Document) Retriever {
	copied := make([]Document, len(docs))
	copy(copied, docs)
	return Retriever{docs: copied}
}

func (r Retriever) Retrieve(query Query) Result {
	queryTerms := uniqueTokens(query.Text)
	if len(queryTerms) == 0 || strings.TrimSpace(query.Workflow) == "" || strings.TrimSpace(query.Scope) == "" {
		return Result{}
	}

	var scored []scoredCitation
	for _, doc := range r.docs {
		if doc.Workflow != query.Workflow || doc.Scope != query.Scope {
			continue
		}

		score := scoreDocument(queryTerms, doc)
		if score == 0 {
			continue
		}

		scored = append(scored, scoredCitation{
			score: score,
			citation: Citation{
				SourceID:     doc.SourceID,
				Title:        doc.Title,
				Workflow:     doc.Workflow,
				Scope:        doc.Scope,
				RevisionDate: doc.RevisionDate,
				CitationRef:  citationRef(doc),
				Snippet:      snippet(doc.Body),
				ContentRole:  ContentRoleData,
			},
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score > scored[j].score
		}
		return scored[i].citation.SourceID < scored[j].citation.SourceID
	})

	limit := query.Limit
	if limit <= 0 || limit > len(scored) {
		limit = len(scored)
	}

	matches := make([]Citation, limit)
	for i := 0; i < limit; i++ {
		matches[i] = scored[i].citation
	}
	return Result{Matches: matches}
}

type scoredCitation struct {
	score    int
	citation Citation
}

func scoreDocument(queryTerms []string, doc Document) int {
	documentTerms := tokenSet(doc.Title + " " + doc.Body)
	score := 0
	for _, term := range queryTerms {
		if _, ok := documentTerms[term]; ok {
			score++
		}
	}
	return score
}

func uniqueTokens(text string) []string {
	seen := make(map[string]struct{})
	var tokens []string
	for _, token := range rawTokens(text) {
		if _, stop := stopWords[token]; stop {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		tokens = append(tokens, token)
	}
	return tokens
}

func tokenSet(text string) map[string]struct{} {
	tokens := make(map[string]struct{})
	for _, token := range rawTokens(text) {
		tokens[token] = struct{}{}
	}
	return tokens
}

func rawTokens(text string) []string {
	parts := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	tokens := make([]string, 0, len(parts))
	for _, part := range parts {
		token := strings.TrimSpace(part)
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func citationRef(doc Document) string {
	return fmt.Sprintf("%s#%s", doc.SourceID, doc.RevisionDate.Format(time.DateOnly))
}

func snippet(body string) string {
	body = strings.TrimSpace(body)
	const maxSnippetLength = 240
	if len(body) <= maxSnippetLength {
		return body
	}
	return strings.TrimSpace(body[:maxSnippetLength])
}

var stopWords = map[string]struct{}{
	"a":       {},
	"about":   {},
	"after":   {},
	"an":      {},
	"and":     {},
	"any":     {},
	"applies": {},
	"before":  {},
	"for":     {},
	"in":      {},
	"is":      {},
	"of":      {},
	"or":      {},
	"the":     {},
	"to":      {},
	"what":    {},
	"where":   {},
	"with":    {},
}
