<template>
  <div class="container">
    <h2>Mapping Rule List</h2>
    <AuthorityStatusBox :status="authorityStatus" />
    <section v-if="canWrite">
      <NewMappingRuleModal @completed="created" />
    </section>
    <section v-if="canRead">
      <BTable
        v-if="mappingRules.length"
        :data="mappingRules"
        :columns="columns"
        variant="primary"
      >
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/mapping_rule/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
      <p v-else>There's no mapping rules.</p>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, watch, defineComponent } from 'vue'
import axios from 'axios'
import { MappingRule } from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow, BTableColumn } from '@/components/bootstrap/BTable.vue'
import NewMappingRuleModal from '@/components/modals/NewMappingRuleModal.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    NewMappingRuleModal,
    AuthorityStatusBox,
  },
  setup() {
    const {
      authorityStatus,
      authorityLoadCompleted,
      canWrite,
      canRead,
    } = useUserAuthority()
    
    const mappingRules = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'priority', label: 'Priority' },
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'rule_type', label: 'Rule Type' },
      { key: 'rule_detail', label: 'Description' },
      { key: 'association_type', label: 'Associate With' },
      { key: 'action', label: '' },
    ])

    const ruleTypeLabel: { [index: number]: string} = {
      1: 'Match exactly with Domain',
      2: 'Match with Whole of Domain',
      3: 'Match with OpenID Connect Group',
      4: 'Match with Email Address',
    }

    const fetchList = (async () => {
      const res = await axios.get('/api/mapping_rules').catch(e => e.response)
      const rules = res.data.mapping_rules as MappingRule[]
      mappingRules.value = rules.map(r => {
        return {
          name: r.name,
          description: r.description,
          rule_type: ruleTypeLabel[r.rule_type],
          rule_detail: r.detail,
          association_type: r.association_type,
          association_id: r.association_id,
          priority: r.priority,
        }
      })
    })

    watch(authorityLoadCompleted, (val) => {
      if (val) {
        if (canRead.value) {
          fetchList()
        }
      }
    })


    const created = (() => {
      fetchList()
    })

    return {
      authorityStatus,
      canWrite,
      canRead,
      columns,
      mappingRules,
      created,
    }
  },
})
</script>
