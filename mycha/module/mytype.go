package module

type Type string


//当前的所有的合法的类型
const  (
	TYPE_DOWNLOADER Type = "downloader"
	TYPE_ANALYZER Type = "analyzer"
	TYPE_PIPELINE Type = "pipeline"
)

//当前类型的类型的映射
var legalTypeLetterMap = map[Type]string{
	TYPE_DOWNLOADER:"D",
	TYPE_ANALYZER:"A",
	TYPE_PIPELINE:"P",
}
//当前类型的简写映射
var legalLetterTypeMap = map[string]Type {
	"D":TYPE_DOWNLOADER,
	"A":TYPE_ANALYZER,
	"P":TYPE_PIPELINE,
}

//类型断言来判断
func CheckType(moduleType Type,module Module) bool {
	if moduleType == "" || module == nil {
		return false
	}
	switch moduleType {
	case TYPE_DOWNLOADER:
		if _,ok := module.(Downloader); ok {
			return true
		}
	case TYPE_ANALYZER:
		if _,ok := module.(Analyzer); ok {
			return true
		}
	case TYPE_PIPELINE:
		if _,ok := module.(Pipeline); ok {
			return true
		}
		}
	return  false
}

//无知道这样的判断到底安全在哪里？
func LegalType(moduleType Type) bool {
	if _,ok := legalTypeLetterMap[moduleType]; ok {
		return true
	}
	return false
}
//根据MID获取到当前组件的类型
func GetType(mid MID) (bool,Type) {
	parts,err := SplitMID(mid)
	if err != nil {
		return false,""
	}
	mt,ok := legalLetterTypeMap[parts[0]]
	return ok,mt
}


//用于获取到字母代号
func getLetter(moduleType Type) (bool,string) {
	var letter string
	var found bool
	for l,t := range legalLetterTypeMap {
		if t== moduleType {
			letter = l
			found = true
			break
		}
	}
	return found,letter
}


// typeToLetter 用于根据给定的组件类型获得其字母代号。
// 若给定的组件类型不合法，则第一个结果值会是false。
func typeToLetter(moduleType Type) (bool, string) {
	switch moduleType {
	case TYPE_DOWNLOADER:
		return true, "D"
	case TYPE_ANALYZER:
		return true, "A"
	case TYPE_PIPELINE:
		return true, "P"
	default:
		return false, ""
	}
}

// letterToType 用于根据字母代号获得对应的组件类型。
// 若给定的字母代号不合法，则第一个结果值会是false。
func letterToType(letter string) (bool, Type) {
	switch letter {
	case "D":
		return true, TYPE_DOWNLOADER
	case "A":
		return true, TYPE_ANALYZER
	case "P":
		return true, TYPE_PIPELINE
	default:
		return false, ""
	}
}






























