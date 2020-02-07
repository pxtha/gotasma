import dayjs from 'dayjs'
export default {
    exceptionDays: (state) => {
        let arr = []
        state.exceptions.forEach(exception => {
            if (exception.date[1] - exception.date[0] === 86400000) {
                arr.push(exception.date[0])
            } else {
                do {
                    arr.push(exception.date[0])
                    exception.date[0] += 86400000
                } while (exception.date[0] < exception.date[1])
            }
        })
        state.exceptionDays = arr
        return arr
    },
    getResourceOfProject: (state) => {
        return state.resources.filter(resource => { return state.project.users.indexOf(resource.id) >= 0 })
    },
    availableResources: (state) => {
        return state.resources.filter(resource => { return state.project.users.indexOf(resource.id) === -1 })
    },
    getStatus(state) {
        return dayjs(state.project.updateDate).format('DD-MM-YYYY hh:mm:ss')
    },
    getHistoryById: (state) => (id) => {
        for (let i = 0; i < state.history.length; i++) {
            if (id === state.history[i].id) {
                return state.history[i]
            }
        }
    }
}