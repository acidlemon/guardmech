import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
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
      path: '/users',
      name: 'Users',
      component: HelloWorld
    }
  ]
})
