package database

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// GameLogDispatcher implementa un despachador asíncrono no bloqueante para eventos de auditoría de gameplay.
// Mantiene un canal bufferizado y un grupo de workers concurrentes que persisten hacia GameLogRepository.
type GameLogDispatcher struct {
	repo        GameLogRepository
	queue       chan *GameLogEntry
	workerCount int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	closed      atomic.Bool
	dropped     atomic.Uint64
	processed   atomic.Uint64
}

// NewGameLogDispatcher inicializa un despachador con la capacidad de buffer y número de workers configurados.
func NewGameLogDispatcher(repo GameLogRepository, bufferSize int, workerCount int) *GameLogDispatcher {
	if bufferSize <= 0 {
		bufferSize = 10000
	}
	if workerCount <= 0 {
		workerCount = 4
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &GameLogDispatcher{
		repo:        repo,
		queue:       make(chan *GameLogEntry, bufferSize),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start lanza las goroutines de los workers para procesar la cola de logs.
func (d *GameLogDispatcher) Start() {
	for i := 0; i < d.workerCount; i++ {
		d.wg.Add(1)
		go d.workerLoop(i)
	}
}

// workerLoop procesa registros de la cola hasta que el canal se cierre.
func (d *GameLogDispatcher) workerLoop(workerID int) {
	defer d.wg.Done()

	for entry := range d.queue {
		if entry == nil {
			continue
		}
		// Contexto con timeout por operación de persistencia
		opCtx, opCancel := context.WithTimeout(d.ctx, 5*time.Second)
		_, err := d.repo.AddLog(opCtx, entry)
		opCancel()

		if err != nil {
			// En caso de que el contexto padre ya esté cancelado durante shutdown
			if !errors.Is(err, context.Canceled) {
				log.Printf("[GameLogDispatcher] Worker %d failed to persist log (Type=%d, CharNo=%d): %v",
					workerID, entry.Type, entry.CharNo, err)
			}
		} else {
			d.processed.Add(1)
		}
	}
}

// Enqueue encola un evento de log de forma no bloqueante.
// Retorna true si fue aceptado, o false si la cola está llena o el despachador está cerrado.
func (d *GameLogDispatcher) Enqueue(entry *GameLogEntry) bool {
	if entry == nil || d.closed.Load() {
		d.dropped.Add(1)
		return false
	}

	select {
	case d.queue <- entry:
		return true
	default:
		d.dropped.Add(1)
		log.Printf("[GameLogDispatcher] Warning: queue full (%d items), dropping log (Type=%d, CharNo=%d)",
			len(d.queue), entry.Type, entry.CharNo)
		return false
	}
}

// EnqueueType1 encola un log de movimiento/mapa (Tipo 1).
func (d *GameLogDispatcher) EnqueueType1(nType int, nCharNo int, sMap string, x, y, z int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:   nType,
		CharNo: nCharNo,
		Map:    sMap,
		MapX:   x,
		MapY:   y,
		MapZ:   z,
	})
}

// EnqueueType2 encola un log de nivel/stats (Tipo 2).
func (d *GameLogDispatcher) EnqueueType2(nType int, nCharNo int, sMap string, x, y, z, int1, int2 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:   nType,
		CharNo: nCharNo,
		Map:    sMap,
		MapX:   x,
		MapY:   y,
		MapZ:   z,
		Int1:   int1,
		Int2:   int2,
	})
}

// EnqueueType3 encola un log de economía/NPC (Tipo 3).
func (d *GameLogDispatcher) EnqueueType3(nType int, nCharNo int, sMap string, x, y, z, targetCharNo, targetID, int1 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:         nType,
		CharNo:       nCharNo,
		Map:          sMap,
		MapX:         x,
		MapY:         y,
		MapZ:         z,
		TargetCharNo: targetCharNo,
		TargetID:     targetID,
		Int1:         int1,
	})
}

// EnqueueType4 encola un log de trade/scripts (Tipo 4).
func (d *GameLogDispatcher) EnqueueType4(nType int, nCharNo int, sMap string, x, y, z, int1, int2, int3 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:   nType,
		CharNo: nCharNo,
		Map:    sMap,
		MapX:   x,
		MapY:   y,
		MapZ:   z,
		Int1:   int1,
		Int2:   int2,
		Int3:   int3,
	})
}

// EnqueueType5 encola un log de almacén (Tipo 5).
func (d *GameLogDispatcher) EnqueueType5(nType int, nCharNo int, int1, int2 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:   nType,
		CharNo: nCharNo,
		Int1:   int1,
		Int2:   int2,
	})
}

