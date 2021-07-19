<template>
  <BModal
    button-title="Attach Role"
    modal-type="confirm-cancel"
    modal-title="Attach Role"
    confirm-title="Attach"
    @proceeded="proceeded"
  >
    <BSelect v-model="selectedRole" label="Role List" :items="roleList" />
  </BModal>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent, SetupContext } from 'vue'
import axios from 'axios'
import { Role } from '@/types/api'

import BModal from '@/components/bootstrap/BModal.vue'
import BSelect, { BSelectItem } from '@/components/bootstrap/BSelect.vue'

type RoleAttachableType = 'group' | 'principal'

type Props = {
  ownerId: string
  ownerType: RoleAttachableType
  attachedRoles: Role[]
}

export default defineComponent({
  name: 'AttachRoleModal',
  components: {
    BModal,
    BSelect,
  },
  props: {
    ownerId: {
      type: String,
      default: '',
    },
    ownerType: {
      type: String as () => RoleAttachableType,
      default: 'principal',
    },
    attachedRoles: {
      type: Array as () => Role[],
      default: () => [],
    },
  },
  emits: ['completed'],
  setup(props: Props, context: SetupContext) {
    const name = ref('')
    const description = ref('')

    const selectedRole = ref('')

    const attachRole = async (roleId: string) => {
      const params = new URLSearchParams({
        role_id: roleId,
      })

      if (!props.ownerType) {
        console.log('[BUG] you forgot to set ownerType')
        return
      }

      const path = `/api/${props.ownerType}/${props.ownerId}/attach_role`
      const res = await axios.post(path, params).catch(e => e.response)

      console.log(res)
      context.emit('completed')
    }

    const roleList = ref<BSelectItem[]>([])
    const fetchRole = (async () => {
      const fetchRoles = await axios.get('/api/roles').catch(e => e.response)
      roleList.value = fetchRoles.data.roles.map((x: Role) => (
        {
          label: x.name,
          value: x.id,
          disabled: props.attachedRoles.find(r => r.id === x.id) ? true : false,
        }
      ))
    })

    onMounted(() => {
      fetchRole()
    })

    const proceeded = (() => {
      console.log(selectedRole.value)
      attachRole(selectedRole.value)
    })

    return {
      name,
      description,
      roleList,
      selectedRole,

      proceeded,
    }
  },
})
</script>
