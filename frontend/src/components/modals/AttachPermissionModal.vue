<template>
  <BModal
    button-title="Attach Permission"
    modal-type="confirm-cancel"
    modal-title="Attach Permission"
    confirm-title="Attach"
    @proceeded="proceeded"
  >
    <BSelect v-model="selectedPermission" label="Permission List" :items="permissionList" />
  </BModal>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent, SetupContext } from 'vue'
import axios from 'axios'
import { Permission } from '@/types/api'

import BModal from '@/components/bootstrap/BModal.vue'
import BSelect, { BSelectItem } from '@/components/bootstrap/BSelect.vue'

type Props = {
  ownerId: string
  attachedPermissions: Permission[]
}

export default defineComponent({
  name: 'AttachPermissionModal',
  components: {
    BModal,
    BSelect,
  },
  props: {
    ownerId: {
      type: String,
      default: '',
    },
    attachedPermissions: {
      type: Array as () => Permission[],
      default: () => [],
    },
  },
  emits: ['completed'],
  setup(props: Props, context: SetupContext) {
    const name = ref('')
    const description = ref('')

    const selectedPermission = ref('')

    const attachPermission = async (permissionId: string) => {
      const params = new URLSearchParams({
        permission_id: permissionId,
      })

      const path = `/api/role/${props.ownerId}/attach_permission`
      const res = await axios.post(path, params).catch(e => e.response)

      console.log(res)
      context.emit('completed')
    }

    const permissionList = ref<BSelectItem[]>([])
    const fetchPermission = (async () => {
      const fetchPermissions = await axios.get('/api/permissions').catch(e => e.response)
      permissionList.value = fetchPermissions.data.permissions.map((x: Permission) => (
        {
          label: x.name,
          value: x.id,
          disabled: props.attachedPermissions.find(r => r.id === x.id) ? true : false,
        }
      ))
    })

    onMounted(() => {
      fetchPermission()
    })

    const proceeded = (() => {
      console.log(selectedPermission.value)
      attachPermission(selectedPermission.value)
    })

    return {
      name,
      description,
      permissionList,
      selectedPermission,

      proceeded,
    }
  },
})
</script>