// EnqueueType6 encola un log de drops/dinero (Tipo 6).
func (d *GameLogDispatcher) EnqueueType6(nType int, nCharNo int, sMap string, x, y, z, int1 int, bigint1 int64) bool {
	return d.Enqueue(&GameLogEntry{
		Type:    nType,
		CharNo:  nCharNo,
		Map:     sMap,
		MapX:    x,
		MapY:    y,
		MapZ:    z,
		Int1:    int1,
		Bigint1: bigint1,
	})
}

// EnqueueType7 encola un log de upgrade/enchant (Tipo 7).
func (d *GameLogDispatcher) EnqueueType7(nType int, nCharNo int, sMap string, x, y, z, targetCharNo, int1, int2 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:         nType,
		CharNo:       nCharNo,
		Map:          sMap,
		MapX:         x,
		MapY:         y,
		MapZ:         z,
		TargetCharNo: targetCharNo,
		Int1:         int1,
		Int2:         int2,
	})
}

// EnqueueType8 encola un log de muerte (Tipo 8).
func (d *GameLogDispatcher) EnqueueType8(nType int, nCharNo int, sMap string, x, y, z, int1 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:   nType,
		CharNo: nCharNo,
		Map:    sMap,
		MapX:   x,
		MapY:   y,
		MapZ:   z,
		Int1:   int1,
	})
}

// EnqueueType9 encola un log de misiones (Tipo 9).
func (d *GameLogDispatcher) EnqueueType9(nType int, nCharNo int, sMap string, x, y, z int, itemKey int64, int1, int2 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:    nType,
		CharNo:  nCharNo,
		Map:     sMap,
		MapX:    x,
		MapY:    y,
		MapZ:    z,
		ItemKey: itemKey,
		Int1:    int1,
		Int2:    int2,
	})
}

// EnqueueTypeA encola un log de gremio (Tipo A).
func (d *GameLogDispatcher) EnqueueTypeA(nType int, nCharNo, targetCharNo, targetID int, itemKey int64, int1, int2, int3 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:         nType,
		CharNo:       nCharNo,
		TargetCharNo: targetCharNo,
		TargetID:     targetID,
		ItemKey:      itemKey,
		Int1:         int1,
		Int2:         int2,
		Int3:         int3,
	})
}

// EnqueueTypeB encola un log de almacén de gremio (Tipo B).
func (d *GameLogDispatcher) EnqueueTypeB(nType int, nCharNo int, sMap string, x, y, z, targetCharNo, targetID int, itemKey int64, int1, int2, int3 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:         nType,
		CharNo:       nCharNo,
		Map:          sMap,
		MapX:         x,
		MapY:         y,
		MapZ:         z,
		TargetCharNo: targetCharNo,
		TargetID:     targetID,
		ItemKey:      itemKey,
		Int1:         int1,
		Int2:         int2,
		Int3:         int3,
	})
}

// EnqueueTypeC encola un log de tiendas Booth (Tipo C).
func (d *GameLogDispatcher) EnqueueTypeC(nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2, int3 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:     nType,
		CharNo:   nCharNo,
		Map:      sMap,
		MapX:     x,
		MapY:     y,
		MapZ:     z,
		TargetID: targetID,
		ItemKey:  itemKey,
		Int1:     int1,
		Int2:     int2,
		Int3:     int3,
	})
}

// EnqueueTypeD encola un log de producción/mix/dismantle (Tipo D).
func (d *GameLogDispatcher) EnqueueTypeD(nType int, nCharNo int, itemKey int64, int1, int2 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:    nType,
		CharNo:  nCharNo,
		ItemKey: itemKey,
		Int1:    int1,
		Int2:    int2,
	})
}

// EnqueueTypeE encola un log de skills (Tipo E).
func (d *GameLogDispatcher) EnqueueTypeE(nType int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:     nType,
		Map:      sMap,
		MapX:     x,
		MapY:     y,
		MapZ:     z,
		TargetID: targetID,
		ItemKey:  itemKey,
		Int1:     int1,
		Int2:     int2,
	})
}

