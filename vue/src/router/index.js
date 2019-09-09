import Vue from 'vue'
import Router from 'vue-router'
import PrincipalList from '@/components/PrincipalList'
import PrincipalSingle from '@/components/PrincipalSingle'
import RoleList from '@/components/RoleList'
import Top from '@/components/Top'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: '/guardmech/admin/',
  routes: [
    {
      path: '/',
      name: 'Top',
      component: Top
    },
    {
      path: '/principals',
      name: 'Principals',
      component: PrincipalList
    },
    {
      path: '/principal/:id',
      name: 'PrincipalSingle',
      component: PrincipalSingle
    },
    {
      path: '/roles',
      name: 'Roles',
      component: RoleList
    }
  ]
})
