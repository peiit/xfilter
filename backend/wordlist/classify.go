package wordlist

// segment
func VerifyWords(text string) int {
	segmenter := GetSegmenter()
	segments := segmenter.Segment([]byte(text))
	var level int

	if len(segments) > 0 {
		max_freq := segments[0].Token().Frequency()
		for _, seg := range segments {
			token := seg.Token()
			if token.Frequency() > max_freq {
				max_freq = token.Frequency()
			}
		}

		if max_freq == 1 {
			level = 0
		} else if max_freq == 2 {
			level = 1 // need to be reviewed
		} else if max_freq == 3 {
			level = 2 // ban
		}
	}

	return level
}
