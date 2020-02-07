import NetworkHelper from './NetworkHelper'

export const getResources = async() => {
    return await NetworkHelper.requestGet('/resources')
}

export const addResource = async(resource) => {
    return await NetworkHelper.requestPost('/resources', resource)
}

export const deleteResource = async(id) => {
    return await NetworkHelper.requestDelete('/resources/' + id)
}

export const editResource = async(resource) => {
    return await NetworkHelper.requestPut('/resources/' + resource.id, resource)
}

export const addProjectToResource = async(payload) => {
    return await NetworkHelper.requestPatch('/resources/' + payload.id, { projects: payload.projects })
}

export const deleteProjectofResource = async(payload) => {
    return await NetworkHelper.requestPatch('/resources/' + payload.id, { projects: payload.info })
}