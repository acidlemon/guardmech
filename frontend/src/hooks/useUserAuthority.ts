import axios from 'axios'
import { ref, computed, onMounted } from 'vue'

export type AuthorityStatus = 'OWNER' | 'WRITABLE' | 'READONLY' | 'NONE' | 'ERROR' | ''

export const useUserAuthority = () => {

  const authorityStatus = ref<AuthorityStatus>('')
  const authorityLoadCompleted = ref(false)

  onMounted(() => {
    axios.get('/api/authority').then((res) => {
      const status = res.data.authority
      if (!status) {
        authorityStatus.value = 'ERROR'
        authorityLoadCompleted.value = true
        return
      }

      authorityStatus.value = status
      authorityLoadCompleted.value = true
    }).catch(e => {
      console.log(e)
      // session expired, need refresh
      location.reload()
    })
  })

  const isOwner = computed(() => authorityStatus.value === 'OWNER')
  const canWrite = computed(() => ['OWNER', 'WRITABLE'].includes(authorityStatus.value))
  const canRead = computed(() => ['OWNER', 'WRITABLE', 'READONLY'].includes(authorityStatus.value))

  return {
    authorityStatus,
    authorityLoadCompleted,
    isOwner,
    canWrite,
    canRead,
  }
}

