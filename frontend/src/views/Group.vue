<template>
  <div>
    <div class="container information">
      <h2>Group Information</h2>
      <BTable
        v-if="group"
        :data="basicRow"
        :columns="basicColumns"
      />

      <h3>Attached Roles</h3>
    </div>

    <template v-if="group">
      <div class="danger-zone">
        <div class="container">
          <h3>Danger Zone</h3>
          <DestructionModal
            button-title="Delete This Group"
            @confirmDelete="onDelete"
          />
        </div>
      </div>
    </template>
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

    const group = ref<Group>()

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
      group.value = res.data.group as Group 

      if (group.value) {
        basicRow.value = [{
          label: 'Name',
          value: group.value.name,
        }, {
          label: 'Description',
          value: group.value.description,
        }]
      } else {
        basicRow.value = []
      }
    })

    const onDelete = (() => {
      deleteGroup()
    })
    const deleteGroup = (async () => {
      const res = await axios.delete('/api/group/' + id)
      console.log(res)
    })

    return {
      group,
      basicRow,
      basicColumns,
      onDelete,
    }
  },
})
</script>

<style scoped>
.information {
  padding-top: 20px;
  padding-bottom: 20px;
}

.danger-zone {
  background: #F8E0E0;
  padding-top: 20px;
  padding-bottom: 30px;
  border-top: dashed 1px #E06060;
}
</style>