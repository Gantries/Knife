package national

import (
	"fmt"
	"sync"
	"testing"
)

func TestFindOrCreateLocalizerConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	const concurrency = 100 // 并发级别
	const languages = 99    // 语言数量

	// 启动多个协程来并发调用FindOrCreateLocalizer
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			lang := fmt.Sprintf("lang%d", index%languages) // 循环使用不同的语言代码
			localizer := FindOrCreateLocalizer(lang)
			if localizer == nil {
				t.Errorf("FindOrCreateLocalizer returned nil for lang%d", index%languages)
			}
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()

	// 检查是否每个语言代码只创建了一个Localizer实例
	localizerCount := 0
	translators.Range(func(key, value interface{}) bool {
		localizerCount++
		return true
	})
	if localizerCount-1 != languages {
		t.Errorf("Expected %d localizers, but got %d", languages, localizerCount)
	}
}
