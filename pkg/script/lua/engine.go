package lua

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// Engine administra el entorno de ejecución de Lua para ZoneService (Kingdom Quests, Dungeons e IA).
type Engine struct {
	mu       sync.Mutex
	state    *lua.LState
	basePath string
}

// Config contiene las rutas y opciones para inicializar el motor de scripts.
type Config struct {
	BasePath string // Ruta base a data/LuaScript
	Debug    bool
}

// NewEngine inicializa una nueva máquina virtual de Lua configurada con las rutas y funciones oficiales.
func NewEngine(cfg Config) (*Engine, error) {
	if cfg.BasePath == "" {
		cfg.BasePath = "./data/LuaScript"
	}

	absPath, err := filepath.Abs(cfg.BasePath)
	if err != nil {
		return nil, fmt.Errorf("error resolviendo ruta de LuaScript: %w", err)
	}

	L := lua.NewState(lua.Options{
		CallStackSize: 1024,
		RegistrySize:  1024 * 20,
	})

	e := &Engine{
		state:    L,
		basePath: absPath,
	}

	// 1. Configurar package.path para resolución de requires anidados oficiales
	packagePath := fmt.Sprintf("%s/?.lua;%s/?/?.lua;%s/KQ/?.lua;%s/ID/?.lua;%s/AIScript/?.lua",
		absPath, absPath, absPath, absPath, absPath)
	
	packageTable := L.GetGlobal("package")
	if packageTable.Type() == lua.LTTable {
		L.SetField(packageTable, "path", lua.LString(packagePath))
	}

	// 2. Registrar tablas globales oficiales de memoria de instancia
	L.SetGlobal("InstanceField", L.NewTable())

	// 3. Registrar funciones C-API oficiales de Zone.exe
	e.registerOfficialAPIs(cfg.Debug)

	// 4. Cargar el archivo común maestro common.lua
	commonPath := filepath.Join(absPath, "common.lua")
	if _, err := os.Stat(commonPath); err == nil {
		if err := L.DoFile(commonPath); err != nil {
			L.Close()
			return nil, fmt.Errorf("error cargando common.lua: %w", err)
		}
	}

	return e, nil
}

// Close libera los recursos y memoria de la máquina virtual de Lua.
func (e *Engine) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.state != nil {
		e.state.Close()
		e.state = nil
	}
}

// DoFile ejecuta un archivo Lua dentro del contexto del motor.
func (e *Engine) DoFile(filename string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	path := filename
	if !filepath.IsAbs(path) {
		candidate := filepath.Join(e.basePath, filename)
		if _, err := os.Stat(candidate); err == nil {
			path = candidate
		}
	}
	return e.state.DoFile(path)
}

// CallFunction invoca una función global de Lua con argumentos y devuelve los valores de retorno.
func (e *Engine) CallFunction(name string, args ...lua.LValue) ([]lua.LValue, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	fn := e.state.GetGlobal(name)
	if fn.Type() != lua.LTFunction {
		return nil, fmt.Errorf("función Lua '%s' no encontrada o no es función", name)
	}

	top := e.state.GetTop()
	if err := e.state.CallByParam(lua.P{
		Fn:      fn,
		NRet:    lua.MultRet,
		Protect: true,
	}, args...); err != nil {
		return nil, fmt.Errorf("error ejecutando función Lua '%s': %w", name, err)
	}

	nRet := e.state.GetTop() - top
	rets := make([]lua.LValue, nRet)
	for i := 1; i <= nRet; i++ {
		rets[i-1] = e.state.Get(top + i)
	}
	e.state.SetTop(top)

	return rets, nil
}

// GetGlobal devuelve el valor de una variable global de Lua.
func (e *Engine) GetGlobal(name string) lua.LValue {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.state.GetGlobal(name)
}

// GetGlobalTableField obtiene el valor de un campo dentro de una tabla global (ej: "ObjectType", "Player").
func (e *Engine) GetGlobalTableField(tableName, fieldName string) (lua.LValue, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	tbl := e.state.GetGlobal(tableName)
	if tbl.Type() != lua.LTTable {
		return lua.LNil, fmt.Errorf("tabla global '%s' no encontrada o no es tabla", tableName)
	}

	val := e.state.GetField(tbl, fieldName)
	return val, nil
}

