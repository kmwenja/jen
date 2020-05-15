package jen

import "bytes"

// ParseFrontmatter extracts frontmatter from any sequence of bytes
func ParseFrontmatter(content []byte) (matter []byte, rest []byte) {
	var (
		beforematter = 1
		inmatter     = 2
		aftermatter  = 3
	)
	status := beforematter

	hyphenCount := 0

	matterBuf := &bytes.Buffer{}
	mdBuf := &bytes.Buffer{}
	tempBuf := &bytes.Buffer{}

	for _, b := range content {
		switch status {
		case beforematter:
			switch b {
			case '-':
				hyphenCount++
				// just in case it turns out it wasn't a separator
				tempBuf.WriteByte(b)
			case '\n':
				if hyphenCount >= 3 {
					status = inmatter
					hyphenCount = 0
					// since we got front matter, reset the content
					// we had already put into the markdown buffer
					mdBuf.Truncate(0)
					tempBuf.Truncate(0)
				} else {
					mdBuf.Write(tempBuf.Bytes())
					tempBuf.Truncate(0)
					mdBuf.WriteByte(b)
				}
			default:
				// assume there's no front matter
				mdBuf.Write(tempBuf.Bytes())
				tempBuf.Truncate(0)
				mdBuf.WriteByte(b)
				hyphenCount = 0
			}
		case inmatter:
			switch b {
			case '-':
				hyphenCount++
				// just in case it turns out it wasn't a separator
				tempBuf.WriteByte(b)
			case '\n':
				if hyphenCount >= 3 {
					status = aftermatter
					hyphenCount = 0
					tempBuf.Truncate(0)
				} else {
					matterBuf.Write(tempBuf.Bytes())
					tempBuf.Truncate(0)
					matterBuf.WriteByte(b)
				}
			default:
				matterBuf.Write(tempBuf.Bytes())
				tempBuf.Truncate(0)
				matterBuf.WriteByte(b)
				hyphenCount = 0
			}
		case aftermatter:
			mdBuf.WriteByte(b)
		}
	}

	// in case any temp stuff was left
	switch status {
	case beforematter, aftermatter:
		mdBuf.Write(tempBuf.Bytes())
	case inmatter:
		matterBuf.Write(tempBuf.Bytes())
	}

	matter = matterBuf.Bytes()
	rest = mdBuf.Bytes()

	return
}
