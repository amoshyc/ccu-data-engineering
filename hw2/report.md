# HW2-2: External Sort

403410034 資工四 黃鈺程

這篇報告只記錄上次繳交後對程式所做的修改與實驗，如果需要我之前的報告，可以請 [這裡](https://gitlab.com/amoshuangyc/ccu/tree/master/data-engineering/hw2)。我將程式碼進行許多修改與重構，並放至與上次不同的 [repo](https://github.com/amoshyc/ccu-data-engineering)。這次我探討了 **平行排序** 對時間的影響，並使用了 **10GB** 的測資來實驗。程式碼仍然是用 **Go** 寫的，用 Go 來寫平行排序真是一大享受。

## 環境

- Go 1.9
- Linux 工作站

工作站每人能使用的空間是有限額的，不過我們發現創建在 `/tmp/` 下的檔案是不被考慮，所以我們測資就放在 `/tmp/` 下，輸出也寫到 `/tmp/` 下。


## 平行排序

對 chunk 排序時，原先是使用 Go 內建的 `sort.Slice`，這次使用平行化的 merge sort 來替代。

```go
func parallelSort(data []string) {
	ch := make(chan []string, 1)
	mergesort(data, ch, 0)
	data = <-ch
}

func mergesort(data []string, out chan []string, dep int) {
	if dep >= 3 {
		sort.Slice(data, func(i, j int) bool {
			return data[i] > data[j]
		})
		out <- data
		return
	}

	N := len(data)
	res1 := make(chan []string, 1)
	res2 := make(chan []string, 1)
	go mergesort(data[:N/2], res1, dep+1)
	go mergesort(data[N/2:], res2, dep+1)

	l, r := <-res1, <-res2
	i, j := 0, 0
	for ix := range data {
		switch {
		case i == len(l):
			data[ix] = r[j]
			j++
		case j == len(r):
			data[ix] = l[i]
			i++
		case l[i] > r[j]:
			data[ix] = l[i]
			i++
		default:
			data[ix] = r[j]
			j++
		}
	}

	out <- data
}
```

我原先還想了另一個版本，是使用回傳值而不是 `chan` 來回傳結果。不過實驗了一下，發現後者會炸記憶體，所以我最後使用了第一個版本：使用 `chan` 來回傳結果。一些實作細節為：不是每次遞迴都要開新的 routine 出來，這反而會比較慢，並使用更多記憶體。正確的做法為只對前幾層開 routine，我這個程式碼是使用前三層，所以最多時會展開成 4 個 routines 同時在跑。


## 結果

### 時間

平行排序對程式的加速是顯著的，以 10GB 資料為列，不使用平行排序的時間為 37 分鐘多，使用後變成 27 分鐘多，簡省了約 10 分鐘，加快了 `37/27 = 1.37` 倍。不管是 `external_merge_sort` 還是 `external_partition_sort`，我兩個程式都跑出來差不多的時間，誤差不到一分鐘。平行排序在小測資時也有明顯的加速，不過增幅就沒這麼誇張。這應該是目前同學中最快的程式，遠遠比沒使用平行排序的 C 快，至於有沒有比使用平行排序的 C 快就不知道了，沒有人去寫這個程式來比較。

### 記憶體使用

記憶體使用量直接跟 chunk 的大小成正比，在 chunk 數為 150 時，記憶體使用量最大為 7%，chunk 數 1000 時，記憶體最大 1.5%。看起來，我的程式的記憶體使用量很穩定，也非常節省。


## 感想

這次重構將程式變得更智慧，api 也更統一了。不管是 `external_merge_sort` 還是 `external_partition_sort`，都不用指定 chunk 的大小，只需指定 chunk 的數量，程式寫出 chunk 時會自動推導每個 chunk 多大。Go 相比 C/C++ 真的讓人很滿意，根據同學使用 Go 的 profiler 的結果，Go 的 IO 速度與 C++ 是相當的，但 Go 在撰寫平行／並行程式時的效率是 C++ 遠遠不及的。

我對我這次實現的外部排序非常滿意，使用VS Code 來寫 Go 也是極為舒爽的，猶如 IDE 中的待遇，程式一存檔，編輯器就顯示出程式哪裡有 CE，也把我的 code format 好，使用 Go 統一的 format 程式 `go fmt`。我本來想放完整程式的，不過因為行數有點多（`external_merge_sort` 297 行，`external_partition_sort` 202 行），所以就請自行參閱我的 [repo](https://github.com/amoshyc/ccu-data-engineering) 吧~