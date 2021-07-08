<template>
  <BModal
    button-title="Create New"
    modal-type="ok-cancel"
    modal-title="Create New Principal"
    @proceeded="proceeded"
  >
    <BInput v-model="name" label="Name" placeholder="new name" />
    <BInput v-model="description" label="Description" placeholder="write something" />
  </BModal>
</template>

<script lang="ts">
import { ref, defineComponent } from 'vue'
import axios from 'axios'
import BModal from '@/components/bootstrap/BModal.vue'
import BInput from '@/components/bootstrap/BInput.vue'


export default defineComponent({
  name: 'NewPrincipalModal',
  components: {
    BModal,
    BInput,
  },
  setup() {
    const name = ref('')
    const description = ref('')

    const createPrincipal = async (name: string, description: string) => {
      const params = new URLSearchParams({
        name,
        description,
      })
      const res = await axios.post('/api/principal', params)

      console.log(res)
    }

    const proceeded = (() => {
      createPrincipal(name.value, description.value)
    })

    return {
      name,
      description,

      proceeded,
    }
  },
})
</script>
