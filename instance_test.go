package wasmtime

import "testing"

func TestInstance(t *testing.T) {
	wasm, err := Wat2Wasm(`
          (module
            (func (export "f"))
            (global (export "g") i32 (i32.const 0))
            (table (export "t") 1 funcref)
            (memory (export "m") 1)
          )
        `)
	if err != nil {
		panic(err)
	}
	store := NewStore(NewEngine())
	module, err := NewModule(store.Engine, wasm)
	if err != nil {
		panic(err)
	}
	instance, err := NewInstance(store, module, []AsExtern{})
	if err != nil {
		panic(err)
	}
	exports := instance.Exports(store)
	if len(exports) != 4 {
		panic("wrong number of exports")
	}
	if exports[0].Func() == nil {
		panic("not a func")
	}
	if exports[0].Global() != nil {
		panic("should be a func")
	}
	if exports[0].Memory() != nil {
		panic("should be a func")
	}
	if exports[0].Table() != nil {
		panic("should be a func")
	}
	if exports[1].Func() != nil {
		panic("should be a global")
	}
	if exports[1].Global() == nil {
		panic("should be a global")
	}
	if exports[1].Memory() != nil {
		panic("should be a func")
	}
	if exports[1].Table() != nil {
		panic("should be a func")
	}
	if exports[2].Table() == nil {
		panic("should be a table")
	}
	if exports[3].Memory() == nil {
		panic("should be a memory")
	}

	f := exports[0].Func()
	g := exports[1].Global()
	table := exports[2].Table()
	m := exports[3].Memory()

	if len(f.Type(store).Params()) != 0 {
		panic("bad params on type")
	}
	if len(exports[0].Type(store).FuncType().Params()) != 0 {
		panic("bad params on type")
	}
	if g.Type(store).Content().Kind() != KindI32 {
		panic("bad global type")
	}
	if exports[1].Type(store).GlobalType().Content().Kind() != KindI32 {
		panic("bad global type")
	}
	if table.Type(store).Element().Kind() != KindFuncref {
		panic("bad table type")
	}
	if exports[2].Type(store).TableType().Element().Kind() != KindFuncref {
		panic("bad table type")
	}
	if m.Type(store).Minimum() != 1 {
		panic("bad memory type")
	}
	if exports[3].Type(store).MemoryType().Minimum() != 1 {
		panic("bad memory type")
	}
}

func TestInstanceBad(t *testing.T) {
	store := NewStore(NewEngine())
	wasm, err := Wat2Wasm(`(module (import "" "" (func)))`)
	assertNoError(err)
	module, err := NewModule(NewEngine(), wasm)
	assertNoError(err)

	// wrong number of imports
	instance, err := NewInstance(store, module, []AsExtern{})
	if instance != nil {
		panic("expected nil instance")
	}
	if err == nil {
		panic("expected an error")
	}

	// wrong types of imports
	f := WrapFunc(store, func(a int32) {})
	instance, err = NewInstance(store, module, []AsExtern{f})
	if instance != nil {
		panic("expected nil instance")
	}
	if err == nil {
		panic("expected an error")
	}
}

func TestInstanceGetFunc(t *testing.T) {
	wasm, err := Wat2Wasm(`
          (module
            (func (export "f") (nop))
            (global (export "g") i32 (i32.const 0))
          )
	`)
	if err != nil {
		panic(err)
	}
	store := NewStore(NewEngine())
	module, err := NewModule(store.Engine, wasm)
	if err != nil {
		panic(err)
	}
	instance, err := NewInstance(store, module, []AsExtern{})
	if err != nil {
		panic(err)
	}

	f := instance.GetFunc(store, "f")
	if f == nil {
		panic("expected a function")
	}
	_, err = f.Call(store)
	if err != nil {
		panic(err)
	}

	f = instance.GetFunc(store, "g")
	if f != nil {
		panic("expected an error")
	}

	f = instance.GetFunc(store, "f2")
	if f != nil {
		panic("expected an error")
	}
}
