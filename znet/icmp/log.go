package icmp

import "log"

func Debug(ps *Statistics) {
	log.Printf("\n--- %s icmp statistics ---\n", ps.Addr)
	log.Printf("ip %s, %d packets transmitted, %d packets received, %v%% packet loss\n", ps.Addr,
		ps.PacketsSent, ps.PacketsReceive, ps.PacketLoss)
	log.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		ps.MinRtt, ps.AvgRtt, ps.MaxRtt, ps.StdDevRtt)
	log.Printf("rtts is %v \n", ps.RttS)
}
