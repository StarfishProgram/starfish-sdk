package sdk

import "os"

type Signal struct {
	chs []chan os.Signal
}

func (s *Signal) Add(chs ...chan os.Signal) {
	s.chs = append(s.chs, chs...)
}

func (s *Signal) Waiting() os.Signal {
	println("程序已就绪")
	sign := Waiting()
	println("收到终止信号", sign)
	for index := range s.chs {
		ch := s.chs[index]
		ch <- sign
	}
	for index := range s.chs {
		ch := s.chs[index]
		<-ch
	}
	println("程序已停止")
	return sign
}

func NewSignal() *Signal {
	return &Signal{chs: []chan os.Signal{}}
}
