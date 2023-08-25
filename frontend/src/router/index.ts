import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import LayoutPage from '@/pages/Layout.vue'

import HomePage from '@/pages/HomePage.vue'
import PrincipalListPage from '@/pages/PrincipalListPage.vue'
import PrincipalPage from '@/pages/PrincipalPage.vue'
import MappingRuleListPage from '@/pages/MappingRuleListPage.vue'
import MappingRulePage from '@/pages/MappingRulePage.vue'
import GroupListPage from '@/pages/GroupListPage.vue'
import GroupPage from '@/pages/GroupPage.vue'
import RoleListPage from '@/pages/RoleListPage.vue'
import RolePage from '@/pages/RolePage.vue'
import PermissionListPage from '@/pages/PermissionListPage.vue'
import PermissionPage from '@/pages/PermissionPage.vue'


const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: LayoutPage,
    redirect: '/home',
    children: [
      {
        path: 'home',
        name: 'Home',
        component: HomePage,
      },
      {
        path: 'principals',
        name: 'Principals',
        component: PrincipalListPage,
      },
      {
        path: 'principal/:id',
        name: 'Principal',
        component: PrincipalPage,
      },
      {
        path: 'mapping_rules',
        name: 'MappingRuleList',
        component: MappingRuleListPage,
      },
      {
        path: 'mapping_rule/:id',
        name: 'MappingRule',
        component: MappingRulePage,
      },
      {
        path: 'groups',
        name: 'GroupList',
        component: GroupListPage,
      },
      {
        path: 'group/:id',
        name: 'Group',
        component: GroupPage,
      },
      {
        path: 'roles',
        name: 'Roles',
        component: RoleListPage,
      },
      {
        path: 'role/:id',
        name: 'Role',
        component: RolePage,
      },
      {
        path: 'permissions',
        name: 'PermissionList',
        component: PermissionListPage,
      },
      {
        path: 'permission/:id',
        name: 'Permission',
        component: PermissionPage,
      },
    ]
  },
]

const curSrc = document.currentScript as HTMLScriptElement
const baseUrl = new URL(curSrc.src)
const base = baseUrl.pathname.substring(0, baseUrl.pathname.indexOf('/js/'))

const router = createRouter({
  history: createWebHistory(base),
  routes
})

export default router
