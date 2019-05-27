package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

/*func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	result := ""
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data := dataRaw.(string)
			result = data
		}),
	}

	ExecutePipeline(hashSignJobs...)
	fmt.Println(result)
}*/

func ExecutePipeline(funcs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	for _, f := range funcs {
		wg.Add(1)
		out := make(chan interface{})
		go func(wg1 *sync.WaitGroup, f1 job, ch1, ch2 chan interface{}) {
			f1(ch1, ch2)
			close(ch2)
			wg1.Done()
		}(wg, f, in, out)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for value := range in {
		d := value.(int)
		data := strconv.Itoa(d)
		tmp := DataSignerMd5(data)
		ch1 := make(chan string)
		ch2 := make(chan string)
		wg.Add(1)
		go func(wg1 *sync.WaitGroup, ch chan interface{}) {
			go func(d string, ch11 chan string) {
				ch11 <- DataSignerCrc32(d)
			}(data, ch1)
			go func(t string, ch22 chan string) {
				ch22 <- DataSignerCrc32(t)
			}(tmp, ch2)
			x1 := <-ch1
			x2 := <-ch2
			ch <- x1 + "~" + x2
			wg.Done()
		}(wg, out)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for value := range in {
		data := value.(string)
		wg.Add(1)
		go func(wg1 *sync.WaitGroup, d string, ch chan interface{}) {
			wg2 := &sync.WaitGroup{}
			mas := make([]string, 6)
			for i := 0; i < 6; i++ {
				wg2.Add(1)
				go func(wg2 *sync.WaitGroup, j int, d string, ar []string) {
					th := strconv.Itoa(j)
					mas[j] = DataSignerCrc32(th + d)
					wg2.Done()
				}(wg2, i, d, mas)
			}
			wg2.Wait()
			ch <- strings.Join(mas, "")
			wg1.Done()
		}(wg, data, out)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	res := ""
	var mas []string
	for value := range in {
		mas = append(mas, value.(string))
	}
	sort.Strings(mas)
	res = strings.Join(mas, "_")
	out <- res
}
