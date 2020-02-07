import NetworkHelper from './NetworkHelper'

export const getProjects = async() => {
    return await NetworkHelper.requestGet('/projects')
}

export const addProject = async(project) => {
    console.log(project)
    return await NetworkHelper.requestPost('/projects', project)
}

export const deleteProject = async(id) => {
    return await NetworkHelper.requestDelete('/projects/' + id)
}

export const editProject = async(project) => {
    return await NetworkHelper.requestPatch('/projects/' + project.id, { name: project.name })
}

export const getProjectByID = async(id) => {
    return await NetworkHelper.requestGet('/projects/' + id)
}

export const highlightProjects = async(payload) => {
    return await NetworkHelper.requestPatch('/projects/' + payload.id, { highlighted: payload.highlighted })
}

export const getHighlightedProjects = async() => {
    return await NetworkHelper.requestGet('/projects?highlighted=true')
}

export const addResourceToProject = async(payload) => {
    return await NetworkHelper.requestPatch('/projects/' + payload.id, { users: payload.newInfo })
}

export const saveProject = async(payload) => {

    let d = new Date()
    let n = d.valueOf()
    
    return await NetworkHelper.requestPut('/projects/' + payload.id, {
        name: payload.name,
        tasks: payload.tasks,
        updateDate: n,
        users: payload.users
    })
    
}

export const deleteUserToProject = async(payload) => {
    return await NetworkHelper.requestPatch('/projects/' + payload.id, { users: payload.info })
}