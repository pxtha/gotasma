import * as Services from '../services'

export default {
    // Resources actions
    getResources({ commit }) {
        Services.getResources()
            .then((response) => {
                commit('GET_RESOURCES', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    addResource({ commit }, payload) {
        Services.addResource(payload)
            .then((response) => {
                commit('ADD_RESOURCE')
            })
            .catch(error => {
                console.log(error)
            })
    },
    deleteResource({ commit }, payload) {
        Services.deleteResource(payload)
            .then((response) => {
                commit('DELETE_RESOURCE')
            })
            .catch(error => {
                console.log(error)
            })
    },
    editResource({ commit }, payload) {
        Services.editResource(payload)
            .then((response) => {
                commit('EDIT_RESOURCE')
            })
            .catch(error => {
                console.log(error)
            })
    },
    addProjectToResource({ commit }, payload) {
        Services.addProjectToResource(payload)
            .then((response) => {
                commit('ADD_PROJECT_TO_RESOURCE')
            })
            .catch(error => {
                console.log(error)
            })
    },
    deleteProjectofResource({ commit }, payload) {
        Services.deleteProjectofResource(payload)
            .then((response) => {
                commit('DELETE_PROJECT_TO_RESOURCE')
            })
            .catch(error => {
                console.log(error)
            })
    },
    // History actions
    addHistory({ commit }, payload) {
        Services.addHistory(payload)
            .then((response) => {
                commit('ADD_HISTORY', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    getHistoryById({ commit }, payload) {
        Services.getHistoryByID(payload)
            .then((response) => {
                commit('GET_HISTORY_BY_ID', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    // Project actions
    saveProject({ commit }, payload) {
        Services.saveProject(payload)
            .then((response) => {
                commit('SAVE_PROJECT', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    getProjects({ commit }) {
        Services.getProjects()
            .then((response) => {
                commit('GET_PROJECTS', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    getHighlightedProjects({ commit }) {
        Services.getHighlightedProjects()
            .then((response) => {
                commit('GET_HIGHLIGHTED_PROJECTS', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    highlightProjects({ commit }, payload) {
        Services.highlightProjects(payload)
            .then((response) => {
                commit('HIGHLIGHT_PROJECT', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    getProjectById({ commit }, payload) {
        Services.getProjectByID(payload)
            .then((response) => {
                commit('GET_PROJECT_BY_ID', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    addProject({ commit }, payload) {
        Services.addProject(payload)
            .then((response) => {
                commit('ADD_PROJECT', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    deleteProject({ commit }, payload) {
        Services.deleteProject(payload)
            .then((response) => {
                commit('DELETE_PROJECT', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    addResourceToProject({ commit }, payload) {
        Services.addResourceToProject(payload)
            .then((response) => {
                commit('ADD_USER_TO_PROJECT', payload)
            })
            .catch(error => {
                console.log(error)
            })
    },
    deleteUserToProject({ commit }, payload) {
        Services.deleteUserToProject(payload)
            .then((response) => {
                commit('DELETE_USER_TO_PROJECT', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    editProject({ commit }, payload) {
        Services.editProject(payload)
            .then((response) => {
                commit('EDIT_PROJECT', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    // Exception actions
    getExceptions({ commit }) {
        Services.getExceptions()
            .then((response) => {
                commit('GET_EXCEPTIONS', response)
            })
            .catch(error => {
                console.log(error)
            })
    },
    addExceptions({ commit }, payload) {
        Services.addExceptions(payload)
            .then((response) => {
                commit('ADD_EXCEPTION')
            })
            .catch(error => {
                console.log(error)
            })
    },
    deleteExceptions({ commit }, payload) {
        Services.deleteExceptions(payload)
            .then((response) => {
                commit('DELETE_EXCEPTION')
            })
            .catch(error => {
                console.log(error)
            })
    },
    // LOCAL ACTIONS
    addTask({ commit }, payload) {
        commit('addTask', payload)
    },
    addSumTask({ commit }, payload) {
        commit('addSumTask', payload)
    },
    addMilestone({ commit }, payload) {
        commit('addMilestone', payload)
    },
    breakTask({ commit }, payload) {
        commit('breakTask', payload)
    },
    deleteThisTask({ commit }, payload) {
        commit('deleteThisTask', payload)
    },
    assignMember({ commit }, payload) {
        commit('assignMember', payload)
    },
    editTask({ commit }, payload) {
        commit('editTask', payload)
    }
}