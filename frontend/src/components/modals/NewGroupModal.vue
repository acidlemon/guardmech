<template>
  <BModal
    button-title="Create New"
    modal-type="confirm-cancel"
    modal-title="Create New Group"
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
  name: 'NewGroupModal',
  components: {
    BModal,
    BInput,
  },
  setup() {
    const name = ref('')
    const description = ref('')

    const createGroup = async (name: string, description: string) => {
      const params = new URLSearchParams({
        name,
        description,
      })
      const res = await axios.post('/api/group', params)

      console.log(res)
    }

    const proceeded = (() => {
      createGroup(name.value, description.value)
    })

    return {
      name,
      description,

      proceeded,
    }
  },
})
</script>
