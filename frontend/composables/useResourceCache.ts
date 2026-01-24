// 共享的資源緩存，避免重複 API 調用
const resourceCache = ref<{
  offerings: any[]
  teachers: Map<number, any>
  rooms: Map<number, any>
  loaded: boolean
}>({
  offerings: [],
  teachers: new Map(),
  rooms: new Map(),
  loaded: false,
})

// 使用 ref 來確保在組件間共享同一個 Promise
const fetchPromise = ref<Promise<void> | null>(null)

export function useResourceCache() {
  const { getCenterId } = useCenterId()

  const fetchAllResources = async () => {
    // 如果已經在加載中，返回現有的 promise
    if (fetchPromise.value) {
      return fetchPromise.value
    }

    // 如果已經加載過，直接返回
    if (resourceCache.value.loaded) {
      return
    }

    // 創建新的 promise
    const promise = (async () => {
      try {
        const api = useApi()
        const centerId = getCenterId()

        console.log('開始載入資源資料...')

        const [offeringsRes, teachersRes, roomsRes] = await Promise.all([
          api.get<{ code: number; datas: any }>(`/admin/offerings`),
          api.get<{ code: number; datas: any[] }>('/teachers'),
          api.get<{ code: number; datas: any[] }>(`/admin/rooms`)
        ])

        // 處理 offerings
        resourceCache.value.offerings = offeringsRes.datas?.offerings || []

        // 處理 teachers
        resourceCache.value.teachers = new Map()
        teachersRes.datas?.forEach((t: any) => {
          resourceCache.value.teachers.set(t.id, t)
        })

        // 處理 rooms
        resourceCache.value.rooms = new Map()
        roomsRes.datas?.forEach((r: any) => {
          resourceCache.value.rooms.set(r.id, r)
        })

        resourceCache.value.loaded = true

        console.log('資源資料載入完成:', {
          offerings: resourceCache.value.offerings.length,
          teachers: resourceCache.value.teachers.size,
          rooms: resourceCache.value.rooms.size,
        })
      } catch (error) {
        console.error('Failed to fetch resources:', error)
        resourceCache.value.loaded = false
        fetchPromise.value = null
        throw error
      }
    })()

    fetchPromise.value = promise
    return promise
  }

  const getTeacherName = (teacherId: number): string => {
    const teacher = resourceCache.value.teachers.get(teacherId)
    return teacher?.name || `老師 ${teacherId}`
  }

  const getRoomName = (roomId: number): string => {
    const room = resourceCache.value.rooms.get(roomId)
    return room?.name || `教室 ${roomId}`
  }

  const getOfferingName = (offeringId: number): string => {
    const offering = resourceCache.value.offerings.find(o => o.id === offeringId)
    return offering?.name || `課程 ${offeringId}`
  }

  return {
    resourceCache,
    fetchAllResources,
    getTeacherName,
    getRoomName,
    getOfferingName,
  }
}
