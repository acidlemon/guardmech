<template>
  <div>
    <h2>Principal Information</h2>

    <template v-if="principal">
      <h3>OpenID Connect Authorization Info</h3>
      <BTable
        v-if="principal.auth_oidc"
        :data="authTableRow"
        :columns="authTableColumns"
      />
      <p v-else>No authorization information.</p>

      <h3>API Keys</h3>
      <BTable
        v-if="principal.auth_oidc"
        :data="authTableRow"
        :columns="authTableColumns"
      />
      <p v-else>No API keys.</p>

      <h3>Groups</h3>
      <BTable
        v-if="principal.auth_oidc"
        :data="authTableRow"
        :columns="authTableColumns"
      />
      <p v-else>No groups.</p>

      <h3>Roles</h3>
      <BTable
        v-if="principal.auth_oidc"
        :data="authTableRow"
        :columns="authTableColumns"
      />
      <p v-else>No roles.</p>

      <h3>Permissions</h3>
      <BTable
        v-if="principal.auth_oidc"
        :data="authTableRow"
        :columns="authTableColumns"
      />
      <p v-else>No permissions.</p>

    </template>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent, SetupContext } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { PrincipalPayload } from '@/types/api'

import BTable, { BTableRow } from '@/components/bootstrap/BTable.vue'

export default defineComponent({
  components: {
    BTable,
  },
  setup(_, context: SetupContext) {
    const router = useRouter()
    const id = router.currentRoute.value.params['id']
    console.log(id)

    const principal = ref<PrincipalPayload>()

    const authTableRow = ref<BTableRow[]>([])
    const authTableColumns = [
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
      const res = await axios.get('/api/principal/' + id)
      const payload = res.data.principal as PrincipalPayload
      principal.value = payload

      if (payload.auth_oidc) {
        authTableRow.value = [{
          label: 'Issuer',
          value: payload.auth_oidc.issuer,
        }, {
          label: 'Subject',
          value: payload.auth_oidc.subject,
        }, {
          label: 'Email',
          value: payload.auth_oidc.email,
        }]
      } else {
        authTableRow.value = []
      }
    })


    return {
      principal,
      authTableRow,
      authTableColumns,
    }

  },
})
</script>
