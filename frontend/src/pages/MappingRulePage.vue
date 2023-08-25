<template>
  <div>
    <div class="container information">
      <h2>Mapping Rule Information</h2>
      <AuthorityStatusBox :status="authorityStatus" />
      <BTable
        v-if="mappingRule"
        :data="basicRow"
        :columns="basicColumns"
      />

      <h3>Matcher</h3>
      <ul>
        <li>{{ ruleTypeString }}</li>
        <li>match with <code>{{ mappingRule.detail }}</code></li>
      </ul>

      <section>
        <!-- TODO いいかんじにする -->
        <h3 v-if="mappingRule.association_type == 'group'">Associated Group</h3>
        <h3 v-if="mappingRule.association_type == 'role'">Associated Role</h3>
        <BTable
          v-if="mappingRule"
          :data="associatedRow"
          :columns="basicColumns"
        />
      </section>

    </div>
    <template v-if="mappingRule">
      <div class="danger-zone">
        <div class="container">
          <h3>Danger Zone</h3>
          <DestructionModal
            button-title="Delete This Mapping Rule"
            @confirmDelete="onDelete"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { ref, computed, watch, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import {
  MappingRule,
  MappingRuleTypeSpecificDomain,
  MappingRuleTypeWholeDomain,
  MappingRuleTypeGroupMember,
  MappingRuleTypeSpecificAddress,
} from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import DestructionModal from '@/components/modals/DestructionModal.vue'
// import BButton from '@/components/bootstrap/BButton.vue'
import BTable from '@/components/bootstrap/BTable.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'
import { BTableRow } from '@/types/bootstrap'

export default defineComponent({
  components: {
    // BButton,
    BTable,
    DestructionModal,
    AuthorityStatusBox,
  },
  setup() {
    const router = useRouter()
    const id = router.currentRoute.value.params['id'] as string

    const {
      authorityStatus,
      authorityLoadCompleted,
      canWrite,
      canRead,
    } = useUserAuthority()

    const mappingRule = ref<MappingRule>({
      id: '',
      name: '',
      description: '',
      rule_type: 0,
      detail: '',
      priority: 0,
      association_type: '',
      association_id: '',
    })

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

    const associatedRow = ref<BTableRow[]>([])

    const ruleTypeString = computed(() => {
      switch (mappingRule.value.rule_type) {
      case MappingRuleTypeSpecificDomain:
        return 'Match exactly with Domain'
      case MappingRuleTypeWholeDomain:
        return 'Match with Whole of Domain'
      case MappingRuleTypeGroupMember:
        return 'Match with OpenID Connect Group'
      case MappingRuleTypeSpecificAddress:
        return 'Match with Email Address'
      }
      return ''
    })

    const fetchMappingRule = (async () => {
      const res = await axios.get('/api/mapping_rule/' + id)
      mappingRule.value = res.data.mapping_rule as MappingRule

      if (!mappingRule.value) {
        basicRow.value = []
        return
      }

      if (mappingRule.value.association_type == 'group') {
        const res2 = await axios.get('/api/group/' + mappingRule.value.association_id)

        const associated = res2.data.group
        associatedRow.value = [{
          label: 'Name',
          value: associated.name,
        }, {
          label: 'Description',
          value: associated.description,
        }]
      } else {
        const res2 = await axios.get('/api/role/' + mappingRule.value.association_id)

        const associated = res2.data.role
        associatedRow.value = [{
          label: 'Name',
          value: associated.name,
        }, {
          label: 'Description',
          value: associated.description,
        }]
      }

      basicRow.value = [{
        label: 'Name',
        value: mappingRule.value.name,
      }, {
        label: 'Description',
        value: mappingRule.value.description,
      }]
    })

    watch(authorityLoadCompleted, (val) => {
      if (val) {
        if (canRead.value) {
          fetchMappingRule()
        }
      }
    })

    const onDelete = (() => {
      console.log('onDelete')
      deleteMappingRule()
    })
    const deleteMappingRule = (async () => {
      const res = await axios.delete('/api/mapping_rule/' + id)
      console.log(res)
      // TODO check

      router.push('/mapping_rules')
    })

    return {
      authorityStatus,
      canWrite,
      canRead,
      mappingRule,
      basicRow,
      basicColumns,
      associatedRow,
      ruleTypeString,
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