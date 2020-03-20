package local

type LocalGrab struct {
	LocalTotal int64
	LocalSold  int64
}

//本地扣库存,返回bool值
func (ticket *LocalGrab) LocalGrabTicket() bool {
	ticket.LocalSold++
	return ticket.LocalSold <= ticket.LocalTotal
}
