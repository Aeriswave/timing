package main

import (
	"fmt"
	"time"
)

// SynchroTimers Таймеры для отслеживания событий временной синхронизации
type SynchroTimers struct {
	roundTimer     *time.Ticker         // таймер синхронизации начала периода|раунда (уведомления о начале)
	stopTimer      *time.Ticker         // таймер контроля активности приложения
	syncNumber     *int                 // счетчик раундов|периодов времени
	syncSpeed      *int                 // действий в раунде|периоде времени
	actionDuration *time.Duration       // длительность 1 действия
	appCheckTimer  *time.Duration       // длительность таймера контроля
	pause          *bool                // индикатор паузы
	roundDuration  func() time.Duration // продолжительность раунда (вычисляемое значение)
}

// Clock методы управления таймером
type Clock interface {
	Start() // Запустить часы
	Stop()  // Остановить часы
	Pause() // Поставить часы на паузу
}

func main() {
	fmt.Println("Привет")
	fmt.Println("Старт", time.Now())
	Clock := Init()
	Clock.Start()

	fmt.Println("Программа завершена")
	return
}

// Init Инициализация
func Init() Clock {
	var syncNum int = 0
	var syncSpeed int = 256
	var actionDuration time.Duration = time.Millisecond * 16
	var appCheckTimer time.Duration = time.Millisecond * 15000
	var pause bool = true
	rd := func() time.Duration {
		rr := time.Duration(actionDuration * time.Duration(syncSpeed))
		return rr
	}

	var t1 Clock = SynchroTimers{
		syncNumber:     &syncNum,
		syncSpeed:      &syncSpeed,
		actionDuration: &actionDuration,
		appCheckTimer:  &appCheckTimer,
		pause:          &pause,
		roundDuration:  rd,
	}

	return t1
}

// Start Запуск
func (t SynchroTimers) Start() {
	t.Pause()
	t.roundTimer = time.NewTicker(t.roundDuration()) // вычисляемое значение
	t.stopTimer = time.NewTicker(*t.appCheckTimer)

	go t.Sync()

	for tt := range t.stopTimer.C {
		t.Pause()
		t.roundTimer.Stop()
		fmt.Println("Синхронизация остановлена", tt)
		t.Stop()
		return
	}
}

// Stop Остановка обработки событий контроля
func (t SynchroTimers) Stop() {
	t.stopTimer.Stop()
	fmt.Println("Таймер контроля остановлен")
	return
}

// Sync Запуск обработки событий контроля
func (t SynchroTimers) Sync() {
	*t.pause = false
	t.Round()
	for tt := range t.roundTimer.C {
		fmt.Println(tt)
		t.Round()
	}
	return
}

// Round Запуск обработки событий синхронизации
func (t SynchroTimers) Round() {
	if !*t.pause {
		fmt.Println("Раунд ", *t.syncNumber)
		// Новый раунд
		//.. код

		// Завершение раунда
		//.. код
		fmt.Println("Раунд ", *t.syncNumber, " завершен")
		*t.syncNumber = (*t.syncNumber + 1) % int((^uint(0))>>4) // с защитой от превышения счётчика
	}
	return
}

// Pause Приостановить обработку событий синхронизации | раундов
func (t SynchroTimers) Pause() {
	*t.pause = true

	return
}
