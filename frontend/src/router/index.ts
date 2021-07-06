import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Layout from '../views/Layout.vue'

import Dashboard from '@/views/Dashboard.vue'
import PrincipalList from '@/views/PrincipalList.vue'
import Principal from '@/views/Principal.vue'
import MappingRuleList from '@/views/MappingRuleList.vue'
import MappingRule from '@/views/MappingRule.vue'
import GroupList from '@/views/GroupList.vue'
import Group from '@/views/Group.vue'
import RoleList from '@/views/RoleList.vue'
import Role from '@/views/Role.vue'
import PermissionList from '@/views/PermissionList.vue'
import Permission from '@/views/Permission.vue'


const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: Dashboard,
      },
      {
        path: 'principals',
        name: 'Principals',
        component: PrincipalList,
      },
      {
        path: 'principal/:id',
        name: 'Principal',
        component: Principal,
      },
      {
        path: 'mapping_rules',
        name: 'MappingRuleList',
        component: MappingRuleList,
      },
      {
        path: 'mapping_rules/:id',
        name: 'MappingRule',
        component: MappingRule,
      },
      {
        path: 'groups',
        name: 'GroupList',
        component: GroupList,
      },
      {
        path: 'group/:id',
        name: 'Group',
        component: Group,
      },
      {
        path: 'roles',
        name: 'Roles',
        component: RoleList,
      },
      {
        path: 'role/:id',
        name: 'Role',
        component: Role,
      },
      {
        path: 'permissions',
        name: 'PermissionList',
        component: PermissionList,
      },
      {
        path: 'permission/:id',
        name: 'Permission',
        component: Permission,
      },
    ]
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