// registerOfficialAPIs registra las funciones nativas que Zone.exe expone hacia los scripts de Lua.
func (e *Engine) registerOfficialAPIs(debug bool) {
	L := e.state

	// cExecCheck: traza de ejecución de fase
	L.SetGlobal("cExecCheck", L.NewFunction(func(L *lua.LState) int {
		if debug {
			str := L.CheckString(1)
			log.Printf("[Lua][ExecCheck] %s", str)
		}
		return 0
	}))

	// cDebugLog: log de depuración
	L.SetGlobal("cDebugLog", L.NewFunction(func(L *lua.LState) int {
		if debug {
			str := L.CheckString(1)
			log.Printf("[Lua][DebugLog] %s", str)
		}
		return 0
	}))

	// cAssertLog: aserciones y alertas críticas
	L.SetGlobal("cAssertLog", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		log.Printf("[Lua][AssertLog] %s", str)
		return 0
	}))

	// cCurSec / cCurrentSecond: tiempo actual del servidor en segundos
	currentTimeFn := L.NewFunction(func(L *lua.LState) int {
		sec := time.Now().Unix()
		L.Push(lua.LNumber(sec))
		return 1
	})
	L.SetGlobal("cCurSec", currentTimeFn)
	L.SetGlobal("cCurrentSecond", currentTimeFn)

	// cRandom / cRandomInt: generación de números aleatorios oficiales
	L.SetGlobal("cRandom", L.NewFunction(func(L *lua.LState) int {
		min := L.OptInt(1, 0)
		max := L.OptInt(2, 100)
		if max <= min {
			L.Push(lua.LNumber(min))
			return 1
		}
		val := min + rand.Intn(max-min+1)
		L.Push(lua.LNumber(val))
		return 1
	}))

	L.SetGlobal("cRandomInt", L.NewFunction(func(L *lua.LState) int {
		max := L.OptInt(1, 100)
		if max <= 0 {
			L.Push(lua.LNumber(0))
			return 1
		}
		L.Push(lua.LNumber(rand.Intn(max)))
		return 1
	}))

	// cNotice / cNoticeString: mensajes de aviso a la zona o mapa
	L.SetGlobal("cNotice", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		if debug {
			log.Printf("[Lua][Notice] %s", msg)
		}
		return 0
	}))

	L.SetGlobal("cNoticeString", L.NewFunction(func(L *lua.LState) int {
		msg := L.CheckString(1)
		if debug {
			log.Printf("[Lua][NoticeString] %s", msg)
		}
		return 0
	}))

	// cDistanceSquar: distancia euclidiana al cuadrado (para evitar raíces cuadradas en CPU)
	L.SetGlobal("cDistanceSquar", L.NewFunction(func(L *lua.LState) int {
		x1 := L.CheckNumber(1)
		y1 := L.CheckNumber(2)
		x2 := L.CheckNumber(3)
		y2 := L.CheckNumber(4)
		dx := x1 - x2
		dy := y1 - y2
		distSq := (dx * dx) + (dy * dy)
		L.Push(lua.LNumber(distSq))
		return 1
	}))

	// Stubs para funciones de gestión de objetos y mapas (para que los scripts compilen y corran sin errores de nil)
	stubNoOp := L.NewFunction(func(L *lua.LState) int {
		return 0
	})
	stubReturnZero := L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(0))
		return 1
	})
	stubEmptyString := L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LString(""))
		return 1
	})

	L.SetGlobal("cEndOfKingdomQuest", stubNoOp)
	L.SetGlobal("cTimerStart", stubNoOp)
	L.SetGlobal("cTimerEnd", stubNoOp)
	L.SetGlobal("cMobRegen_XY", stubReturnZero)
	L.SetGlobal("cKillObject", stubNoOp)
	L.SetGlobal("cObjectHP", stubReturnZero)
	L.SetGlobal("cSetObjectHP", stubNoOp)

	// Funciones para soporte completo de scripts de KQ y MobBehavior
	L.SetGlobal("cWhoIsTarget", stubEmptyString)
	L.SetGlobal("cWhoKillMe", stubEmptyString)
	L.SetGlobal("cGetName", stubEmptyString)
	L.SetGlobal("cPermillage", stubReturnZero)
	L.SetGlobal("cChat", stubNoOp)
	L.SetGlobal("cMobRegen", stubNoOp)
	L.SetGlobal("cMobSummon", stubNoOp)
	L.SetGlobal("cHowManyMob", stubReturnZero)
	L.SetGlobal("cKQRestMinute", stubReturnZero)
	L.SetGlobal("cCurrentHP", stubReturnZero)
	L.SetGlobal("cCurrentMap", stubEmptyString)
	L.SetGlobal("cMyName", stubEmptyString)
	L.SetGlobal("cMyHandle", stubReturnZero)
	L.SetGlobal("cShout", stubNoOp)
	L.SetGlobal("cBroadcast", stubNoOp)
	L.SetGlobal("cWaitSec", stubNoOp)
	L.SetGlobal("cWaitMin", stubNoOp)
	L.SetGlobal("cReward", stubNoOp)
	L.SetGlobal("cEndOfEvent", stubNoOp)
	L.SetGlobal("cLinkTo", stubNoOp)
	L.SetGlobal("cMapClear", stubNoOp)
	L.SetGlobal("cSkillBlast", stubNoOp)
	L.SetGlobal("cRetRate", stubReturnZero)
	L.SetGlobal("CallNamedRoutine", stubNoOp)
}
