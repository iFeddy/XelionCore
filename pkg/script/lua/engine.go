package lua

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

type cachedLuaProto struct {
	size    int64
	modTime int64
	proto   *lua.FunctionProto
}

var (
	luaProtoCacheMu sync.RWMutex
	luaProtoCache   = make(map[string]cachedLuaProto)
)

// ClearProtoCache limpia el caché global de bytecode compilado *lua.FunctionProto.
func ClearProtoCache() {
	luaProtoCacheMu.Lock()
	luaProtoCache = make(map[string]cachedLuaProto)
	luaProtoCacheMu.Unlock()
}

// CompileFileCached lee y compila un archivo Lua a *lua.FunctionProto una sola vez en todo el proceso.
// El *lua.FunctionProto resultante es inmutable y seguro para compartir entre múltiples lua.LState.
func CompileFileCached(path string) (*lua.FunctionProto, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	fi, statErr := os.Stat(absPath)
	if statErr == nil {
		size := fi.Size()
		modTime := fi.ModTime().UnixNano()

		luaProtoCacheMu.RLock()
		if entry, ok := luaProtoCache[absPath]; ok && entry.size == size && entry.modTime == modTime {
			luaProtoCacheMu.RUnlock()
			return entry.proto, nil
		}
		luaProtoCacheMu.RUnlock()

		src, err := os.ReadFile(absPath)
		if err != nil {
			return nil, err
		}
		chunk, err := parse.Parse(bytes.NewReader(src), absPath)
		if err != nil {
			return nil, err
		}
		proto, err := lua.Compile(chunk, absPath)
		if err != nil {
			return nil, err
		}

		luaProtoCacheMu.Lock()
		luaProtoCache[absPath] = cachedLuaProto{
			size:    size,
			modTime: modTime,
			proto:   proto,
		}
		luaProtoCacheMu.Unlock()
		return proto, nil
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	chunk, err := parse.Parse(bytes.NewReader(src), path)
	if err != nil {
		return nil, err
	}
	return lua.Compile(chunk, path)
}

// Engine administra el entorno de ejecución de Lua para ZoneService (Kingdom Quests, Dungeons e IA).
type Engine struct {
	mu       sync.Mutex
	state    *lua.LState
	basePath string
	bridge   WorldBridge
	debug    bool
}

// Config contiene las rutas y opciones para inicializar el motor de scripts.
type Config struct {
	BasePath string      // Ruta base a data/LuaScript
	Debug    bool        // Habilita trazado de depuración
	Bridge   WorldBridge // Enlace opcional a las entidades del juego
}

// NewEngine inicializa una nueva máquina virtual de Lua configurada con las rutas y las 274 funciones C-API oficiales.
func NewEngine(cfg Config) (*Engine, error) {
	if cfg.BasePath == "" {
		cfg.BasePath = "./data/LuaScript"
	}
	if cfg.Bridge == nil {
		cfg.Bridge = &DefaultBridge{}
	}

	absPath, err := filepath.Abs(cfg.BasePath)
	if err != nil {
		return nil, fmt.Errorf("error resolviendo ruta de LuaScript: %w", err)
	}

	L := lua.NewState(lua.Options{
		CallStackSize:       128,
		RegistrySize:        1024,
		RegistryMaxSize:     1024 * 20,
		RegistryGrowStep:    1024,
		MinimizeStackMemory: true,
	})

	e := &Engine{
		state:    L,
		basePath: absPath,
		bridge:   cfg.Bridge,
		debug:    cfg.Debug,
	}

	// 1. Configurar package.path para resolución de requires anidados oficiales
	packagePath := fmt.Sprintf("%s/?.lua;%s/?/?.lua;%s/KQ/?.lua;%s/ID/?.lua;%s/AIScript/?.lua;%s/Scenario/?.lua",
		absPath, absPath, absPath, absPath, absPath, absPath)

	packageTable := L.GetGlobal("package")
	if packageTable.Type() == lua.LTTable {
		L.SetField(packageTable, "path", lua.LString(packagePath))
	}

	// 2. Registrar tablas globales oficiales de memoria de instancia
	L.SetGlobal("InstanceField", L.NewTable())

	// 3. Registrar el 100% de las 274 funciones oficiales C-API de Zone.pdb
	e.registerAllOfficialAPIs()

	// 4. Cargar el archivo común maestro common.lua usando el caché de bytecode *lua.FunctionProto
	commonPath := filepath.Join(absPath, "common.lua")
	if _, err := os.Stat(commonPath); err == nil {
		if err := e.doFileLocked(commonPath); err != nil {
			L.Close()
			return nil, fmt.Errorf("error cargando common.lua: %w", err)
		}
	}

	return e, nil
}

func (e *Engine) doFileLocked(path string) error {
	proto, err := CompileFileCached(path)
	if err != nil {
		return err
	}
	fn := e.state.NewFunctionFromProto(proto)
	e.state.Push(fn)
	return e.state.PCall(0, 0, nil)
}

// SetBridge enlaza una nueva implementación de WorldBridge (por ejemplo al arrancar un mapa).
func (e *Engine) SetBridge(bridge WorldBridge) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if bridge != nil {
		e.bridge = bridge
	}
}

// GetBridge retorna el WorldBridge actualmente enlazado.
func (e *Engine) GetBridge() WorldBridge {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.bridge
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

// DoFile ejecuta un archivo Lua dentro del contexto del motor reutilizando su bytecode compilado.
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
	return e.doFileLocked(path)
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

// CallTableFunction invoca una función dentro de una tabla global (ej: KQConditionOfHero.Main).
func (e *Engine) CallTableFunction(tableName, funcName string, args ...lua.LValue) ([]lua.LValue, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	tbl := e.state.GetGlobal(tableName)
	if tbl.Type() != lua.LTTable {
		return nil, fmt.Errorf("tabla global '%s' no encontrada o no es tabla", tableName)
	}
	fn := e.state.GetField(tbl, funcName)
	if fn.Type() != lua.LTFunction {
		return nil, fmt.Errorf("función '%s.%s' no encontrada o no es función", tableName, funcName)
	}

	top := e.state.GetTop()
	if err := e.state.CallByParam(lua.P{
		Fn:      fn,
		NRet:    lua.MultRet,
		Protect: true,
	}, args...); err != nil {
		return nil, fmt.Errorf("error ejecutando '%s.%s': %w", tableName, funcName, err)
	}

	nRet := e.state.GetTop() - top
	rets := make([]lua.LValue, nRet)
	for i := 1; i <= nRet; i++ {
		rets[i-1] = e.state.Get(top + i)
	}
	e.state.SetTop(top)
	return rets, nil
}

// SetGlobal asigna un valor global en la VM de Lua (p. ej. InitValue para scripts KQ).
func (e *Engine) SetGlobal(name string, val lua.LValue) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.state.SetGlobal(name, val)
}
