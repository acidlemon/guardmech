<template>
  <BModal
    button-title="Attach Group"
    modal-type="confirm-cancel"
    modal-title="Attach Group"
    confirm-title="Attach"
    @proceeded="proceeded"
  >
    <BSelect v-model="selectedGroup" label="Group List" :items="groupList" />
  </BModal>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent, SetupContext } from 'vue'
import axios from 'axios'
import { Group } from '@/types/api'

import BModal from '@/components/bootstrap/BModal.vue'
import BSelect, { BSelectItem } from '@/components/bootstrap/BSelect.vue'

type Props = {
  ownerId: string
  attachedGroups: Group[]
}

export default defineComponent({
  name: 'AttachGroupModal',
  components: {
    BModal,
    BSelect,
  },
  props: {
    ownerId: {
      type: String,
      default: '',
    },
    attachedGroups: {
      type: Array as () => Group[],
      default: () => [],
    },
  },
  emits: ['completed'],
  setup(props: Props, context: SetupContext) {
    const name = ref('')
    const description = ref('')

    const selectedGroup = ref('')

    const attachGroup = async (groupId: string) => {
      const params = new URLSearchParams({
        group_id: groupId,
      })

      const path = `/api/principal/${props.ownerId}/attach_group`
      const res = await axios.post(path, params).catch(e => e.response)

      console.log(res)
      context.emit('completed')
    }

    const groupList = ref<BSelectItem[]>([])
    const fetchGroup = (async () => {
      const fetchGroups = await axios.get('/api/groups').catch(e => e.response)
      groupList.value = fetchGroups.data.groups.map((x: Group) => (
        {
          label: x.name,
          value: x.id,
          disabled: props.attachedGroups.find(r => r.id === x.id) ? true : false,
        }
      ))
    })

    onMounted(() => {
      fetchGroup()
    })

    const proceeded = (() => {
      console.log(selectedGroup.value)
      attachGroup(selectedGroup.value)
    })

    return {
      name,
      description,
      groupList,
      selectedGroup,

      proceeded,
    }
  },
})
</script>
