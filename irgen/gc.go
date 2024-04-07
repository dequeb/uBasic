package irgen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (m *Module) genGC() {
	i8 := types.I8
	i64 := types.I64
	ptrType := types.NewPointer(i8)

	// %struct.GarbageCollector type definition
	garbageCollectorType := types.NewStruct(ptrType, i8, ptrType, i64)
	m.NewTypeDef("struct.GarbageCollector", garbageCollectorType)

	// gc global
	gcGlobal := m.NewGlobal("gc", garbageCollectorType)
	gcGlobal.Linkage = enum.LinkageExternal
	gcGlobal.Align = 8

	garbageCollectorPtrType := types.NewPointer(garbageCollectorType)

	// gc_start
	gcStartParam1 := ir.NewParam("gc", garbageCollectorPtrType)
	gcStartParam1.Attrs = append(gcStartParam1.Attrs, enum.ParamAttrNoUndef)
	gcStartParam2 := ir.NewParam("base_stack", ptrType)
	gcStartParam2.Attrs = append(gcStartParam2.Attrs, enum.ParamAttrNoUndef)
	gcStart := m.NewFunc("gc_start", types.Void, gcStartParam1, gcStartParam2)
	// gc_stop
	gcStopParam1 := ir.NewParam("gc", garbageCollectorPtrType)
	gcStopParam1.Attrs = append(gcStopParam1.Attrs, enum.ParamAttrNoUndef)
	gcStop := m.NewFunc("gc_stop", i64, gcStopParam1)
	// gc_malloc
	gcMallocParam1 := ir.NewParam("gc", garbageCollectorPtrType)
	gcMallocParam1.Attrs = append(gcMallocParam1.Attrs, enum.ParamAttrNoUndef)
	gcMallocParam2 := ir.NewParam("size", ptrType)
	gcMallocParam2.Attrs = append(gcMallocParam2.Attrs, enum.ParamAttrNoUndef)
	gcMalloc := m.NewFunc("gc_malloc", ptrType, gcMallocParam1, gcMallocParam2)
	_ = gcStart
	_ = gcStop
	_ = gcMalloc
}

func (m *Module) GCmalloc(f *Function, size value.Value) value.Value {
	// get gc object
	gcGlobal := m.LookupGlobal("gc")
	// call gc_malloc
	gcMalloc := m.LookupFunction("gc_malloc")
	malloc := f.currentBlock.NewCall(gcMalloc, gcGlobal, size)
	return malloc
}

func (m *Module) GCstart(f *Function, ptr value.Value) {
	// get gc object
	gcGlobal := m.LookupGlobal("gc")
	// call gc_start
	gcStart := m.LookupFunction("gc_start")
	f.currentBlock.NewCall(gcStart, gcGlobal, ptr)
}

func (m *Module) GCstop(f *Function) {
	// get gc object
	gcGlobal := m.LookupGlobal("gc")
	// call gc_stop
	gcStop := m.LookupFunction("gc_stop")
	f.currentBlock.NewCall(gcStop, gcGlobal)
}
