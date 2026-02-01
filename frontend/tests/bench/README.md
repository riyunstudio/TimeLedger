# 效能基準測試

本目錄包含循環展開和衝突檢測的效能基準測試，確保在大量數據下的效能表現。

## 測試檔案

| 檔案 | 說明 |
|:---|:---|
| `recurrence-bench.test.ts` | 循環展開效能測試 |
| `conflict-detection-bench.test.ts` | 衝突檢測效能測試 |

## 執行基準測試

```bash
# 執行所有基準測試
npm run test:bench

# 或使用 vitest 直接執行
npx vitest run tests/bench/
```

## 測試結果

### 循環展開測試 (`recurrence-bench.test.ts`)

| 測試項目 | 數據量 | 執行時間 | 目標 | 狀態 |
|:---|:---:|:---:|:---:|:---:|
| 展開 500 個循環事件 | 500 | ~2.6ms | <10ms | ✅ |
| 展開 1000 個循環事件 | 1000 | ~5.2ms | <20ms | ✅ |
| 展開 500 個每日事件 | 500 | ~3.5ms | <10ms | ✅ |
| 展開混合頻率事件 | 500 | ~4.9ms | <10ms | ✅ |

### 衝突檢測測試 (`conflict-detection-bench.test.ts`)

| 測試項目 | 數據量 | 執行時間 | 目標 | 狀態 |
|:---|:---:|:---:|:---:|:---:|
| 優化版檢測 500 個時段 | 500 | ~93ms | <100ms | ✅ |
| 優化版檢測 1000 個時段 | 1000 | ~405ms | <500ms | ✅ |
| 基本版檢測 500 個時段 | 500 | ~340ms | <500ms | ✅ |
| 基本版檢測 1000 個時段 | 1000 | ~1370ms | <2000ms | ✅ |

**注意**：優化版使用索引結構，效能顯著優於基本版。建議在生產環境中使用優化版 (`detectAllConflictsOptimized`)。

## 測試場景

### 循環展開測試 (`recurrence-bench.test.ts`)

- **單事件展開**：測試單一循環事件的展開效能
- **多事件展開**：測試 50/100/200/500/1000 個事件的展開效能
- **大規模場景**：模擬月檢視的混合頻率事件展開
- **邊界條件**：空陣列、無循環規則、大間隔事件

### 衝突檢測測試 (`conflict-detection-bench.test.ts`)

- **基本檢測**：O(n²) 算法的衝突檢測效能
- **優化檢測**：使用索引的 O(n) 算法衝突檢測效能
- **版本比較**：基本版 vs 優化版的效能差異
- **多日期場景**：測試跨多天的衝突檢測

## 新增測試案例

如需新增效能測試，請遵循以下模式：

```typescript
it('處理 N 個項目', () => {
  const startTime = performance.now()
  myFunction(generateTestData(N))
  const duration = performance.now() - startTime
  
  console.log(`${N} 個項目處理時間: ${duration.toFixed(2)}ms`)
  expect(duration).toBeLessThan(THRESHOLD_MS)
})
```

## 閾值設定指南

- **循環展開**：<10ms（線性複雜度）
- **衝突檢測優化版**：<100ms（使用索引）
- **衝突檢測基本版**：<500ms（O(n²) 複雜度，僅供參考）

## 注意事項

1. 測試使用 `performance.now()` 測量執行時間
2. 衝突檢測的優化版本使用 Map 索引，顯著提升效能
3. 基本版衝突檢測為 O(n²) 複雜度，不建議用於生產環境
4. 測試結果會輸出到控制台，方便監控效能趨勢
