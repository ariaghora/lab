package main

func (r *RaptorCfg) DeleteRow(rowNum int) {
	removeAtIndex(r.Rows, rowNum)
	r.CurrentFileNumRows -= 1
	if r.CY+r.CurrentRowOffset > r.CurrentFileNumRows-1 {
		r.CY -= 1
	}
	line := r.Rows[r.CY+r.CurrentRowOffset].Chars
	if r.CX+r.CurrentColOffset > len(line)-1 {
		r.CX = r.CurrentColOffset + len(line) - 1
		if r.CX < 0 {
			r.CX += 1
		}
	}
}
