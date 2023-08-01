package pkg

import "log"

func Log(ps *Statistics) {
	log.Printf("\n--- %s ping statistics ---\n", ps.Addr)
	log.Printf("ip %s, %d packets transmitted, %d packets received, %v%% packet loss\n", ps.Addr,
		ps.PacketsSent, ps.PacketsRev, ps.PacketLoss)
	log.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		ps.MinRtt, ps.AvgRtt, ps.MaxRtt, ps.StdDevRtt)
	log.Printf("rtts is %v \n", ps.RttS)
}
