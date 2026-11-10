package zone

import (
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// MapGameLoop gobierna los ciclos determinísticos del mapa (100ms, 1s, 30s) según la arquitectura oficial.
type MapGameLoop struct {
	zMap     *ZoneMap
	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// NewMapGameLoop instancia un bucle de juego para el mapa indicado.
func NewMapGameLoop(zMap *ZoneMap) *MapGameLoop {
	return &MapGameLoop{
		zMap:     zMap,
		stopChan: make(chan struct{}),
	}
}

// Start inicia las rutinas periódicas de simulación del mapa.
func (gl *MapGameLoop) Start() {
	gl.mu.Lock()
	if gl.running {
		gl.mu.Unlock()
		return
	}
	gl.running = true
	gl.mu.Unlock()

	gl.wg.Add(1)
	go gl.run()
}

// Stop finaliza limpiamente el bucle de juego del mapa.
func (gl *MapGameLoop) Stop() {
	gl.mu.Lock()
	if !gl.running {
		gl.mu.Unlock()
		return
	}
	gl.running = false
	close(gl.stopChan)
	gl.mu.Unlock()

	gl.wg.Wait()
}

func (gl *MapGameLoop) run() {
	defer gl.wg.Done()

	tick100ms := time.NewTicker(100 * time.Millisecond)
	defer tick100ms.Stop()

	tick1s := time.NewTicker(1 * time.Second)
	defer tick1s.Stop()

	tick30s := time.NewTicker(30 * time.Second)
	defer tick30s.Stop()

	for {
		select {
		case <-gl.stopChan:
			return
		case <-tick100ms.C:
			gl.tick100ms()
		case <-tick1s.C:
			gl.tick1s()
		case <-tick30s.C:
			gl.tick30s()
		}
	}
}

// tick100ms ejecuta la lógica de alta frecuencia (IA de monstruos, barras de casteo, etc.).
func (gl *MapGameLoop) tick100ms() {
	// Lógica de alta frecuencia extensible
}

// tick1s ejecuta la regeneración natural y en tienda (Tent) de HP y SP cada segundo.
func (gl *MapGameLoop) tick1s() {
	gl.zMap.mu.RLock()
	players := make([]*ZonePlayer, 0, len(gl.zMap.players))
	for _, p := range gl.zMap.players {
		players = append(players, p)
	}
	gl.zMap.mu.RUnlock()

	for _, p := range players {
		if p.Char == nil || p.Char.CurHP <= 0 {
			continue
		}

		hpChanged := false
		spChanged := false

		// Estado 4 = Resting en Tent: regeneración rápida de HP y SP (10% por segundo)
		if p.State == 4 {
			regenHP := p.MaxHP / 10
			if regenHP < 5 {
				regenHP = 5
			}
			if p.Char.CurHP < p.MaxHP {
				p.Char.CurHP += regenHP
				if p.Char.CurHP > p.MaxHP {
					p.Char.CurHP = p.MaxHP
				}
				hpChanged = true
			}

			regenSP := p.MaxSP / 10
			if regenSP < 5 {
				regenSP = 5
			}
			if p.Char.CurSP < p.MaxSP {
				p.Char.CurSP += regenSP
				if p.Char.CurSP > p.MaxSP {
					p.Char.CurSP = p.MaxSP
				}
				spChanged = true
			}
		}

		if hpChanged && p.Client != nil {
			hpUpdate := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
			hpUpdate.WriteUint32(uint32(p.Char.CurHP))
			hpUpdate.WriteUint16(0)
			_ = p.Client.SendPacket(hpUpdate)
		}

		if spChanged && p.Client != nil {
			spUpdate := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
			spUpdate.WriteUint32(uint32(p.Char.CurSP))
			_ = p.Client.SendPacket(spUpdate)
		}
	}
}

// tick30s limpia ítems del suelo expirados y estados periódicos lentos.
func (gl *MapGameLoop) tick30s() {
	gl.zMap.CleanupExpiredDrops(2 * time.Minute)
}
