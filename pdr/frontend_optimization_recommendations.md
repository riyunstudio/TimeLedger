# 🚀 TimeLedger 前端優化與重構建議報告

本報告針對 `TimeLedger` 前端專案，從**可維護性**、**高復用性**、**效能**與**使用者體驗**四個維度，提供了 11 項具體的優化建議與重構計畫。

## 1. 核心瓶頸：重構「上帝組件」(God Component) - ScheduleGrid.vue
*   **現況分析**：`ScheduleGrid.vue` 超過 1300 行，承載了週導航、篩選過濾、手機/桌機視圖切換、拖拽邏輯及多個 Modal 彈窗。這導致邏輯耦合嚴重，難以測試與維護。
*   **建議方向**：依據職責拆分組件，將 UI 與 業務邏輯分離。
*   **Cursor 指令**：
    > `Refactor ScheduleGrid.vue: 1. Create a specialized directory components/Schedule/. 2. Extract CalendarHeader, WeekGrid, ScheduleCard, and overlap dialogs into separate SFC files. 3. Use props and emits for communication.`

## 2. 狀態管理去中心化：拆分「上帝 Store」- teacher.ts
*   **現況分析**：`useTeacherStore` 高達 760 行，融合了認證、排課、個人行程、技能、證照、通知等不同領域的資料。
*   **建議方向**：依據 Domain 拆分為多個專注的小 Store（如 `useScheduleStore`, `useProfileStore`, `useNotificationStore`）。
*   **Cursor 指令**：
    > `Refactor stores/teacher.ts: Split it into useProfileStore.ts (for skills/certs/profile), useScheduleStore.ts (for schedules/exceptions), and useNotificationStore.ts. Update all component imports.`

## 3. 採用目錄特徵分層 (Feature-based Structure)
*   **現況分析**：`components/` 目錄偏向扁平，隨著組件增加會變得難以導航。
*   **建議方向**：將組件按功能模組分組（例如 `components/auth/`, `components/scheduling/`, `components/shared/`）。
*   **Cursor 指令**：
    > `Reorganize components/ folder: Group components into feature directories: Admin, Teacher, Navigation, Scheduling, and Base. Ensure all auto-imports in Nuxt are preserved.`

## 4. 業務邏輯抽離至 Composables (邏輯複用)
*   **現況分析**：部分覆雜邏輯（如 `teacher.ts` 中的週期性行程展開邏輯）被硬編碼在 Store 中。
*   **建議方向**：將純數據處理逻辑 (Pure Logic) 抽離至 Composables，便於跨組件/Store 複用並進行單元測試。
*   **Cursor 指令**：
    > `Extract recurrence expansion logic from stores/teacher.ts (lines 280-340) into a new composable composables/useRecurrence.ts. Add unit tests for this logic.`

## 5. 強化 TypeScript 型別嚴謹度
*   **現況分析**：專案中多處使用 `any`（例如 `ScheduleGrid.vue` 的 `schedules?: any[]`），這抵消了 TS 的維護優勢。
*   **建議方向**：對 API 響應數據與 Store 狀態定義嚴格的 Interface。
*   **Cursor 指令**
    > `Audit types/index.ts and all files in components/: Replace 'any' with specific interfaces. Define missing models for ScheduleItem and API responses.`

## 6. 標準化彈窗管理 (Standardized Modals)
*   **現況分析**：多個 Modal（AddSkill, AddCertificate 等）存在高度重覆的樣版代碼 (Boilerplate) 與 Teleport 邏輯。
*   **建議方向**：建立一個 `BaseModal` 核心組件，或使用專屬的 Modal Store 管理彈窗狀態。
*   **Cursor 指令**：
    > `Create a components/base/BaseModal.vue with consistent transitions, backdrop, and close logic. Refactor AddSkillModal.vue and AddCertificateModal.vue to extend BaseModal.`

## 7. 高頻 UI 效能優化 (Schedule Grid Performance)
*   **現況分析**：排課網格在課程量大時，樣式計算（getScheduleStyle）與 Reactivity 可能造成渲染延遲。
*   **建議方向**：對複雜的樣式計算使用 `computed` 緩存，並在 `DynamicScroller` 內部使用 `v-memo` 優化。
*   **Cursor 指令**：
    > `Optimize ScheduleGrid.vue rendering: 1. Memoize getScheduleStyle results. 2. Ensure only visible cards are being re-evaluated during scroll. 3. Check for unnecessary re-renders in the grid template.`

## 8. 統一錯誤處理與使用者反饋 (UX Improvement)
*   **現況分析**：雖然 `useApi` 有基礎處理，但前端 UI 對不同錯誤（如：權限重疊或網絡超時）的反饋方式不夠統一。
*   **建議方向**：實作一個全域的錯誤攔截 UI 跟隨系統，提供更具引導性的錯誤提示（Actionable Feedback）。
*   **Cursor 指令**：
    > `Enhance UI feedback: Create a global error alerting system that maps API error codes to user-friendly messages and actions, integrated with useToast or GlobalAlert.`

## 9. 資產優化與 SVG 圖示標準化
*   **現況分析**：組件內直接嵌入大量內聯 SVG 代碼，影響代碼可讀性且無法統一管理樣式。
*   **建議方向**：將所有圖示標準化，或使用 `nuxt-icon` 模組統一管理。
*   **Cursor 指令**：
    > `Extract all inline SVGs in components/ into a library or use nuxt-icon. Ensure icons support consistent sizing and color through Tailwind classes.`

## 10. 引入自動化效能審計 (Continuous Maintainability)
*   **現況分析**：目前缺乏量化的效能指標，難以防止後續開發造成的效能退化。
*   **建議方向**：在測試流程中加入 ` Lighthouse CI` 或 `Vitest Bench` 針對核心邏輯（如排課衝突檢測）進行效能測試。
*   **Cursor 指令**：
    > `Add a performance benchmark test using Vitest in tests/bench/ for the recurrence expansion and conflict detection logic to ensure it stays below 10ms for 500+ items.`

## 11. 實作單元測試與組件測試
*   **現況分析**：雖然有 `playwright` 和 `vitest` 配置，但現有測試較少，尤其是針對複雜狀態邏輯的測試。
*   **建議方向**：針對 Store 中的業務邏輯和關鍵 UI 組件實作測試。
*   **Cursor 指令**：
    > `Generate Vitest tests for stores/schedule.ts (once refactored) and components/Schedule/WeekGrid.vue to ensure correct coordinate mapping and event handling.`

---

## 🛠️ 推薦重構執行順序 (Task List)

1.  **[High Priority]** 拆分 `useTeacherStore` - 解決數據架構混亂。
2.  **[High Priority]** 拆分 `ScheduleGrid.vue` - 解決 UI 開發瓶頸。
3.  **[Medium Priority]** 標準化 `BaseModal` 與圖示 - 統一視覺與開發規範。
4.  **[Medium Priority]** 型別補完 (Kill the 'any') - 增強系統穩定性。
5.  **[Low Priority]** 實施效能測試與 Composables 邏輯抽離。
