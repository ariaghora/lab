package main

func (r *RaptorCfg) DeleteRow(rowNum int) {
	removeAtIndex(r.Rows, rowNum)
	r.CurrentFileNumRows -= 1
}
