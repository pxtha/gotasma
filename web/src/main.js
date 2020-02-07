// Import ES6 Promise
import 'es6-promise/auto'

// Import System requirements
import Vue from 'vue'
import VueRouter from 'vue-router'
import VueDateNow from 'vue-date-now'
import VModal from 'vue-js-modal'
import VeeValidate from 'vee-validate'
import routes from './routes'
import store from './store'
// Import Views - Top level
import AppView from './components/App.vue'
import { count, momentNormalDate, momentDetailDate } from './Utils/filter'

// Routing logic
var router = new VueRouter({
    routes: routes,
    mode: 'history',
    scrollBehavior: function(to, from, savedPosition) {
        return savedPosition || { x: 0, y: 0 }
    }
})

Vue.use(VueRouter)
Vue.use(VeeValidate)
Vue.use(VueDateNow)
Vue.use(VModal, { dialog: true, dynamic: true })

Vue.filter('count', count)
Vue.filter('momentNormalDate', momentNormalDate)
Vue.filter('momentDetailDate', momentDetailDate)

export const EventBus = new Vue()

new Vue({
    router,
    store,
    render: h => h(AppView)
}).$mount('#root')

// Start out app!
// eslint-disable-next-line no-new
// change this. demo
window.bugsnagClient = window.bugsnag('02fe1c2caaf5874c50b6ee19534f5932')
window.bugsnagClient.use(window.bugsnag__vue(Vue))