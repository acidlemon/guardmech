<template>
  <div>
    <div class="container information">
      <h2>Principal Information</h2>

      <template v-if="principal">
        <h3>OpenID Connect Authorization Info</h3>
        <BTable
          v-if="principal.auth_oidc"
          :data="authTableRows"
          :columns="authTableColumns"
        />
        <p v-else>No authorization information.</p>

        <h3>API Keys</h3>
        <BTable
          v-if="apiKeyTableRows.length"
          :data="apiKeyTableRows"
          :columns="standardColumns"
        />
        <p v-else>No API keys.</p>

        <h3>Attached Groups</h3>
        <AttachGroupModal
          :owner-id="principal.principal.id"
          :attached-groups="attachedGroups"
          @completed="needRefresh"
        />
        <BTable
          v-if="groupTableRows.length"
          :data="groupTableRows"
          :columns="standardColumns"
        >
          <template #cell(action)="data">
            <!-- TODO: confirm modal -->
            <BButton variant="danger" @click="onDetachGroup(data.row.id)">Detach</BButton>
          </template>
        </BTable>
        <p v-else>No groups.</p>

        <h3>Attached Roles</h3>
        <AttachRoleModal
          owner-type="principal"
          :owner-id="principal.principal.id"
          :attached-roles="attachedRoles"
          @completed="needRefresh"
        />
        <BTable
          v-if="roleTableRows.length"
          :data="roleTableRows"
          :columns="standardColumns"
        >
          <template #cell(action)="data">
            <!-- TODO: confirm modal -->
            <BButton variant="danger" @click="onDetachRole(data.row.id)">Detach</BButton>
          </template>
        </BTable>
        <p v-else>No roles.</p>

        <h3>Having Permissions</h3>
        <BTable
          v-if="permissionTableRows.length"
          :data="permissionTableRows"
          :columns="standardColumns"
        />
        <p v-else>No permissions.</p>
      </template>
    </div>
    <template v-if="principal">
      <div class="danger-zone">
        <div class="container">
          <h3>Danger Zone</h3>
          <DestructionModal
            button-title="Delete This Principal"
            @confirmDelete="onDelete"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { ref, computed, onMounted, defineComponent, SetupContext } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { PrincipalPayload, Group, Role } from '@/types/api'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow } from '@/components/bootstrap/BTable.vue'
import AttachGroupModal from '@/components/modals/AttachGroupModal.vue'
import AttachRoleModal from '@/components/modals/AttachRoleModal.vue'
import DestructionModal from '@/components/modals/DestructionModal.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    AttachGroupModal,
    AttachRoleModal,
    DestructionModal,
  },
  setup(_, context: SetupContext) {
    const router = useRouter()
    const id = router.currentRoute.value.params['id']

    const principal = ref<PrincipalPayload>()

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
    const authTableRows = ref<BTableRow[]>([])

    const standardColumns = [
      {
        key: 'name',
        label: 'Name',
      },
      {
        key: 'description',
        label: 'Description',
      },
      {
        key: 'action',
        label: '',
      },
    ]
    const apiKeyTableRows = ref<BTableRow[]>([])
    const groupTableRows = ref<BTableRow[]>([])
    const roleTableRows = ref<BTableRow[]>([])
    const permissionTableRows = ref<BTableRow[]>([])

    const fetchPrincipal = (async () => {
      const res = await axios.get('/api/principal/' + id)
      const payload = res.data.principal as PrincipalPayload
      principal.value = payload

      if (payload.auth_oidc) {
        authTableRows.value = [{
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
        authTableRows.value = []
      }

      apiKeyTableRows.value = payload.auth_apikeys
      groupTableRows.value = payload.groups
      roleTableRows.value = payload.roles
      permissionTableRows.value = payload.permissions
    })

    onMounted(() => {
      fetchPrincipal()
    })

    const attachedGroups = computed<Group[]>(() => principal.value ? principal.value.groups : [])
    const attachedRoles = computed<Role[]>(() => principal.value ? principal.value.roles : [])

    const onDelete = (() => {
      deletePrincipal()
    })

    const needRefresh = (() => {
      fetchPrincipal()
    })

    const onDetachGroup = ((groupId: string) => {
      detachGroup(groupId)
    })

    const detachGroup = (async (groupId: string) => {
      console.log(groupId)
      const params = new URLSearchParams({
        group_id: groupId,
      })
      const res = await axios.post('/api/principal/' + id + '/detach_group', params)
      console.log(res)
      fetchPrincipal()
    })

    const onDetachRole = ((roleId: string) => {
      detachRole(roleId)
    })

    const detachRole = (async (roleId: string) => {
      console.log(roleId)
      const params = new URLSearchParams({
        role_id: roleId,
      })
      const res = await axios.post('/api/principal/' + id + '/detach_role', params)
      console.log(res)
      fetchPrincipal()
    })

    const deletePrincipal = (async () => {
      const res = await axios.delete('/api/principal/' + id)
      console.log(res)
      
      //window.location.href = '../principals'
      router.push('/principals')
    })

    return {
      principal,
      authTableRows,
      authTableColumns,
      standardColumns,
      apiKeyTableRows,
      groupTableRows, 
      roleTableRows,
      permissionTableRows,
      attachedGroups,
      attachedRoles,
      onDelete,
      needRefresh,
      onDetachGroup,
      onDetachRole,
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