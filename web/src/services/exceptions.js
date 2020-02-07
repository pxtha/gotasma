import NetworkHelper from './NetworkHelper'

export const getExceptions = async() => {
    return await NetworkHelper.requestGet('/exceptions')
}

export const addExceptions = async(exception) => {
    return await NetworkHelper.requestPost('/exceptions', exception)
}

export const deleteExceptions = async(id) => {
    return await NetworkHelper.requestDelete('/exceptions/' + id)
}