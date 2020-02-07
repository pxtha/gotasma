export default {
    // =====================RESOURCES================================
    GET_RESOURCES(state, resources) {
        state.resources = resources
    },
    ADD_RESOURCE() {
        console.log('add resource done')
    },
    DELETE_RESOURCE() {
        console.log('delete resource done')
    },
    EDIT_RESOURCE() {
        console.log('edit resource done')
    },
    ADD_PROJECT_TO_RESOURCE() {
        console.log('add project to resource')
    },
    DELETE_PROJECT_TO_RESOURCE() {
        console.log('delete project to resource')
    },
    // === HIstory ===
    ADD_HISTORY() {
        console.log('History updated')
    },
    GET_HISTORY_BY_ID(state, history) {
        state.history = history
    },
    // ==========================PROJECTs=============================
    SAVE_PROJECT() {
        console.log('Saved')
    },
    GET_PROJECTS(state, projects) {
        state.projects = projects
    },
    GET_HIGHLIGHTED_PROJECTS(state, projects) {
        state.highlightedProjects = projects
    },
    HIGHLIGHT_PROJECT() {
        console.log('highlighted one project')
    },
    ADD_PROJECT() {
        console.log('added one project')
    },
    DELETE_PROJECT() {
        console.log('deleted one project')
    },
    GET_PROJECT_BY_ID(state, project) {
        state.project = project
    },
    ADD_USER_TO_PROJECT(state, payload) {
        console.log(payload)
    },
    DELETE_USER_TO_PROJECT() {
        console.log('delete one user to project')
    },
    EDIT_PROJECT() {
        console.log('project edited')
    },
    // =============================EXCEPTIONS===========================
    GET_EXCEPTIONS(state, exceptions) {
        state.exceptions = exceptions
    },
    ADD_EXCEPTION() {
        console.log('add exception done')
    },
    DELETE_EXCEPTION() {
        console.log('delete exception done')
    },
    // ===========local action=====
    addTask: (state, newTaskInfo) => {
        state.project.tasks.push({
            parentId: newTaskInfo.parentId,
            id: newTaskInfo.id,
            label: newTaskInfo.label,
            start: (newTaskInfo.start).valueOf(),
            duration: newTaskInfo.duration * 86400000,
            type: newTaskInfo.type,
            parents: newTaskInfo.parents
        })
    },
    addSumTask: (state, newTaskInfo) => {
        state.project.tasks.push({
            id: newTaskInfo.id,
            label: newTaskInfo.label,
            start: (newTaskInfo.start).valueOf(),
            duration: newTaskInfo.duration * 86400000,
            type: newTaskInfo.type
        })
    },
    addMilestone: (state, newTaskInfo) => {
        state.project.tasks.push({
            id: newTaskInfo.id,
            parentId: newTaskInfo.parentId,
            label: newTaskInfo.label,
            start: (newTaskInfo.start).valueOf(),
            duration: 86400000,
            type: 'milestone'
        })
    },
    breakTask: (state, breakTaskInfo) => {
        state.project.tasks.splice(state.project.tasks.findIndex(task => task.id === breakTaskInfo.adjacentId) + 1, 0, {
            parentId: breakTaskInfo.parentId,
            id: breakTaskInfo.id,
            label: breakTaskInfo.label,
            start: breakTaskInfo.start,
            duration: breakTaskInfo.duration,
            type: breakTaskInfo.type,
            effort: breakTaskInfo.effort
        })
    },
    deleteThisTask(state, idTaskDelete) {
        for (let i = 0; i < state.project.tasks.length; i++) {
            if (state.project.tasks[i].id === idTaskDelete) {
                if (state.project.tasks[i].children.length === 0) {
                    state.project.tasks.splice(state.project.tasks.findIndex(deleteTask => deleteTask.id === idTaskDelete), 1)
                } else {
                    state.project.tasks[i].allChildren.forEach(child => {
                        state.project.tasks.splice(state.project.tasks.findIndex(deleteTask => deleteTask.id === child), 1)
                    })
                    state.project.tasks.splice(state.project.tasks.findIndex(deleteTask => deleteTask.id === idTaskDelete), 1)
                }
                break
            }
        }
    },
    assignMember(state, userInfo) {
        let newUser = []
        userInfo.user.forEach(element => {
            newUser.push(element.id)
        })
        for (let i = 0; i < state.project.tasks.length; i++) {
            if (state.project.tasks[i].id === userInfo.taskId) {
                state.project.tasks[i].user = newUser
                break
            }
        }
    },
    editTask(state, newTaskInfo) {
        for (let i = 0; i < state.project.tasks.length; i++) {
            let task = state.project.tasks[i]
            if (task.id === newTaskInfo.id) {
                task.label = newTaskInfo.label
                task.startTime = newTaskInfo.start.valueOf()
                task.start = newTaskInfo.start.valueOf()
                task.duration = newTaskInfo.estimateDuration * 86400000
                task.effort = newTaskInfo.effort

                task.estimateDuration = task.duration
                let timeStart = new Date(task.startTime)
                let calculateTimeChart = task.startTime
                let dayofWeek = (timeStart.getDay())
                let durationDays = task.duration / 86400000
                let isHoliday = false
                for (let i = 0; i < durationDays; i++) {
                    for (let j = 0; j < state.exceptionDays.length; j++) {
                        if (calculateTimeChart === state.exceptionDays[j]) {
                            isHoliday = true
                            break
                        }
                    }
                    if (isHoliday) {
                        task.duration += 86400000
                        isHoliday = false
                        durationDays++
                        if (dayofWeek === 6) {
                            dayofWeek = 0
                        } else {
                            dayofWeek += 1
                        }
                    } else if (dayofWeek === 6) {
                        dayofWeek = 0
                        task.duration += 86400000
                        durationDays++
                    } else if (dayofWeek === 0) {
                        dayofWeek += 1
                        task.duration += 86400000
                        durationDays++
                    } else {
                        dayofWeek += 1
                    }
                    calculateTimeChart += 86400000
                }
                // task.duration = task.duration
                task.endTime = task.startTime + task.duration
            }
        }
    }
}