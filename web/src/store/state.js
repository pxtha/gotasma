export default {
    adminInfo: {
        id: 'admin',
        fullName: 'Coffee Cat',
        email: 'coffaacat@gmail.com',
        country: 'Heaven',
        avatar: '/static/img/stock/giphy.gif',
        aboutAvatar: '/static/img/stock/giphy.gif',
        aboutBackground: '/static/img/about.jpg',
        role: 'Task Planning'
    },
    // gantt chart header default options
    headerOptions: {
        title: {
            label: 'Project',
            html: false
        },
        locale: {
            name: 'en',
            Now: 'Now',
            'X-Scale': 'Width',
            'Y-Scale': 'Height',
            'Task list width': 'Task width',
            'Before/After': 'Expand',
            'Display task list': 'Hide Task'
        }
    },
    resources: [],
    users: [],
    projects: [],
    project: {},
    exceptions: [],
    highlightedProjects: [],
    exceptionDays: [],
    chartText: false,
    history: []
}