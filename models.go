package main

type Word struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

type Exceptions struct {
	Word string `storm:"id" json:"word"`
}

type Subtitles struct {
	ID    int    `storm:"id,increment" json:"id"`
	Name  string `json:"name"`
	Words []Word `json:"words"`
}

func (w *Word) translate() {
	w.Translation = translator.translate(w.Word)
}

func (s *Subtitles) refresh() {
	words := make([]Word, 0, len(s.Words))
	for i := range s.Words {
		if !isException(s.Words[i].Word) {
			words = append(words, s.Words[i])
		}
	}
	s.Words = words
	db.Save(s)
}
