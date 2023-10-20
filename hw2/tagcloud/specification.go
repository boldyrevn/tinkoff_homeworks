package tagcloud

type TagCloud struct {
	cloud []TagStat
}

type TagStat struct {
	Tag             string
	OccurrenceCount int
}

func New() TagCloud {
	return TagCloud{[]TagStat{}}
}

func (t *TagCloud) AddTag(tag string) {
	for i := range t.cloud {
		if t.cloud[i].Tag == tag {
			t.cloud[i].OccurrenceCount += 1
			for i > 0 && t.cloud[i-1].OccurrenceCount < t.cloud[i].OccurrenceCount {
				t.cloud[i-1], t.cloud[i] = t.cloud[i], t.cloud[i-1]
			}
			return
		}
	}
	t.cloud = append(t.cloud, TagStat{tag, 1})
}

func (t *TagCloud) TopN(n int) []TagStat {
	if n > len(t.cloud) {
		n = len(t.cloud)
	}
	topArray := make([]TagStat, n)
	copy(topArray, t.cloud[:n])
	return topArray
}
