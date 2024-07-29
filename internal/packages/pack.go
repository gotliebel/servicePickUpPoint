package packages

var PackageMap = map[string]PackagingStrategy{
	"bag":      NewBag(),
	"box":      NewBox(),
	"wrapping": NewWrapping(),
}