// EnqueueTypeF encola un log de MiniHouse (Tipo F).
func (d *GameLogDispatcher) EnqueueTypeF(nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2, int3 int, bigint1 int64) bool {
	return d.Enqueue(&GameLogEntry{
		Type:     nType,
		CharNo:   nCharNo,
		Map:      sMap,
		MapX:     x,
		MapY:     y,
		MapZ:     z,
		TargetID: targetID,
		ItemKey:  itemKey,
		Int1:     int1,
		Int2:     int2,
		Int3:     int3,
		Bigint1:  bigint1,
	})
}

// EnqueueTypeG encola un log de Pet/Mover (Tipo G).
func (d *GameLogDispatcher) EnqueueTypeG(nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2 int, bigint1 int64) bool {
	return d.Enqueue(&GameLogEntry{
		Type:     nType,
		CharNo:   nCharNo,
		Map:      sMap,
		MapX:     x,
		MapY:     y,
		MapZ:     z,
		TargetID: targetID,
		ItemKey:  itemKey,
		Int1:     int1,
		Int2:     int2,
		Bigint1:  bigint1,
	})
}

// EnqueueTypeH encola un log de Kingdom Quest (Tipo H).
func (d *GameLogDispatcher) EnqueueTypeH(nType int, nCharNo int, sMap string, x, y, z, targetCharNo int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:         nType,
		CharNo:       nCharNo,
		Map:          sMap,
		MapX:         x,
		MapY:         y,
		MapZ:         z,
		TargetCharNo: targetCharNo,
	})
}

// EnqueueTypeI encola un log de casino/dados (Tipo I).
func (d *GameLogDispatcher) EnqueueTypeI(nType int, nCharNo int, sMap string, x, y, z, targetCharNo, targetID int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:         nType,
		CharNo:       nCharNo,
		Map:          sMap,
		MapX:         x,
		MapY:         y,
		MapZ:         z,
		TargetCharNo: targetCharNo,
		TargetID:     targetID,
	})
}

// EnqueueTypeJ encola un log de bodas/anillos (Tipo J).
func (d *GameLogDispatcher) EnqueueTypeJ(nType int, nCharNo int, sMap string, x, y, z, targetID int, itemKey int64, int1, int2 int) bool {
	return d.Enqueue(&GameLogEntry{
		Type:     nType,
		CharNo:   nCharNo,
		Map:      sMap,
		MapX:     x,
		MapY:     y,
		MapZ:     z,
		TargetID: targetID,
		ItemKey:  itemKey,
		Int1:     int1,
		Int2:     int2,
	})
}

// EnqueueTypeK encola un log de Academia de gremio (Tipo K).
func (d *GameLogDispatcher) EnqueueTypeK(nType int, nCharNo int, sMap string, x, y, z, int1, int2 int, bigint1 int64) bool {
	return d.Enqueue(&GameLogEntry{
		Type:    nType,
		CharNo:  nCharNo,
		Map:     sMap,
		MapX:    x,
		MapY:    y,
		MapZ:    z,
		Int1:    int1,
		Int2:    int2,
		Bigint1: bigint1,
	})
}

// EnqueueTypeL encola un log de reportes/sanciones GM (Tipo L).
func (d *GameLogDispatcher) EnqueueTypeL(nType int, nCharNo, targetCharNo, targetID int, itemKey int64, int1, int2, int3 int, bigint1 int64) bool {
	return d.Enqueue(&GameLogEntry{
		Type:         nType,
		CharNo:       nCharNo,
		TargetCharNo: targetCharNo,
		TargetID:     targetID,
		ItemKey:      itemKey,
		Int1:         int1,
		Int2:         int2,
		Int3:         int3,
		Bigint1:      bigint1,
	})
}

// Flush espera a que todos los elementos encolados actualmente sean procesados o hasta que ctx expire.
func (d *GameLogDispatcher) Flush(ctx context.Context) error {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		if len(d.queue) == 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// Close realiza un apagado ordenado del despachador:
// 1. Marca closed para no aceptar nuevos logs.
// 2. Cierra el canal.
// 3. Espera a que los workers terminen de drenar todos los logs pendientes o se agote el timeout del contexto.
func (d *GameLogDispatcher) Close(ctx context.Context) error {
	if d.closed.Swap(true) {
		return nil // Ya cerrado
	}

	close(d.queue)

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		d.cancel()
		return nil
	case <-ctx.Done():
		d.cancel() // Cancela operaciones en curso
		return ctx.Err()
	}
}

// Stats retorna estadísticas operativas del despachador.
func (d *GameLogDispatcher) Stats() (queueLen int, processed uint64, dropped uint64) {
	return len(d.queue), d.processed.Load(), d.dropped.Load()
}

