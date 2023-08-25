<template>
  <BModal
    button-title="Create New"
    modal-type="confirm-cancel"
    modal-title="Create New Mapping Rule"
    @proceeded="proceeded"
  >
    
    <BInput v-model="name" label="Name" placeholder="new name" />
    <BInput v-model="description" label="Description" placeholder="write something" />
    <BSelect v-model="ruleType" label="Rule Type" :items="ruleTypeItems" />

    <template v-if="ruleType === '1'">
      <BInput v-model="detail" label="Email Address Domain" placeholder="example.com" pre-text="@" />
    </template>
    <template v-else-if="ruleType === '2'">
      <BInput v-model="detail" label="Email Address Suffix" placeholder="example.com" />
    </template>
    <template v-else-if="ruleType === '3'">
      <BInput v-model="detail" label="OpneID Connect Providing Group" placeholder="group@example.com" />
    </template>
    <template v-else-if="ruleType === '4'">
      <BInput v-model="detail" label="Email Address" placeholder="john@example.com" />
    </template>

    <BCheckList v-model="targetType" :items="targetItems" check-type="radio" name="association-type" label="Association Type" />
    <BSelect v-model="associationId" label="Candidate List" :items="targetList" />

  </BModal>
</template>

<script lang="ts">
import { ref, computed, onMounted, defineComponent, SetupContext } from 'vue'
import axios from 'axios'
import BModal from '@/components/bootstrap/BModal.vue'
import BInput from '@/components/bootstrap/BInput.vue'
import BCheckList from '@/components/bootstrap/BCheckList.vue'
import BSelect, { BSelectItem } from '@/components/bootstrap/BSelect.vue'
import { Group, Role } from '@/types/api'

export default defineComponent({
  name: 'NewMappingRuleModal',
  components: {
    BModal,
    BInput,
    BCheckList,
    BSelect,
  },
  emits: ['completed'],
  setup(_, context: SetupContext) {
    const name = ref('')
    const description = ref('')
    const ruleType = ref('1')
    const detail = ref('')
    const associationId = ref('')

    const createMappingRule = async (
      name: string,
      description: string,
      rule_type: string,
      detail: string,
      association_type: string,
      association_id: string,
    ) => {
      const params = new URLSearchParams({
        name,
        description,
        rule_type,
        detail,
        association_type,
        association_id,
      })
      const res = await axios.post('/api/mapping_rule', params).catch(e => e.response)

      console.log(res)
      context.emit('completed')
    }

    const proceeded = (() => {
      createMappingRule(name.value, description.value, ruleType.value, detail.value, targetType.value, associationId.value)
    })

    const ruleTypeItems: BSelectItem[] = [
      {
        label: 'Match exactly with Domain',
        tips: 'Specific Email Address Domain (end with @example.com)',
        value: '1',
      },
      {
        label: 'Match with Whole of Domain',
        tips: 'Whole Email Address Domain (end with example.com, including @sub.example.com)',
        value: '2',
      },
      {
        label: 'Match with OpenID Connect Group',
        tips: 'OpenID Connect Provider\'s Group',
        value: '3',
      },
      {
        label: 'Match with Email Address',
        tips: 'Specific Email Address (like john@example.com)',
        value: '4',
      },
    ]

    const groupList = ref<BSelectItem[]>([])
    const roleList = ref<BSelectItem[]>([])

    const fetchGroupAndRole = (async () => {
      const fetchGroups = axios.get('/api/groups').catch(e => e.response)
      const fetchRoles = axios.get('/api/roles').catch(e => e.response)

      const [groups, roles] = await Promise.all([fetchGroups, fetchRoles])

      groupList.value = groups.data.groups.map((x: Group) => ({ label: x.name, value: x.id }))
      roleList.value = roles.data.roles.map((x: Role) => ({ label: x.name, value: x.id }))
    })

    onMounted(() => {
      fetchGroupAndRole()
    })

    const targetType = ref('group')
    const targetItems = [
      {
        label: 'Group',
        value: 'group',
        defaultChecked: true,
      },
      {
        label: 'Role',
        value: 'role',
        defaultChecked: false,
      },
    ]
    const targetList = computed<BSelectItem[]>(() => {
      if (targetType.value === 'group') 
        return groupList.value

      return roleList.value
    })

    return {
      name,
      description,
      ruleType,
      detail,
      associationId,
      ruleTypeItems,
      targetType,
      targetItems,
      targetList,
      proceeded,
    }
  },
})
</script>
