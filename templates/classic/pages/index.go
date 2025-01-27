package classic_pages

import (
	"github.com/eriicafes/tmpl"
	"github.com/eriicafes/tmplist/db"
)

type Index struct {
	Layout
	Topics []db.Topic
	Search string
}

func (i Index) Template() (string, any) {
	return tmpl.Tmpl("classic/pages/index", i.Layout, i).Template()
}

func (i Index) EmptyCells() []int {
	maxEmptyCells := 7
	if len(i.Topics) >= maxEmptyCells {
		return nil
	}
	var cells []int
	for i := range maxEmptyCells - len(i.Topics) {
		cells = append(cells, i)
	}
	return cells
}

func (i Index) Gradient(id int) string {
	// based on the id, return a tailwind gradient color
	gradients := []string{
		"from-purple-500 to-pink-500",
		"from-green-400 to-blue-500",
		"from-yellow-400 to-red-500",
		"from-blue-400 to-indigo-500",
		"from-red-400 to-yellow-500",
		"from-pink-400 to-purple-500",
		"from-indigo-400 to-green-500",
		"from-teal-400 to-cyan-500",
		"from-orange-400 to-red-500",
		"from-lime-400 to-green-500",
		"from-amber-400 to-yellow-500",
		"from-emerald-400 to-teal-500",
		"from-sky-400 to-blue-500",
		"from-rose-400 to-pink-500",
		"from-fuchsia-400 to-purple-500",
		"from-violet-400 to-indigo-500",
		"from-cyan-400 to-teal-500",
		"from-lime-500 to-green-600",
		"from-amber-500 to-orange-600",
		"from-emerald-500 to-teal-600",
		"from-sky-500 to-blue-600",
		"from-rose-500 to-pink-600",
		"from-fuchsia-500 to-purple-600",
		"from-violet-500 to-indigo-600",
		"from-cyan-500 to-teal-600",
		"from-lime-600 to-green-700",
		"from-amber-600 to-orange-700",
		"from-emerald-600 to-teal-700",
		"from-sky-600 to-blue-700",
		"from-rose-600 to-pink-700",
		"from-fuchsia-600 to-purple-700",
		"from-violet-600 to-indigo-700",
		"from-cyan-600 to-teal-700",
		"from-lime-700 to-green-800",
		"from-amber-700 to-orange-800",
		"from-emerald-700 to-teal-800",
		"from-sky-700 to-blue-800",
		"from-rose-700 to-pink-800",
		"from-fuchsia-700 to-purple-800",
		"from-violet-700 to-indigo-800",
		"from-cyan-700 to-teal-800",
		"from-lime-800 to-green-900",
		"from-amber-800 to-orange-900",
		"from-emerald-800 to-teal-900",
		"from-sky-800 to-blue-900",
		"from-rose-800 to-pink-900",
		"from-fuchsia-800 to-purple-900",
		"from-violet-800 to-indigo-900",
		"from-cyan-800 to-teal-900",
	}
	return gradients[id%len(gradients)]
}
