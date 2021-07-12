<template>
  <div>
    <h2>Single Group</h2>

    <h3>Basic Information</h3>
    <BTable
      v-if="data"
      :data="basicRow"
      :columns="basicColumns"
    />

    <h3>Attached Roles</h3>

    <h3>Danger Zone</h3>
    <DestructionModal
      button-title="Delete This Group"
      @confirmDelete="onDelete"
    />

  </div>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

import { Group } from '@/types/api'
import DestructionModal from '@/components/modals/DestructionModal.vue'
import BTable, { BTableRow } from '@/components/bootstrap/BTable.vue'

export default defineComponent({
  components: {
    BTable,
    DestructionModal,
  },
  setup() {
    const router = useRouter()
    const id = router.currentRoute.value.params['id'] as string

    const data = ref<Group>()

    const basicRow = ref<BTableRow[]>([])
    const basicColumns = [
      {
        key: 'label',
        label: 'Item',
      },
      {
        key: 'value',
        label: 'Data',
      },
    ]

    onMounted(async () => {
      const res = await axios.get('/api/group/' + id)
      data.value = res.data.group as Group 

      if (data.value) {
        basicRow.value = [{
          label: 'Name',
          value: data.value.name,
        }, {
          label: 'Description',
          value: data.value.description,
        }]
      } else {
        basicRow.value = []
      }
    })

    const onDelete = (() => {
      deleteGroup()
    })
    const deleteGroup = (async () => {
      console.log('delete')
      console.log(id)
      const res = await axios.delete('/api/group/' + id)
      console.log(res)
    })

    return {
      data,
      basicRow,
      basicColumns,
      onDelete,
    }
  },
})
</script>
