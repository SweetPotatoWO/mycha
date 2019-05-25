package module

import (
	"mycha/errors"
	"sync"


)

//注册器的接口
type Registrar interface {
	 Register(module Module) (bool,error) //用于注册组件实例
	 Unregister(mid MID) (bool,error)  //用于注销实例

	 //用于获取到某个实例 这里用了负载均衡的方式？ 不是很懂
	 Get(moduleType Type)  (Module,error)
	 //获取到全部的某种类型实例
	 GetAllByType(moduleType Type)(map[MID]Module,error)
	 //获取到全部组件
	 GetAll() map[MID]Module
	 //清除
	 Clear()
}

// NewRegistrar 用于创建一个组件注册器的实例。
func NewRegistrar() Registrar {
	return &myRegistrar{
		moduleTypeMap: map[Type]map[MID]Module{},
	}
}

//实现的结构体
type myRegistrar struct {
	//代表组件类型与对应的类型
	moduleTypeMap map[Type]map[MID]Module
	rwlock sync.RWMutex  //读写锁
}

func (registrar *myRegistrar) Register(module Module) (bool,error) {
	if module == nil {
		return false,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"错误的参数")
	}
	mid := module.ID()
	parts,err := SplitMID(mid)
	if err != nil {
		return  false,err
	}

	moduleType := legalLetterTypeMap[parts[0]]
	if !CheckType(moduleType,module) {
		return false,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"不符合类型的组件")
	}

	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()
	modules := registrar.moduleTypeMap[moduleType]
	if modules == nil {
		modules = map[MID]Module{}
	}
	if _,ok := modules[mid]; ok {
		return false,nil
	}
	modules[mid] = module
	registrar.moduleTypeMap[moduleType] = modules
	return true, nil
}


func (registrar *myRegistrar) Unregister(mid MID) (bool,error) {
	parts,err := SplitMID(mid)
	if err != nil {
		return false,nil
	}
	moduleType := legalLetterTypeMap[parts[0]]
	var deleted bool
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()
	if modules,ok:= registrar.moduleTypeMap[moduleType];ok {
		if _,ok := modules[mid]; ok {
			delete(modules,mid)
			deleted = true
		}
	}

	return deleted,nil
}

//在某个类型的组件数组中 随机的给予一个空闲的组件
func (registrar *myRegistrar) Get(moduleType Type) (Module,error) {
	modules, err := registrar.GetAllByType(moduleType)
	if err != nil {
		return nil,err
	}
	minScore := uint32(0)
	var selectedModule Module
	for _,module := range modules {
		SetScore(module)
		if err != nil {
			return nil,err
		}
		score:= module.Score()
		if minScore == 0 || score<minScore {
			selectedModule = module
			minScore = score
		}
	}
	return selectedModule,nil
}


func (registrar *myRegistrar) GetAllByType(moduleType Type) (map[MID]Module,error) {

	if !LegalType(moduleType) {
		return nil,errors.NewCrawlerError(errors.ERROR_TYPE_PARAMETER,"不存在的类型")
	}
	result := map[MID]Module{}
	registrar.rwlock.RLock()
	defer  registrar.rwlock.RUnlock()
	modules,err := registrar.moduleTypeMap[moduleType]

	if err || len(modules) != 0  {
		for mid,m := range  modules {
			result[mid] = m
		}
		return result,nil
	}
	return result,ErrNotFoundModulInstance

}


func (registrar *myRegistrar) GetAll() map[MID]Module {
	result := map[MID]Module{}
	registrar.rwlock.RLock()
	defer registrar.rwlock.RUnlock()
	for _,modules := range registrar.moduleTypeMap {
		for mid,module := range  modules {
			result[mid] = module
		}
	}
	return result
}


func(registrar *myRegistrar) Clear() {
	registrar.rwlock.Lock()
	defer registrar.rwlock.Unlock()
	registrar.moduleTypeMap = map[Type]map[MID]Module{}
}

























