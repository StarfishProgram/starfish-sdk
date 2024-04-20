package starfish_sdk

import "os"

type Signal struct {
	chs []chan os.Signal
}

func (s *Signal) AddChans(chs ...chan os.Signal) {
	s.chs = append(s.chs, chs...)
}

func (s *Signal) Waiting() os.Signal {
	Log().Info("程序已就绪")
	sign := Waiting()
	Log().Info("收到终止信号", sign)
	for index := range s.chs {
		ch := s.chs[index]
		ch <- sign
	}
	for index := range s.chs {
		ch := s.chs[index]
		<-ch
	}
	Log().Info("程序已停止")
	return sign
}

func SignalNew() *Signal {
	return &Signal{chs: []chan os.Signal{}}
}
