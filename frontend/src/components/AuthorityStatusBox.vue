<template>
  <BAlert :variant="variant">
    {{ message }}
  </BAlert>
</template>

<script lang="ts">
import { ref, watch, toRef, defineComponent } from 'vue'
import { AuthorityStatus } from '@/hooks/useUserAuthority'

import BAlert from '@/components/bootstrap/BAlert.vue'

type Props = {
  status: AuthorityStatus
}

export default defineComponent({
  name: 'AuthorityStatusBox',
  components: {
    BAlert,
  },
  props: {
    status: {
      type: String as () => AuthorityStatus,
      default: 'NONE',
    }
  },
  setup(props: Props) {
    const message = ref('')
    const variant = ref('')

    const status = toRef(props, 'status')
    watch(status, val => {
      switch (val) {
      case 'OWNER':
        message.value = 'You have owner permission of guardmech. You can operate Guardmech-related membership.'
        variant.value = 'primary'
        break
      case 'WRITABLE':
        message.value = 'You have writable permission of guardmech. You can operate memberships without Guardmech-related.'
        variant.value = 'success'
        break
      case 'READONLY':
        message.value = 'You have read only permission of guardmech. You can only view memberships.'
        variant.value = 'warning'
        break
      case 'NONE':
        message.value = 'You have no permission.'
        variant.value = 'secondary'
        break
      case 'ERROR':
        message.value = 'There is nothing to show because something error occurred.'
        variant.value = 'danger'
        break
      }
    })

    return {
      message,
      variant,
    }
  },
})
</script>
