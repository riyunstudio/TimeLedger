import { defineStore } from 'pinia'
import type { TeacherSkill, TeacherCertificate, Teacher, Hashtag } from '~/types'
import { withLoading } from '~/utils/loadingHelper'

export interface BackgroundImage {
  path: string
  filename: string
  url: string
}

export const useProfileStore = defineStore('profile', () => {
  // 資料狀態
  const skills = ref<TeacherSkill[]>([])
  const certificates = ref<TeacherCertificate[]>([])
  const profile = ref<Teacher | null>(null)
  const backgrounds = ref<BackgroundImage[]>([])

  // Loading 狀態
  const isFetching = ref(false)
  const isUpdating = ref(false)
  const isCreatingSkill = ref(false)
  const isUpdatingSkill = ref(false)
  const isDeletingSkill = ref(false)
  const isCreatingCertificate = ref(false)
  const isDeletingCertificate = ref(false)
  const isFetchingBackgrounds = ref(false)
  const isUploadingBackground = ref(false)
  const isDeletingBackground = ref(false)

  // 技能操作
  const fetchSkills = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const response = await api.get<{ code: number; message: string; datas: TeacherSkill[] }>('/teacher/me/skills')
        skills.value = response.datas || []
      } catch (error) {
        console.error('Failed to fetch skills:', error)
        throw error
      }
    })
  }

  const createSkill = async (data: { category: string; skill_name: string; hashtag_ids?: number[] }) => {
    return withLoading(isCreatingSkill, async () => {
      const api = useApi()
      const response = await api.post<{ code: number; message: string; datas: TeacherSkill }>('/teacher/me/skills', data)
      skills.value.push(response.datas)
      return response.datas
    })
  }

  const updateSkill = async (skillId: number, data: { category: string; skill_name: string; hashtags?: string[] }) => {
    return withLoading(isUpdatingSkill, async () => {
      const api = useApi()
      const response = await api.put<{ code: number; message: string; datas: TeacherSkill }>(`/teacher/me/skills/${skillId}`, data)
      const index = skills.value.findIndex(s => s.id === skillId)
      if (index !== -1) {
        skills.value[index] = response.datas
      }
      return response.datas
    })
  }

  const deleteSkill = async (skillId: number) => {
    return withLoading(isDeletingSkill, async () => {
      const api = useApi()
      await api.delete(`/teacher/me/skills/${skillId}`)
      skills.value = skills.value.filter(s => s.id !== skillId)
    })
  }

  // 證照操作
  const fetchCertificates = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const response = await api.get<{ code: number; message: string; datas: TeacherCertificate[] }>('/teacher/me/certificates')
        certificates.value = response.datas || []
      } catch (error) {
        console.error('Failed to fetch certificates:', error)
        throw error
      }
    })
  }

  const createCertificate = async (data: {
    name: string
    file_url?: string
    issued_at?: string
  }) => {
    return withLoading(isCreatingCertificate, async () => {
      const api = useApi()
      const response = await api.post<{ code: number; message: string; datas: TeacherCertificate }>('/teacher/me/certificates', data)
      certificates.value.push(response.datas)
      return response.datas
    })
  }

  const deleteCertificate = async (certId: number) => {
    return withLoading(isDeletingCertificate, async () => {
      const api = useApi()
      await api.delete(`/teacher/me/certificates/${certId}`)
      certificates.value = certificates.value.filter(c => c.id !== certId)
    })
  }

  // 個人檔案操作
  const fetchProfile = async () => {
    return withLoading(isFetching, async () => {
      try {
        const api = useApi()
        const response = await api.get<{ code: number; message: string; datas: Teacher }>('/teacher/me/profile')
        profile.value = response.datas || null
        return response.datas
      } catch (error) {
        console.error('Failed to fetch profile:', error)
        throw error
      }
    })
  }

  const updateProfile = async (data: Partial<Teacher>) => {
    return withLoading(isUpdating, async () => {
      const api = useApi()
      const response = await api.put<{ code: number; message: string; datas: Teacher }>('/teacher/me/profile', data)
      profile.value = response.datas || null
      return response.datas
    })
  }

  // 標籤相關方法
  const searchHashtags = async (query: string): Promise<Hashtag[]> => {
    if (!query || query.length < 1) {
      return []
    }

    try {
      const api = useApi()
      const response = await api.get<{ code: number; datas: Hashtag[] }>('/hashtags/search', { q: query })
      return response.datas || []
    } catch (error) {
      console.error('Failed to search hashtags:', error)
      return []
    }
  }

  const createHashtag = async (name: string): Promise<Hashtag | null> => {
    if (!name || name.length < 2) {
      return null
    }

    try {
      const api = useApi()
      const response = await api.post<{ code: number; datas: Hashtag }>('/hashtags', { name: '#' + name })
      return response.datas || null
    } catch (error) {
      console.error('Failed to create hashtag:', error)
      return null
    }
  }

  const processHashtag = async (tagName: string): Promise<number | null> => {
    if (!tagName || tagName.length < 2) {
      return null
    }

    try {
      const searchResults = await searchHashtags(tagName)
      const existing = searchResults.find(h => h.name === tagName || h.name === '#' + tagName)

      if (existing) {
        return existing.id
      }

      const newHashtag = await createHashtag(tagName)
      return newHashtag?.id || null
    } catch (error) {
      console.error('Failed to process hashtag:', tagName, error)
      return null
    }
  }

  // 背景圖片操作
  const fetchBackgrounds = async () => {
    return withLoading(isFetchingBackgrounds, async () => {
      try {
        const api = useApi()
        const config = useRuntimeConfig()
        const response = await api.get<string[]>('/teacher/me/backgrounds')
        // 轉換為完整的 URL 陣列
        backgrounds.value = (response || []).map((path: string) => ({
          path,
          filename: path.split('/').pop() || path,
          url: `${config.public.apiBase}/${path}`,
        }))
      } catch (error) {
        console.error('Failed to fetch backgrounds:', error)
        throw error
      }
    })
  }

  const uploadBackground = async (file: File): Promise<BackgroundImage> => {
    return withLoading(isUploadingBackground, async () => {
      const api = useApi()
      const config = useRuntimeConfig()
      const response = await api.upload<{ image_path: string }>('/teacher/me/backgrounds', file)

      const newImage: BackgroundImage = {
        path: response.image_path,
        filename: response.image_path.split('/').pop() || response.image_path,
        url: `${config.public.apiBase}/${response.image_path}`,
      }
      backgrounds.value.push(newImage)
      return newImage
    })
  }

  const deleteBackground = async (path: string) => {
    return withLoading(isDeletingBackground, async () => {
      const api = useApi()
      await api.delete(`/teacher/me/backgrounds?path=${encodeURIComponent(path)}`)
      backgrounds.value = backgrounds.value.filter(b => b.path !== path)
    })
  }

  return {
    // 資料狀態
    skills,
    certificates,
    profile,
    backgrounds,

    // Loading 狀態
    isFetching,
    isUpdating,
    isCreatingSkill,
    isUpdatingSkill,
    isDeletingSkill,
    isCreatingCertificate,
    isDeletingCertificate,
    isFetchingBackgrounds,
    isUploadingBackground,
    isDeletingBackground,

    // 技能方法
    fetchSkills,
    createSkill,
    updateSkill,
    deleteSkill,

    // 證照方法
    fetchCertificates,
    createCertificate,
    deleteCertificate,

    // 個人檔案方法
    fetchProfile,
    updateProfile,

    // 標籤方法
    searchHashtags,
    createHashtag,
    processHashtag,

    // 背景圖片方法
    fetchBackgrounds,
    uploadBackground,
    deleteBackground,
  }
}, {
  persist: {
    key: 'timeledger-profile',
    paths: ['skills', 'certificates', 'profile'],
  },
})
