<template>
  <div>
    <div class="container information">
      <h2>Principal Information</h2>
      <AuthorityStatusBox :status="authorityStatus" />

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
            <template v-if="data && data.row.isAttached">
              <BButton variant="danger" @click="onDetachRole(data.row.id)">Detach</BButton>
            </template>
            <template v-else-if="data && data.row.from">
              from {{ data.row.from }}
            </template>
          </template>
        </BTable>
        <p v-else>No roles.</p>

        <h3>Having Permissions</h3>
        <BTable
          v-if="permissionTableRows.length"
          :data="permissionTableRows"
          :columns="standardColumns"
        >
          <template #cell(action)="data">
            <template v-if="data && data.row.from">
              from {{ data.row.from }}
            </template>
          </template>
        </BTable>
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
import { ref, computed, watch, defineComponent, SetupContext } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { PrincipalPayload, Group, Role } from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable from '@/components/bootstrap/BTable.vue'
import AttachGroupModal from '@/components/modals/AttachGroupModal.vue'
import AttachRoleModal from '@/components/modals/AttachRoleModal.vue'
import DestructionModal from '@/components/modals/DestructionModal.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'
import { BTableRow } from '@/types/bootstrap'

export default defineComponent({
  components: {
    BButton,
    BTable,
    AttachGroupModal,
    AttachRoleModal,
    DestructionModal,
    AuthorityStatusBox,
  },
  setup(_, context: SetupContext) {
    const router = useRouter()
    const id = router.currentRoute.value.params['id']

    const {
      authorityStatus,
      authorityLoadCompleted,
      canWrite,
      canRead,
    } = useUserAuthority()

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
      assignPrincipal(payload)
    })

    const assignPrincipal = ((payload: PrincipalPayload) => {
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
      roleTableRows.value = payload.having_roles.map(x => {
        if (payload.attached_roles.find(a => a.id === x.id)) {
          return {
            ...x,
            isAttached: true,
          }
        }
        // groupにある
        console.log(payload.groups[0].attached_roles)
        const groups = payload.groups.filter(g => g.attached_roles.find(r => r.id === x.id) ? true : false)
        console.log(groups)

        return {
          ...x,
          from: groups.map(g => g.name).join(),
        }
      })
      permissionTableRows.value = payload.having_permissions.map(x => {
        const roles = payload.having_roles.filter(r => r.attached_permissions.find(p => p.id === x.id) ? true : false)

        return {
          ...x,
          from: roles.map(r => r.name).join(),
        }
      })

      console.log(roleTableRows.value)
    })

    watch(authorityLoadCompleted, (val) => {
      if (val) {
        if (canRead.value) {
          fetchPrincipal()
        }
      }
    })

    const attachedGroups = computed<Group[]>(() => principal.value ? principal.value.groups : [])
    const attachedRoles = computed<Role[]>(() => principal.value ? principal.value.attached_roles : [])

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
      const params = new URLSearchParams({
        group_id: groupId,
      })
      const res = await axios.post('/api/principal/' + id + '/detach_group', params)
      console.log(res)

      assignPrincipal(res.data.principal)
    })

    const onDetachRole = ((roleId: string) => {
      detachRole(roleId)
    })

    const detachRole = (async (roleId: string) => {
      const params = new URLSearchParams({
        role_id: roleId,
      })
      const res = await axios.post('/api/principal/' + id + '/detach_role', params)
      console.log(res)

      assignPrincipal(res.data.principal)
    })

    const deletePrincipal = (async () => {
      const res = await axios.delete('/api/principal/' + id)
      console.log(res)
      // TODO check
      
      router.push('/principals')
    })

    return {
      authorityStatus,
      canWrite,
      canRead,
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
      needRefresh,
      onDelete,
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